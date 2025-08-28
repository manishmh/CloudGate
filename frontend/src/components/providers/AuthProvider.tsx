"use client";

import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

type AuthContextType = {
  isAuthenticated: boolean;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (data: {
    email: string;
    username: string;
    first_name: string;
    last_name: string;
    password: string;
  }) => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [token, setToken] = useState<string | null>(null);
  const isAuthenticated = !!token;

  useEffect(() => {
    if (typeof window !== "undefined") {
      const t = localStorage.getItem("auth_token");
      if (t) setToken(t);
    }
  }, []);

  const login = async (email: string, password: string) => {
    const res = await fetch(
      `${
        process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8081"
      }/auth/login`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
        credentials: "include",
      }
    );
    if (!res.ok) throw new Error("Login failed");
    const data = await res.json();
    setToken(data.access_token);
    if (typeof window !== "undefined") {
      localStorage.setItem("auth_token", data.access_token);
      localStorage.setItem("refresh_token", data.refresh_token);
    }
  };

  const logout = async () => {
    const refresh =
      typeof window !== "undefined"
        ? localStorage.getItem("refresh_token")
        : null;
    if (refresh) {
      await fetch(
        `${
          process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8081"
        }/auth/logout`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ refresh_token: refresh }),
          credentials: "include",
        }
      ).catch(() => {});
    }
    setToken(null);
    if (typeof window !== "undefined") {
      localStorage.removeItem("auth_token");
      localStorage.removeItem("refresh_token");
    }
  };

  const register: AuthContextType["register"] = async (data) => {
    const res = await fetch(
      `${
        process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8081"
      }/auth/register`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
        credentials: "include",
      }
    );
    if (!res.ok) {
      const msg = await res.text().catch(() => "Registration failed");
      throw new Error(msg || "Registration failed");
    }
  };

  const value = useMemo(
    () => ({ isAuthenticated, token, login, register, logout }),
    [isAuthenticated, token]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextType => {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
};
