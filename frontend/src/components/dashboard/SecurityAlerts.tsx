"use client";

import type { SecurityAlert } from "@/lib/api";
import {
  IoAlertCircle,
  IoCheckmarkCircle,
  IoInformationCircle,
  IoWarning,
} from "react-icons/io5";

interface SecurityAlertsProps {
  alerts: SecurityAlert[];
  onUpdateStatus: (alertId: string, status: SecurityAlert["status"]) => void;
}

const getSeverityConfig = (severity: SecurityAlert["severity"]) => {
  switch (severity) {
    case "critical":
      return {
        icon: IoAlertCircle,
        color: "text-red-600",
        bgColor: "bg-red-50",
      };
    case "high":
      return {
        icon: IoWarning,
        color: "text-orange-600",
        bgColor: "bg-orange-50",
      };
    case "medium":
      return {
        icon: IoInformationCircle,
        color: "text-yellow-600",
        bgColor: "bg-yellow-50",
      };
    case "low":
    default:
      return {
        icon: IoCheckmarkCircle,
        color: "text-blue-600",
        bgColor: "bg-blue-50",
      };
  }
};

export default function SecurityAlerts({
  alerts,
  onUpdateStatus,
}: SecurityAlertsProps) {
  if (!alerts || alerts.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow p-6 text-center">
        <IoCheckmarkCircle className="mx-auto h-12 w-12 text-green-500" />
        <h3 className="mt-2 text-sm font-medium text-gray-900">All Clear</h3>
        <p className="mt-1 text-sm text-gray-500">
          No security alerts at this time.
        </p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-6">
        <h3 className="text-lg font-medium text-gray-900">Security Alerts</h3>
        <p className="mt-1 text-sm text-gray-500">
          Real-time threats and suspicious activities detected in your account.
        </p>
      </div>
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th
                scope="col"
                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Alert
              </th>
              <th
                scope="col"
                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Severity
              </th>
              <th
                scope="col"
                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Date
              </th>
              <th
                scope="col"
                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Status
              </th>
              <th scope="col" className="relative px-6 py-3">
                <span className="sr-only">Actions</span>
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {alerts.map((alert) => {
              const severityConfig = getSeverityConfig(alert.severity);
              return (
                <tr key={alert.id}>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-full">
                        <severityConfig.icon
                          className={`h-6 w-6 ${severityConfig.color}`}
                          aria-hidden="true"
                        />
                      </div>
                      <div className="ml-4">
                        <div className="text-sm font-medium text-gray-900">
                          {alert.description}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span
                      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${severityConfig.bgColor} ${severityConfig.color}`}
                    >
                      {alert.severity}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {new Date(alert.created_at).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 capitalize">
                    {alert.status}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    {alert.status === "new" && (
                      <button
                        onClick={() => onUpdateStatus(alert.id, "acknowledged")}
                        className="text-indigo-600 hover:text-indigo-900"
                      >
                        Acknowledge
                      </button>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
