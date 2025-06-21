"use client";

// Developer Note: To manually clear login tracking for testing, run in browser console:
// localStorage.removeItem(`login_event_${new Date().toDateString()}`)

import DashboardLayout from "@/components/DashboardLayout";
import MFASetup from "@/components/dashboard/MFASetup";
import SecurityAlerts from "@/components/dashboard/SecurityAlerts";
import SecurityEnhancements from "@/components/dashboard/SecurityEnhancements";
import { SECURITY_FEATURES } from "@/constants";
import {
  apiClient,
  type SecurityAlert,
  type SecurityEvent,
  type TrustedDevice,
  type UserSettings,
} from "@/lib/api";
import { useCallback, useEffect, useState } from "react";
import {
  IoCheckmarkCircle,
  IoDesktop,
  IoInformationCircle,
  IoPhonePortrait,
  IoRefresh,
  IoShieldCheckmark,
  IoTabletLandscape,
  IoTrash,
  IoWarning,
} from "react-icons/io5";
import { toast } from "sonner";

// Device fingerprinting utility
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
    (navigator as Navigator & { deviceMemory?: number }).deviceMemory ||
      "unknown",
  ].join("|");

  // Simple hash function
  let hash = 0;
  for (let i = 0; i < fingerprint.length; i++) {
    const char = fingerprint.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash = hash & hash; // Convert to 32bit integer
  }
  return Math.abs(hash).toString(16);
};

const getDeviceInfo = () => {
  const userAgent = navigator.userAgent;
  let deviceType = "desktop";
  let deviceName = "Unknown Device";
  let browser = "Unknown Browser";
  let os = "Unknown OS";

  // Detect device type
  if (/Mobile|Android|iPhone|iPad/.test(userAgent)) {
    if (/iPad/.test(userAgent)) {
      deviceType = "tablet";
      deviceName = "iPad";
    } else if (/iPhone/.test(userAgent)) {
      deviceType = "mobile";
      deviceName = "iPhone";
    } else if (/Android/.test(userAgent)) {
      deviceType = /Mobile/.test(userAgent) ? "mobile" : "tablet";
      deviceName = deviceType === "mobile" ? "Android Phone" : "Android Tablet";
    }
  } else {
    deviceType = "desktop";
    if (/Windows/.test(userAgent)) {
      deviceName = "Windows PC";
    } else if (/Mac/.test(userAgent)) {
      deviceName = "Mac";
    } else if (/Linux/.test(userAgent)) {
      deviceName = "Linux PC";
    }
  }

  // Detect browser
  if (/Chrome/.test(userAgent) && !/Edge/.test(userAgent)) {
    browser = "Chrome";
  } else if (/Firefox/.test(userAgent)) {
    browser = "Firefox";
  } else if (/Safari/.test(userAgent) && !/Chrome/.test(userAgent)) {
    browser = "Safari";
  } else if (/Edge/.test(userAgent)) {
    browser = "Edge";
  }

  // Detect OS
  if (/Windows/.test(userAgent)) {
    os = "Windows";
  } else if (/Mac/.test(userAgent)) {
    os = "macOS";
  } else if (/Linux/.test(userAgent)) {
    os = "Linux";
  } else if (/Android/.test(userAgent)) {
    os = "Android";
  } else if (/iOS|iPhone|iPad/.test(userAgent)) {
    os = "iOS";
  }

  return { deviceType, deviceName, browser, os };
};

// Get user's location (approximate)
const getUserLocation = async (): Promise<string> => {
  try {
    // Try to get location from geolocation API
    if (navigator.geolocation) {
      return new Promise((resolve) => {
        navigator.geolocation.getCurrentPosition(
          (position) => {
            resolve(
              `${position.coords.latitude.toFixed(
                2
              )}, ${position.coords.longitude.toFixed(2)}`
            );
          },
          () => {
            // Fallback to timezone-based location
            resolve(Intl.DateTimeFormat().resolvedOptions().timeZone);
          },
          { timeout: 5000 }
        );
      });
    }
  } catch (error) {
    console.warn("Could not get location:", error);
  }

  // Fallback to timezone
  return Intl.DateTimeFormat().resolvedOptions().timeZone;
};

export default function SecurityPage() {
  const [loading, setLoading] = useState(false);
  const [securityEvents, setSecurityEvents] = useState<SecurityEvent[]>([]);
  const [trustedDevices, setTrustedDevices] = useState<TrustedDevice[]>([]);
  const [userSettings, setUserSettings] = useState<UserSettings | null>(null);
  const [securityAlerts, setSecurityAlerts] = useState<SecurityAlert[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [deviceRegistered, setDeviceRegistered] = useState(false);
  const [hasCreatedLoginEvent, setHasCreatedLoginEvent] = useState(false);

  const registerCurrentDevice = useCallback(async () => {
    // Check if device was already registered using fingerprint
    const fingerprint = generateDeviceFingerprint();
    const storageKey = `device_registered_${fingerprint}`;
    const deviceAlreadyRegistered = localStorage.getItem(storageKey);

    if (deviceRegistered || deviceAlreadyRegistered) return;

    try {
      const deviceInfo = getDeviceInfo();
      const location = await getUserLocation();

      await apiClient.registerDevice({
        device_name: deviceInfo.deviceName,
        device_type: deviceInfo.deviceType,
        browser: deviceInfo.browser,
        os: deviceInfo.os,
        fingerprint: fingerprint,
        location: location,
      });

      setDeviceRegistered(true);
      // Store device registration to prevent duplicates
      localStorage.setItem(storageKey, "true");
    } catch (error) {
      console.warn("Failed to register device:", error);
    }
  }, [deviceRegistered]);

  const createLoginEvent = useCallback(async () => {
    // Check if login event was already created in this session
    const sessionKey = `login_event_${new Date().toDateString()}`;
    const loginEventCreated = localStorage.getItem(sessionKey);

    if (hasCreatedLoginEvent || loginEventCreated) {
      console.log("Login event already created today, skipping...");
      return;
    }

    try {
      const location = await getUserLocation();

      await apiClient.createSecurityEvent({
        event_type: "login",
        description: "User logged in to CloudGate",
        severity: "low",
        location: location,
        risk_score: 0.1,
      });

      console.log("Login event created successfully");
      setHasCreatedLoginEvent(true);
      // Store in localStorage with today's date to prevent multiple logins per day
      localStorage.setItem(sessionKey, "true");
    } catch (error) {
      console.warn("Failed to create login event:", error);
    }
  }, [hasCreatedLoginEvent]);

  const initializeSecurityPage = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      // Register current device and create login event
      await registerCurrentDevice();
      await createLoginEvent();

      // Load security data
      await loadSecurityData();
    } catch (err) {
      console.error("Failed to initialize security page:", err);
      setError("Failed to initialize security page. Please try again.");
    } finally {
      setLoading(false);
    }
  }, [registerCurrentDevice, createLoginEvent]);

  useEffect(() => {
    // Clean up old login events from localStorage (older than 7 days)
    const cleanupOldLoginEvents = () => {
      const keys = Object.keys(localStorage);
      const currentDate = new Date();

      keys.forEach((key) => {
        if (
          key.startsWith("login_event_") ||
          key.startsWith("device_registered_")
        ) {
          if (key.startsWith("login_event_")) {
            const dateStr = key.replace("login_event_", "");
            const eventDate = new Date(dateStr);
            const daysDiff =
              (currentDate.getTime() - eventDate.getTime()) /
              (1000 * 3600 * 24);

            if (daysDiff > 7) {
              localStorage.removeItem(key);
            }
          }
          // For device registration, keep it unless manually cleared
        }
      });
    };

    cleanupOldLoginEvents();
    initializeSecurityPage();
  }, [initializeSecurityPage]);

  const loadSecurityData = async () => {
    try {
      // Load all security data in parallel
      const [
        eventsResponse,
        devicesResponse,
        settingsResponse,
        alertsResponse,
      ] = await Promise.all([
        apiClient.getSecurityEvents(10),
        apiClient.getTrustedDevices(),
        apiClient.getUserSettings(),
        apiClient.getSecurityAlerts(),
      ]);

      setSecurityEvents(eventsResponse.events);
      setTrustedDevices(devicesResponse.devices);
      setUserSettings(settingsResponse.settings);
      setSecurityAlerts(alertsResponse.alerts);
    } catch (err) {
      console.error("Failed to load security data:", err);
      setError("Failed to load security data. Please try again.");
    }
  };

  const handleRefresh = () => {
    initializeSecurityPage();
  };

  const handleSettingChange = async (
    setting: string,
    value: boolean | number
  ) => {
    if (!userSettings) return;

    try {
      const response = await apiClient.updateSingleSetting(setting, value);
      const updatedSettings = response.settings;
      setUserSettings(updatedSettings);

      toast.success(`Security setting updated successfully`);

      // Create a security event for settings change
      await apiClient.createSecurityEvent({
        event_type: "settings_change",
        description: `Security setting '${setting}' changed to '${value}'`,
        severity: "low",
        risk_score: 0.2,
      });

      // Refresh events to show the new one
      const eventsResponse = await apiClient.getSecurityEvents(10);
      setSecurityEvents(eventsResponse.events);
    } catch (err) {
      console.error("Failed to update setting:", err);
      toast.error("Failed to update setting. Please try again.");
      setError("Failed to update setting. Please try again.");
    }
  };

  const revokeDevice = async (deviceId: string) => {
    try {
      // Find device before revoking for the security event
      const revokedDevice = trustedDevices.find((d) => d.id === deviceId);

      await apiClient.revokeDevice(deviceId);

      // Remove device from local state immediately
      setTrustedDevices((devices) => devices.filter((d) => d.id !== deviceId));

      toast.success(
        `Device "${revokedDevice?.device_name || "Unknown"}" has been revoked`
      );

      // Create security event for device revocation
      await apiClient.createSecurityEvent({
        event_type: "device_revoked",
        description: `Device '${
          revokedDevice?.device_name || "Unknown"
        }' was revoked`,
        severity: "medium",
        risk_score: 0.4,
      });

      // Refresh events
      const eventsResponse = await apiClient.getSecurityEvents(10);
      setSecurityEvents(eventsResponse.events);
    } catch (err) {
      console.error("Failed to revoke device:", err);
      toast.error("Failed to revoke device. Please try again.");
      setError("Failed to revoke device. Please try again.");
    }
  };

  const trustDevice = async (deviceId: string) => {
    try {
      const deviceToTrust = trustedDevices.find((d) => d.id === deviceId);

      await apiClient.trustDevice(deviceId);
      setTrustedDevices((devices) =>
        devices.map((d) => (d.id === deviceId ? { ...d, trusted: true } : d))
      );

      toast.success(
        `Device "${deviceToTrust?.device_name || "Unknown"}" is now trusted`
      );

      // Create security event for device trust
      await apiClient.createSecurityEvent({
        event_type: "device_trusted",
        description: `Device '${
          deviceToTrust?.device_name || "Unknown"
        }' was marked as trusted`,
        severity: "low",
        risk_score: 0.1,
      });

      // Refresh events
      const eventsResponse = await apiClient.getSecurityEvents(10);
      setSecurityEvents(eventsResponse.events);
    } catch (err) {
      console.error("Failed to trust device:", err);
      toast.error("Failed to trust device. Please try again.");
      setError("Failed to trust device. Please try again.");
    }
  };

  const handleUpdateAlertStatus = async (
    alertId: string,
    status: SecurityAlert["status"]
  ) => {
    try {
      await apiClient.updateSecurityAlertStatus(alertId, status);
      setSecurityAlerts(
        securityAlerts.map((alert) =>
          alert.id === alertId ? { ...alert, status } : alert
        )
      );

      toast.success(`Alert ${status} successfully`);
    } catch (err) {
      console.error("Failed to update alert status:", err);
      toast.error("Failed to update alert status");
    }
  };

  const handleMFAStatusChange = useCallback(
    (enabled: boolean) => {
      if (userSettings) {
        setUserSettings({ ...userSettings, two_factor_enabled: enabled });
      }
    },
    [userSettings]
  );

  const getEventIcon = (type: string) => {
    switch (type) {
      case "login":
        return "🔓";
      case "logout":
        return "🔒";
      case "failed_login":
        return "⚠️";
      case "password_change":
        return "🔑";
      case "mfa_enabled":
      case "mfa_disabled":
        return "🛡️";
      case "device_registered":
        return "📱";
      case "suspicious_activity":
        return "🚨";
      default:
        return "ℹ️";
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity.toLowerCase()) {
      case "high":
        return "bg-red-100 text-red-800";
      case "medium":
        return "bg-yellow-100 text-yellow-800";
      case "low":
        return "bg-green-100 text-green-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const getDeviceIcon = (deviceType: string) => {
    switch (deviceType.toLowerCase()) {
      case "mobile":
      case "phone":
        return <IoPhonePortrait className="h-5 w-5" />;
      case "tablet":
        return <IoTabletLandscape className="h-5 w-5" />;
      case "desktop":
      case "computer":
      default:
        return <IoDesktop className="h-5 w-5" />;
    }
  };

  const getSecurityScore = () => {
    if (!userSettings) return "Medium";

    let score = 0;
    if (userSettings.two_factor_enabled) score += 40;
    if (userSettings.login_notifications) score += 20;
    if (userSettings.suspicious_activity_alerts) score += 20;
    if (userSettings.session_timeout <= 30) score += 10;
    if (userSettings.password_expiry_days <= 90) score += 10;

    if (score >= 80) return "High";
    if (score >= 50) return "Medium";
    return "Low";
  };

  const getActiveSessionsCount = () => {
    return trustedDevices.filter(
      (d) =>
        d.trusted &&
        new Date(d.last_seen) > new Date(Date.now() - 24 * 60 * 60 * 1000)
    ).length;
  };

  const getRecentAlertsCount = () => {
    const sevenDaysAgo = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000);
    return securityEvents.filter(
      (e) =>
        new Date(e.created_at) > sevenDaysAgo &&
        (e.severity === "medium" || e.severity === "high")
    ).length;
  };

  const refreshAction = (
    <button
      onClick={handleRefresh}
      disabled={loading}
      className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 cursor-pointer"
    >
      <IoRefresh className={`h-4 w-4 mr-2 ${loading ? "animate-spin" : ""}`} />
      {loading ? "Refreshing..." : "Refresh"}
    </button>
  );

  if (loading && !userSettings) {
    return (
      <DashboardLayout>
        <div className="space-y-6">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-white rounded-lg shadow p-6">
              <div className="animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
                <div className="h-8 bg-gray-200 rounded w-1/2"></div>
              </div>
            </div>
          ))}
        </div>
      </DashboardLayout>
    );
  }

  if (error) {
    return (
      <DashboardLayout>
        <div className="bg-red-50 border border-red-200 rounded-lg p-6">
          <div className="flex items-center">
            <IoWarning className="h-6 w-6 text-red-500 mr-3" />
            <div>
              <h3 className="text-lg font-medium text-red-800">
                Error Loading Security Data
              </h3>
              <p className="text-red-600">{error}</p>
              <button
                onClick={loadSecurityData}
                className="mt-2 text-red-700 hover:text-red-900 font-medium cursor-pointer"
              >
                Try Again
              </button>
            </div>
          </div>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      {/* Security Overview */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-green-100 rounded-lg">
              <IoShieldCheckmark className="h-6 w-6 text-green-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">
                Security Score
              </p>
              <p className="text-2xl font-semibold text-green-600">
                {getSecurityScore()}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-blue-100 rounded-lg">
              <IoInformationCircle className="h-6 w-6 text-blue-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">
                Active Sessions
              </p>
              <p className="text-2xl font-semibold text-gray-900">
                {getActiveSessionsCount()}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-yellow-100 rounded-lg">
              <IoWarning className="h-6 w-6 text-yellow-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">
                Alerts (7 days)
              </p>
              <p className="text-2xl font-semibold text-gray-900">
                {getRecentAlertsCount()}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Security Features Status */}
      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <h3 className="text-lg font-medium text-gray-900 mb-4">
          Security Features
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {SECURITY_FEATURES.map((feature) => (
            <div key={feature.id} className="flex items-center space-x-3">
              <div
                className={`flex-shrink-0 w-10 h-10 bg-${feature.color}-100 rounded-full flex items-center justify-center`}
              >
                <span className={`text-${feature.color}-600 text-lg`}>✓</span>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-900">
                  {feature.title}
                </p>
                <p className="text-sm text-gray-500">{feature.description}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* MFA Setup */}
      <MFASetup onMFAStatusChange={handleMFAStatusChange} />

      {/* Security Alerts Section */}
      <div className="mb-8">
        <SecurityAlerts
          alerts={securityAlerts}
          onUpdateStatus={handleUpdateAlertStatus}
        />
      </div>

      {/* Security Enhancements */}
      <div className="mb-8">
        <SecurityEnhancements />
      </div>

      {/* Trusted Devices */}
      <div className="bg-white rounded-lg shadow p-6 mb-8 mt-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">
          Trusted Devices
        </h3>
        <div className="space-y-4">
          {trustedDevices.length === 0 ? (
            <div className="text-center py-8">
              <div className="text-gray-400 text-4xl mb-4">📱</div>
              <h4 className="text-lg font-medium text-gray-900 mb-2">
                No trusted devices
              </h4>
              <p className="text-gray-500">
                Devices you trust will appear here.
              </p>
            </div>
          ) : (
            trustedDevices.map((device) => (
              <div
                key={device.id}
                className="flex items-center justify-between p-4 border border-gray-200 rounded-lg"
              >
                <div className="flex items-center space-x-4">
                  <div className="text-gray-500">
                    {getDeviceIcon(device.device_type)}
                  </div>
                  <div>
                    <h4 className="text-sm font-medium text-gray-900">
                      {device.device_name}
                    </h4>
                    <div className="text-xs text-gray-500 space-y-1">
                      <div>
                        {device.browser} on {device.os}
                      </div>
                      <div>
                        Last seen: {new Date(device.last_seen).toLocaleString()}
                      </div>
                      <div>Location: {device.location || "Unknown"}</div>
                    </div>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  {device.trusted ? (
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      <IoCheckmarkCircle className="h-3 w-3 mr-1" />
                      Trusted
                    </span>
                  ) : (
                    <button
                      onClick={() => trustDevice(device.id)}
                      className="text-blue-600 hover:text-blue-700 text-sm font-medium cursor-pointer"
                    >
                      Trust Device
                    </button>
                  )}
                  <button
                    onClick={() => revokeDevice(device.id)}
                    className="text-red-600 hover:text-red-700 p-1 cursor-pointer"
                    title="Remove device"
                  >
                    <IoTrash className="h-4 w-4" />
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      {/* Security Settings */}
      {userSettings && (
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h3 className="text-lg font-medium text-gray-900 mb-4">
            Security Settings
          </h3>
          <div className="space-y-6">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-gray-900">
                  Login Notifications
                </h4>
                <p className="text-sm text-gray-500">
                  Get notified when someone logs into your account
                </p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  checked={userSettings.login_notifications}
                  onChange={(e) =>
                    handleSettingChange("login_notifications", e.target.checked)
                  }
                  className="sr-only peer"
                />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>

            <div className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-gray-900">
                  Suspicious Activity Alerts
                </h4>
                <p className="text-sm text-gray-500">
                  Get alerted about unusual account activity
                </p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  checked={userSettings.suspicious_activity_alerts}
                  onChange={(e) =>
                    handleSettingChange(
                      "suspicious_activity_alerts",
                      e.target.checked
                    )
                  }
                  className="sr-only peer"
                />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>

            <div className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-gray-900">
                  Session Timeout
                </h4>
                <p className="text-sm text-gray-500">
                  Automatically log out after inactivity (minutes)
                </p>
              </div>
              <select
                value={userSettings.session_timeout}
                onChange={(e) =>
                  handleSettingChange(
                    "session_timeout",
                    parseInt(e.target.value)
                  )
                }
                className="px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 text-black/80 focus:border-blue-500"
              >
                <option value={15}>15 minutes</option>
                <option value={30}>30 minutes</option>
                <option value={60}>1 hour</option>
                <option value={120}>2 hours</option>
                <option value={240}>4 hours</option>
              </select>
            </div>

            <div className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-gray-900">
                  Password Expiry
                </h4>
                <p className="text-sm text-gray-500">
                  Require password change every (days)
                </p>
              </div>
              <select
                value={userSettings.password_expiry_days}
                onChange={(e) =>
                  handleSettingChange(
                    "password_expiry_days",
                    parseInt(e.target.value)
                  )
                }
                className="px-3 py-2 border text-black/80 border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value={30}>30 days</option>
                <option value={60}>60 days</option>
                <option value={90}>90 days</option>
                <option value={180}>180 days</option>
                <option value={365}>1 year</option>
              </select>
            </div>
          </div>
        </div>
      )}

      {/* Recent Security Events */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">
            Recent Security Events
          </h3>
          <p className="text-sm text-gray-500">Monitor your account activity</p>
        </div>
        <div className="divide-y divide-gray-200">
          {securityEvents.length === 0 ? (
            <div className="px-6 py-12 text-center">
              <div className="text-gray-400 text-4xl mb-4">🛡️</div>
              <h4 className="text-lg font-medium text-gray-900 mb-2">
                No security events
              </h4>
              <p className="text-gray-500">
                Security events will appear here as they occur.
              </p>
            </div>
          ) : (
            securityEvents.map((event) => (
              <div key={event.id} className="px-6 py-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <div className="text-2xl mr-4">
                      {getEventIcon(event.event_type)}
                    </div>
                    <div>
                      <h4 className="text-sm font-medium text-gray-900">
                        {event.description}
                      </h4>
                      <div className="flex items-center space-x-4 mt-1">
                        <span className="text-xs text-gray-500">
                          {new Date(event.created_at).toLocaleString()}
                        </span>
                        <span className="text-xs text-gray-500">
                          {event.location || "Unknown location"}
                        </span>
                        <span className="text-xs text-gray-500">
                          IP: {event.ip_address}
                        </span>
                        {event.risk_score > 0 && (
                          <span className="text-xs text-gray-500">
                            Risk: {event.risk_score}/10
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                  <span
                    className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getSeverityColor(
                      event.severity
                    )}`}
                  >
                    {event.severity.charAt(0).toUpperCase() +
                      event.severity.slice(1)}
                  </span>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </DashboardLayout>
  );
}
