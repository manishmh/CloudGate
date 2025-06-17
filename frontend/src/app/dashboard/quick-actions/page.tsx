"use client";

import DashboardLayout from "@/components/DashboardLayout";
import Link from "next/link";
import { useState } from "react";
import { IoHelpCircle, IoRefresh, IoShare } from "react-icons/io5";

interface QuickAction {
  id: string;
  title: string;
  description: string;
  icon: string;
  action: () => void;
  category: "account" | "security" | "apps" | "support";
  external?: boolean;
  href?: string;
}

export default function QuickActionsPage() {
  const [loading, setLoading] = useState(false);

  const handleExportData = () => {
    // Simulate data export
    const data = {
      user: "demo@cloudgate.com",
      exportDate: new Date().toISOString(),
      applications: ["Google Workspace", "Slack", "GitHub"],
      connections: 3,
    };

    const blob = new Blob([JSON.stringify(data, null, 2)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "cloudgate-data-export.json";
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);

    alert("Data exported successfully!");
  };

  const handleGenerateReport = () => {
    alert("Security report generated! Check your email for the report.");
  };

  const handleClearCache = () => {
    localStorage.clear();
    sessionStorage.clear();
    alert("Cache cleared successfully!");
  };

  const handleTestConnections = async () => {
    setLoading(true);
    // Simulate connection testing
    await new Promise((resolve) => setTimeout(resolve, 2000));
    setLoading(false);
    alert("All connections tested successfully!");
  };

  const quickActions: QuickAction[] = [
    {
      id: "profile-settings",
      title: "Profile Settings",
      description: "Update your personal information and preferences",
      icon: "👤",
      action: () => {},
      category: "account",
      href: "/profile",
    },
    {
      id: "change-password",
      title: "Change Password",
      description: "Update your account password",
      icon: "🔑",
      action: () => {
        window.open(
          `${
            process.env.NEXT_PUBLIC_KEYCLOAK_URL || "http://localhost:8080"
          }/realms/cloudgate/account/`,
          "_blank"
        );
      },
      category: "security",
      external: true,
    },
    {
      id: "export-data",
      title: "Export Data",
      description: "Download your account data and settings",
      icon: "📥",
      action: handleExportData,
      category: "account",
    },
    {
      id: "security-report",
      title: "Generate Security Report",
      description: "Get a detailed security analysis of your account",
      icon: "📊",
      action: handleGenerateReport,
      category: "security",
    },
    {
      id: "test-connections",
      title: "Test App Connections",
      description: "Verify all your application connections are working",
      icon: "🔗",
      action: handleTestConnections,
      category: "apps",
    },
    {
      id: "clear-cache",
      title: "Clear Cache",
      description: "Clear browser cache and stored data",
      icon: "🧹",
      action: handleClearCache,
      category: "account",
    },
    {
      id: "privacy-policy",
      title: "Privacy Policy",
      description: "Review our privacy policy and data handling",
      icon: "🔒",
      action: () => {},
      category: "support",
      href: "/privacy-policy",
    },
    {
      id: "terms-service",
      title: "Terms of Service",
      description: "Read our terms and conditions",
      icon: "📋",
      action: () => {},
      category: "support",
      href: "/terms",
    },
    {
      id: "contact-support",
      title: "Contact Support",
      description: "Get help from our support team",
      icon: "💬",
      action: () => {
        window.open(
          "mailto:support@cloudgate.com?subject=CloudGate Support Request",
          "_blank"
        );
      },
      category: "support",
      external: true,
    },
    {
      id: "documentation",
      title: "Documentation",
      description: "Browse our help documentation",
      icon: "📚",
      action: () => {
        alert("Documentation coming soon!");
      },
      category: "support",
    },
    {
      id: "api-keys",
      title: "API Keys",
      description: "Manage your API keys and tokens",
      icon: "🔐",
      action: () => {
        alert("API key management coming soon!");
      },
      category: "account",
    },
    {
      id: "activity-log",
      title: "Activity Log",
      description: "View your recent account activity",
      icon: "📝",
      action: () => {},
      category: "security",
      href: "/dashboard/security",
    },
  ];

  const categories = [
    { id: "account", name: "Account", icon: "👤", color: "blue" },
    { id: "security", name: "Security", icon: "🔒", color: "green" },
    { id: "apps", name: "Applications", icon: "📱", color: "purple" },
    { id: "support", name: "Support", color: "yellow", icon: "💬" },
  ];

  const handleRefresh = () => {
    setLoading(true);
    setTimeout(() => setLoading(false), 1000);
  };

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
      title="Quick Actions"
      description="Common tasks and shortcuts for your account"
      actions={refreshAction}
    >
      {/* Categories Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {categories.map((category) => {
          const categoryActions = quickActions.filter(
            (action) => action.category === category.id
          );
          return (
            <div key={category.id} className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className={`p-2 bg-${category.color}-100 rounded-lg`}>
                  <span className="text-2xl">{category.icon}</span>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">
                    {category.name}
                  </p>
                  <p className="text-2xl font-semibold text-gray-900">
                    {categoryActions.length}
                  </p>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Quick Actions by Category */}
      {categories.map((category) => {
        const categoryActions = quickActions.filter(
          (action) => action.category === category.id
        );

        return (
          <div
            key={category.id}
            className="bg-white rounded-lg shadow p-6 mb-8"
          >
            <div className="flex items-center mb-6">
              <span className="text-2xl mr-3">{category.icon}</span>
              <h3 className="text-lg font-medium text-gray-900">
                {category.name}
              </h3>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {categoryActions.map((action) => {
                const commonContent = (
                  <>
                    <div className="text-2xl mr-4 flex-shrink-0">
                      {action.icon}
                    </div>
                    <div className="flex-1">
                      <h4 className="font-medium text-gray-900 mb-1">
                        {action.title}
                        {action.external && (
                          <span className="ml-1 text-xs text-gray-500">↗</span>
                        )}
                      </h4>
                      <p className="text-sm text-gray-500">
                        {action.description}
                      </p>
                      {loading && action.id === "test-connections" && (
                        <div className="flex items-center mt-2">
                          <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600 mr-2"></div>
                          <span className="text-xs text-blue-600">
                            Testing...
                          </span>
                        </div>
                      )}
                    </div>
                  </>
                );

                if (action.href) {
                  return (
                    <Link
                      key={action.id}
                      href={action.href}
                      className="flex items-start p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors text-left"
                    >
                      {commonContent}
                    </Link>
                  );
                }

                return (
                  <button
                    key={action.id}
                    onClick={action.action}
                    disabled={loading && action.id === "test-connections"}
                    className="flex items-start p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors text-left disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {commonContent}
                  </button>
                );
              })}
            </div>
          </div>
        );
      })}

      {/* Help Section */}
      <div className="bg-blue-50 rounded-lg p-6">
        <div className="flex items-start">
          <IoHelpCircle className="h-6 w-6 text-blue-600 mt-1 mr-3 flex-shrink-0" />
          <div>
            <h3 className="text-lg font-medium text-blue-900 mb-2">
              Need Help?
            </h3>
            <p className="text-blue-700 mb-4">
              Can&apos;t find what you&apos;re looking for? Our support team is
              here to help you with any questions or issues.
            </p>
            <div className="flex flex-wrap gap-3">
              <button
                onClick={() =>
                  window.open("mailto:support@cloudgate.com", "_blank")
                }
                className="inline-flex items-center px-4 py-2 border border-blue-300 rounded-md text-sm font-medium text-blue-700 bg-white hover:bg-blue-50 transition-colors"
              >
                <IoShare className="h-4 w-4 mr-2" />
                Email Support
              </button>
              <button
                onClick={() => alert("Live chat coming soon!")}
                className="inline-flex items-center px-4 py-2 border border-blue-300 rounded-md text-sm font-medium text-blue-700 bg-white hover:bg-blue-50 transition-colors"
              >
                💬 Live Chat
              </button>
              <button
                onClick={() => alert("Knowledge base coming soon!")}
                className="inline-flex items-center px-4 py-2 border border-blue-300 rounded-md text-sm font-medium text-blue-700 bg-white hover:bg-blue-50 transition-colors"
              >
                📚 Knowledge Base
              </button>
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
