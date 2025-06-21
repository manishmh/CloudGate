"use client";

import { DASHBOARD_NAV_ITEMS } from "@/constants";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { setSidebarOpen } from "@/store/slices/sidebarSlice";
import { useKeycloak } from "@react-keycloak/web";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import {
  HiCog,
  HiDocumentText,
  HiHome,
  HiLightningBolt,
  HiLink,
  HiLogout,
  HiShieldCheck,
  HiUser,
  HiViewGrid,
  HiX,
} from "react-icons/hi";
import { MdOutlineSecurity } from "react-icons/md";

const iconMap = {
  HiHome,
  HiViewGrid,
  HiLink,
  HiShieldCheck,
  MdOutlineSecurity,
  HiLightningBolt,
  HiDocumentText,
};

export default function Sidebar() {
  const { keycloak } = useKeycloak();
  const pathname = usePathname();
  const dispatch = useAppDispatch();
  const { isOpen: sidebarOpen } = useAppSelector((state) => state.sidebar);
  const [isHovered, setIsHovered] = useState(false);
  const [profilePicture, setProfilePicture] = useState<string | null>(null);

  // Determine if sidebar should be expanded
  // On desktop: expand if open (clicked) OR hovered (when not open)
  // On mobile: expand only if open
  const shouldExpand = sidebarOpen || (isHovered && !sidebarOpen);

  useEffect(() => {
    // Load profile picture
    if (keycloak?.tokenParsed?.sub) {
      const savedPicture = localStorage.getItem(
        `profile_picture_${keycloak.tokenParsed.sub}`
      );
      if (savedPicture) {
        setProfilePicture(savedPicture);
      }
    }
  }, [keycloak]);

  const handleLogout = () => {
    if (keycloak) {
      keycloak.logout({
        redirectUri: `${window.location.origin}/login`,
      });
    }
  };

  const handleCloseSidebar = () => {
    dispatch(setSidebarOpen(false));
  };

  const getUserDisplayName = () => {
    if (keycloak?.tokenParsed) {
      return (
        keycloak.tokenParsed.preferred_username ||
        keycloak.tokenParsed.name ||
        keycloak.tokenParsed.email ||
        "User"
      );
    }
    return "User";
  };

  const getUserEmail = () => {
    if (keycloak?.tokenParsed) {
      return keycloak.tokenParsed.email || "";
    }
    return "";
  };

  return (
    <>
      {/* Overlay for mobile */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
          onClick={handleCloseSidebar}
        />
      )}

      {/* Sidebar */}
      <div
        className={`fixed left-0 top-0 h-full bg-white border-r border-gray-200 shadow-lg transition-all duration-300 ease-in-out ${
          shouldExpand ? "w-64" : "w-16"
        } ${
          sidebarOpen
            ? "translate-x-0 z-50"
            : "-translate-x-full lg:translate-x-0 lg:z-40"
        }`}
        onMouseEnter={() => !sidebarOpen && setIsHovered(true)}
        onMouseLeave={() => !sidebarOpen && setIsHovered(false)}
        style={{
          zIndex: sidebarOpen ? 50 : isHovered ? 45 : 40,
        }}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 py-5 border-b border-gray-200">
          <div className="flex items-center space-x-3">
            <div className="h-8 w-8 bg-blue-600 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg
                className="h-5 w-5 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                />
              </svg>
            </div>
            {shouldExpand && (
              <span className="font-bold text-gray-900 text-lg">CloudGate</span>
            )}
          </div>
          {shouldExpand && (
            <button
              onClick={handleCloseSidebar}
              className="p-1 rounded-md hover:bg-gray-100 lg:hidden cursor-pointer"
            >
              <HiX className="h-5 w-5 text-gray-500" />
            </button>
          )}
        </div>

        {/* User Banner */}
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center space-x-3">
            <div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center flex-shrink-0 overflow-hidden">
              {profilePicture ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={profilePicture}
                  alt="Profile"
                  className="h-full w-full object-cover"
                />
              ) : (
                <HiUser className="h-6 w-6 text-blue-600" />
              )}
            </div>
            {shouldExpand && (
              <div className="min-w-0 flex-1">
                <p className="text-sm font-medium text-gray-900 truncate">
                  {getUserDisplayName()}
                </p>
                <p className="text-xs text-gray-500 truncate">
                  {getUserEmail()}
                </p>
              </div>
            )}
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-4 space-y-2">
          {DASHBOARD_NAV_ITEMS.map((item) => {
            const isActive = pathname === item.href;
            const IconComponent = iconMap[item.icon as keyof typeof iconMap];
            return (
              <Link
                key={item.name}
                href={item.href}
                className={`flex items-center space-x-3 px-2 py-2  rounded-lg text-sm font-medium transition-colors duration-200 cursor-pointer 
                ${sidebarOpen || isHovered ? "justify-start" : "justify-center"}
                ${
                  isActive
                    ? "bg-blue-50 text-blue-700 border-r-2 border-blue-700"
                    : "text-gray-700 hover:text-gray-900 hover:bg-gray-100"
                }
                `}
                title={!shouldExpand ? item.description : undefined}
              >
                <IconComponent
                  className={`h-5 w-5 flex-shrink-0 ${
                    isActive ? "text-blue-700" : "text-gray-400"
                  }`}
                />
                {shouldExpand && <span>{item.name}</span>}
              </Link>
            );
          })}
        </nav>

        {/* Bottom Section */}
        <div className="p-4 border-t border-gray-200 space-y-2">
          <Link
            href="/dashboard/profile"
            className="flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900 transition-colors duration-200 cursor-pointer"
            title={!shouldExpand ? "Profile Settings" : undefined}
          >
            <HiUser className="h-5 w-5 text-gray-400 flex-shrink-0" />
            {shouldExpand && <span>Profile</span>}
          </Link>
          <Link
            href="/dashboard/settings"
            className="flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900 transition-colors duration-200 cursor-pointer"
            title={!shouldExpand ? "Settings" : undefined}
          >
            <HiCog className="h-5 w-5 text-gray-400 flex-shrink-0" />
            {shouldExpand && <span>Settings</span>}
          </Link>
          <button
            onClick={handleLogout}
            className="w-full flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium text-red-600 hover:bg-red-50 transition-colors duration-200 cursor-pointer"
            title={!shouldExpand ? "Logout" : undefined}
          >
            <HiLogout className="h-5 w-5 flex-shrink-0" />
            {shouldExpand && <span>Logout</span>}
          </button>
        </div>
      </div>
    </>
  );
}
