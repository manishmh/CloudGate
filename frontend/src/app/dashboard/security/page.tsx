"use client";

import DashboardLayout from "@/components/DashboardLayout";
import MFASetup from "@/components/dashboard/MFASetup";
import SecurityEnhancements from "@/components/dashboard/SecurityEnhancements";
import { SECURITY_FEATURES } from "@/constants";
import { useState } from "react";
import {
  IoInformationCircle,
  IoRefresh,
  IoShieldCheckmark,
  IoWarning,
} from "react-icons/io5";

interface SecurityEvent {
  id: string;
  type: "login" | "logout" | "failed_login" | "password_change" | "mfa_enabled";
  description: string;
  timestamp: string;
  location: string;
  ip: string;
  severity: "low" | "medium" | "high";
}

export default function SecurityPage() {
  const [loading, setLoading] = useState(false);
  const [mfaEnabled, setMfaEnabled] = useState(false);
  const [securityEvents] = useState<SecurityEvent[]>([
    {
      id: "1",
      type: "login",
      description: "Successful login",
      timestamp: "2024-01-20T10:30:00Z",
      location: "New York, US",
      ip: "192.168.1.100",
      severity: "low",
    },
    {
      id: "2",
      type: "mfa_enabled",
      description: "Multi-factor authentication enabled",
      timestamp: "2024-01-19T15:45:00Z",
      location: "New York, US",
      ip: "192.168.1.100",
      severity: "low",
    },
    {
      id: "3",
      type: "failed_login",
      description: "Failed login attempt",
      timestamp: "2024-01-18T08:20:00Z",
      location: "Unknown",
      ip: "203.0.113.1",
      severity: "medium",
    },
    {
      id: "4",
      type: "password_change",
      description: "Password changed successfully",
      timestamp: "2024-01-17T14:15:00Z",
      location: "New York, US",
      ip: "192.168.1.100",
      severity: "low",
    },
  ]);

  const [securitySettings, setSecuritySettings] = useState({
    sessionTimeout: 30,
    loginNotifications: true,
    suspiciousActivityAlerts: true,
    passwordExpiry: 90,
  });

  const handleRefresh = () => {
    setLoading(true);
    // Simulate refresh
    setTimeout(() => setLoading(false), 1000);
  };

  const handleSettingChange = (setting: string, value: boolean | number) => {
    setSecuritySettings((prev) => ({
      ...prev,
      [setting]: value,
    }));
  };

  const getEventIcon = (type: SecurityEvent["type"]) => {
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
        return "🛡️";
      default:
        return "ℹ️";
    }
  };

  const getSeverityColor = (severity: SecurityEvent["severity"]) => {
    switch (severity) {
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

  const refreshAction = (
    <button
      onClick={handleRefresh}
      disabled={loading}
      className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
    >
      <IoRefresh className="h-4 w-4 mr-2" />
      {loading ? "Refreshing..." : "Refresh"}
    </button>
  );

  return (
    <DashboardLayout
      title="Security Center"
      description="Monitor and manage your security settings"
      actions={refreshAction}
    >
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
                {mfaEnabled ? "High" : "Medium"}
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
              <p className="text-2xl font-semibold text-gray-900">2</p>
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
              <p className="text-2xl font-semibold text-gray-900">1</p>
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
      <MFASetup onMFAStatusChange={setMfaEnabled} />

      {/* Security Enhancements */}
      <SecurityEnhancements />

      {/* Security Settings */}
      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <h3 className="text-lg font-medium text-gray-900 mb-4">
          Additional Security Settings
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
                checked={securitySettings.loginNotifications}
                onChange={(e) =>
                  handleSettingChange("loginNotifications", e.target.checked)
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
                checked={securitySettings.suspiciousActivityAlerts}
                onChange={(e) =>
                  handleSettingChange(
                    "suspiciousActivityAlerts",
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
              value={securitySettings.sessionTimeout}
              onChange={(e) =>
                handleSettingChange("sessionTimeout", parseInt(e.target.value))
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
              value={securitySettings.passwordExpiry}
              onChange={(e) =>
                handleSettingChange("passwordExpiry", parseInt(e.target.value))
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

      {/* Recent Security Events */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">
            Recent Security Events
          </h3>
          <p className="text-sm text-gray-500">Monitor your account activity</p>
        </div>
        <div className="divide-y divide-gray-200">
          {securityEvents.map((event) => (
            <div key={event.id} className="px-6 py-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <div className="text-2xl mr-4">
                    {getEventIcon(event.type)}
                  </div>
                  <div>
                    <h4 className="text-sm font-medium text-gray-900">
                      {event.description}
                    </h4>
                    <div className="flex items-center space-x-4 mt-1">
                      <span className="text-xs text-gray-500">
                        {new Date(event.timestamp).toLocaleString()}
                      </span>
                      <span className="text-xs text-gray-500">
                        {event.location}
                      </span>
                      <span className="text-xs text-gray-500">
                        IP: {event.ip}
                      </span>
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
          ))}
        </div>
      </div>
    </DashboardLayout>
  );
}
