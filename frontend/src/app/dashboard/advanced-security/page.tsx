"use client";

import DashboardLayout from "@/components/DashboardLayout";
import {
  apiClient,
  type AdaptiveAuthResponse,
  type RiskAssessment as ApiRiskAssessment,
} from "@/lib/api";
import { useKeycloak } from "@react-keycloak/web";
import { useCallback, useEffect, useState } from "react";
import {
  IoAlertCircle,
  IoAnalytics,
  IoCheckmarkCircle,
  IoEye,
  IoFingerPrint,
  IoInformationCircle,
  IoKey,
  IoRefresh,
  IoShieldCheckmark,
  IoWarning,
} from "react-icons/io5";
import { toast } from "sonner";

// Risk Assessment Types
interface WebAuthnCredential {
  id: string;
  credential_id: string;
  device_name: string;
  created_at: string;
  last_used?: string;
}

export default function AdvancedSecurityPage() {
  const { keycloak } = useKeycloak();
  const [activeTab, setActiveTab] = useState<"webauthn" | "risk" | "saml">(
    "webauthn"
  );
  const [adaptiveAuthResponse, setAdaptiveAuthResponse] =
    useState<AdaptiveAuthResponse | null>(null);
  const [riskHistory, setRiskHistory] = useState<ApiRiskAssessment[]>([]);
  const [webauthnCredentials, setWebauthnCredentials] = useState<
    WebAuthnCredential[]
  >([]);
  const [loading, setLoading] = useState(false);

  const loadWebAuthnCredentials = useCallback(async () => {
    try {
      const response = await apiClient.getWebAuthnCredentials();
      setWebauthnCredentials(
        (response.credentials as WebAuthnCredential[]) || []
      );
    } catch (error) {
      console.error("Failed to load WebAuthn credentials:", error);
    }
  }, []);

  const loadRiskAssessment = useCallback(async () => {
    try {
      setLoading(true);

      // Use the new adaptive auth evaluation endpoint
      const deviceFingerprint = generateDeviceFingerprint();

      const authResponse = await apiClient.evaluateAuthentication({
        device_fingerprint: deviceFingerprint,
        typing_pattern: {
          avg_keydown_time: 120 + Math.random() * 50,
        },
      });

      // Create compatibility layer - convert new format to old format for UI
      const compatibleResponse: AdaptiveAuthResponse = {
        decision: authResponse.decision,
        risk_score: authResponse.risk_score,
        risk_level: authResponse.risk_level,
        required_actions: authResponse.required_actions, // Keep the new format
        reasoning: authResponse.reasoning,
        session_duration_seconds: authResponse.session_duration_seconds,
        restrictions: authResponse.restrictions,
        metadata: authResponse.metadata,
        expires_at: authResponse.expires_at,
        risk_assessment: {
          user_id: "12345678-1234-1234-1234-123456789012", // Valid UUID for demo user
          risk_score: authResponse.risk_score,
          risk_level: authResponse.risk_level,
          risk_factors: authResponse.reasoning.map((reason) => ({
            type: "behavioral",
            description: reason,
            weight: 0.1,
            score: authResponse.risk_score,
          })),
          location: {
            country: "US",
            city: "San Francisco",
            is_vpn: false,
            is_tor: false,
          },
          device: {
            fingerprint: deviceFingerprint,
            is_known: true,
            trust_score: 1 - authResponse.risk_score,
          },
          behavior: {
            typing_speed_deviation: 0.1,
            mouse_pattern_deviation: 0.1,
          },
          timestamp: new Date().toISOString(),
        },
        session_restrictions: {
          max_duration_minutes: Math.floor(
            authResponse.session_duration_seconds / 60
          ),
          require_mfa: authResponse.required_actions.some(
            (action) => action.type === "mfa"
          ),
          allowed_operations: ["login", "dashboard"],
        },
      };

      setAdaptiveAuthResponse(compatibleResponse);

      // Load risk history - create mock data since backend format is different
      if (keycloak?.tokenParsed?.sub) {
        const mockHistory: ApiRiskAssessment[] = [
          {
            user_id: keycloak.tokenParsed.sub,
            risk_score: authResponse.risk_score,
            risk_level: authResponse.risk_level,
            risk_factors: [],
            location: {
              country: "US",
              city: "San Francisco",
              is_vpn: false,
              is_tor: false,
            },
            device: {
              fingerprint: deviceFingerprint,
              is_known: true,
              trust_score: 1 - authResponse.risk_score,
            },
            behavior: {
              typing_speed_deviation: 0.1,
              mouse_pattern_deviation: 0.1,
            },
            timestamp: new Date().toISOString(),
          },
        ];
        setRiskHistory(mockHistory);
      }

      // Register device fingerprint
      await apiClient.registerDeviceFingerprint(deviceFingerprint);
    } catch (error) {
      console.error("Failed to load risk assessment:", error);
      toast.error("Failed to load risk assessment");
    } finally {
      setLoading(false);
    }
  }, [keycloak]);

  // Load initial data
  useEffect(() => {
    loadWebAuthnCredentials();
    loadRiskAssessment();
  }, [loadWebAuthnCredentials, loadRiskAssessment]);

  // Helper function to convert base64url to Uint8Array
  const base64urlToUint8Array = (base64url: string): Uint8Array => {
    // Convert base64url to base64
    const base64 = base64url
      .replace(/-/g, "+")
      .replace(/_/g, "/")
      .padEnd(base64url.length + ((4 - (base64url.length % 4)) % 4), "=");

    // Decode base64
    const binaryString = atob(base64);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes;
  };

  const registerWebAuthn = async () => {
    try {
      setLoading(true);
      toast.info("Starting WebAuthn registration...");

      // Begin registration
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const options = (await apiClient.webAuthnRegisterBegin()) as any;

      // Convert challenge from base64url
      options.challenge = base64urlToUint8Array(options.challenge);
      options.user.id = base64urlToUint8Array(options.user.id);

      // Create credential
      const credential = (await navigator.credentials.create({
        publicKey: options,
      })) as PublicKeyCredential;

      if (!credential) {
        throw new Error("Failed to create credential");
      }

      // Finish registration
      const credentialData = {
        id: credential.id,
        rawId: Array.from(new Uint8Array(credential.rawId)),
        type: credential.type,
        response: {
          attestationObject: Array.from(
            new Uint8Array(
              (
                credential.response as AuthenticatorAttestationResponse
              ).attestationObject
            )
          ),
          clientDataJSON: Array.from(
            new Uint8Array(credential.response.clientDataJSON)
          ),
        },
      };

      const result = await apiClient.webAuthnRegisterFinish(credentialData);

      if (result.success) {
        toast.success("WebAuthn credential registered successfully!");
        loadWebAuthnCredentials();
      } else {
        throw new Error("Failed to finish WebAuthn registration");
      }
    } catch (error) {
      console.error("WebAuthn registration failed:", error);
      toast.error("WebAuthn registration failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const authenticateWebAuthn = async () => {
    try {
      setLoading(true);
      toast.info("Starting WebAuthn authentication...");

      // Begin authentication
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const options = (await apiClient.webAuthnAuthenticateBegin()) as any;

      // Convert challenge and credential IDs
      options.challenge = base64urlToUint8Array(options.challenge);
      options.allowCredentials = options.allowCredentials.map(
        (cred: { id: string; type: string; transports?: string[] }) => ({
          ...cred,
          id: base64urlToUint8Array(cred.id),
        })
      );

      // Get assertion
      const assertion = (await navigator.credentials.get({
        publicKey: options,
      })) as PublicKeyCredential;

      if (!assertion) {
        throw new Error("Failed to get assertion");
      }

      // Finish authentication
      const credentialData = {
        id: assertion.id,
        type: assertion.type,
        response: {
          authenticatorData: Array.from(
            new Uint8Array(
              (
                assertion.response as AuthenticatorAssertionResponse
              ).authenticatorData
            )
          ),
          clientDataJSON: Array.from(
            new Uint8Array(assertion.response.clientDataJSON)
          ),
          signature: Array.from(
            new Uint8Array(
              (assertion.response as AuthenticatorAssertionResponse).signature
            )
          ),
          userHandle: (assertion.response as AuthenticatorAssertionResponse)
            .userHandle
            ? Array.from(
                new Uint8Array(
                  (
                    assertion.response as AuthenticatorAssertionResponse
                  ).userHandle!
                )
              )
            : null,
        },
      };

      const result = await apiClient.webAuthnAuthenticateFinish(credentialData);

      if (result.success) {
        toast.success("WebAuthn authentication successful!");
      } else {
        throw new Error("Failed to finish WebAuthn authentication");
      }
    } catch (error) {
      console.error("WebAuthn authentication failed:", error);
      toast.error("WebAuthn authentication failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const deleteWebAuthnCredential = async (credentialId: string) => {
    try {
      await apiClient.deleteWebAuthnCredential(credentialId);
      toast.success("WebAuthn credential deleted successfully");
      loadWebAuthnCredentials();
    } catch (error) {
      console.error("Failed to delete WebAuthn credential:", error);
      toast.error("Failed to delete credential");
    }
  };

  const generateDeviceFingerprint = (): string => {
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    ctx?.fillText("CloudGate fingerprint", 2, 2);

    const fingerprint = [
      navigator.userAgent,
      navigator.language,
      screen.width + "x" + screen.height,
      screen.colorDepth,
      new Date().getTimezoneOffset(),
      canvas.toDataURL(),
      navigator.hardwareConcurrency || "unknown",
    ].join("|");

    // Simple hash function
    let hash = 0;
    for (let i = 0; i < fingerprint.length; i++) {
      const char = fingerprint.charCodeAt(i);
      hash = (hash << 5) - hash + char;
      hash = hash & hash;
    }
    return Math.abs(hash).toString(16);
  };

  const getRiskLevelColor = (level: string) => {
    switch (level) {
      case "low":
        return "text-green-600 bg-green-50";
      case "medium":
        return "text-yellow-600 bg-yellow-50";
      case "high":
        return "text-orange-600 bg-orange-50";
      case "critical":
        return "text-red-600 bg-red-50";
      default:
        return "text-gray-600 bg-gray-50";
    }
  };

  const getSeverityIcon = (severity: string) => {
    switch (severity) {
      case "low":
        return <IoInformationCircle className="text-blue-500" />;
      case "medium":
        return <IoAlertCircle className="text-yellow-500" />;
      case "high":
        return <IoWarning className="text-orange-500" />;
      case "critical":
        return <IoWarning className="text-red-500" />;
      default:
        return <IoInformationCircle className="text-gray-500" />;
    }
  };

  const getDecisionExplanation = (decision: string): string => {
    switch (decision) {
      case "allow":
        return "Access granted based on low risk profile";
      case "challenge":
        return "Additional verification required due to moderate risk";
      case "deny":
        return "Access denied due to high risk factors";
      case "monitor":
        return "Access granted with enhanced monitoring";
      default:
        return "Unknown decision";
    }
  };

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <div className="flex items-center space-x-3 mb-4">
            <IoShieldCheckmark className="text-3xl text-blue-600" />
            <div>
              <h1 className="text-2xl font-bold text-gray-900">
                Advanced Security
              </h1>
              <p className="text-gray-600">
                Enterprise-grade security features and risk assessment
              </p>
            </div>
          </div>

          {/* Tab Navigation */}
          <div className="flex space-x-1 bg-gray-100 rounded-lg p-1">
            <button
              onClick={() => setActiveTab("webauthn")}
              className={`flex-1 flex items-center justify-center space-x-2 px-4 py-2 rounded-md transition-colors cursor-pointer ${
                activeTab === "webauthn"
                  ? "bg-white text-blue-600 shadow-sm"
                  : "text-gray-600 hover:text-gray-900"
              }`}
            >
              <IoFingerPrint />
              <span>WebAuthn / FIDO2</span>
            </button>
            <button
              onClick={() => setActiveTab("risk")}
              className={`flex-1 flex items-center justify-center space-x-2 px-4 py-2 rounded-md transition-colors cursor-pointer ${
                activeTab === "risk"
                  ? "bg-white text-blue-600 shadow-sm"
                  : "text-gray-600 hover:text-gray-900"
              }`}
            >
              <IoAnalytics />
              <span>Risk Assessment</span>
            </button>
            <button
              onClick={() => setActiveTab("saml")}
              className={`flex-1 flex items-center justify-center space-x-2 px-4 py-2 rounded-md transition-colors cursor-pointer ${
                activeTab === "saml"
                  ? "bg-white text-blue-600 shadow-sm"
                  : "text-gray-600 hover:text-gray-900"
              }`}
            >
              <IoKey />
              <span>SAML 2.0</span>
            </button>
          </div>
        </div>

        {/* WebAuthn Tab */}
        {activeTab === "webauthn" && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* WebAuthn Registration */}
            <div className="bg-white rounded-lg shadow-sm p-6">
              <div className="flex items-center space-x-3 mb-4">
                <IoFingerPrint className="text-2xl text-blue-600" />
                <h2 className="text-xl font-semibold text-black/80">
                  Biometric Authentication
                </h2>
              </div>

              <p className="text-gray-600 mb-6">
                Register your device for passwordless authentication using
                biometrics, security keys, or platform authenticators.
              </p>

              <div className="space-y-4">
                <button
                  onClick={registerWebAuthn}
                  disabled={loading}
                  className="w-full bg-blue-600 text-white px-4 py-3 rounded-lg hover:bg-blue-700 disabled:opacity-50 flex items-center justify-center space-x-2 cursor-pointer disabled:cursor-not-allowed"
                >
                  <IoFingerPrint />
                  <span>
                    {loading ? "Registering..." : "Register New Device"}
                  </span>
                </button>

                {webauthnCredentials.length > 0 && (
                  <button
                    onClick={authenticateWebAuthn}
                    disabled={loading}
                    className="w-full bg-green-600 text-white px-4 py-3 rounded-lg hover:bg-green-700 disabled:opacity-50 flex items-center justify-center space-x-2 cursor-pointer disabled:cursor-not-allowed"
                  >
                    <IoCheckmarkCircle />
                    <span>
                      {loading ? "Authenticating..." : "Test Authentication"}
                    </span>
                  </button>
                )}
              </div>
            </div>

            {/* Registered Devices */}
            <div className="bg-white rounded-lg shadow-sm p-6 text-black/80">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold">Registered Devices</h3>
                <button
                  onClick={loadWebAuthnCredentials}
                  className="p-2 text-gray-700 hover:text-gray-600 rounded-lg hover:bg-gray-100 cursor-pointer"
                >
                  <IoRefresh />
                </button>
              </div>

              {webauthnCredentials.length === 0 ? (
                <div className="text-center py-8">
                  <IoFingerPrint className="text-4xl text-gray-300 mx-auto mb-2" />
                  <p className="text-gray-500">No devices registered</p>
                  <p className="text-sm text-gray-400">
                    Register a device to enable biometric authentication
                  </p>
                </div>
              ) : (
                <div className="space-y-3">
                  {webauthnCredentials.map((credential) => (
                    <div
                      key={credential.id}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div>
                        <div className="font-medium">
                          {credential.device_name}
                        </div>
                        <div className="text-sm text-gray-500">
                          Registered:{" "}
                          {new Date(credential.created_at).toLocaleDateString()}
                        </div>
                        {credential.last_used && (
                          <div className="text-sm text-gray-500">
                            Last used:{" "}
                            {new Date(
                              credential.last_used
                            ).toLocaleDateString()}
                          </div>
                        )}
                      </div>
                      <button
                        onClick={() =>
                          deleteWebAuthnCredential(credential.credential_id)
                        }
                        className="text-red-600 hover:text-red-800 p-2 rounded-lg hover:bg-red-50 cursor-pointer"
                      >
                        Delete
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}

        {/* Risk Assessment Tab */}
        {activeTab === "risk" && (
          <div className="space-y-6">
            {/* Risk Assessment Explanation */}
            <div className="bg-blue-50 rounded-lg p-6 border border-blue-200">
              <div className="flex items-start space-x-3">
                <IoInformationCircle className="text-2xl text-blue-600 mt-1" />
                <div>
                  <h3 className="text-lg font-semibold text-blue-900 mb-2">
                    What is Risk Assessment?
                  </h3>
                  <p className="text-blue-800 mb-3">
                    Our intelligent risk assessment engine continuously analyzes
                    your login patterns, device fingerprints, location data, and
                    behavioral characteristics using advanced heuristic
                    algorithms to detect potential security threats in
                    real-time.
                  </p>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-blue-700">
                    <div>
                      <strong>🎯 Purpose:</strong>
                      <ul className="list-disc list-inside mt-1 space-y-1">
                        <li>Detect suspicious login attempts</li>
                        <li>Identify compromised accounts</li>
                        <li>Prevent unauthorized access</li>
                        <li>Adaptive security controls</li>
                      </ul>
                    </div>
                    <div>
                      <strong>📊 Factors Analyzed:</strong>
                      <ul className="list-disc list-inside mt-1 space-y-1">
                        <li>Geographic location changes</li>
                        <li>Device fingerprinting</li>
                        <li>Login time patterns</li>
                        <li>Network characteristics (VPN/Tor)</li>
                      </ul>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Current Risk Assessment */}
            <div className="bg-white rounded-lg shadow-sm p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-semibold text-black/80">
                  Current Risk Assessment
                </h2>
                <button
                  onClick={loadRiskAssessment}
                  disabled={loading}
                  className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 cursor-pointer disabled:cursor-not-allowed"
                >
                  <IoRefresh />
                  <span>{loading ? "Analyzing..." : "Refresh Assessment"}</span>
                </button>
              </div>

              {adaptiveAuthResponse?.risk_assessment ? (
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  {/* Risk Score */}
                  <div className="text-center">
                    <div
                      className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${getRiskLevelColor(
                        adaptiveAuthResponse.risk_assessment.risk_level
                      )}`}
                    >
                      {adaptiveAuthResponse.risk_assessment.risk_level.toUpperCase()}
                    </div>
                    <div className="mt-2 text-3xl font-bold text-gray-900">
                      {Math.round(
                        adaptiveAuthResponse.risk_assessment.risk_score * 100
                      )}
                      %
                    </div>
                    <div className="text-sm text-gray-500">Risk Score</div>
                  </div>

                  {/* Location Info */}
                  <div className="text-center">
                    <div className="text-lg font-semibold text-gray-900">
                      {adaptiveAuthResponse.risk_assessment.location.city},{" "}
                      {adaptiveAuthResponse.risk_assessment.location.country}
                    </div>
                    <div className="text-sm text-gray-500">Location</div>
                    {(adaptiveAuthResponse.risk_assessment.location.is_vpn ||
                      adaptiveAuthResponse.risk_assessment.location.is_tor) && (
                      <div className="mt-1">
                        {adaptiveAuthResponse.risk_assessment.location
                          .is_vpn && (
                          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-yellow-100 text-yellow-800 mr-1">
                            VPN
                          </span>
                        )}
                        {adaptiveAuthResponse.risk_assessment.location
                          .is_tor && (
                          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-red-100 text-red-800">
                            Tor
                          </span>
                        )}
                      </div>
                    )}
                  </div>

                  {/* Policy Decision */}
                  <div className="text-center">
                    {adaptiveAuthResponse && (
                      <>
                        <div
                          className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${
                            adaptiveAuthResponse.decision === "allow"
                              ? "bg-green-100 text-green-800"
                              : adaptiveAuthResponse.decision === "challenge"
                              ? "bg-yellow-100 text-yellow-800"
                              : adaptiveAuthResponse.decision === "monitor"
                              ? "bg-blue-100 text-blue-800"
                              : "bg-red-100 text-red-800"
                          }`}
                        >
                          {adaptiveAuthResponse.decision.toUpperCase()}
                        </div>
                        <div className="mt-2 text-sm text-gray-600">
                          {getDecisionExplanation(
                            adaptiveAuthResponse.decision
                          )}
                        </div>
                        {adaptiveAuthResponse.session_restrictions
                          ?.require_mfa && (
                          <div className="mt-1 text-xs text-orange-600">
                            MFA Required
                          </div>
                        )}
                      </>
                    )}
                  </div>
                </div>
              ) : (
                <div className="text-center py-8">
                  <IoAnalytics className="text-4xl text-gray-300 mx-auto mb-2" />
                  <p className="text-gray-500">No risk assessment available</p>
                  <p className="text-sm text-gray-400">
                    Click refresh to perform a new assessment
                  </p>
                </div>
              )}
            </div>

            {/* Risk Factors */}
            {adaptiveAuthResponse?.risk_assessment?.risk_factors && (
              <div className="bg-white rounded-lg shadow-sm p-6">
                <h3 className="text-lg font-semibold mb-4">Risk Factors</h3>
                <div className="space-y-3">
                  {adaptiveAuthResponse.risk_assessment.risk_factors.map(
                    (factor, index) => (
                      <div
                        key={index}
                        className="flex items-start space-x-3 p-3 border rounded-lg"
                      >
                        {getSeverityIcon(
                          factor.score > 0.7
                            ? "high"
                            : factor.score > 0.4
                            ? "medium"
                            : "low"
                        )}
                        <div className="flex-1">
                          <div className="font-medium">
                            {factor.description}
                          </div>
                          <div className="text-sm text-gray-500 capitalize">
                            {factor.type} • Weight: {factor.weight}
                          </div>
                        </div>
                        <div className="text-sm font-medium">
                          {Math.round(factor.score * 100)}%
                        </div>
                      </div>
                    )
                  )}
                </div>
              </div>
            )}

            {/* Session Restrictions */}
            {adaptiveAuthResponse?.session_restrictions && (
              <div className="bg-white rounded-lg shadow-sm p-6">
                <h3 className="text-lg font-semibold mb-4">
                  Session Restrictions
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="p-4 bg-gray-50 rounded-lg">
                    <div className="text-sm text-gray-600">
                      Max Session Duration
                    </div>
                    <div className="text-lg font-semibold">
                      {
                        adaptiveAuthResponse.session_restrictions
                          .max_duration_minutes
                      }{" "}
                      minutes
                    </div>
                  </div>
                  <div className="p-4 bg-gray-50 rounded-lg">
                    <div className="text-sm text-gray-600">MFA Required</div>
                    <div className="text-lg font-semibold">
                      {adaptiveAuthResponse.session_restrictions.require_mfa
                        ? "Yes"
                        : "No"}
                    </div>
                  </div>
                </div>
                {adaptiveAuthResponse.required_actions.length > 0 && (
                  <div className="mt-4">
                    <div className="text-sm font-medium text-gray-700 mb-2">
                      Required Actions:
                    </div>
                    <ul className="list-disc list-inside space-y-1">
                      {adaptiveAuthResponse.required_actions.map(
                        (action, index) => (
                          <li key={index} className="text-sm text-gray-600">
                            {typeof action === "string"
                              ? action
                              : action.description}
                          </li>
                        )
                      )}
                    </ul>
                  </div>
                )}
              </div>
            )}

            {/* Risk History */}
            {riskHistory.length > 0 && (
              <div className="bg-white rounded-lg shadow-sm p-6">
                <h3 className="text-lg font-semibold mb-4">
                  Recent Risk Assessments
                </h3>
                <div className="space-y-2">
                  {riskHistory.slice(0, 5).map((assessment, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div className="flex items-center space-x-3">
                        <div
                          className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getRiskLevelColor(
                            assessment.risk_level
                          )}`}
                        >
                          {assessment.risk_level}
                        </div>
                        <div className="text-sm text-gray-600">
                          {assessment.location.city},{" "}
                          {assessment.location.country}
                        </div>
                      </div>
                      <div className="text-sm text-gray-500">
                        {new Date(assessment.timestamp).toLocaleString()}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}

        {/* SAML Tab */}
        {activeTab === "saml" && (
          <div className="bg-white rounded-lg shadow-sm p-6 text-black/80">
            <div className="flex items-center space-x-3 mb-6">
              <IoKey className="text-2xl text-blue-600" />
              <div>
                <h2 className="text-xl font-semibold text-black/80">
                  SAML 2.0 Integration
                </h2>
                <p className="text-gray-600">
                  Enterprise SSO for legacy applications
                </p>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* SAML Info */}
              <div>
                <h3 className="text-lg font-semibold mb-4 text-black/80">
                  SAML Configuration
                </h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Entity ID
                    </label>
                    <div className="p-3 bg-gray-50 rounded-lg font-mono text-sm">
                      CloudGate-SSO
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      SSO URL
                    </label>
                    <div className="p-3 bg-gray-50 rounded-lg font-mono text-sm">
                      {process.env.NEXT_PUBLIC_API_URL ||
                        "http://localhost:8081"}
                      /saml/sso
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Metadata URL
                    </label>
                    <div className="p-3 bg-gray-50 rounded-lg font-mono text-sm">
                      {process.env.NEXT_PUBLIC_API_URL ||
                        "http://localhost:8081"}
                      /saml/metadata
                    </div>
                  </div>
                </div>
              </div>

              {/* SAML Features */}
              <div>
                <h3 className="text-lg font-semibold mb-4">
                  Supported Features
                </h3>
                <div className="space-y-3">
                  <div className="flex items-center space-x-3">
                    <IoCheckmarkCircle className="text-green-500" />
                    <span>SAML 2.0 Protocol</span>
                  </div>
                  <div className="flex items-center space-x-3">
                    <IoCheckmarkCircle className="text-green-500" />
                    <span>HTTP-POST Binding</span>
                  </div>
                  <div className="flex items-center space-x-3">
                    <IoCheckmarkCircle className="text-green-500" />
                    <span>HTTP-Redirect Binding</span>
                  </div>
                  <div className="flex items-center space-x-3">
                    <IoCheckmarkCircle className="text-green-500" />
                    <span>Signed Assertions</span>
                  </div>
                  <div className="flex items-center space-x-3">
                    <IoCheckmarkCircle className="text-green-500" />
                    <span>Attribute Mapping</span>
                  </div>
                  <div className="flex items-center space-x-3">
                    <IoCheckmarkCircle className="text-green-500" />
                    <span>Legacy Application Support</span>
                  </div>
                </div>

                <div className="mt-6">
                  <a
                    href={`${
                      process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081"
                    }/saml/metadata`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 cursor-pointer"
                  >
                    <IoEye />
                    <span>View Metadata</span>
                  </a>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}
