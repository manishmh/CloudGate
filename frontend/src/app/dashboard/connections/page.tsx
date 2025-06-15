"use client";

import DashboardLayout from "@/components/DashboardLayout";
import { useState } from "react";
import { IoRefresh } from "react-icons/io5";

interface AppConnection {
  name: string;
  status: "connected" | "disconnected";
  icon: string;
  description: string;
  connectUrl: string;
  connectedAt?: string;
  lastUsed?: string;
}

export default function ConnectionsPage() {
  const [connections, setConnections] = useState<AppConnection[]>([
    {
      name: "Google Workspace",
      status: "disconnected",
      icon: "🔍",
      description: "Access Gmail, Drive, Calendar, and more",
      connectUrl: "/oauth/google/connect",
    },
    {
      name: "Microsoft 365",
      status: "disconnected",
      icon: "🏢",
      description: "Access Outlook, OneDrive, Teams, and more",
      connectUrl: "/oauth/microsoft/connect",
    },
    {
      name: "Slack",
      status: "connected",
      icon: "💬",
      description: "Access your Slack workspaces",
      connectUrl: "/oauth/slack/connect",
      connectedAt: "2024-01-15",
      lastUsed: "2024-01-20",
    },
    {
      name: "GitHub",
      status: "connected",
      icon: "🐙",
      description: "Access your repositories and organizations",
      connectUrl: "/oauth/github/connect",
      connectedAt: "2024-01-10",
      lastUsed: "2024-01-19",
    },
    {
      name: "Salesforce",
      status: "disconnected",
      icon: "☁️",
      description: "Customer relationship management platform",
      connectUrl: "/oauth/salesforce/connect",
    },
    {
      name: "Jira",
      status: "connected",
      icon: "🎯",
      description: "Issue tracking and project management",
      connectUrl: "/oauth/jira/connect",
      connectedAt: "2024-01-12",
      lastUsed: "2024-01-18",
    },
    {
      name: "Confluence",
      status: "disconnected",
      icon: "📝",
      description: "Team workspace and knowledge management",
      connectUrl: "/oauth/confluence/connect",
    },
    {
      name: "Dropbox",
      status: "disconnected",
      icon: "📦",
      description: "Cloud storage and file synchronization",
      connectUrl: "/oauth/dropbox/connect",
    },
  ]);

  const [loading, setLoading] = useState(false);

  const handleConnect = async (connectionName: string) => {
    setLoading(true);
    // Simulate connection process
    await new Promise((resolve) => setTimeout(resolve, 1500));

    setConnections((prev) =>
      prev.map((conn) =>
        conn.name === connectionName
          ? {
              ...conn,
              status: "connected" as const,
              connectedAt: new Date().toISOString().split("T")[0],
              lastUsed: new Date().toISOString().split("T")[0],
            }
          : conn
      )
    );
    setLoading(false);
    alert(`Successfully connected to ${connectionName}!`);
  };

  const handleDisconnect = async (connectionName: string) => {
    if (
      !confirm(`Are you sure you want to disconnect from ${connectionName}?`)
    ) {
      return;
    }

    setLoading(true);
    // Simulate disconnection process
    await new Promise((resolve) => setTimeout(resolve, 1000));

    setConnections((prev) =>
      prev.map((conn) =>
        conn.name === connectionName
          ? {
              ...conn,
              status: "disconnected" as const,
              connectedAt: undefined,
              lastUsed: undefined,
            }
          : conn
      )
    );
    setLoading(false);
    alert(`Disconnected from ${connectionName}`);
  };

  const handleRefresh = () => {
    setLoading(true);
    // Simulate refresh
    setTimeout(() => setLoading(false), 1000);
  };

  const connectedCount = connections.filter(
    (c) => c.status === "connected"
  ).length;
  const totalCount = connections.length;

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
      title="App Connections"
      description={`${connectedCount} of ${totalCount} applications connected`}
      actions={refreshAction}
    >
      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-green-100 rounded-lg">
              <span className="text-2xl">🔗</span>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Connected</p>
              <p className="text-2xl font-semibold text-gray-900">
                {connectedCount}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-gray-100 rounded-lg">
              <span className="text-2xl">⏸️</span>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Available</p>
              <p className="text-2xl font-semibold text-gray-900">
                {totalCount - connectedCount}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-blue-100 rounded-lg">
              <span className="text-2xl">📊</span>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total</p>
              <p className="text-2xl font-semibold text-gray-900">
                {totalCount}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Connected Applications */}
      <div className="bg-white rounded-lg shadow mb-8">
        <div className="px-6 py-4 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">
            Connected Applications
          </h3>
          <p className="text-sm text-gray-500">
            Applications you&apos;re currently connected to
          </p>
        </div>
        <div className="divide-y divide-gray-200">
          {connections.filter((conn) => conn.status === "connected").length ===
          0 ? (
            <div className="px-6 py-12 text-center">
              <div className="text-gray-400 text-4xl mb-4">🔌</div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">
                No connected applications
              </h3>
              <p className="text-gray-500">
                Connect to applications below to get started.
              </p>
            </div>
          ) : (
            connections
              .filter((conn) => conn.status === "connected")
              .map((connection) => (
                <div key={connection.name} className="px-6 py-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center">
                      <div className="text-3xl mr-4">{connection.icon}</div>
                      <div>
                        <h4 className="text-lg font-medium text-gray-900">
                          {connection.name}
                        </h4>
                        <p className="text-sm text-gray-500">
                          {connection.description}
                        </p>
                        <div className="flex items-center space-x-4 mt-1">
                          {connection.connectedAt && (
                            <span className="text-xs text-gray-400">
                              Connected:{" "}
                              {new Date(
                                connection.connectedAt
                              ).toLocaleDateString()}
                            </span>
                          )}
                          {connection.lastUsed && (
                            <span className="text-xs text-gray-400">
                              Last used:{" "}
                              {new Date(
                                connection.lastUsed
                              ).toLocaleDateString()}
                            </span>
                          )}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center space-x-3">
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        ✓ Connected
                      </span>
                      <button
                        onClick={() => handleDisconnect(connection.name)}
                        disabled={loading}
                        className="bg-red-600 text-white px-4 py-2 rounded-md text-sm hover:bg-red-700 transition-colors disabled:opacity-50"
                      >
                        Disconnect
                      </button>
                    </div>
                  </div>
                </div>
              ))
          )}
        </div>
      </div>

      {/* Available Applications */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">
            Available Applications
          </h3>
          <p className="text-sm text-gray-500">
            Applications you can connect to
          </p>
        </div>
        <div className="divide-y divide-gray-200">
          {connections.filter((conn) => conn.status === "disconnected")
            .length === 0 ? (
            <div className="px-6 py-12 text-center">
              <div className="text-gray-400 text-4xl mb-4">✅</div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">
                All applications connected
              </h3>
              <p className="text-gray-500">
                You&apos;ve connected to all available applications.
              </p>
            </div>
          ) : (
            connections
              .filter((conn) => conn.status === "disconnected")
              .map((connection) => (
                <div key={connection.name} className="px-6 py-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center">
                      <div className="text-3xl mr-4">{connection.icon}</div>
                      <div>
                        <h4 className="text-lg font-medium text-gray-900">
                          {connection.name}
                        </h4>
                        <p className="text-sm text-gray-500">
                          {connection.description}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center space-x-3">
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                        ○ Not Connected
                      </span>
                      <button
                        onClick={() => handleConnect(connection.name)}
                        disabled={loading}
                        className="bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700 transition-colors disabled:opacity-50"
                      >
                        Connect
                      </button>
                    </div>
                  </div>
                </div>
              ))
          )}
        </div>
      </div>
    </DashboardLayout>
  );
}
