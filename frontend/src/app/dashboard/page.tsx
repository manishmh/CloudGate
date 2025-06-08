"use client";

import {
  APP_LAUNCH_URLS,
  DEMO_CONFIG,
  ERROR_MESSAGES,
  FALLBACK_SAAS_APPS,
  LOADING_MESSAGES,
  SECURITY_FEATURES,
  STATUS_CONFIG,
  SUCCESS_MESSAGES,
} from "@/constants";
import {
  apiClient,
  type AppConnectionResponse,
  type AppLaunchResponse,
  type SaaSApplication,
} from "@/lib/api";
import { useKeycloak } from "@react-keycloak/web";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { IoRefresh } from "react-icons/io5";
import Link from 'next/link';

interface AppConnection {
  name: string;
  status: 'connected' | 'disconnected';
  icon: string;
  description: string;
  connectUrl: string;
}

export default function Dashboard() {
  const { keycloak, initialized } = useKeycloak();
  const router = useRouter();
  const [apps, setApps] = useState<SaaSApplication[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [connectingApps, setConnectingApps] = useState<Set<string>>(new Set());
  const [launchingApps, setLaunchingApps] = useState<Set<string>>(new Set());
  const [profilePicture, setProfilePicture] = useState<string | null>(null);
  const [connections, setConnections] = useState<AppConnection[]>([
    {
      name: 'Google Workspace',
      status: 'disconnected',
      icon: '🔍',
      description: 'Access Gmail, Drive, Calendar, and more',
      connectUrl: '/oauth/google/connect'
    },
    {
      name: 'Microsoft 365',
      status: 'disconnected',
      icon: '🏢',
      description: 'Access Outlook, OneDrive, Teams, and more',
      connectUrl: '/oauth/microsoft/connect'
    },
    {
      name: 'Slack',
      status: 'disconnected',
      icon: '💬',
      description: 'Access your Slack workspaces',
      connectUrl: '/oauth/slack/connect'
    },
    {
      name: 'GitHub',
      status: 'disconnected',
      icon: '🐙',
      description: 'Access your repositories and organizations',
      connectUrl: '/oauth/github/connect'
    }
  ]);
  const [user, setUser] = useState<any>(null);

  useEffect(() => {
    loadApps();
    // Load profile picture
    if (keycloak?.tokenParsed?.sub) {
      const savedPicture = localStorage.getItem(
        `profile_picture_${keycloak.tokenParsed.sub}`
      );
      if (savedPicture) {
        setProfilePicture(savedPicture);
      }
    }

    // Check for connection success from URL params
    const urlParams = new URLSearchParams(window.location.search);
    const connected = urlParams.get('connected');
    const email = urlParams.get('email');
    
    if (connected && email) {
      // Update the connection status
      setConnections(prev => prev.map(conn => 
        conn.name.toLowerCase().includes(connected) 
          ? { ...conn, status: 'connected' as const }
          : conn
      ));
      
      // Show success message
      alert(`Successfully connected to ${connected}! Email: ${email}`);
      
      // Clean up URL
      window.history.replaceState({}, document.title, '/dashboard');
    }

    // Load user info (mock for now)
    setUser({
      name: 'Test User',
      email: 'test@cloudgate.com',
      keycloakId: 'test-user-id'
    });
    setLoading(false);
  }, [keycloak]);

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

        if (response.method === "redirect") {
          window.location.href = response.launch_url;
        } else if (response.method === "popup") {
          window.open(response.launch_url, "_blank", "width=800,height=600");
        }
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
      STATUS_CONFIG.available;
    return (
      <span
        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${config.color}`}
      >
        {config.text}
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
          className="w-full bg-green-600 hover:bg-green-700 disabled:bg-green-400 text-white py-2 px-4 rounded-lg text-sm font-medium transition-colors cursor-pointer"
        >
          {isLaunching ? (
            <span className="flex items-center justify-center">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
              {LOADING_MESSAGES.LAUNCHING}
            </span>
          ) : (
            "Launch"
          )}
        </button>
      );
    } else if (app.status === "pending") {
      return (
        <button
          disabled
          className="w-full bg-yellow-400 text-yellow-800 py-2 px-4 rounded-lg text-sm font-medium cursor-not-allowed"
        >
          Pending
        </button>
      );
    } else {
      return (
        <button
          onClick={() => handleConnect(app.id)}
          disabled={isConnecting}
          className="w-full bg-blue-600 cursor-pointer hover:bg-blue-700 disabled:bg-blue-400 text-white py-2 px-4 rounded-lg text-sm font-medium transition-colors"
        >
          {isConnecting ? (
            <span className="flex items-center justify-center">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
              {LOADING_MESSAGES.CONNECTING}
            </span>
          ) : (
            "Connect"
          )}
        </button>
      );
    }
  };

  const handleDisconnect = async (app: AppConnection) => {
    // TODO: Implement disconnect functionality
    setConnections(prev => prev.map(conn => 
      conn.name === app.name 
        ? { ...conn, status: 'disconnected' as const }
        : conn
    ));
  };

  if (!initialized) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!keycloak?.authenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">
            Access Denied
          </h1>
          <p className="text-gray-600 mb-6">
            Please log in to access the dashboard.
          </p>
          <button
            onClick={() => keycloak?.login()}
            className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-md font-medium cursor-pointer"
          >
            Login
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                CloudGate Dashboard
              </h1>
              <p className="text-gray-600">
                Welcome, {keycloak?.tokenParsed?.preferred_username || "User"}
              </p>
            </div>
            <div className="flex items-center space-x-4">
              <button
                onClick={handleRefresh}
                className="p-2 text-gray-600 cursor-pointer hover:text-gray-900 hover:bg-gray-100 rounded-full transition-colors"
                title="Refresh"
              >
                <IoRefresh />
              </button>
              <button
                onClick={() => router.push("/profile")}
                className="flex items-center cursor-pointer justify-center w-10 h-10 rounded-full bg-blue-600 hover:bg-blue-700 transition-colors"
                title="Profile"
              >
                {profilePicture ? (
                  <Image
                    src={profilePicture}
                    alt="Profile"
                    width={40}
                    height={40}
                    className="w-10 h-10 rounded-full object-cover"
                    unoptimized={true}
                  />
                ) : (
                  <span className="text-white font-semibold text-sm">
                    {keycloak?.tokenParsed?.given_name?.charAt(0) || "U"}
                    {keycloak?.tokenParsed?.family_name?.charAt(0) || ""}
                  </span>
                )}
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {/* Security Status */}
        <div className="mb-8">
          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg leading-6 font-medium text-gray-900 mb-4">
                Security Status
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {SECURITY_FEATURES.map((feature) => (
                  <div key={feature.id} className="flex items-center space-x-3">
                    <div
                      className={`flex-shrink-0 w-8 h-8 bg-${feature.color}-100 rounded-full flex items-center justify-center`}
                    >
                      <span className={`text-${feature.color}-600 text-sm`}>
                        ✓
                      </span>
                    </div>
                    <div>
                      <p className="text-sm font-medium text-gray-900">
                        {feature.title}
                      </p>
                      <p className="text-sm text-gray-500">
                        {feature.description}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

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

        {/* SaaS Applications */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-lg font-semibold text-gray-900">
              SaaS Applications
            </h2>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-500">
                {apps.length} applications available
              </span>
              <button
                onClick={handleRefresh}
                disabled={loading}
                className="text-blue-600 cursor-pointer hover:text-blue-700 text-sm font-medium disabled:opacity-50"
              >
                {loading ? "Refreshing..." : "Refresh"}
              </button>
            </div>
          </div>

          {loading ? (
            <div className="flex items-center justify-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              <span className="ml-2 text-gray-600">
                {LOADING_MESSAGES.LOADING_APPS}
              </span>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {apps.map((app) => (
                <div
                  key={app.id}
                  className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
                >
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center">
                      <span className="text-2xl mr-3">{app.icon}</span>
                      <div>
                        <h3 className="font-medium text-gray-900">
                          {app.name}
                        </h3>
                        <p className="text-xs text-gray-500 capitalize">
                          {app.category}
                        </p>
                      </div>
                    </div>
                    {getStatusBadge(app.status)}
                  </div>

                  <p className="text-sm text-gray-600 mb-4">
                    {app.description}
                  </p>

                  {getActionButton(app)}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* App Connections */}
        <div className="mt-8 bg-white shadow overflow-hidden sm:rounded-md">
          <div className="px-4 py-5 sm:px-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900">
              SaaS Application Connections
            </h3>
            <p className="mt-1 max-w-2xl text-sm text-gray-500">
              Connect to your favorite SaaS applications for seamless single sign-on.
            </p>
          </div>
          <ul className="divide-y divide-gray-200">
            {connections.map((app) => (
              <li key={app.name}>
                <div className="px-4 py-4 flex items-center justify-between">
                  <div className="flex items-center">
                    <div className="text-2xl mr-4">{app.icon}</div>
                    <div>
                      <h4 className="text-lg font-medium text-gray-900">{app.name}</h4>
                      <p className="text-sm text-gray-500">{app.description}</p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-3">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      app.status === 'connected' 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-gray-100 text-gray-800'
                    }`}>
                      {app.status === 'connected' ? '✓ Connected' : '○ Not Connected'}
                    </span>
                    {app.status === 'connected' ? (
                      <button
                        onClick={() => handleDisconnect(app)}
                        className="bg-red-600 text-white px-4 py-2 rounded-md text-sm hover:bg-red-700 transition-colors"
                      >
                        Disconnect
                      </button>
                    ) : (
                      <button
                        onClick={() => handleConnect(app.name.toLowerCase())}
                        className="bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700 transition-colors"
                      >
                        Connect
                      </button>
                    )}
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>

        {/* Quick Actions */}
        <div className="mt-8 bg-white shadow overflow-hidden sm:rounded-md">
          <div className="px-4 py-5 sm:px-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900">
              Quick Actions
            </h3>
          </div>
          <div className="px-4 py-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <Link
              href="/profile"
              className="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <div className="text-2xl mr-3">👤</div>
              <div>
                <h4 className="font-medium text-gray-900">Profile Settings</h4>
                <p className="text-sm text-gray-500">Manage your account</p>
              </div>
            </Link>
            
            <Link
              href="/privacy-policy"
              className="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <div className="text-2xl mr-3">🔒</div>
              <div>
                <h4 className="font-medium text-gray-900">Privacy Policy</h4>
                <p className="text-sm text-gray-500">View our privacy policy</p>
              </div>
            </Link>
            
            <Link
              href="/terms"
              className="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <div className="text-2xl mr-3">📋</div>
              <div>
                <h4 className="font-medium text-gray-900">Terms of Service</h4>
                <p className="text-sm text-gray-500">View our terms</p>
              </div>
            </Link>
            
            <button
              onClick={() => window.location.href = 'http://localhost:8080/realms/cloudgate/account/'}
              className="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors text-left"
            >
              <div className="text-2xl mr-3">🔑</div>
              <div>
                <h4 className="font-medium text-gray-900">Keycloak Account</h4>
                <p className="text-sm text-gray-500">Manage authentication</p>
              </div>
            </button>
          </div>
        </div>
      </main>
    </div>
  );
}
