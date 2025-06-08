import KeycloakProvider from "@/components/providers/KeycloakProvider";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "CloudGate SSO Portal",
  description: "Enterprise-grade Single Sign-On Portal with adaptive security",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <KeycloakProvider>{children}</KeycloakProvider>
      </body>
    </html>
  );
}
