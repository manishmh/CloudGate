import Link from "next/link";
import {
  HiArrowRight,
  HiClock,
  HiExclamationCircle,
  HiLightningBolt,
  HiLink,
  HiShieldCheck,
  HiViewGrid,
} from "react-icons/hi";

interface ActivityItem {
  id: string;
  type: "login" | "app_launch" | "connection" | "security";
  description: string;
  timestamp: string;
  icon: string;
  severity?: "info" | "warning" | "success";
}

interface RecentActivityProps {
  activities: ActivityItem[];
}

const iconMap = {
  HiShieldCheck,
  HiViewGrid,
  HiLink,
  HiExclamationCircle,
  HiLightningBolt,
};

export default function RecentActivity({ activities }: RecentActivityProps) {
  return (
    <div>
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium text-gray-900">
              Recent Activity
            </h3>
            <Link
              href="/dashboard/security"
              className="text-sm text-blue-600 hover:text-blue-800 font-medium cursor-pointer"
            >
              View all <HiArrowRight className="inline h-4 w-4 ml-1" />
            </Link>
          </div>
        </div>
        <div className="p-6">
          <div className="space-y-4">
            {activities.slice(0, 5).map((activity) => {
              const IconComponent =
                iconMap[activity.icon as keyof typeof iconMap];
              return (
                <div key={activity.id} className="flex items-start space-x-3">
                  <div
                    className={`p-1 rounded-full ${
                      activity.severity === "success"
                        ? "bg-green-100"
                        : activity.severity === "warning"
                        ? "bg-yellow-100"
                        : "bg-blue-100"
                    }`}
                  >
                    <IconComponent
                      className={`h-4 w-4 ${
                        activity.severity === "success"
                          ? "text-green-600"
                          : activity.severity === "warning"
                          ? "text-yellow-600"
                          : "text-blue-600"
                      }`}
                    />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-gray-900">
                      {activity.description}
                    </p>
                    <div className="flex items-center mt-1">
                      <HiClock className="h-3 w-3 text-gray-400 mr-1" />
                      <p className="text-xs text-gray-500">
                        {activity.timestamp}
                      </p>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
