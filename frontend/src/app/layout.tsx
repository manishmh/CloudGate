"use client";

import { store } from "@/store";
import { ReactKeycloakProvider } from "@react-keycloak/web";
import Keycloak from "keycloak-js";
import { Inter } from "next/font/google";
import { Provider } from "react-redux";
import { Toaster } from "sonner";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

// Keycloak configuration with extensive logging
const keycloakConfig = {
  url: process.env.NEXT_PUBLIC_KEYCLOAK_URL || "http://localhost:8080",
  realm: process.env.NEXT_PUBLIC_KEYCLOAK_REALM || "cloudgate",
  clientId: process.env.NEXT_PUBLIC_KEYCLOAK_CLIENT_ID || "cloudgate-frontend",
};

console.log("🔧 Keycloak Configuration:", {
  url: keycloakConfig.url,
  realm: keycloakConfig.realm,
  clientId: keycloakConfig.clientId,
  environment: process.env.NODE_ENV,
  timestamp: new Date().toISOString(),
});

const keycloak = new Keycloak(keycloakConfig);

// Keycloak initialization options with debugging
const keycloakInitOptions = {
  onLoad: "check-sso" as const,
  checkLoginIframe: false,
  pkceMethod: "S256" as const,
  enableLogging: true,
  flow: "standard" as const,
  responseMode: "fragment" as const,
  messageReceiveTimeout: 10000,
  silentCheckSsoRedirectUri:
    typeof window !== "undefined"
      ? `${window.location.origin}/silent-check-sso.html`
      : undefined,
};

console.log("🚀 Keycloak Init Options:", keycloakInitOptions);

// Add event listeners for debugging
if (typeof window !== "undefined") {
  keycloak.onReady = (authenticated) => {
    console.log("🔐 Keycloak Ready:", {
      authenticated,
      timestamp: new Date().toISOString(),
    });
  };

  keycloak.onAuthSuccess = () => {
    console.log("✅ Keycloak Auth Success:", {
      token: keycloak.token ? "present" : "missing",
      refreshToken: keycloak.refreshToken ? "present" : "missing",
      timestamp: new Date().toISOString(),
    });
  };

  keycloak.onAuthError = (error) => {
    console.error("❌ Keycloak Auth Error:", {
      error,
      timestamp: new Date().toISOString(),
    });
  };

  keycloak.onAuthLogout = () => {
    console.log("🚪 Keycloak Logout:", { timestamp: new Date().toISOString() });
  };

  keycloak.onTokenExpired = () => {
    console.warn("⏰ Keycloak Token Expired:", {
      timestamp: new Date().toISOString(),
    });
  };
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className} suppressHydrationWarning>
        <Provider store={store}>
          <ReactKeycloakProvider
            authClient={keycloak}
            initOptions={keycloakInitOptions}
          >
            {children}
            <Toaster position="top-right" richColors />
          </ReactKeycloakProvider>
        </Provider>
      </body>
    </html>
  );
}
