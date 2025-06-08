"use client";

import { keycloak, keycloakInitOptions } from "@/lib/keycloak";
import { ReactKeycloakProvider } from "@react-keycloak/web";
import { ReactNode } from "react";

interface KeycloakProviderProps {
  children: ReactNode;
}

const LoadingComponent = (
  <div className="flex items-center justify-center min-h-screen">
    <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
  </div>
);

export default function KeycloakProvider({ children }: KeycloakProviderProps) {
  return (
    <ReactKeycloakProvider
      authClient={keycloak}
      initOptions={keycloakInitOptions}
      LoadingComponent={LoadingComponent}
    >
      {children}
    </ReactKeycloakProvider>
  );
}
