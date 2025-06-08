"use client";

import { useKeycloak } from "@react-keycloak/ssr";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function HomePage() {
  const { keycloak, initialized } = useKeycloak();
  const router = useRouter();

  useEffect(() => {
    if (initialized) {
      if (keycloak?.authenticated) {
        router.push("/dashboard");
      } else {
        router.push("/login");
      }
    }
  }, [initialized, keycloak?.authenticated, router]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600 mx-auto mb-4"></div>
        <p className="text-white">Loading CloudGate SSO Portal...</p>
      </div>
    </div>
  );
}
