"use client";

import {
  apiClient,
  type MFASetupResponse,
  type MFAStatusResponse,
} from "@/lib/api";
import { useEffect, useState } from "react";
import {
  IoCopy,
  IoEye,
  IoEyeOff,
  IoShieldCheckmark,
  IoWarning,
} from "react-icons/io5";

interface MFASetupProps {
  onMFAStatusChange?: (enabled: boolean) => void;
}

export default function MFASetup({ onMFAStatusChange }: MFASetupProps) {
  const [mfaStatus, setMfaStatus] = useState<MFAStatusResponse | null>(null);
  const [setupData, setSetupData] = useState<MFASetupResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [step, setStep] = useState<"status" | "setup" | "verify" | "backup">(
    "status"
  );
  const [verificationCode, setVerificationCode] = useState("");
  const [showBackupCodes, setShowBackupCodes] = useState(false);
  const [newBackupCodes, setNewBackupCodes] = useState<string[]>([]);

  useEffect(() => {
    loadMFAStatus();
  }, []);

  const loadMFAStatus = async () => {
    try {
      setLoading(true);
      setError(null);
      const status = await apiClient.getMFAStatus();
      setMfaStatus(status);
      onMFAStatusChange?.(status.enabled);
    } catch (err) {
      setError("Failed to load MFA status");
      console.error("MFA status error:", err);
    } finally {
      setLoading(false);
    }
  };

  const startMFASetup = async () => {
    try {
      setLoading(true);
      setError(null);
      const setup = await apiClient.setupMFA();
      setSetupData(setup);
      setStep("setup");
    } catch (err) {
      setError("Failed to start MFA setup");
      console.error("MFA setup error:", err);
    } finally {
      setLoading(false);
    }
  };

  const verifyMFASetup = async () => {
    if (!verificationCode.trim()) {
      setError("Please enter a verification code");
      return;
    }

    try {
      setLoading(true);
      setError(null);
      await apiClient.verifyMFASetup(verificationCode);
      setSuccess("MFA enabled successfully!");
      setStep("backup");
      await loadMFAStatus();
    } catch (err) {
      setError("Invalid verification code. Please try again.");
      console.error("MFA verification error:", err);
    } finally {
      setLoading(false);
    }
  };

  const disableMFA = async () => {
    if (!verificationCode.trim()) {
      setError("Please enter a verification code to disable MFA");
      return;
    }

    try {
      setLoading(true);
      setError(null);
      await apiClient.disableMFA(verificationCode);
      setSuccess("MFA disabled successfully");
      setVerificationCode("");
      await loadMFAStatus();
      setStep("status");
    } catch (err) {
      setError("Invalid verification code. Please try again.");
      console.error("MFA disable error:", err);
    } finally {
      setLoading(false);
    }
  };

  const regenerateBackupCodes = async () => {
    if (!verificationCode.trim()) {
      setError("Please enter a verification code to regenerate backup codes");
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const response = await apiClient.regenerateBackupCodes(verificationCode);
      setNewBackupCodes(response.backup_codes);
      setSuccess("Backup codes regenerated successfully");
      setVerificationCode("");
      await loadMFAStatus();
    } catch (err) {
      setError("Invalid verification code. Please try again.");
      console.error("Backup codes regeneration error:", err);
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    setSuccess("Copied to clipboard!");
    setTimeout(() => setSuccess(null), 2000);
  };

  const copyAllBackupCodes = () => {
    const codes = setupData?.backup_codes || newBackupCodes;
    const codesText = codes.join("\n");
    copyToClipboard(codesText);
  };

  if (loading && !mfaStatus) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
          <div className="h-8 bg-gray-200 rounded w-1/2"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className="text-lg font-medium text-gray-900">
            Multi-Factor Authentication
          </h3>
          <p className="text-sm text-gray-500">
            Add an extra layer of security to your account with TOTP
          </p>
        </div>
        {mfaStatus?.enabled && (
          <div className="flex items-center text-green-600">
            <IoShieldCheckmark className="h-5 w-5 mr-2" />
            <span className="text-sm font-medium">Enabled</span>
          </div>
        )}
      </div>

      {error && (
        <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg flex items-center">
          <IoWarning className="h-5 w-5 text-red-500 mr-2" />
          <span className="text-red-700">{error}</span>
        </div>
      )}

      {success && (
        <div className="mb-4 p-4 bg-green-50 border border-green-200 rounded-lg flex items-center">
          <IoShieldCheckmark className="h-5 w-5 text-green-500 mr-2" />
          <span className="text-green-700">{success}</span>
        </div>
      )}

      {step === "status" && (
        <div>
          {mfaStatus?.enabled ? (
            <div className="space-y-4">
              <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                <div className="flex items-center">
                  <IoShieldCheckmark className="h-6 w-6 text-green-500 mr-3" />
                  <div>
                    <h4 className="text-sm font-medium text-green-800">
                      MFA is Active
                    </h4>
                    <p className="text-sm text-green-600">
                      Your account is protected with multi-factor authentication
                    </p>
                    {mfaStatus.setup_date && (
                      <p className="text-xs text-green-500 mt-1">
                        Enabled on{" "}
                        {new Date(mfaStatus.setup_date).toLocaleDateString()}
                      </p>
                    )}
                  </div>
                </div>
              </div>

              <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div>
                  <h4 className="text-sm font-medium text-gray-900">
                    Backup Codes
                  </h4>
                  <p className="text-sm text-gray-500">
                    {mfaStatus.backup_codes_remaining} backup codes remaining
                  </p>
                </div>
                <button
                  onClick={() => setStep("backup")}
                  className="text-blue-600 hover:text-blue-700 text-sm font-medium"
                >
                  Manage
                </button>
              </div>

              <div className="flex space-x-3">
                <button
                  onClick={() => setStep("backup")}
                  className="flex-1 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors"
                >
                  Regenerate Backup Codes
                </button>
                <button
                  onClick={() => setStep("verify")}
                  className="flex-1 bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 transition-colors"
                >
                  Disable MFA
                </button>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                <div className="flex items-center">
                  <IoWarning className="h-6 w-6 text-yellow-500 mr-3" />
                  <div>
                    <h4 className="text-sm font-medium text-yellow-800">
                      MFA Not Enabled
                    </h4>
                    <p className="text-sm text-yellow-600">
                      Your account is not protected with multi-factor
                      authentication
                    </p>
                  </div>
                </div>
              </div>

              <button
                onClick={startMFASetup}
                disabled={loading}
                className="w-full max-w-3xs bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:bg-blue-400 transition-colors"
              >
                {loading ? "Setting up..." : "Enable MFA"}
              </button>
            </div>
          )}
        </div>
      )}

      {step === "setup" && setupData && (
        <div className="space-y-6">
          <div className="text-center">
            <h4 className="text-lg font-medium text-gray-900 mb-2">
              Scan QR Code
            </h4>
            <p className="text-sm text-gray-500 mb-4">
              Use your authenticator app (Google Authenticator, Authy, etc.) to
              scan this QR code
            </p>

            <div className="inline-block p-4 bg-white border-2 border-gray-200 rounded-lg">
              <img
                src={setupData.qr_code_data_url}
                alt="MFA QR Code"
                className="w-48 h-48 mx-auto"
              />
            </div>
          </div>

          <div className="bg-gray-50 rounded-lg p-4">
            <h5 className="text-sm font-medium text-gray-900 mb-2">
              Manual Entry
            </h5>
            <p className="text-xs text-gray-500 mb-2">
              If you can&apos;t scan the QR code, enter this secret manually:
            </p>
            <div className="flex items-center space-x-2">
              <code className="flex-1 text-black/70 bg-white px-3 py-2 rounded border text-sm font-mono">
                {setupData.secret}
              </code>
              <button
                onClick={() => copyToClipboard(setupData.secret)}
                className="p-2 text-gray-500 hover:text-gray-700"
                title="Copy secret"
              >
                <IoCopy className="h-4 w-4" />
              </button>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Enter verification code from your authenticator app:
            </label>
            <input
              type="text"
              value={verificationCode}
              onChange={(e) =>
                setVerificationCode(
                  e.target.value.replace(/\D/g, "").slice(0, 6)
                )
              }
              placeholder="123456"
              className="w-full text-black px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              maxLength={6}
            />
          </div>

          <div className="flex space-x-3 max-w-md">
            <button
              onClick={() => setStep("status")}
              className="flex-1 bg-gray-600 text-white px-4 py-2 rounded-md hover:bg-gray-700 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={verifyMFASetup}
              disabled={loading || verificationCode.length !== 6}
              className="flex-1 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:bg-blue-400 transition-colors"
            >
              {loading ? "Verifying..." : "Verify & Enable"}
            </button>
          </div>
        </div>
      )}

      {step === "verify" && (
        <div className="space-y-4">
          <div className="bg-red-50 border border-red-200 rounded-lg p-4">
            <div className="flex items-center">
              <IoWarning className="h-6 w-6 text-red-500 mr-3" />
              <div>
                <h4 className="text-sm font-medium text-red-800">
                  Disable MFA
                </h4>
                <p className="text-sm text-red-600">
                  This will remove multi-factor authentication from your account
                </p>
              </div>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Enter verification code to confirm:
            </label>
            <input
              type="text"
              value={verificationCode}
              onChange={(e) =>
                setVerificationCode(
                  e.target.value.replace(/\D/g, "").slice(0, 6)
                )
              }
              placeholder="123456"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              maxLength={6}
            />
          </div>

          <div className="flex space-x-3">
            <button
              onClick={() => setStep("status")}
              className="flex-1 bg-gray-600 text-white px-4 py-2 rounded-md hover:bg-gray-700 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={disableMFA}
              disabled={loading || verificationCode.length !== 6}
              className="flex-1 bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 disabled:bg-red-400 transition-colors"
            >
              {loading ? "Disabling..." : "Disable MFA"}
            </button>
          </div>
        </div>
      )}

      {step === "backup" && (
        <div className="space-y-6">
          <div>
            <h4 className="text-lg font-medium text-gray-900 mb-2">
              Backup Codes
            </h4>
            <p className="text-sm text-gray-500 mb-4">
              Save these backup codes in a secure location. You can use them to
              access your account if you lose your authenticator device.
            </p>
          </div>

          {(setupData?.backup_codes || newBackupCodes.length > 0) && (
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="flex items-center justify-between mb-3">
                <h5 className="text-sm font-medium text-gray-900">
                  Your Backup Codes
                </h5>
                <div className="flex space-x-2">
                  <button
                    onClick={copyAllBackupCodes}
                    className="text-blue-600 hover:text-blue-700 text-sm font-medium flex items-center"
                  >
                    <IoCopy className="h-4 w-4 mr-1" />
                    Copy All
                  </button>
                  <button
                    onClick={() => setShowBackupCodes(!showBackupCodes)}
                    className="text-gray-600 hover:text-gray-700 text-sm font-medium flex items-center"
                  >
                    {showBackupCodes ? (
                      <IoEyeOff className="h-4 w-4" />
                    ) : (
                      <IoEye className="h-4 w-4" />
                    )}
                  </button>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-2">
                {(setupData?.backup_codes || newBackupCodes).map(
                  (code, index) => (
                    <div
                      key={index}
                      className="bg-white px-3 py-2 rounded border"
                    >
                      <code className="text-sm font-mono">
                        {showBackupCodes ? code : "••••••••••"}
                      </code>
                    </div>
                  )
                )}
              </div>

              <div className="mt-3 p-3 bg-yellow-50 border border-yellow-200 rounded">
                <p className="text-xs text-yellow-700">
                  ⚠️ Each backup code can only be used once. Store them securely
                  and don&apos;t share them.
                </p>
              </div>
            </div>
          )}

          {mfaStatus?.enabled && (
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Enter verification code to regenerate backup codes:
                </label>
                <input
                  type="text"
                  value={verificationCode}
                  onChange={(e) =>
                    setVerificationCode(
                      e.target.value.replace(/\D/g, "").slice(0, 6)
                    )
                  }
                  placeholder="123456"
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  maxLength={6}
                />
              </div>

              <button
                onClick={regenerateBackupCodes}
                disabled={loading || verificationCode.length !== 6}
                className="w-full bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:bg-blue-400 transition-colors"
              >
                {loading ? "Regenerating..." : "Regenerate Backup Codes"}
              </button>
            </div>
          )}

          <button
            onClick={() => setStep("status")}
            className="w-full bg-gray-600 text-white px-4 py-2 rounded-md hover:bg-gray-700 transition-colors"
          >
            Done
          </button>
        </div>
      )}
    </div>
  );
}
