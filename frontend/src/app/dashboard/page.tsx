"use client";

import DashboardLayout from "@/components/DashboardLayout";
import {
  DashboardMetrics,
  QuickAccess,
  QuickActions,
  RecentActivity,
} from "@/components/dashboard";
import {
  DASHBOARD_QUICK_ACTIONS,
  DEFAULT_APP_CONNECTIONS,
  DEFAULT_DASHBOARD_METRICS,
  DEFAULT_RECENT_ACTIVITY,
  ERROR_MESSAGES,
  FALLBACK_SAAS_APPS,
} from "@/constants";
import { apiClient } from "@/lib/api";
import { useKeycloak } from "@react-keycloak/web";
import { useCallback, useEffect, useState } from "react";
import { HiRefresh } from "react-icons/hi";

interface AppConnection {
  name: string;
  status: "connected" | "disconnected";
  icon: string;
  description: string;
  connect_url: string;
  last_used?: string;
}

interface ActivityItem {
  id: string;
  type: "login" | "app_launch" | "connection" | "security";
  description: string;
  timestamp: string;
  icon: string;
  severity?: "info" | "warning" | "success";
}

interface DashboardMetricsType {
  totalApps: number;
  connectedApps: number;
  recentLogins: number;
  securityScore: number;
  lastActivity: string;
}

export default function Dashboard() {
  const { keycloak } = useKeycloak();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [connections, setConnections] = useState<AppConnection[]>([
    ...DEFAULT_APP_CONNECTIONS,
  ]);
  const [recentActivity] = useState<ActivityItem[]>([
    ...DEFAULT_RECENT_ACTIVITY,
  ]);
  const [metrics, setMetrics] = useState<DashboardMetricsType>({
    ...DEFAULT_DASHBOARD_METRICS,
  });

  const updateMetrics = useCallback(() => {
    const connectedCount = connections.filter(
      (c) => c.status === "connected"
    ).length;
    setMetrics((prev) => ({
      ...prev,
      totalApps: connections.length,
      connectedApps: connectedCount,
    }));
  }, [connections]);

  const loadApps = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      // Try to load dashboard data from API
      try {
        const response = await apiClient.getDashboardData();
        if (response.success && response.data) {
          // Update state with API data
          setConnections(response.data.connections);
          setMetrics({
            totalApps: response.data.metrics.total_apps,
            connectedApps: response.data.metrics.connected_apps,
            recentLogins: response.data.metrics.recent_logins,
            securityScore: response.data.metrics.security_score,
            lastActivity: response.data.metrics.last_activity,
          });
          console.log("Dashboard data loaded from API:", response.data);
        }
      } catch (apiError) {
        console.warn(
          "Failed to load dashboard data from API, using fallback:",
          apiError
        );
        // Fallback to default data
        setConnections([...DEFAULT_APP_CONNECTIONS]);
        setMetrics({ ...DEFAULT_DASHBOARD_METRICS });
      }

      // Also try to load apps
      const appsResponse = await apiClient.getApps();
      console.log("Apps loaded:", appsResponse.apps);
    } catch (err) {
      console.error("Failed to load apps:", err);
      setError(ERROR_MESSAGES.NETWORK_ERROR);
      console.log("Using fallback apps:", FALLBACK_SAAS_APPS);
    } finally {
      setLoading(false);
    }
  }, []);

  // Initial load - runs only once on mount
  useEffect(() => {
    loadApps();
  }, [loadApps]);

  // Handle URL params for connection success - runs only once on mount
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const connected = urlParams.get("connected");
    const email = urlParams.get("email");

    if (connected && email) {
      // Update the connection status
      setConnections((prev) =>
        prev.map((conn) =>
          conn.name.toLowerCase().includes(connected)
            ? { ...conn, status: "connected" as const, last_used: "Just now" }
            : conn
        )
      );

      // Clean up URL
      window.history.replaceState({}, document.title, "/dashboard");
    }
  }, []); // Empty dependency array - runs only once

  // Update metrics when connections change
  useEffect(() => {
    updateMetrics();
  }, [updateMetrics]);

  const handleRefresh = () => {
    loadApps();
  };

  const getUserDisplayName = () => {
    if (keycloak?.tokenParsed) {
      return (
        keycloak.tokenParsed.preferred_username ||
        keycloak.tokenParsed.name ||
        keycloak.tokenParsed.email ||
        "User"
      );
    }
    return "User";
  };

  const refreshAction = (
    <button
      onClick={handleRefresh}
      disabled={loading}
      className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 cursor-pointer"
    >
      <HiRefresh className={`h-4 w-4 mr-2 ${loading ? "animate-spin" : ""}`} />
      Refresh
    </button>
  );

  return (
    <DashboardLayout
      title={`Welcome back, ${getUserDisplayName()}`}
      description="Your secure single sign-on dashboard"
      actions={refreshAction}
    >
      {/* Error Message */}
      {error && (
        <div className="mb-6 bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex">
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800">Error</h3>
              <p className="text-sm text-red-700 mt-1">{error}</p>
            </div>
            <button
              onClick={() => setError(null)}
              className="ml-auto text-red-400 hover:text-red-600 cursor-pointer"
            >
              ×
            </button>
          </div>
        </div>
      )}

      {/* Dashboard Metrics */}
      <DashboardMetrics metrics={metrics} />

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Quick Access */}
        <QuickAccess connections={connections} />

        {/* Recent Activity */}
        <RecentActivity activities={recentActivity} />
      </div>

      {/* Quick Actions */}
      <QuickActions actions={DASHBOARD_QUICK_ACTIONS} />
    </DashboardLayout>
  );
}
