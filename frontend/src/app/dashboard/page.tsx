"use client";

import DashboardLayout from "@/components/DashboardLayout";
import {
  DashboardMetrics,
  QuickAccess,
  QuickActions,
  RecentActivity,
} from "@/components/dashboard";
import { useAuth } from "@/components/providers/AuthProvider";
import {
  DASHBOARD_QUICK_ACTIONS,
  DEFAULT_APP_CONNECTIONS,
  DEFAULT_DASHBOARD_METRICS,
  ERROR_MESSAGES,
  FALLBACK_SAAS_APPS,
} from "@/constants";
import { apiClient } from "@/lib/api";

import { useCallback, useEffect, useState } from "react";

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

// Helper function to map security event types to activity types and icons
const mapSecurityEventToActivity = (event: any): ActivityItem => {
  let activityType: ActivityItem["type"] = "security";
  let icon = "HiShieldCheck";
  let severity: ActivityItem["severity"] = "info";

  // Map event types to appropriate activity types and icons
  switch (event.event_type.toLowerCase()) {
    case "login":
    case "authentication":
    case "login_success":
    case "login_failure":
      activityType = "login";
      icon =
        event.event_type.includes("failure") ||
        event.event_type.includes("failed")
          ? "HiExclamationCircle"
          : "HiShieldCheck";
      severity =
        event.event_type.includes("failure") ||
        event.event_type.includes("failed")
          ? "warning"
          : "success";
      break;
    case "app_launch":
    case "application_access":
      activityType = "app_launch";
      icon = "HiViewGrid";
      severity = "info";
      break;
    case "connection":
    case "oauth_connection":
    case "app_connection":
      activityType = "connection";
      icon = "HiLink";
      severity = "success";
      break;
    default:
      activityType = "security";
      icon = "HiShieldCheck";
      severity =
        event.severity === "warning" || event.severity === "error"
          ? "warning"
          : event.severity === "critical"
          ? "warning"
          : "success";
  }

  // Format timestamp to relative time
  const formatTimestamp = (dateString: string): string => {
    const now = new Date();
    const eventDate = new Date(dateString);
    const diffInMinutes = Math.floor(
      (now.getTime() - eventDate.getTime()) / (1000 * 60)
    );

    if (diffInMinutes < 1) return "Just now";
    if (diffInMinutes < 60)
      return `${diffInMinutes} minute${diffInMinutes === 1 ? "" : "s"} ago`;

    const diffInHours = Math.floor(diffInMinutes / 60);
    if (diffInHours < 24)
      return `${diffInHours} hour${diffInHours === 1 ? "" : "s"} ago`;

    const diffInDays = Math.floor(diffInHours / 24);
    return `${diffInDays} day${diffInDays === 1 ? "" : "s"} ago`;
  };

  return {
    id: event.id,
    type: activityType,
    description: event.description,
    timestamp: formatTimestamp(event.created_at),
    icon,
    severity,
  };
};

export default function Dashboard() {
  const { isAuthenticated } = useAuth();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [connections, setConnections] = useState<AppConnection[]>([
    ...DEFAULT_APP_CONNECTIONS,
  ]);
  const [recentActivity, setRecentActivity] = useState<ActivityItem[]>([]);
  const [loadingActivity, setLoadingActivity] = useState(true);
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

  const loadRecentActivity = useCallback(async () => {
    try {
      setLoadingActivity(true);

      // Fetch recent security events (limit to 10 to get enough for 5 recent activities)
      const response = await apiClient.getSecurityEvents(10);

      if (response.events && response.events.length > 0) {
        // Transform security events to activity items
        const activities = response.events
          .map(mapSecurityEventToActivity)
          .slice(0, 5); // Only show 5 most recent

        setRecentActivity(activities);

        // Update last activity time in metrics if we have activities
        if (activities.length > 0) {
          setMetrics((prev) => ({
            ...prev,
            lastActivity: activities[0].timestamp,
          }));
        }
      } else {
        // Fallback to empty array if no events
        setRecentActivity([]);
      }
    } catch (error) {
      console.warn(
        "Failed to load recent activity from API, using empty array:",
        error
      );
      setRecentActivity([]);
    } finally {
      setLoadingActivity(false);
    }
  }, []);

  const loadApps = useCallback(async () => {
    try {
      setError(null);
      setLoading(true);

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

  // Initial load - only after authentication is ready
  useEffect(() => {
    if (!isAuthenticated) return;
    loadApps();
    loadRecentActivity();
  }, [isAuthenticated, loadApps, loadRecentActivity]);

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

      // Reload recent activity to include the new connection
      loadRecentActivity();
    }
  }, [loadRecentActivity]); // Include loadRecentActivity in dependencies

  // Update metrics when connections change
  useEffect(() => {
    updateMetrics();
  }, [updateMetrics]);

  if (loading) {
    return (
      <DashboardLayout>
        <div className="space-y-6">
          {/* Metrics Skeleton */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="bg-white rounded-lg shadow p-6">
                <div className="animate-pulse">
                  <div className="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
                  <div className="h-8 bg-gray-200 rounded w-1/3"></div>
                </div>
              </div>
            ))}
          </div>

          {/* Content Skeleton */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {[1, 2].map((i) => (
              <div key={i} className="bg-white rounded-lg shadow p-6">
                <div className="animate-pulse">
                  <div className="h-4 bg-gray-200 rounded w-1/3 mb-4"></div>
                  <div className="space-y-3">
                    {[1, 2, 3].map((j) => (
                      <div key={j} className="h-4 bg-gray-200 rounded"></div>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
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
        {loadingActivity ? (
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">
              Recent Activity
            </h3>
            <div className="space-y-4">
              {[1, 2, 3, 4, 5].map((i) => (
                <div
                  key={i}
                  className="animate-pulse flex items-start space-x-3"
                >
                  <div className="h-6 w-6 bg-gray-200 rounded-full"></div>
                  <div className="flex-1 space-y-2">
                    <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                    <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        ) : (
          <RecentActivity activities={recentActivity} />
        )}
      </div>

      {/* Quick Actions */}
      <QuickActions actions={DASHBOARD_QUICK_ACTIONS} />
    </DashboardLayout>
  );
}
