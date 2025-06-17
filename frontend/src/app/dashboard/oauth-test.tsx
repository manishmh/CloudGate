"use client";

import { useEffect, useState } from "react";

interface AppConnection {
  name: string;
  status: "connected" | "disconnected";
  icon: string;
  description: string;
  connectUrl: string;
}

export default function OAuthTest() {
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
      status: "disconnected",
      icon: "💬",
      description: "Access your Slack workspaces",
      connectUrl: "/oauth/slack/connect",
    },
    {
      name: "GitHub",
      status: "disconnected",
      icon: "🐙",
      description: "Access your repositories and organizations",
      connectUrl: "/oauth/github/connect",
    },
    {
      name: "Trello",
      status: "disconnected",
      icon: "📋",
      description: "Access your boards and cards (OAuth 1.0a)",
      connectUrl: "/oauth/trello/connect",
    },
  ]);

  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  useEffect(() => {
    // Check for connection success from URL params
    const urlParams = new URLSearchParams(window.location.search);
    const connected = urlParams.get("connected");
    const email = urlParams.get("email");

    if (connected && email) {
      // Update the connection status
      setConnections((prev) =>
        prev.map((conn) =>
          conn.name.toLowerCase().includes(connected)
            ? { ...conn, status: "connected" as const }
            : conn
        )
      );

      setMessage(`✅ Successfully connected to ${connected}! Email: ${email}`);

      // Clean up URL
      window.history.replaceState({}, document.title, "/dashboard/oauth-test");
    }
  }, []);

  const handleConnect = async (app: AppConnection) => {
    setLoading(true);
    setMessage(`🔄 Connecting to ${app.name}...`);

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081"}${
          app.connectUrl
        }`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        if (data.auth_url) {
          setMessage(`🚀 Redirecting to ${app.name} OAuth...`);
          // Redirect to OAuth provider
          window.location.href = data.auth_url;
        } else {
          throw new Error("No auth URL received");
        }
      } else {
        const error = await response.json();
        throw new Error(error.error || "Failed to initiate OAuth");
      }
    } catch (error) {
      console.error("OAuth connection error:", error);
      setMessage(`❌ Failed to connect to ${app.name}: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  const testBackendConnection = async () => {
    setLoading(true);
    setMessage("🔄 Testing backend connection...");

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081"}/health`
      );
      if (response.ok) {
        const data = await response.json();
        setMessage(`✅ Backend is healthy: ${data.status}`);
      } else {
        throw new Error("Backend not responding");
      }
    } catch (error) {
      setMessage(`❌ Backend connection failed: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white shadow rounded-lg p-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">
            OAuth Integration Test
          </h1>

          {/* Status Message */}
          {message && (
            <div className="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
              <p className="text-blue-800">{message}</p>
            </div>
          )}

          {/* Backend Test */}
          <div className="mb-8">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              Backend Connection Test
            </h2>
            <button
              onClick={testBackendConnection}
              disabled={loading}
              className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:bg-blue-400 transition-colors"
            >
              {loading ? "Testing..." : "Test Backend"}
            </button>
          </div>

          {/* OAuth Connections */}
          <div>
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              OAuth Provider Connections
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {connections.map((app) => (
                <div
                  key={app.name}
                  className="border border-gray-200 rounded-lg p-4"
                >
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center">
                      <span className="text-2xl mr-3">{app.icon}</span>
                      <div>
                        <h3 className="font-medium text-gray-900">
                          {app.name}
                        </h3>
                        <p className="text-sm text-gray-500">
                          {app.description}
                        </p>
                      </div>
                    </div>
                    <span
                      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        app.status === "connected"
                          ? "bg-green-100 text-green-800"
                          : "bg-gray-100 text-gray-800"
                      }`}
                    >
                      {app.status === "connected"
                        ? "✓ Connected"
                        : "○ Not Connected"}
                    </span>
                  </div>

                  <div className="flex space-x-2">
                    {app.status === "connected" ? (
                      <button
                        onClick={() => {
                          setConnections((prev) =>
                            prev.map((conn) =>
                              conn.name === app.name
                                ? { ...conn, status: "disconnected" as const }
                                : conn
                            )
                          );
                          setMessage(`🔌 Disconnected from ${app.name}`);
                        }}
                        className="flex-1 bg-red-600 text-white px-3 py-2 rounded-md text-sm hover:bg-red-700 transition-colors"
                      >
                        Disconnect
                      </button>
                    ) : (
                      <button
                        onClick={() => handleConnect(app)}
                        disabled={loading}
                        className="flex-1 bg-blue-600 text-white px-3 py-2 rounded-md text-sm hover:bg-blue-700 disabled:bg-blue-400 transition-colors"
                      >
                        {loading ? "Connecting..." : "Connect"}
                      </button>
                    )}

                    <button
                      onClick={async () => {
                        setMessage(`🔍 Testing ${app.name} endpoint...`);
                        try {
                          const response = await fetch(
                            `${
                              process.env.NEXT_PUBLIC_API_URL ||
                              "http://localhost:8081"
                            }${app.connectUrl}`
                          );
                          const data = await response.json();
                          if (response.ok) {
                            setMessage(
                              `✅ ${
                                app.name
                              } endpoint working! Auth URL: ${data.auth_url?.substring(
                                0,
                                50
                              )}...`
                            );
                          } else {
                            setMessage(
                              `⚠️ ${app.name} endpoint error: ${data.error}`
                            );
                          }
                        } catch (error) {
                          setMessage(
                            `❌ ${app.name} endpoint failed: ${error}`
                          );
                        }
                      }}
                      className="bg-gray-600 text-white px-3 py-2 rounded-md text-sm hover:bg-gray-700 transition-colors"
                    >
                      Test
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Debug Info */}
          <div className="mt-8 p-4 bg-gray-50 rounded-lg">
            <h3 className="font-medium text-gray-900 mb-2">
              Debug Information
            </h3>
            <div className="text-sm text-gray-600 space-y-1">
              <p>
                <strong>Backend URL:</strong>{" "}
                {process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081"}
              </p>
              <p>
                <strong>Frontend URL:</strong>{" "}
                {process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000"}
              </p>
              <p>
                <strong>Keycloak URL:</strong>{" "}
                {process.env.NEXT_PUBLIC_KEYCLOAK_URL ||
                  "http://localhost:8080"}
              </p>
              <p>
                <strong>Current URL:</strong>{" "}
                {typeof window !== "undefined" ? window.location.href : "N/A"}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
