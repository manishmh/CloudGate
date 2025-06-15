"use client";

import { store } from "@/store";
import { ReactKeycloakProvider } from "@react-keycloak/web";
import Keycloak from "keycloak-js";
import { Inter } from "next/font/google";
import { Provider } from "react-redux";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

// Keycloak configuration
const keycloakConfig = {
  url: process.env.NEXT_PUBLIC_KEYCLOAK_URL || "http://localhost:8080",
  realm: process.env.NEXT_PUBLIC_KEYCLOAK_REALM || "cloudgate",
  clientId: process.env.NEXT_PUBLIC_KEYCLOAK_CLIENT_ID || "cloudgate-frontend",
};

const keycloak = new Keycloak(keycloakConfig);

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className} suppressHydrationWarning>
        <Provider store={store}>
          <ReactKeycloakProvider authClient={keycloak}>
            {children}
          </ReactKeycloakProvider>
        </Provider>
      </body>
    </html>
  );
}
