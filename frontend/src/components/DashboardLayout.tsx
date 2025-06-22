"use client";

import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { setSidebarOpen } from "@/store/slices/sidebarSlice";
import { useKeycloak } from "@react-keycloak/web";
import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";
import { HiHome, HiMenuAlt2 } from "react-icons/hi";
import Sidebar from "./Sidebar";

interface DashboardLayoutProps {
  children: React.ReactNode;
}

const breadcrumbNameMap: { [key: string]: string } = {
  "/dashboard": "Dashboard",
  "/dashboard/applications": "Applications",
  "/dashboard/connections": "Connections",
  "/dashboard/security": "Security",
  "/dashboard/advanced-security": "Advanced Security",
  "/dashboard/audit-log": "Audit Log",
  "/dashboard/settings": "Settings",
  "/dashboard/profile": "User Profile",
  "/dashboard/quick-actions": "Quick Actions",
  "/dashboard/oauth-test": "OAuth Test",
};

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const { keycloak, initialized } = useKeycloak();
  const router = useRouter();
  const pathname = usePathname();
  const dispatch = useAppDispatch();
  const { isOpen: sidebarOpen } = useAppSelector((state) => state.sidebar);

  const pathParts = pathname.split("/").filter((part) => part);

  useEffect(() => {
    if (initialized && !keycloak?.authenticated) {
      router.push("/login");
    }
  }, [initialized, keycloak?.authenticated, router]);

  // Load sidebar state from localStorage on mount
  useEffect(() => {
    const savedSidebarState = localStorage.getItem("sidebarOpen");
    if (savedSidebarState === "true") {
      dispatch(setSidebarOpen(true));
    }
  }, [dispatch]);

  // Save sidebar state to localStorage when it changes
  useEffect(() => {
    localStorage.setItem("sidebarOpen", sidebarOpen.toString());
  }, [sidebarOpen]);

  if (!initialized) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (!keycloak?.authenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Redirecting to login...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Sidebar />

      {/* Main Content */}
      <div
        className={`transition-all duration-300 ${
          sidebarOpen ? "lg:pl-64" : "lg:pl-20"
        }`}
      >
        {/* Top Bar with integrated Header */}
        <div className="sticky top-0 z-30 bg-white border-b border-gray-200 shadow-sm">
          <div className="flex h-16 items-center">
            <button
              type="button"
              className="border-r border-gray-200 px-4 text-gray-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 lg:hidden"
              onClick={() => dispatch(setSidebarOpen(true))}
            >
              <span className="sr-only">Open sidebar</span>
              <HiMenuAlt2 className="h-6 w-6" aria-hidden="true" />
            </button>

            <div className="flex flex-1 justify-between px-4">
              <div className="flex flex-col justify-center">
                <h1 className="text-xl font-bold tracking-tight text-gray-900">
                  {breadcrumbNameMap[pathname] || "Dashboard"}
                </h1>
                <nav className="flex" aria-label="Breadcrumb">
                  <ol role="list" className="flex items-center space-x-2">
                    <li>
                      <div>
                        <a
                          href="/dashboard"
                          className="text-gray-400 hover:text-gray-500"
                        >
                          <HiHome
                            className="h-4 w-4 flex-shrink-0"
                            aria-hidden="true"
                          />
                          <span className="sr-only">Home</span>
                        </a>
                      </div>
                    </li>
                    {pathParts.map((part, index) => {
                      const href =
                        "/" + pathParts.slice(0, index + 1).join("/");
                      const isLast = index === pathParts.length - 1;
                      const name =
                        breadcrumbNameMap[href] ||
                        part.charAt(0).toUpperCase() + part.slice(1);

                      return (
                        <li key={href}>
                          <div className="flex items-center">
                            <svg
                              className="h-4 w-4 flex-shrink-0 text-gray-300"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                              aria-hidden="true"
                            >
                              <path d="M5.555 17.776l8-16 .894.448-8 16-.894-.448z" />
                            </svg>
                            <a
                              href={href}
                              className={`ml-2 text-xs font-medium ${
                                isLast
                                  ? "text-gray-500"
                                  : "text-gray-400 hover:text-gray-500"
                              }`}
                              aria-current={isLast ? "page" : undefined}
                            >
                              {name}
                            </a>
                          </div>
                        </li>
                      );
                    })}
                  </ol>
                </nav>
              </div>
              <div className="ml-4 flex items-center md:ml-6">
                {/* Profile dropdown */}
              </div>
            </div>
          </div>
        </div>

        {/* Page Content */}
        <main className="flex-1">
          <div className="py-6">
            <div className="mx-auto max-w-7xl px-4 sm:px-6 md:px-8">
              {children}
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
