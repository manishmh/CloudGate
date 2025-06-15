import Link from "next/link";
import { HiShieldCheck, HiViewGrid, HiLink, HiExclamationCircle, HiLightningBolt } from "react-icons/hi";

interface QuickAction {
  id: string;
  title: string;
  description: string;
  href: string;
  icon: string;
  gradient: string;
}

interface QuickActionsProps {
  actions: readonly QuickAction[];
}

const iconMap = {
  HiShieldCheck,
  HiViewGrid,
  HiLink,
  HiExclamationCircle,
  HiLightningBolt,
};

export default function QuickActions({ actions }: QuickActionsProps) {
  return (
    <div className="mt-8">
      <div className="bg-white rounded-2xl shadow-sm border border-gray-100">
        <div className="px-8 py-6 border-b border-gray-100">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-xl font-semibold text-gray-900">Quick Actions</h3>
              <p className="text-gray-600 text-sm mt-1">Common tasks and shortcuts</p>
            </div>
            <Link
              href="/dashboard/quick-actions"
              className="text-sm text-blue-600 hover:text-blue-800 font-medium cursor-pointer"
            >
              View all
            </Link>
          </div>
        </div>
        <div className="p-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {actions.map((action) => {
              const IconComponent = iconMap[action.icon as keyof typeof iconMap];
              return (
                <Link
                  key={action.id}
                  href={action.href}
                  className="group relative flex flex-col items-center p-8 border border-gray-100 rounded-2xl hover:bg-gray-50 hover:border-gray-200 transition-all duration-200 cursor-pointer overflow-hidden"
                >
                  {/* Gradient Background */}
                  <div className={`absolute inset-0 bg-gradient-to-br ${action.gradient} opacity-0 group-hover:opacity-5 transition-opacity duration-300`}></div>
                  
                  <div className="relative z-10 text-center flex flex-col items-center">
                    <div className={`w-12 h-12 bg-gradient-to-br ${action.gradient} rounded-2xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-200`}>
                      <IconComponent className="h-6 w-6 text-white" />
                    </div>
                    <h4 className="text-lg font-semibold text-gray-900 mb-2">
                      {action.title}
                    </h4>
                    <p className="text-sm text-gray-600">
                      {action.description}
                    </p>
                  </div>
                </Link>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
} 