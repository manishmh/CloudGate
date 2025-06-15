import {
  HiLink,
  HiShieldCheck,
  HiTrendingUp,
  HiViewGrid,
} from "react-icons/hi";

interface DashboardMetrics {
  totalApps: number;
  connectedApps: number;
  recentLogins: number;
  securityScore: number;
  lastActivity: string;
}

interface DashboardMetricsProps {
  metrics: DashboardMetrics;
}

export default function DashboardMetrics({ metrics }: DashboardMetricsProps) {
  const metricCards = [
    {
      title: "Total Apps",
      value: metrics.totalApps,
      icon: HiViewGrid,
      color: "blue",
      bgColor: "bg-blue-100",
      iconColor: "text-blue-600",
    },
    {
      title: "Connected",
      value: metrics.connectedApps,
      icon: HiLink,
      color: "green",
      bgColor: "bg-green-100",
      iconColor: "text-green-600",
    },
    {
      title: "Security Score",
      value: `${metrics.securityScore}%`,
      icon: HiShieldCheck,
      color: "yellow",
      bgColor: "bg-yellow-100",
      iconColor: "text-yellow-600",
    },
    {
      title: "Recent Logins",
      value: metrics.recentLogins,
      icon: HiTrendingUp,
      color: "purple",
      bgColor: "bg-purple-100",
      iconColor: "text-purple-600",
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      {metricCards.map((metric) => (
        <div key={metric.title} className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className={`p-2 ${metric.bgColor} rounded-lg`}>
              <metric.icon className={`h-6 w-6 ${metric.iconColor}`} />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">
                {metric.title}
              </p>
              <p className="text-2xl font-semibold text-gray-900">
                {metric.value}
              </p>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
