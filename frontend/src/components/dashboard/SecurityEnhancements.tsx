"use client";

import { useEffect, useState } from "react";
import {
  IoAlert,
  IoLocation,
  IoPhonePortrait,
  IoRefresh,
  IoShieldCheckmark,
  IoWarning,
} from "react-icons/io5";

interface DeviceInfo {
  id: string;
  name: string;
  type: "desktop" | "mobile" | "tablet";
  browser: string;
  os: string;
  location: string;
  lastSeen: string;
  trusted: boolean;
  fingerprint: string;
}

interface SecurityEvent {
  id: string;
  type: "login" | "suspicious_location" | "new_device" | "failed_mfa";
  description: string;
  timestamp: string;
  location: string;
  ip: string;
  riskScore: number;
  severity: "low" | "medium" | "high";
}

interface RiskAssessment {
  score: number;
  level: "low" | "medium" | "high";
  factors: string[];
  recommendations: string[];
}

export default function SecurityEnhancements() {
  const [devices, setDevices] = useState<DeviceInfo[]>([]);
  const [securityEvents, setSecurityEvents] = useState<SecurityEvent[]>([]);
  const [riskAssessment, setRiskAssessment] = useState<RiskAssessment | null>(
    null
  );
  const [loading, setLoading] = useState(false);
  const [currentDevice, setCurrentDevice] = useState<DeviceInfo | null>(null);

  useEffect(() => {
    loadSecurityData();
    generateDeviceFingerprint();
  }, []);

  const loadSecurityData = async () => {
    setLoading(true);

    // Simulate loading security data
    setTimeout(() => {
      setDevices([
        {
          id: "1",
          name: "Chrome on Windows",
          type: "desktop",
          browser: "Chrome 120.0",
          os: "Windows 11",
          location: "New York, US",
          lastSeen: "2024-01-20T10:30:00Z",
          trusted: true,
          fingerprint: "fp_abc123def456",
        },
        {
          id: "2",
          name: "Safari on iPhone",
          type: "mobile",
          browser: "Safari 17.0",
          os: "iOS 17.2",
          location: "New York, US",
          lastSeen: "2024-01-19T15:45:00Z",
          trusted: true,
          fingerprint: "fp_xyz789ghi012",
        },
        {
          id: "3",
          name: "Firefox on Linux",
          type: "desktop",
          browser: "Firefox 121.0",
          os: "Ubuntu 22.04",
          location: "Unknown",
          lastSeen: "2024-01-18T08:20:00Z",
          trusted: false,
          fingerprint: "fp_jkl345mno678",
        },
      ]);

      setSecurityEvents([
        {
          id: "1",
          type: "login",
          description: "Successful login from trusted device",
          timestamp: "2024-01-20T10:30:00Z",
          location: "New York, US",
          ip: "192.168.1.100",
          riskScore: 0.1,
          severity: "low",
        },
        {
          id: "2",
          type: "new_device",
          description: "Login from new device detected",
          timestamp: "2024-01-18T08:20:00Z",
          location: "Unknown",
          ip: "203.0.113.1",
          riskScore: 0.7,
          severity: "medium",
        },
        {
          id: "3",
          type: "suspicious_location",
          description: "Login from unusual location",
          timestamp: "2024-01-17T22:15:00Z",
          location: "Moscow, RU",
          ip: "198.51.100.1",
          riskScore: 0.9,
          severity: "high",
        },
      ]);

      setRiskAssessment({
        score: 0.3,
        level: "low",
        factors: [
          "Trusted device usage",
          "Consistent location patterns",
          "MFA enabled",
          "Recent password change",
        ],
        recommendations: [
          "Continue using trusted devices",
          "Review device list regularly",
          "Enable login notifications",
        ],
      });

      setLoading(false);
    }, 1000);
  };

  const generateDeviceFingerprint = async () => {
    try {
      // Basic device fingerprinting
      const canvas = document.createElement("canvas");
      const ctx = canvas.getContext("2d");
      if (ctx) {
        ctx.textBaseline = "top";
        ctx.font = "14px Arial";
        ctx.fillText("Device fingerprint", 2, 2);
      }

      const fingerprint = {
        userAgent: navigator.userAgent,
        language: navigator.language,
        platform: navigator.platform,
        screenResolution: `${screen.width}x${screen.height}`,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
        canvasFingerprint: canvas.toDataURL(),
        cookieEnabled: navigator.cookieEnabled,
        doNotTrack: navigator.doNotTrack,
        hardwareConcurrency: navigator.hardwareConcurrency,
        maxTouchPoints: navigator.maxTouchPoints,
      };

      // Generate a simple hash
      const fingerprintString = JSON.stringify(fingerprint);
      const hash = await crypto.subtle.digest(
        "SHA-256",
        new TextEncoder().encode(fingerprintString)
      );
      const hashArray = Array.from(new Uint8Array(hash));
      const hashHex = hashArray
        .map((b) => b.toString(16).padStart(2, "0"))
        .join("");

      setCurrentDevice({
        id: "current",
        name: `${getBrowserName()} on ${getOSName()}`,
        type: getDeviceType(),
        browser: getBrowserName(),
        os: getOSName(),
        location: "Current Location",
        lastSeen: new Date().toISOString(),
        trusted: true,
        fingerprint: `fp_${hashHex.substring(0, 12)}`,
      });
    } catch (error) {
      console.error("Error generating device fingerprint:", error);
    }
  };

  const getBrowserName = (): string => {
    const userAgent = navigator.userAgent;
    if (userAgent.includes("Chrome")) return "Chrome";
    if (userAgent.includes("Firefox")) return "Firefox";
    if (userAgent.includes("Safari")) return "Safari";
    if (userAgent.includes("Edge")) return "Edge";
    return "Unknown Browser";
  };

  const getOSName = (): string => {
    const userAgent = navigator.userAgent;
    if (userAgent.includes("Windows")) return "Windows";
    if (userAgent.includes("Mac")) return "macOS";
    if (userAgent.includes("Linux")) return "Linux";
    if (userAgent.includes("Android")) return "Android";
    if (userAgent.includes("iOS")) return "iOS";
    return "Unknown OS";
  };

  const getDeviceType = (): "desktop" | "mobile" | "tablet" => {
    const userAgent = navigator.userAgent;
    if (/tablet|ipad/i.test(userAgent)) return "tablet";
    if (/mobile|android|iphone/i.test(userAgent)) return "mobile";
    return "desktop";
  };

  const trustDevice = (deviceId: string) => {
    setDevices((prev) =>
      prev.map((device) =>
        device.id === deviceId ? { ...device, trusted: true } : device
      )
    );
  };

  const revokeDevice = (deviceId: string) => {
    setDevices((prev) => prev.filter((device) => device.id !== deviceId));
  };

  const getDeviceIcon = (type: string) => {
    switch (type) {
      case "mobile":
        return "📱";
      case "tablet":
        return "📱";
      case "desktop":
        return "💻";
      default:
        return "🖥️";
    }
  };

  const getRiskColor = (level: string) => {
    switch (level) {
      case "low":
        return "text-green-600 bg-green-100";
      case "medium":
        return "text-yellow-600 bg-yellow-100";
      case "high":
        return "text-red-600 bg-red-100";
      default:
        return "text-gray-600 bg-gray-100";
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case "low":
        return "bg-green-100 text-green-800";
      case "medium":
        return "bg-yellow-100 text-yellow-800";
      case "high":
        return "bg-red-100 text-red-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  if (loading) {
    return (
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
    );
  }

  return (
    <div className="space-y-6">
      {/* Risk Assessment */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h3 className="text-lg font-medium text-gray-900">
              Risk Assessment
            </h3>
            <p className="text-sm text-gray-500">
              Real-time security risk analysis based on your activity
            </p>
          </div>
          <button
            onClick={loadSecurityData}
            className="p-2 text-gray-500 hover:text-gray-700"
            title="Refresh assessment"
          >
            <IoRefresh className="h-5 w-5" />
          </button>
        </div>

        {riskAssessment && (
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <div
                  className={`px-3 py-1 rounded-full text-sm font-medium ${getRiskColor(
                    riskAssessment.level
                  )}`}
                >
                  {riskAssessment.level.toUpperCase()} RISK
                </div>
                <span className="text-2xl font-bold text-gray-900">
                  {Math.round(riskAssessment.score * 100)}%
                </span>
              </div>
              <IoShieldCheckmark className="h-8 w-8 text-green-500" />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 className="text-sm font-medium text-gray-900 mb-2">
                  Risk Factors
                </h4>
                <ul className="space-y-1">
                  {riskAssessment.factors.map((factor, index) => (
                    <li
                      key={index}
                      className="text-sm text-gray-600 flex items-center"
                    >
                      <span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
                      {factor}
                    </li>
                  ))}
                </ul>
              </div>

              <div>
                <h4 className="text-sm font-medium text-gray-900 mb-2">
                  Recommendations
                </h4>
                <ul className="space-y-1">
                  {riskAssessment.recommendations.map((rec, index) => (
                    <li
                      key={index}
                      className="text-sm text-gray-600 flex items-center"
                    >
                      <span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
                      {rec}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Device Management */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h3 className="text-lg font-medium text-gray-900">
              Trusted Devices
            </h3>
            <p className="text-sm text-gray-500">
              Manage devices that have access to your account
            </p>
          </div>
          <IoPhonePortrait className="h-6 w-6 text-gray-400" />
        </div>

        {/* Current Device */}
        {currentDevice && (
          <div className="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <span className="text-2xl">
                  {getDeviceIcon(currentDevice.type)}
                </span>
                <div>
                  <h4 className="text-sm font-medium text-blue-900">
                    {currentDevice.name} (Current Device)
                  </h4>
                  <p className="text-xs text-blue-600">
                    Fingerprint: {currentDevice.fingerprint}
                  </p>
                </div>
              </div>
              <div className="text-blue-600">
                <IoShieldCheckmark className="h-5 w-5" />
              </div>
            </div>
          </div>
        )}

        <div className="space-y-4">
          {devices.map((device) => (
            <div
              key={device.id}
              className="flex items-center justify-between p-4 border border-gray-200 rounded-lg"
            >
              <div className="flex items-center space-x-3">
                <span className="text-2xl">{getDeviceIcon(device.type)}</span>
                <div>
                  <h4 className="text-sm font-medium text-gray-900">
                    {device.name}
                  </h4>
                  <div className="flex items-center space-x-4 mt-1">
                    <span className="text-xs text-gray-500">
                      <IoLocation className="inline h-3 w-3 mr-1" />
                      {device.location}
                    </span>
                    <span className="text-xs text-gray-500">
                      Last seen:{" "}
                      {new Date(device.lastSeen).toLocaleDateString()}
                    </span>
                  </div>
                  <p className="text-xs text-gray-400 mt-1">
                    Fingerprint: {device.fingerprint}
                  </p>
                </div>
              </div>

              <div className="flex items-center space-x-2">
                {device.trusted ? (
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                    <IoShieldCheckmark className="h-3 w-3 mr-1" />
                    Trusted
                  </span>
                ) : (
                  <div className="flex space-x-2">
                    <button
                      onClick={() => trustDevice(device.id)}
                      className="text-blue-600 hover:text-blue-700 text-xs font-medium"
                    >
                      Trust
                    </button>
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                      <IoWarning className="h-3 w-3 mr-1" />
                      Untrusted
                    </span>
                  </div>
                )}

                <button
                  onClick={() => revokeDevice(device.id)}
                  className="text-red-600 hover:text-red-700 text-xs font-medium"
                >
                  Revoke
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Security Events */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h3 className="text-lg font-medium text-gray-900">
              Security Events
            </h3>
            <p className="text-sm text-gray-500">
              Recent security-related activities and alerts
            </p>
          </div>
          <IoAlert className="h-6 w-6 text-gray-400" />
        </div>

        <div className="space-y-4">
          {securityEvents.map((event) => (
            <div
              key={event.id}
              className="flex items-center justify-between p-4 border border-gray-200 rounded-lg"
            >
              <div className="flex items-center space-x-3">
                <div
                  className={`w-3 h-3 rounded-full ${
                    event.severity === "high"
                      ? "bg-red-500"
                      : event.severity === "medium"
                      ? "bg-yellow-500"
                      : "bg-green-500"
                  }`}
                ></div>
                <div>
                  <h4 className="text-sm font-medium text-gray-900">
                    {event.description}
                  </h4>
                  <div className="flex items-center space-x-4 mt-1">
                    <span className="text-xs text-gray-500">
                      {new Date(event.timestamp).toLocaleString()}
                    </span>
                    <span className="text-xs text-gray-500">
                      <IoLocation className="inline h-3 w-3 mr-1" />
                      {event.location}
                    </span>
                    <span className="text-xs text-gray-500">
                      IP: {event.ip}
                    </span>
                  </div>
                </div>
              </div>

              <div className="flex items-center space-x-2">
                <span
                  className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getSeverityColor(
                    event.severity
                  )}`}
                >
                  {event.severity.toUpperCase()}
                </span>
                <span className="text-xs text-gray-500">
                  Risk: {Math.round(event.riskScore * 100)}%
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
