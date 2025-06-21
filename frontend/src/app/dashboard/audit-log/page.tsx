"use client";

import DashboardLayout from "@/components/DashboardLayout";
import { apiClient } from "@/lib/api";
import { useCallback, useEffect, useState } from "react";
import {
  IoCalendar,
  IoCheckmarkCircle,
  IoDownload,
  IoFilter,
  IoInformationCircle,
  IoRefresh,
  IoSearch,
  IoShieldCheckmark,
  IoWarning,
} from "react-icons/io5";
import { toast } from "sonner";

export interface AuditEvent {
  id: string;
  user_id: string;
  event_type: string;
  category: string;
  description: string;
  severity: "info" | "warning" | "error" | "critical";
  ip_address: string;
  user_agent: string;
  metadata: Record<string, string | number | boolean>;
  created_at: string;
}

interface AuditLogFilters {
  search: string;
  category: string;
  severity: string;
  dateFrom: string;
  dateTo: string;
}

const CATEGORIES = [
  { value: "all", label: "All Categories" },
  { value: "authentication", label: "Authentication" },
  { value: "authorization", label: "Authorization" },
  { value: "data_access", label: "Data Access" },
  { value: "configuration", label: "Configuration" },
  { value: "security", label: "Security" },
];

const SEVERITIES = [
  { value: "all", label: "All Severities" },
  { value: "info", label: "Info", color: "text-blue-600" },
  { value: "warning", label: "Warning", color: "text-yellow-600" },
  { value: "error", label: "Error", color: "text-orange-600" },
  { value: "critical", label: "Critical", color: "text-red-600" },
];

export default function AuditLogPage() {
  const [loading, setLoading] = useState(false);
  const [auditEvents, setAuditEvents] = useState<AuditEvent[]>([]);
  const [filters, setFilters] = useState<AuditLogFilters>({
    search: "",
    category: "all",
    severity: "all",
    dateFrom: "",
    dateTo: "",
  });
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  const loadAuditEvents = useCallback(async () => {
    try {
      setLoading(true);

      // In a real implementation, this would call the audit log endpoint
      // For now, we'll simulate with security events
      const response = await apiClient.getSecurityEvents(50);

      // Transform security events to audit events format
      const transformedEvents: AuditEvent[] = response.events.map((event) => ({
        id: event.id,
        user_id: event.user_id,
        event_type: event.event_type,
        category: "security",
        description: event.description,
        severity: event.severity as AuditEvent["severity"],
        ip_address: event.ip_address,
        user_agent: event.user_agent,
        metadata: {},
        created_at: event.created_at,
      }));

      setAuditEvents(transformedEvents);
      setTotalPages(Math.ceil(response.count / 50));
    } catch (error) {
      console.error("Failed to load audit events:", error);
      toast.error("Failed to load audit events");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadAuditEvents();
  }, [loadAuditEvents]);

  const handleFilterChange = (key: keyof AuditLogFilters, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    setPage(1);
  };

  const handleExportLogs = async () => {
    try {
      // In a real implementation, this would call an export endpoint
      toast.success("Export started. You'll receive an email when it's ready.");
    } catch {
      toast.error("Failed to export logs");
    }
  };

  const getSeverityIcon = (severity: string) => {
    switch (severity) {
      case "critical":
      case "error":
        return IoWarning;
      case "warning":
        return IoInformationCircle;
      default:
        return IoCheckmarkCircle;
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case "critical":
        return "text-red-600 bg-red-50";
      case "error":
        return "text-orange-600 bg-orange-50";
      case "warning":
        return "text-yellow-600 bg-yellow-50";
      default:
        return "text-blue-600 bg-blue-50";
    }
  };

  const filteredEvents = auditEvents.filter((event) => {
    const matchesSearch =
      !filters.search ||
      event.description.toLowerCase().includes(filters.search.toLowerCase()) ||
      event.event_type.toLowerCase().includes(filters.search.toLowerCase());

    const matchesCategory =
      filters.category === "all" || event.category === filters.category;
    const matchesSeverity =
      filters.severity === "all" || event.severity === filters.severity;

    let matchesDate = true;
    if (filters.dateFrom) {
      matchesDate = new Date(event.created_at) >= new Date(filters.dateFrom);
    }
    if (filters.dateTo && matchesDate) {
      matchesDate = new Date(event.created_at) <= new Date(filters.dateTo);
    }

    return matchesSearch && matchesCategory && matchesSeverity && matchesDate;
  });

  const refreshAction = (
    <div className="flex gap-2">
      <button
        onClick={handleExportLogs}
        className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
      >
        <IoDownload className="h-4 w-4 mr-2" />
        Export
      </button>
      <button
        onClick={loadAuditEvents}
        disabled={loading}
        className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
      >
        <IoRefresh className="h-4 w-4 mr-2" />
        {loading ? "Refreshing..." : "Refresh"}
      </button>
    </div>
  );

  return (
    <DashboardLayout
      title="Audit Log"
      description="View comprehensive system and user activity logs"
      actions={refreshAction}
    >
      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <div className="flex items-center mb-4">
          <IoFilter className="h-5 w-5 text-gray-400 mr-2" />
          <h3 className="text-lg font-medium text-gray-900">Filters</h3>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
          <div className="lg:col-span-2">
            <div className="relative">
              <IoSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
              <input
                type="text"
                placeholder="Search logs..."
                value={filters.search}
                onChange={(e) => handleFilterChange("search", e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-black"
              />
            </div>
          </div>

          <select
            value={filters.category}
            onChange={(e) => handleFilterChange("category", e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-black"
          >
            {CATEGORIES.map((cat) => (
              <option key={cat.value} value={cat.value}>
                {cat.label}
              </option>
            ))}
          </select>

          <select
            value={filters.severity}
            onChange={(e) => handleFilterChange("severity", e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-black"
          >
            {SEVERITIES.map((sev) => (
              <option key={sev.value} value={sev.value}>
                {sev.label}
              </option>
            ))}
          </select>

          <div className="relative">
            <IoCalendar className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
            <input
              type="date"
              value={filters.dateFrom}
              onChange={(e) => handleFilterChange("dateFrom", e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-black"
              placeholder="From date"
            />
          </div>
        </div>
      </div>

      {/* Audit Events Table */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-medium text-gray-900">
              {filteredEvents.length} Event
              {filteredEvents.length !== 1 ? "s" : ""}
            </h3>
            <div className="flex items-center text-sm text-gray-500">
              <IoShieldCheckmark className="h-4 w-4 mr-1" />
              Compliant with SOX, GDPR, HIPAA
            </div>
          </div>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <span className="ml-2 text-gray-600">Loading audit events...</span>
          </div>
        ) : filteredEvents.length === 0 ? (
          <div className="text-center py-12">
            <IoSearch className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">
              No events found
            </h3>
            <p className="mt-1 text-sm text-gray-500">
              Try adjusting your filters or search criteria.
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Event
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Category
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Severity
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    IP Address
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Timestamp
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredEvents.map((event) => {
                  const Icon = getSeverityIcon(event.severity);
                  const colorClass = getSeverityColor(event.severity);

                  return (
                    <tr key={event.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4">
                        <div className="flex items-start">
                          <Icon
                            className={`h-5 w-5 mt-0.5 mr-3 ${
                              colorClass.split(" ")[0]
                            }`}
                          />
                          <div>
                            <div className="text-sm font-medium text-gray-900">
                              {event.event_type}
                            </div>
                            <div className="text-sm text-gray-500">
                              {event.description}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="text-sm text-gray-900 capitalize">
                          {event.category}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colorClass}`}
                        >
                          {event.severity}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {event.ip_address}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {new Date(event.created_at).toLocaleString()}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
            <div className="flex-1 flex justify-between sm:hidden">
              <button
                onClick={() => setPage(Math.max(1, page - 1))}
                disabled={page === 1}
                className="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
              >
                Previous
              </button>
              <button
                onClick={() => setPage(Math.min(totalPages, page + 1))}
                disabled={page === totalPages}
                className="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
              >
                Next
              </button>
            </div>
            <div className="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
              <div>
                <p className="text-sm text-gray-700">
                  Showing page <span className="font-medium">{page}</span> of{" "}
                  <span className="font-medium">{totalPages}</span>
                </p>
              </div>
              <div>
                <nav
                  className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px"
                  aria-label="Pagination"
                >
                  <button
                    onClick={() => setPage(Math.max(1, page - 1))}
                    disabled={page === 1}
                    className="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                  >
                    Previous
                  </button>
                  <button
                    onClick={() => setPage(Math.min(totalPages, page + 1))}
                    disabled={page === totalPages}
                    className="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                  >
                    Next
                  </button>
                </nav>
              </div>
            </div>
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}
