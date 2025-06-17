"use client";

import DashboardLayout from "@/components/DashboardLayout";
import OAuthConnectionDashboard from "@/components/dashboard/OAuthConnectionDashboard";

export default function ConnectionsPage() {
  return (
    <DashboardLayout
      title="App Connections"
      description="Manage your connected applications and services"
    >
      <OAuthConnectionDashboard />
    </DashboardLayout>
  );
}
