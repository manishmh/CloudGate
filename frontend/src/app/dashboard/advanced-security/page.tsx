"use client";

import DashboardLayout from "@/components/DashboardLayout";
import { useEffect, useState } from "react";
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
interface RiskAssessment {
  user_id: string;
  ip_address: string;
  location: {
    country: string;
    city: string;
    is_vpn: boolean;
    is_tor: boolean;
  };
  risk_score: number;
  risk_level: string;
  risk_factors: RiskFactor[];
  recommendations: string[];
  timestamp: string;
}

interface RiskFactor {
  type: string;
  description: string;
  severity: string;
  weight: number;
  score: number;
}

interface PolicyDecision {
  action: string;
  confidence: number;
  required_mfa: string[];
  explanation: string;
  session_limits: {
    max_duration_minutes: number;
    idle_timeout_minutes: number;
  };
}

interface WebAuthnCredential {
  id: string;
  credential_id: string;
  device_name: string;
  created_at: string;
  last_used?: string;
}

export default function AdvancedSecurityPage() {
  const [activeTab, setActiveTab] = useState<"webauthn" | "risk" | "saml">(
    "webauthn"
  );
  const [riskAssessment, setRiskAssessment] = useState<RiskAssessment | null>(
    null
  );
  const [policyDecision, setPolicyDecision] = useState<PolicyDecision | null>(
    null
  );
  const [webauthnCredentials, setWebauthnCredentials] = useState<
    WebAuthnCredential[]
  >([]);
  const [loading, setLoading] = useState(false);

  // Load initial data
  useEffect(() => {
    loadWebAuthnCredentials();
    loadRiskAssessment();
  }, []);

  const loadWebAuthnCredentials = async () => {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";
      const response = await fetch(`${apiUrl}/webauthn/credentials`, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      if (response.ok) {
        const data = await response.json();
        setWebauthnCredentials(data.credentials || []);
      }
    } catch (error) {
      console.error("Failed to load WebAuthn credentials:", error);
    }
  };

  const loadRiskAssessment = async () => {
    try {
      setLoading(true);

      // Perform risk assessment
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";
      const assessResponse = await fetch(`${apiUrl}/risk/assess`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({
          device_fingerprint: generateDeviceFingerprint(),
          typing_pattern: {
            avg_keydown_time: 120 + Math.random() * 50,
          },
        }),
      });

      if (assessResponse.ok) {
        const assessment = await assessResponse.json();
        setRiskAssessment(assessment);

        // Get policy decision
        const policyResponse = await fetch(`${apiUrl}/risk/policy`, {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });

        if (policyResponse.ok) {
          const policy = await policyResponse.json();
          setPolicyDecision(policy);
        }
      }

      // Load risk history
      const historyResponse = await fetch(`${apiUrl}/risk/history?limit=10`, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });

      if (historyResponse.ok) {
        const history = await historyResponse.json();
        // Risk history loaded successfully but not used in current UI
        console.log(
          "Risk history loaded:",
          history.assessments?.length || 0,
          "items"
        );
      }
    } catch (error) {
      console.error("Failed to load risk assessment:", error);
      toast.error("Failed to load risk assessment");
    } finally {
      setLoading(false);
    }
  };

  const registerWebAuthn = async () => {
    try {
      setLoading(true);
      toast.info("Starting WebAuthn registration...");

      // Begin registration
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";
      const beginResponse = await fetch(`${apiUrl}/webauthn/register/begin`, {
        method: "POST",
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });

      if (!beginResponse.ok) {
        throw new Error("Failed to begin WebAuthn registration");
      }

      const options = await beginResponse.json();

      // Convert challenge from base64url
      options.challenge = new Uint8Array(
        Buffer.from(options.challenge, "base64url")
      );
      options.user.id = new Uint8Array(
        Buffer.from(options.user.id, "base64url")
      );

      // Create credential
      const credential = (await navigator.credentials.create({
        publicKey: options,
      })) as PublicKeyCredential;

      if (!credential) {
        throw new Error("Failed to create credential");
      }

      // Finish registration
      const finishResponse = await fetch(`${apiUrl}/webauthn/register/finish`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({
          credential: {
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
          },
        }),
      });

      if (finishResponse.ok) {
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
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";
      const beginResponse = await fetch(
        `${apiUrl}/webauthn/authenticate/begin`,
        {
          method: "POST",
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        }
      );

      if (!beginResponse.ok) {
        throw new Error("Failed to begin WebAuthn authentication");
      }

      const options = await beginResponse.json();

      // Convert challenge and credential IDs
      options.challenge = new Uint8Array(
        Buffer.from(options.challenge, "base64url")
      );
      options.allowCredentials = options.allowCredentials.map(
        (cred: { id: string; type: string; transports?: string[] }) => ({
          ...cred,
          id: new Uint8Array(Buffer.from(cred.id, "base64url")),
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
      const finishResponse = await fetch(
        `${apiUrl}/webauthn/authenticate/finish`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
          body: JSON.stringify({
            credential: {
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
                    (
                      assertion.response as AuthenticatorAssertionResponse
                    ).signature
                  )
                ),
                userHandle: (
                  assertion.response as AuthenticatorAssertionResponse
                ).userHandle
                  ? Array.from(
                      new Uint8Array(
                        (
                          assertion.response as AuthenticatorAssertionResponse
                        ).userHandle!
                      )
                    )
                  : null,
              },
            },
          }),
        }
      );

      if (finishResponse.ok) {
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
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";
      const response = await fetch(
        `${apiUrl}/webauthn/credentials/${credentialId}`,
        {
          method: "DELETE",
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        }
      );

      if (response.ok) {
        toast.success("WebAuthn credential deleted successfully");
        loadWebAuthnCredentials();
      } else {
        throw new Error("Failed to delete credential");
      }
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

              {riskAssessment ? (
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  {/* Risk Score */}
                  <div className="text-center">
                    <div
                      className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${getRiskLevelColor(
                        riskAssessment.risk_level
                      )}`}
                    >
                      {riskAssessment.risk_level.toUpperCase()}
                    </div>
                    <div className="mt-2 text-3xl font-bold text-gray-900">
                      {Math.round(riskAssessment.risk_score * 100)}%
                    </div>
                    <div className="text-sm text-gray-500">Risk Score</div>
                  </div>

                  {/* Location Info */}
                  <div className="text-center">
                    <div className="text-lg font-semibold text-gray-900">
                      {riskAssessment.location.city},{" "}
                      {riskAssessment.location.country}
                    </div>
                    <div className="text-sm text-gray-500">Location</div>
                    {(riskAssessment.location.is_vpn ||
                      riskAssessment.location.is_tor) && (
                      <div className="mt-1">
                        {riskAssessment.location.is_vpn && (
                          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-yellow-100 text-yellow-800 mr-1">
                            VPN
                          </span>
                        )}
                        {riskAssessment.location.is_tor && (
                          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-red-100 text-red-800">
                            Tor
                          </span>
                        )}
                      </div>
                    )}
                  </div>

                  {/* Policy Decision */}
                  <div className="text-center">
                    {policyDecision && (
                      <>
                        <div
                          className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${
                            policyDecision.action === "allow"
                              ? "bg-green-100 text-green-800"
                              : policyDecision.action === "step_up"
                              ? "bg-yellow-100 text-yellow-800"
                              : "bg-red-100 text-red-800"
                          }`}
                        >
                          {policyDecision.action.toUpperCase()}
                        </div>
                        <div className="mt-2 text-sm text-gray-600">
                          {policyDecision.explanation}
                        </div>
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
            {riskAssessment && riskAssessment.risk_factors && (
              <div className="bg-white rounded-lg shadow-sm p-6">
                <h3 className="text-lg font-semibold mb-4">Risk Factors</h3>
                <div className="space-y-3">
                  {riskAssessment.risk_factors.map((factor, index) => (
                    <div
                      key={index}
                      className="flex items-start space-x-3 p-3 border rounded-lg"
                    >
                      {getSeverityIcon(factor.severity)}
                      <div className="flex-1">
                        <div className="font-medium">{factor.description}</div>
                        <div className="text-sm text-gray-500 capitalize">
                          {factor.type} • {factor.severity} severity • Weight:{" "}
                          {factor.weight}
                        </div>
                      </div>
                      <div className="text-sm font-medium">
                        {Math.round(factor.score * 100)}%
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Recommendations */}
            {riskAssessment && riskAssessment.recommendations && (
              <div className="bg-white rounded-lg shadow-sm p-6">
                <h3 className="text-lg font-semibold mb-4">
                  Security Recommendations
                </h3>
                <div className="space-y-2">
                  {riskAssessment.recommendations.map(
                    (recommendation, index) => (
                      <div key={index} className="flex items-start space-x-3">
                        <IoInformationCircle className="text-blue-500 mt-0.5" />
                        <span className="text-gray-700">{recommendation}</span>
                      </div>
                    )
                  )}
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
                      http://localhost:8081/saml/sso
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Metadata URL
                    </label>
                    <div className="p-3 bg-gray-50 rounded-lg font-mono text-sm">
                      http://localhost:8081/saml/metadata
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
