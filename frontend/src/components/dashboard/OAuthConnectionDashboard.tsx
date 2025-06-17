"use client";

import {
  apiClient,
  type ConnectionStats,
  type EnhancedConnection,
} from "@/lib/api";
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

export default function OAuthConnectionDashboard() {
  const [connections, setConnections] = useState<EnhancedConnection[]>([]);
  const [stats, setStats] = useState<ConnectionStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [selectedConnection, setSelectedConnection] = useState<string | null>(
    null
  );
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadConnections();

    if (autoRefresh) {
      const interval = setInterval(loadConnections, 30000); // Refresh every 30 seconds
      return () => clearInterval(interval);
    }
  }, [autoRefresh]);

  // Auto-expand first connection details
  useEffect(() => {
    if (connections.length > 0 && !selectedConnection) {
      setSelectedConnection(connections[0].id);
    }
  }, [connections, selectedConnection]);

  const loadConnections = async () => {
    try {
      setLoading(true);
      setError(null);

      // Load connections and stats in parallel
      const [connectionsResponse, statsResponse] = await Promise.all([
        apiClient.getConnections(),
        apiClient.getConnectionStats(),
      ]);

      setConnections(connectionsResponse.connections);
      setStats(statsResponse);
    } catch (err) {
      console.error("Failed to load connections:", err);
      setError("Failed to load connection data. Please try again.");
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
    try {
      await apiClient.testConnection(connectionId);
      // Reload connections to get updated health data
      await loadConnections();
    } catch (err) {
      console.error("Failed to test connection:", err);
      setError("Failed to test connection. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const getConnectionIcon = (provider: string) => {
    const iconMap: Record<string, string> = {
      google: "🔍",
      microsoft: "🏢",
      slack: "💬",
      github: "🐙",
      trello: "📋",
      salesforce: "☁️",
      jira: "🎯",
      notion: "📝",
      dropbox: "📦",
    };
    return iconMap[provider] || "🔗";
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

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-6">
        <div className="flex items-center">
          <IoCloseCircle className="h-6 w-6 text-red-500 mr-3" />
          <div>
            <h3 className="text-lg font-medium text-red-800">
              Error Loading Connections
            </h3>
            <p className="text-red-600">{error}</p>
            <button
              onClick={loadConnections}
              className="mt-2 text-red-700 hover:text-red-900 font-medium"
            >
              Try Again
            </button>
          </div>
        </div>
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
                  {stats.total_connections}
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
                  {stats.active_connections}
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
                  {stats.failed_connections}
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
                  {stats.average_response_time}ms
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
                  {stats.uptime_percentage.toFixed(1)}%
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
          {connections.length === 0 ? (
            <div className="px-6 py-12 text-center">
              <div className="text-gray-400 text-4xl mb-4">🔗</div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">
                No connections found
              </h3>
              <p className="text-gray-500">
                Connect to applications to see them here.
              </p>
            </div>
          ) : (
            connections.map((connection) => (
              <div key={connection.id} className="p-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-4">
                    <div className="text-3xl">
                      {getConnectionIcon(connection.provider)}
                    </div>
                    <div>
                      <div className="flex items-center space-x-2">
                        <h4 className="text-lg font-medium text-gray-900">
                          {connection.app_name}
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
                        Provider: {connection.provider}
                      </p>
                      {connection.user_email && (
                        <p className="text-xs text-gray-400">
                          Connected as: {connection.user_email}
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
                            Response Time: {connection.health.response_time}ms
                          </div>
                          <div>Uptime: {connection.health.uptime}%</div>
                          <div>Errors: {connection.health.error_count}</div>
                          <div>
                            Last Check:{" "}
                            {connection.health.last_check
                              ? new Date(
                                  connection.health.last_check
                                ).toLocaleString()
                              : "Never"}
                          </div>
                        </div>
                      </div>

                      <div>
                        <h5 className="text-sm font-medium text-gray-900 mb-2">
                          Usage Statistics
                        </h5>
                        <div className="space-y-1 text-sm text-gray-600">
                          <div>Usage Count: {connection.usage_count}</div>
                          <div>
                            Data Transferred: {connection.data_transferred}
                          </div>
                          {connection.last_used && (
                            <div>
                              Last Used:{" "}
                              {new Date(connection.last_used).toLocaleString()}
                            </div>
                          )}
                          <div>App ID: {connection.app_id}</div>
                        </div>
                      </div>

                      <div>
                        <h5 className="text-sm font-medium text-gray-900 mb-2">
                          Connection Info
                        </h5>
                        <div className="space-y-1 text-sm text-gray-600">
                          <div>Status: {connection.status}</div>
                          <div>Provider: {connection.provider}</div>
                          <div>
                            Connected:{" "}
                            {new Date(
                              connection.connected_at
                            ).toLocaleDateString()}
                          </div>
                          <div>
                            Updated:{" "}
                            {new Date(
                              connection.updated_at
                            ).toLocaleDateString()}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
