"use client";

import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { setSidebarOpen } from "@/store/slices/sidebarSlice";
import { useKeycloak } from "@react-keycloak/web";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { HiMenuAlt2 } from "react-icons/hi";
import Sidebar from "./Sidebar";
import Header from "./dashboard/Header";

interface DashboardLayoutProps {
  children: React.ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
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
        {/* Top Bar */}
        <div className="sticky top-0 z-30 flex h-16 flex-shrink-0 border-b border-gray-200 bg-white">
          <button
            type="button"
            className="border-r border-gray-200 px-4 text-gray-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 lg:hidden"
            onClick={() => dispatch(setSidebarOpen(true))}
          >
            <span className="sr-only">Open sidebar</span>
            <HiMenuAlt2 className="h-6 w-6" aria-hidden="true" />
          </button>
          <div className="flex flex-1 justify-between px-4">
            <div className="flex flex-1">
              <Header />
            </div>
            <div className="ml-4 flex items-center md:ml-6">
              {/* Profile dropdown */}
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
