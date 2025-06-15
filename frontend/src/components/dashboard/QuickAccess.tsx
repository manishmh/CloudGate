import Link from "next/link";
import { HiArrowRight, HiEye, HiLink, HiPlus } from "react-icons/hi";

interface AppConnection {
  name: string;
  status: "connected" | "disconnected";
  icon: string;
  description: string;
  connect_url: string;
  last_used?: string;
}

interface QuickAccessProps {
  connections: AppConnection[];
}

export default function QuickAccess({ connections }: QuickAccessProps) {
  const connectedApps = connections.filter(
    (conn) => conn.status === "connected"
  );

  return (
    <div className="lg:col-span-2">
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium text-gray-900">Quick Access</h3>
            <Link
              href="/dashboard/applications"
              className="text-sm text-blue-600 hover:text-blue-800 font-medium cursor-pointer"
            >
              View all <HiArrowRight className="inline h-4 w-4 ml-1" />
            </Link>
          </div>
        </div>
        <div className="p-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {connectedApps.slice(0, 4).map((connection) => (
              <div
                key={connection.name}
                className="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer"
              >
                <div className="text-2xl mr-3">{connection.icon}</div>
                <div className="flex-1">
                  <h4 className="text-sm font-medium text-gray-900">
                    {connection.name}
                  </h4>
                  <p className="text-xs text-gray-500">
                    Last used: {connection.last_used}
                  </p>
                </div>
                <button className="text-blue-600 hover:text-blue-800">
                  <HiEye className="h-4 w-4" />
                </button>
              </div>
            ))}
          </div>

          {connectedApps.length === 0 && (
            <div className="text-center py-8">
              <HiLink className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h4 className="text-sm font-medium text-gray-900 mb-2">
                No connected applications
              </h4>
              <p className="text-sm text-gray-500 mb-4">
                Connect your first application to get started
              </p>
              <Link
                href="/dashboard/applications"
                className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 cursor-pointer"
              >
                <HiPlus className="h-4 w-4 mr-2" />
                Connect Apps
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
