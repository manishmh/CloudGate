"use client";

import { usePathname } from "next/navigation";
import { HiHome } from "react-icons/hi";

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

const Header = () => {
  const pathname = usePathname();
  const pathParts = pathname.split("/").filter((part) => part);

  return (
    <header className="bg-white shadow-sm">
      <div className="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:px-8">
        <div className="flex flex-col">
          <h1 className="text-2xl font-bold tracking-tight text-gray-900">
            {breadcrumbNameMap[pathname] || "Dashboard"}
          </h1>
          <nav className="flex mt-2" aria-label="Breadcrumb">
            <ol role="list" className="flex items-center space-x-2">
              <li>
                <div>
                  <a
                    href="/dashboard"
                    className="text-gray-400 hover:text-gray-500"
                  >
                    <HiHome
                      className="h-5 w-5 flex-shrink-0"
                      aria-hidden="true"
                    />
                    <span className="sr-only">Home</span>
                  </a>
                </div>
              </li>
              {pathParts.map((part, index) => {
                const href = "/" + pathParts.slice(0, index + 1).join("/");
                const isLast = index === pathParts.length - 1;
                const name =
                  breadcrumbNameMap[href] ||
                  part.charAt(0).toUpperCase() + part.slice(1);

                return (
                  <li key={href}>
                    <div className="flex items-center">
                      <svg
                        className="h-5 w-5 flex-shrink-0 text-gray-300"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                        aria-hidden="true"
                      >
                        <path d="M5.555 17.776l8-16 .894.448-8 16-.894-.448z" />
                      </svg>
                      <a
                        href={href}
                        className={`ml-2 text-sm font-medium ${
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
      </div>
    </header>
  );
};

export default Header;
