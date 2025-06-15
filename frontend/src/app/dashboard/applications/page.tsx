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

  useEffect(() => {
    loadApps();
  }, []);

  const loadApps = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiClient.getApps();
      setApps(response.apps);
    } catch (err) {
      console.error("Failed to load apps:", err);
      setError(ERROR_MESSAGES.NETWORK_ERROR);
      // Use fallback data
      setApps([...FALLBACK_SAAS_APPS]);
    } finally {
      setLoading(false);
    }
  };

  const handleConnect = async (appId: string) => {
    try {
      setConnectingApps((prev) => new Set(prev).add(appId));

      if (DEMO_CONFIG.SIMULATE_OAUTH) {
        // Simulate OAuth flow for demo
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
        // Real OAuth flow
        const response: AppConnectionResponse = await apiClient.connectApp(
          appId
        );

        // Redirect to OAuth provider
        window.location.href = response.auth_url;
      }
    } catch (err) {
      console.error("Failed to connect app:", err);
      setError(ERROR_MESSAGES.CONNECTION_FAILED);
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

      // Check if app is connected
      const app = apps.find((a) => a.id === appId);
      if (!app || app.status !== "connected") {
        setError(ERROR_MESSAGES.CONNECTION_FAILED);
        return;
      }

      if (DEMO_CONFIG.SIMULATE_OAUTH) {
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
        // Real launch flow
        const response: AppLaunchResponse = await apiClient.launchApp(appId);
        window.open(response.launch_url, "_blank", "noopener,noreferrer");
      }
    } catch (err) {
      console.error("Failed to launch app:", err);
      setError(ERROR_MESSAGES.LAUNCH_FAILED);
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
  const filteredApps = apps.filter((app) => {
    const matchesSearch =
      app.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      app.description.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesCategory =
      selectedCategory === "all" || app.category === selectedCategory;
    return matchesSearch && matchesCategory;
  });

  // Get unique categories
  const categories = [
    "all",
    ...Array.from(new Set(apps.map((app) => app.category))),
  ];

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
      title="Applications"
      description="Manage your SaaS applications and connections"
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
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          </div>
          <div className="sm:w-48">
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
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
            {filteredApps.map((app) => (
              <div
                key={app.id}
                className="border border-gray-200 rounded-lg p-6 hover:shadow-md transition-shadow"
              >
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center">
                    <span className="text-3xl mr-3">{app.icon}</span>
                    <div>
                      <h4 className="font-medium text-gray-900">{app.name}</h4>
                      <p className="text-xs text-gray-500 capitalize">
                        {app.category}
                      </p>
                    </div>
                  </div>
                  {getStatusBadge(app.status)}
                </div>

                <p className="text-sm text-gray-600 mb-6 line-clamp-3">
                  {app.description}
                </p>

                {getActionButton(app)}
              </div>
            ))}
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}
