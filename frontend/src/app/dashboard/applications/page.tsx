"use client";

import DashboardLayout from "@/components/DashboardLayout";
import {
  APP_LAUNCH_URLS,
  DEMO_CONFIG,
  ERROR_MESSAGES,
  FALLBACK_SAAS_APPS,
  LOADING_MESSAGES,
  STATUS_CONFIG,
  SUCCESS_MESSAGES,
} from "@/constants";
import {
  apiClient,
  type AppConnectionResponse,
  type AppLaunchResponse,
  type SaaSApplication,
} from "@/lib/api";

import { useEffect, useState } from "react";
import { IoRefresh, IoSearch } from "react-icons/io5";

export default function ApplicationsPage() {
  const [apps, setApps] = useState<SaaSApplication[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [connectingApps, setConnectingApps] = useState<Set<string>>(new Set());
  const [launchingApps, setLaunchingApps] = useState<Set<string>>(new Set());
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("all");
  const [backendAvailable, setBackendAvailable] = useState(false);

  useEffect(() => {
    checkBackendAndLoadApps();

    // Check for OAuth callback success
    const urlParams = new URLSearchParams(window.location.search);
    const connected = urlParams.get("connected");
    const email = urlParams.get("email");

    if (connected && email) {
      // OAuth callback success - update the app status
      handleOAuthCallback(connected, email);

      // Clean up URL
      window.history.replaceState({}, document.title, window.location.pathname);
    }
  }, []);

  const checkBackendAndLoadApps = async () => {
    try {
      setLoading(true);
      setError(null);

      // Check if backend is available
      const isAvailable = await apiClient.isBackendAvailable();
      setBackendAvailable(isAvailable);

      if (isAvailable) {
        // Load apps from backend
        const response = await apiClient.getApps();
        setApps(response.apps);
      } else {
        // Use fallback data if backend is not available
        console.warn("Backend not available, using fallback data");
        setApps([...FALLBACK_SAAS_APPS]);
      }
    } catch (err) {
      console.error("Failed to load apps:", err);
      setError(ERROR_MESSAGES.NETWORK_ERROR);
      // Use fallback data on error
      setApps([...FALLBACK_SAAS_APPS]);
      setBackendAvailable(false);
    } finally {
      setLoading(false);
    }
  };

  const handleOAuthCallback = (provider: string, email: string) => {
    // Update the app status to connected
    setApps((prevApps) =>
      prevApps.map((app) =>
        app.id === provider ||
        app.id === `${provider}-workspace` ||
        app.id === `${provider}-365`
          ? {
              ...app,
              status: "connected" as const,
              connection_details: {
                user_email: email,
                connected_at: new Date().toISOString(),
              },
            }
          : app
      )
    );

    // Show success message
    setError(null);
    // You could show a success toast here instead
    console.log(`Successfully connected to ${provider} with email: ${email}`);
  };

  const loadApps = async () => {
    await checkBackendAndLoadApps();
  };

  const handleConnect = async (appId: string) => {
    try {
      setConnectingApps((prev) => new Set(prev).add(appId));
      setError(null);

      if (DEMO_CONFIG.SIMULATE_OAUTH || !backendAvailable) {
        // Simulate OAuth flow for demo or when backend is unavailable
        await new Promise((resolve) =>
          setTimeout(resolve, DEMO_CONFIG.AUTO_CONNECT_DELAY)
        );

        // Update app status to connected
        setApps((prevApps) =>
          prevApps.map((app) =>
            app.id === appId ? { ...app, status: "connected" as const } : app
          )
        );

        if (DEMO_CONFIG.SHOW_DEMO_ALERTS) {
          alert(SUCCESS_MESSAGES.APP_CONNECTED);
        }
      } else {
        // Real OAuth flow with backend
        const response: AppConnectionResponse = await apiClient.connectApp(
          appId
        );

        // Store the app ID for callback handling
        sessionStorage.setItem("oauth_app_id", appId);

        // Open OAuth in a popup window
        const popup = window.open(
          response.auth_url,
          `oauth_${appId}`,
          "width=600,height=700,scrollbars=yes,resizable=yes"
        );

        if (!popup) {
          // Fallback to redirect if popup is blocked
          window.location.href = response.auth_url;
          return;
        }

        // Monitor popup for completion
        const checkClosed = setInterval(() => {
          if (popup.closed) {
            clearInterval(checkClosed);
            // Refresh the page to check for updates
            setTimeout(() => {
              loadApps();
            }, 1000);
          }
        }, 1000);

        // Handle popup messages (if callback sends postMessage)
        const handleMessage = (event: MessageEvent) => {
          if (event.origin !== window.location.origin) return;

          if (event.data.type === "oauth_success") {
            popup.close();
            clearInterval(checkClosed);
            handleOAuthCallback(event.data.provider, event.data.email);
            window.removeEventListener("message", handleMessage);
          } else if (event.data.type === "oauth_error") {
            popup.close();
            clearInterval(checkClosed);
            setError(`OAuth failed: ${event.data.error}`);
            window.removeEventListener("message", handleMessage);
          }
        };

        window.addEventListener("message", handleMessage);

        // Cleanup after 5 minutes
        setTimeout(() => {
          if (!popup.closed) {
            popup.close();
          }
          clearInterval(checkClosed);
          window.removeEventListener("message", handleMessage);
        }, 300000);
      }
    } catch (err) {
      console.error("Failed to connect app:", err);

      // Handle OAuth not configured error more gracefully
      const errorMessage =
        err instanceof Error ? err.message : "Please try again.";
      if (errorMessage.includes("OAuth not configured")) {
        const appName = apps.find((app) => app.id === appId)?.name || appId;
        setError(
          `${appName} OAuth is not configured yet. Please contact your administrator to set up OAuth credentials for this application.`
        );
      } else if (
        errorMessage.includes("Failed to exchange authorization code")
      ) {
        setError(
          `OAuth authentication failed for ${appId}. This might be due to configuration issues or expired tokens. Please try again or contact support.`
        );
      } else {
        setError(`Failed to connect to ${appId}. ${errorMessage}`);
      }
    } finally {
      setConnectingApps((prev) => {
        const newSet = new Set(prev);
        newSet.delete(appId);
        return newSet;
      });
    }
  };

  const handleLaunch = async (appId: string) => {
    try {
      setLaunchingApps((prev) => new Set(prev).add(appId));
      setError(null);

      // Check if app is connected
      const app = apps.find((a) => a.id === appId);
      if (!app || app.status !== "connected") {
        setError("Application must be connected before launching.");
        return;
      }

      if (DEMO_CONFIG.SIMULATE_OAUTH || !backendAvailable) {
        // Simulate launch delay
        await new Promise((resolve) =>
          setTimeout(resolve, DEMO_CONFIG.LAUNCH_DELAY)
        );

        // Get launch URL from constants
        const launchURL =
          APP_LAUNCH_URLS[appId as keyof typeof APP_LAUNCH_URLS];

        if (launchURL) {
          // Open in new tab
          window.open(launchURL, "_blank", "noopener,noreferrer");

          if (DEMO_CONFIG.SHOW_DEMO_ALERTS) {
            alert(SUCCESS_MESSAGES.APP_LAUNCHED);
          }
        } else {
          throw new Error(ERROR_MESSAGES.APP_NOT_FOUND);
        }
      } else {
        // Real launch flow with backend
        try {
          const response: AppLaunchResponse = await apiClient.launchApp(appId);
          window.open(response.launch_url, "_blank", "noopener,noreferrer");
        } catch (launchError) {
          // Fallback to direct URL if backend launch fails
          const launchURL =
            APP_LAUNCH_URLS[appId as keyof typeof APP_LAUNCH_URLS];
          if (launchURL) {
            window.open(launchURL, "_blank", "noopener,noreferrer");
          } else {
            throw launchError;
          }
        }
      }

      // Update last used timestamp
      setApps((prevApps) =>
        prevApps.map((a) =>
          a.id === appId
            ? {
                ...a,
                connection_details: {
                  ...a.connection_details,
                  last_used: new Date().toISOString(),
                },
              }
            : a
        )
      );
    } catch (err) {
      console.error("Failed to launch app:", err);
      setError(
        `Failed to launch ${appId}. ${
          err instanceof Error ? err.message : "Please try again."
        }`
      );
    } finally {
      setLaunchingApps((prev) => {
        const newSet = new Set(prev);
        newSet.delete(appId);
        return newSet;
      });
    }
  };

  const handleRefresh = () => {
    loadApps();
  };

  const getStatusBadge = (status: string) => {
    const config =
      STATUS_CONFIG[status as keyof typeof STATUS_CONFIG] ||
      STATUS_CONFIG.default;
    return (
      <span
        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${config.className}`}
      >
        {config.label}
      </span>
    );
  };

  const getActionButton = (app: SaaSApplication) => {
    const isConnecting = connectingApps.has(app.id);
    const isLaunching = launchingApps.has(app.id);

    // List of providers that are not configured
    const unconfiguredProviders = ["salesforce"]; // Only Salesforce is not configured
    const isUnconfigured = unconfiguredProviders.includes(app.id);

    if (app.status === "connected") {
      return (
        <button
          onClick={() => handleLaunch(app.id)}
          disabled={isLaunching}
          className="w-full bg-green-600 hover:bg-green-700 disabled:bg-green-400 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors cursor-pointer"
        >
          {isLaunching ? "Launching..." : "Launch"}
        </button>
      );
    }

    if (isUnconfigured) {
      return (
        <button
          disabled
          className="w-full bg-gray-400 text-white px-4 py-2 rounded-md text-sm font-medium cursor-not-allowed"
          title="OAuth not configured for this provider"
        >
          Not Available
        </button>
      );
    }

    return (
      <button
        onClick={() => handleConnect(app.id)}
        disabled={isConnecting}
        className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors cursor-pointer"
      >
        {isConnecting ? "Connecting..." : "Connect"}
      </button>
    );
  };

  // Filter apps based on search term and category
  const filteredApps = apps
    .filter((app) => {
      const matchesSearch =
        app.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        app.description.toLowerCase().includes(searchTerm.toLowerCase());
      const matchesCategory =
        selectedCategory === "all" || app.category === selectedCategory;
      return matchesSearch && matchesCategory;
    })
    .sort((a, b) => {
      // Sort by name to maintain consistent order
      return a.name.localeCompare(b.name);
    });

  // Get unique categories
  const categories = [
    "all",
    ...Array.from(new Set(apps.map((app) => app.category))),
  ];

  return (
    <DashboardLayout>
      {/* Backend Status Indicator */}
      {!backendAvailable && (
        <div className="mb-6 bg-yellow-50 border border-yellow-200 rounded-md p-4">
          <div className="flex">
            <div className="ml-3">
              <h3 className="text-sm font-medium text-yellow-800">Demo Mode</h3>
              <p className="text-sm text-yellow-700 mt-1">
                Backend is not available. Using demo data. OAuth connections
                will not persist.
              </p>
            </div>
          </div>
        </div>
      )}

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
              className="ml-auto text-red-400 hover:text-red-600"
            >
              ×
            </button>
          </div>
        </div>
      )}

      {/* Search and Filter */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1">
            <div className="relative">
              <IoSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
              <input
                type="text"
                placeholder="Search applications..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-black"
              />
            </div>
          </div>
          <div className="sm:w-48">
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="w-full px-3 py-2 border text-black/80 text-sm h-full border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              {categories.map((category) => (
                <option key={category} value={category}>
                  {category === "all"
                    ? "All Categories"
                    : category.charAt(0).toUpperCase() + category.slice(1)}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Applications Grid */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex justify-between items-center mb-6">
          <h3 className="text-lg font-medium text-gray-900">
            {filteredApps.length} Application
            {filteredApps.length !== 1 ? "s" : ""}
            {searchTerm && ` matching "${searchTerm}"`}
            {selectedCategory !== "all" && ` in ${selectedCategory}`}
          </h3>
          {backendAvailable && (
            <div className="flex items-center text-sm text-green-600">
              <div className="w-2 h-2 bg-green-500 rounded-full mr-2"></div>
              Connection established
            </div>
          )}
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <span className="ml-2 text-gray-600">
              {LOADING_MESSAGES.LOADING_APPS}
            </span>
          </div>
        ) : filteredApps.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-gray-400 text-6xl mb-4">📱</div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              No applications found
            </h3>
            <p className="text-gray-500">
              {searchTerm || selectedCategory !== "all"
                ? "Try adjusting your search or filter criteria."
                : "No applications are available at the moment."}
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {filteredApps.map((app) => {
              const unconfiguredProviders = ["salesforce"]; // Only Salesforce is not configured
              const isUnconfigured = unconfiguredProviders.includes(app.id);

              return (
                <div
                  key={app.id}
                  className={`border rounded-lg p-6 hover:shadow-md transition-shadow flex flex-col h-full ${
                    isUnconfigured
                      ? "border-gray-300 bg-gray-50"
                      : "border-gray-200"
                  }`}
                >
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center">
                      <span
                        className={`text-3xl mr-3 ${
                          isUnconfigured ? "opacity-50" : ""
                        }`}
                      >
                        {app.icon}
                      </span>
                      <div>
                        <h4
                          className={`font-medium ${
                            isUnconfigured ? "text-gray-600" : "text-gray-900"
                          }`}
                        >
                          {app.name}
                          {isUnconfigured && (
                            <span className="ml-2 text-xs text-orange-600 font-normal">
                              (Config Required)
                            </span>
                          )}
                        </h4>
                        <p className="text-xs text-gray-500 capitalize">
                          {app.category}
                        </p>
                      </div>
                    </div>
                    {getStatusBadge(app.status)}
                  </div>

                  <p className="text-sm text-gray-600 mb-4 line-clamp-3 flex-grow">
                    {app.description}
                  </p>

                  {/* Connection Details */}
                  {app.status === "connected" && app.connection_details && (
                    <div className="mb-4 p-3 bg-green-50 rounded-md">
                      <p className="text-xs text-green-700">
                        Connected as: {app.connection_details.user_email}
                      </p>
                      {app.connection_details.last_used && (
                        <p className="text-xs text-green-600 mt-1">
                          Last used:{" "}
                          {new Date(
                            app.connection_details.last_used
                          ).toLocaleDateString()}
                        </p>
                      )}
                    </div>
                  )}

                  <div className="mt-auto">{getActionButton(app)}</div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}
