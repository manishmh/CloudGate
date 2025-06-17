"use client";

import { type SaaSApplication } from "@/lib/api";
import { useEffect, useState } from "react";
import {
  IoCheckmarkCircle,
  IoCloseCircle,
  IoEye,
  IoRefresh,
  IoSettings,
  IoShieldCheckmark,
  IoTime,
  IoTrendingUp,
  IoWarning,
} from "react-icons/io5";

interface ConnectionHealth {
  status: "healthy" | "warning" | "error";
  lastCheck: string;
  responseTime: number;
  uptime: number;
  errorCount: number;
}

interface ConnectionStats {
  totalConnections: number;
  activeConnections: number;
  failedConnections: number;
  averageResponseTime: number;
  uptimePercentage: number;
}

interface EnhancedConnection extends SaaSApplication {
  health: ConnectionHealth;
  lastUsed?: string;
  usageCount: number;
  dataTransferred: string;
}

export default function OAuthConnectionDashboard() {
  const [connections, setConnections] = useState<EnhancedConnection[]>([]);
  const [stats, setStats] = useState<ConnectionStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [selectedConnection, setSelectedConnection] = useState<string | null>(
    null
  );
  const [autoRefresh, setAutoRefresh] = useState(true);

  useEffect(() => {
    loadConnections();

    if (autoRefresh) {
      const interval = setInterval(loadConnections, 30000); // Refresh every 30 seconds
      return () => clearInterval(interval);
    }
  }, [autoRefresh]);

  const loadConnections = async () => {
    try {
      setLoading(true);

      // Simulate enhanced connection data with health monitoring
      const mockConnections: EnhancedConnection[] = [
        {
          id: "google-workspace",
          name: "Google Workspace",
          icon: "🔍",
          description: "Access Gmail, Drive, Calendar, and more",
          category: "productivity",
          protocol: "oauth2",
          status: "connected",
          created_at: "2024-01-15T10:00:00Z",
          updated_at: "2024-01-20T10:30:00Z",
          connection_details: {
            user_email: "user@example.com",
            connected_at: "2024-01-15T10:00:00Z",
            last_used: "2024-01-20T10:30:00Z",
          },
          health: {
            status: "healthy",
            lastCheck: new Date().toISOString(),
            responseTime: 120,
            uptime: 99.9,
            errorCount: 0,
          },
          lastUsed: "2024-01-20T10:30:00Z",
          usageCount: 156,
          dataTransferred: "2.3 GB",
        },
        {
          id: "microsoft-365",
          name: "Microsoft 365",
          icon: "🏢",
          description: "Access Outlook, OneDrive, Teams, and more",
          category: "productivity",
          protocol: "oauth2",
          status: "connected",
          created_at: "2024-01-10T14:00:00Z",
          updated_at: "2024-01-19T16:45:00Z",
          connection_details: {
            user_email: "user@company.com",
            connected_at: "2024-01-10T14:00:00Z",
            last_used: "2024-01-19T16:45:00Z",
          },
          health: {
            status: "warning",
            lastCheck: new Date().toISOString(),
            responseTime: 850,
            uptime: 97.2,
            errorCount: 3,
          },
          lastUsed: "2024-01-19T16:45:00Z",
          usageCount: 89,
          dataTransferred: "1.7 GB",
        },
        {
          id: "slack",
          name: "Slack",
          icon: "💬",
          description: "Access your Slack workspaces",
          category: "communication",
          protocol: "oauth2",
          status: "error",
          created_at: "2024-01-12T09:00:00Z",
          updated_at: "2024-01-18T12:00:00Z",
          health: {
            status: "error",
            lastCheck: new Date().toISOString(),
            responseTime: 0,
            uptime: 85.5,
            errorCount: 12,
          },
          lastUsed: "2024-01-18T12:00:00Z",
          usageCount: 45,
          dataTransferred: "890 MB",
        },
        {
          id: "github",
          name: "GitHub",
          icon: "🐙",
          description: "Access your repositories and organizations",
          category: "development",
          protocol: "oauth2",
          status: "connected",
          created_at: "2024-01-08T11:00:00Z",
          updated_at: "2024-01-20T09:15:00Z",
          connection_details: {
            user_email: "dev@example.com",
            connected_at: "2024-01-08T11:00:00Z",
            last_used: "2024-01-20T09:15:00Z",
          },
          health: {
            status: "healthy",
            lastCheck: new Date().toISOString(),
            responseTime: 95,
            uptime: 99.8,
            errorCount: 0,
          },
          lastUsed: "2024-01-20T09:15:00Z",
          usageCount: 234,
          dataTransferred: "4.1 GB",
        },
        {
          id: "trello",
          name: "Trello",
          icon: "📋",
          description: "Manage your boards and projects",
          category: "productivity",
          protocol: "oauth1",
          status: "available",
          created_at: "2024-01-01T00:00:00Z",
          updated_at: "2024-01-01T00:00:00Z",
          health: {
            status: "healthy",
            lastCheck: new Date().toISOString(),
            responseTime: 200,
            uptime: 100,
            errorCount: 0,
          },
          usageCount: 0,
          dataTransferred: "0 MB",
        },
      ];

      setConnections(mockConnections);

      // Calculate stats
      const totalConnections = mockConnections.length;
      const activeConnections = mockConnections.filter(
        (c) => c.status === "connected"
      ).length;
      const failedConnections = mockConnections.filter(
        (c) => c.status === "error"
      ).length;
      const avgResponseTime =
        mockConnections
          .filter((c) => c.status === "connected")
          .reduce((sum, c) => sum + c.health.responseTime, 0) /
          activeConnections || 0;
      const avgUptime =
        mockConnections
          .filter((c) => c.status === "connected")
          .reduce((sum, c) => sum + c.health.uptime, 0) / activeConnections ||
        0;

      setStats({
        totalConnections,
        activeConnections,
        failedConnections,
        averageResponseTime: Math.round(avgResponseTime),
        uptimePercentage: Math.round(avgUptime * 10) / 10,
      });
    } catch (error) {
      console.error("Failed to load connections:", error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "connected":
        return <IoCheckmarkCircle className="h-5 w-5 text-green-500" />;
      case "error":
        return <IoCloseCircle className="h-5 w-5 text-red-500" />;
      case "warning":
        return <IoWarning className="h-5 w-5 text-yellow-500" />;
      default:
        return <IoTime className="h-5 w-5 text-gray-400" />;
    }
  };

  const getHealthColor = (status: string) => {
    switch (status) {
      case "healthy":
        return "text-green-600 bg-green-100";
      case "warning":
        return "text-yellow-600 bg-yellow-100";
      case "error":
        return "text-red-600 bg-red-100";
      default:
        return "text-gray-600 bg-gray-100";
    }
  };

  const testConnection = async (connectionId: string) => {
    setLoading(true);
    // Simulate connection test
    setTimeout(() => {
      setConnections((prev) =>
        prev.map((conn) =>
          conn.id === connectionId
            ? {
                ...conn,
                health: {
                  ...conn.health,
                  lastCheck: new Date().toISOString(),
                  responseTime: Math.floor(Math.random() * 500) + 50,
                },
              }
            : conn
        )
      );
      setLoading(false);
    }, 2000);
  };

  if (loading && connections.length === 0) {
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
      {/* Connection Statistics */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
          <div className="bg-white rounded-lg shadow p-4">
            <div className="flex items-center">
              <div className="p-2 bg-blue-100 rounded-lg">
                <IoSettings className="h-5 w-5 text-blue-600" />
              </div>
              <div className="ml-3">
                <p className="text-sm font-medium text-gray-600">Total</p>
                <p className="text-xl font-semibold text-gray-900">
                  {stats.totalConnections}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-4">
            <div className="flex items-center">
              <div className="p-2 bg-green-100 rounded-lg">
                <IoCheckmarkCircle className="h-5 w-5 text-green-600" />
              </div>
              <div className="ml-3">
                <p className="text-sm font-medium text-gray-600">Active</p>
                <p className="text-xl font-semibold text-green-600">
                  {stats.activeConnections}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-4">
            <div className="flex items-center">
              <div className="p-2 bg-red-100 rounded-lg">
                <IoCloseCircle className="h-5 w-5 text-red-600" />
              </div>
              <div className="ml-3">
                <p className="text-sm font-medium text-gray-600">Failed</p>
                <p className="text-xl font-semibold text-red-600">
                  {stats.failedConnections}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-4">
            <div className="flex items-center">
              <div className="p-2 bg-purple-100 rounded-lg">
                <IoTrendingUp className="h-5 w-5 text-purple-600" />
              </div>
              <div className="ml-3">
                <p className="text-sm font-medium text-gray-600">
                  Avg Response
                </p>
                <p className="text-xl font-semibold text-purple-600">
                  {stats.averageResponseTime}ms
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-4">
            <div className="flex items-center">
              <div className="p-2 bg-indigo-100 rounded-lg">
                <IoShieldCheckmark className="h-5 w-5 text-indigo-600" />
              </div>
              <div className="ml-3">
                <p className="text-sm font-medium text-gray-600">Uptime</p>
                <p className="text-xl font-semibold text-indigo-600">
                  {stats.uptimePercentage}%
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Connection Management */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-medium text-gray-900">
                OAuth Connections
              </h3>
              <p className="text-sm text-gray-500">
                Monitor and manage your application connections
              </p>
            </div>
            <div className="flex items-center space-x-3">
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={autoRefresh}
                  onChange={(e) => setAutoRefresh(e.target.checked)}
                  className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
                <span className="ml-2 text-sm text-gray-600">Auto-refresh</span>
              </label>
              <button
                onClick={loadConnections}
                disabled={loading}
                className="p-2 text-gray-500 hover:text-gray-700 disabled:opacity-50"
                title="Refresh connections"
              >
                <IoRefresh
                  className={`h-5 w-5 ${loading ? "animate-spin" : ""}`}
                />
              </button>
            </div>
          </div>
        </div>

        <div className="divide-y divide-gray-200">
          {connections.map((connection) => (
            <div key={connection.id} className="p-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <div className="text-3xl">{connection.icon}</div>
                  <div>
                    <div className="flex items-center space-x-2">
                      <h4 className="text-lg font-medium text-gray-900">
                        {connection.name}
                      </h4>
                      {getStatusIcon(connection.status)}
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getHealthColor(
                          connection.health.status
                        )}`}
                      >
                        {connection.health.status.toUpperCase()}
                      </span>
                    </div>
                    <p className="text-sm text-gray-500">
                      {connection.description}
                    </p>
                    {connection.connection_details?.user_email && (
                      <p className="text-xs text-gray-400">
                        Connected as: {connection.connection_details.user_email}
                      </p>
                    )}
                  </div>
                </div>

                <div className="flex items-center space-x-4">
                  <button
                    onClick={() =>
                      setSelectedConnection(
                        selectedConnection === connection.id
                          ? null
                          : connection.id
                      )
                    }
                    className="text-blue-600 hover:text-blue-700 text-sm font-medium flex items-center"
                  >
                    <IoEye className="h-4 w-4 mr-1" />
                    Details
                  </button>
                  <button
                    onClick={() => testConnection(connection.id)}
                    disabled={loading}
                    className="text-green-600 hover:text-green-700 text-sm font-medium disabled:opacity-50"
                  >
                    Test
                  </button>
                </div>
              </div>

              {selectedConnection === connection.id && (
                <div className="mt-4 p-4 bg-gray-50 rounded-lg">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                      <h5 className="text-sm font-medium text-gray-900 mb-2">
                        Health Metrics
                      </h5>
                      <div className="space-y-1 text-sm text-gray-600">
                        <div>
                          Response Time: {connection.health.responseTime}ms
                        </div>
                        <div>Uptime: {connection.health.uptime}%</div>
                        <div>Errors: {connection.health.errorCount}</div>
                        <div>
                          Last Check:{" "}
                          {new Date(
                            connection.health.lastCheck
                          ).toLocaleString()}
                        </div>
                      </div>
                    </div>

                    <div>
                      <h5 className="text-sm font-medium text-gray-900 mb-2">
                        Usage Statistics
                      </h5>
                      <div className="space-y-1 text-sm text-gray-600">
                        <div>Usage Count: {connection.usageCount}</div>
                        <div>
                          Data Transferred: {connection.dataTransferred}
                        </div>
                        {connection.lastUsed && (
                          <div>
                            Last Used:{" "}
                            {new Date(connection.lastUsed).toLocaleString()}
                          </div>
                        )}
                        <div>Protocol: {connection.protocol.toUpperCase()}</div>
                      </div>
                    </div>

                    <div>
                      <h5 className="text-sm font-medium text-gray-900 mb-2">
                        Connection Info
                      </h5>
                      <div className="space-y-1 text-sm text-gray-600">
                        <div>Status: {connection.status}</div>
                        <div>Category: {connection.category}</div>
                        <div>
                          Connected:{" "}
                          {new Date(connection.created_at).toLocaleDateString()}
                        </div>
                        <div>
                          Updated:{" "}
                          {new Date(connection.updated_at).toLocaleDateString()}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
