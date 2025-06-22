"use client";

import DashboardLayout from "@/components/DashboardLayout";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import {
  setError,
  setLoading,
  setSaved,
  updateSetting,
  updateSettings,
} from "@/store/slices/settingsSlice";
import { useEffect, useState } from "react";

export default function SettingsPage() {
  const dispatch = useAppDispatch();
  const { settings, loading, error, lastSaved } = useAppSelector(
    (state) => state.settings
  );
  const [localLoading, setLocalLoading] = useState(false);

  // Load settings from backend on mount
  useEffect(() => {
    const loadSettings = async () => {
      try {
        dispatch(setLoading(true));
        // TODO: Replace with actual API call
        // const response = await fetch('/api/user/settings');
        // const data = await response.json();
        // dispatch(setSettings(data.settings));
      } catch {
        dispatch(setError("Failed to load settings"));
      } finally {
        dispatch(setLoading(false));
      }
    };

    loadSettings();
  }, [dispatch]);

  const handleSettingChange = (key: string, value: unknown) => {
    dispatch(updateSetting({ key: key as keyof typeof settings, value }));
  };

  const handleSave = async () => {
    try {
      setLocalLoading(true);
      dispatch(setError(null));

      // TODO: Replace with actual API call
      // const response = await fetch('/api/user/settings', {
      //   method: 'PUT',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify(settings)
      // });
      //
      // if (!response.ok) {
      //   throw new Error('Failed to save settings');
      // }

      // Simulate save operation
      await new Promise((resolve) => setTimeout(resolve, 1500));

      dispatch(setSaved());
    } catch {
      dispatch(setError("Failed to save settings"));
    } finally {
      setLocalLoading(false);
    }
  };

  const handleReset = async () => {
    if (
      confirm("Are you sure you want to reset all settings to default values?")
    ) {
      try {
        setLocalLoading(true);
        dispatch(setError(null));

        // TODO: Replace with actual API call
        // const response = await fetch('/api/user/settings/reset', {
        //   method: 'POST'
        // });
        //
        // if (!response.ok) {
        //   throw new Error('Failed to reset settings');
        // }
        //
        // const data = await response.json();
        // dispatch(setSettings(data.settings));

        // Simulate reset operation
        await new Promise((resolve) => setTimeout(resolve, 1000));

        // Reset to default values
        dispatch(
          updateSettings({
            language: "en",
            timezone: "America/New_York",
            dateFormat: "MM/DD/YYYY",
            emailNotifications: true,
            pushNotifications: false,
            securityAlerts: true,
            appUpdates: true,
            weeklyReports: false,
            defaultView: "dashboard",
            itemsPerPage: 10,
            autoRefresh: true,
            refreshInterval: 30,
            analyticsOptIn: true,
            shareUsageData: false,
            personalizedAds: false,
            apiAccess: false,
            webhookUrl: "",
            maxApiCalls: 1000,
          })
        );

        dispatch(setSaved());
      } catch {
        dispatch(setError("Failed to reset settings"));
      } finally {
        setLocalLoading(false);
      }
    }
  };

  if (loading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      {/* Error Message */}
      {error && (
        <div className="mb-6 bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex justify-between items-center">
            <p className="text-sm text-red-800">{error}</p>
            <button
              onClick={() => dispatch(setError(null))}
              className="text-red-400 hover:text-red-600 cursor-pointer"
            >
              ×
            </button>
          </div>
        </div>
      )}

      {/* Appearance Settings */}
      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <h3 className="text-lg font-medium text-gray-900 mb-6">Appearance</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Language
            </label>
            <select
              value={settings.language}
              onChange={(e) => handleSettingChange("language", e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
            >
              <option value="en">English</option>
              <option value="es">Spanish</option>
              <option value="fr">French</option>
              <option value="de">German</option>
              <option value="ja">Japanese</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Timezone
            </label>
            <select
              value={settings.timezone}
              onChange={(e) => handleSettingChange("timezone", e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
            >
              <option value="America/New_York">Eastern Time</option>
              <option value="America/Chicago">Central Time</option>
              <option value="America/Denver">Mountain Time</option>
              <option value="America/Los_Angeles">Pacific Time</option>
              <option value="Europe/London">London</option>
              <option value="Europe/Paris">Paris</option>
              <option value="Asia/Tokyo">Tokyo</option>
              <option value="Asia/Shanghai">Shanghai</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Date Format
            </label>
            <select
              value={settings.dateFormat}
              onChange={(e) =>
                handleSettingChange("dateFormat", e.target.value)
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
            >
              <option value="MM/DD/YYYY">MM/DD/YYYY</option>
              <option value="DD/MM/YYYY">DD/MM/YYYY</option>
              <option value="YYYY-MM-DD">YYYY-MM-DD</option>
              <option value="DD MMM YYYY">DD MMM YYYY</option>
            </select>
          </div>
        </div>
      </div>

      {/* Notifications Settings */}
      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <h3 className="text-lg font-medium text-gray-900 mb-6">
          Notifications
        </h3>
        <div className="space-y-6">
          {[
            {
              key: "emailNotifications",
              label: "Email Notifications",
              description: "Receive notifications via email",
            },
            {
              key: "pushNotifications",
              label: "Push Notifications",
              description: "Receive browser push notifications",
            },
            {
              key: "securityAlerts",
              label: "Security Alerts",
              description: "Get notified about security events",
            },
            {
              key: "appUpdates",
              label: "Application Updates",
              description: "Notifications about app status changes",
            },
            {
              key: "weeklyReports",
              label: "Weekly Reports",
              description: "Receive weekly activity summaries",
            },
          ].map((item) => (
            <div key={item.key} className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-gray-900">
                  {item.label}
                </h4>
                <p className="text-sm text-gray-500">{item.description}</p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  checked={
                    settings[item.key as keyof typeof settings] as boolean
                  }
                  onChange={(e) =>
                    handleSettingChange(item.key, e.target.checked)
                  }
                  className="sr-only peer"
                />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
          ))}
        </div>
      </div>

      {/* Dashboard Settings */}
      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <h3 className="text-lg font-medium text-gray-900 mb-6">Dashboard</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Default View
            </label>
            <select
              value={settings.defaultView}
              onChange={(e) =>
                handleSettingChange("defaultView", e.target.value)
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
            >
              <option value="dashboard">Dashboard</option>
              <option value="applications">Applications</option>
              <option value="connections">Connections</option>
              <option value="security">Security</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Items Per Page
            </label>
            <select
              value={settings.itemsPerPage}
              onChange={(e) =>
                handleSettingChange("itemsPerPage", parseInt(e.target.value))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
            >
              <option value={5}>5</option>
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
            </select>
          </div>

          <div className="md:col-span-2">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-gray-900">
                  Auto Refresh
                </h4>
                <p className="text-sm text-gray-500">
                  Automatically refresh dashboard data
                </p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  checked={settings.autoRefresh}
                  onChange={(e) =>
                    handleSettingChange("autoRefresh", e.target.checked)
                  }
                  className="sr-only peer"
                />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
          </div>

          {settings.autoRefresh && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Refresh Interval (seconds)
              </label>
              <select
                value={settings.refreshInterval}
                onChange={(e) =>
                  handleSettingChange(
                    "refreshInterval",
                    parseInt(e.target.value)
                  )
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
              >
                <option value={15}>15 seconds</option>
                <option value={30}>30 seconds</option>
                <option value={60}>1 minute</option>
                <option value={300}>5 minutes</option>
              </select>
            </div>
          )}
        </div>
      </div>

      {/* Privacy Settings */}
      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <h3 className="text-lg font-medium text-gray-900 mb-6">Privacy</h3>
        <div className="space-y-6">
          {[
            {
              key: "analyticsOptIn",
              label: "Analytics",
              description: "Help improve CloudGate by sharing usage analytics",
            },
            {
              key: "shareUsageData",
              label: "Usage Data",
              description:
                "Share anonymized usage data for product improvement",
            },
            {
              key: "personalizedAds",
              label: "Personalized Ads",
              description: "Show personalized advertisements",
            },
          ].map((item) => (
            <div key={item.key} className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-gray-900">
                  {item.label}
                </h4>
                <p className="text-sm text-gray-500">{item.description}</p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  checked={
                    settings[item.key as keyof typeof settings] as boolean
                  }
                  onChange={(e) =>
                    handleSettingChange(item.key, e.target.checked)
                  }
                  className="sr-only peer"
                />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
          ))}
        </div>
      </div>

      {/* Integration Settings */}
      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <h3 className="text-lg font-medium text-gray-900 mb-6">Integration</h3>
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <div>
              <h4 className="text-sm font-medium text-gray-900">API Access</h4>
              <p className="text-sm text-gray-500">
                Enable API access for third-party integrations
              </p>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={settings.apiAccess}
                onChange={(e) =>
                  handleSettingChange("apiAccess", e.target.checked)
                }
                className="sr-only peer"
              />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            </label>
          </div>

          {settings.apiAccess && (
            <>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Webhook URL
                </label>
                <input
                  type="url"
                  value={settings.webhookUrl}
                  onChange={(e) =>
                    handleSettingChange("webhookUrl", e.target.value)
                  }
                  placeholder="https://your-app.com/webhook"
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
                />
                <p className="mt-1 text-sm text-gray-500">
                  URL to receive webhook notifications
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Max API Calls per Hour
                </label>
                <select
                  value={settings.maxApiCalls}
                  onChange={(e) =>
                    handleSettingChange("maxApiCalls", parseInt(e.target.value))
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
                >
                  <option value={100}>100</option>
                  <option value={500}>500</option>
                  <option value={1000}>1,000</option>
                  <option value={5000}>5,000</option>
                  <option value={10000}>10,000</option>
                </select>
              </div>
            </>
          )}
        </div>
      </div>
    </DashboardLayout>
  );
}
