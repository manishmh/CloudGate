"use client";

import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { setSidebarOpen, toggleSidebar } from "@/store/slices/sidebarSlice";
import { useKeycloak } from "@react-keycloak/web";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { HiMenuAlt2 } from "react-icons/hi";
import Sidebar from "./Sidebar";

interface DashboardLayoutProps {
  children: React.ReactNode;
  title?: string;
  description?: string;
  actions?: React.ReactNode;
}

export default function DashboardLayout({
  children,
  title,
  description,
  actions,
}: DashboardLayoutProps) {
  const { keycloak, initialized } = useKeycloak();
  const router = useRouter();
  const dispatch = useAppDispatch();
  const { isOpen: sidebarOpen } = useAppSelector((state) => state.sidebar);

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

  const handleSidebarToggle = () => {
    dispatch(toggleSidebar());
  };

  const handleMobileSidebarOpen = () => {
    dispatch(setSidebarOpen(true));
  };

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
          sidebarOpen ? "lg:pl-64" : "lg:pl-16"
        }`}
      >
        {/* Top Bar */}
        <div className="sticky top-0 z-30 bg-white border-b border-gray-200 px-4 py-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <button
                onClick={handleMobileSidebarOpen}
                className="p-2 rounded-md hover:bg-gray-100 lg:hidden cursor-pointer"
              >
                <HiMenuAlt2 className="h-5 w-5 text-gray-600" />
              </button>
              <button
                onClick={handleSidebarToggle}
                className="hidden lg:block p-2 rounded-md hover:bg-gray-100 cursor-pointer"
              >
                <HiMenuAlt2 className="h-5 w-5 text-gray-600" />
              </button>
              <div>
                {title && (
                  <h1 className="text-xl font-semibold text-gray-900">
                    {title}
                  </h1>
                )}
                {description && (
                  <p className="text-sm text-gray-600">{description}</p>
                )}
              </div>
            </div>
            {actions && (
              <div className="flex items-center space-x-2">{actions}</div>
            )}
          </div>
        </div>

        {/* Page Content */}
        <main className="p-6">{children}</main>
      </div>
    </div>
  );
}
