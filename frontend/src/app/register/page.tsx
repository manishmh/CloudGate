"use client";

import { useAuth } from "@/components/providers/AuthProvider";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function RegisterPage() {
  const { register, login } = useAuth();
  const router = useRouter();
  const [form, setForm] = useState({
    email: "",
    username: "",
    first_name: "",
    last_name: "",
    password: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e?: React.FormEvent<HTMLFormElement>) => {
    if (e) e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      await register(form);
      await login(form.email, form.password);
      router.replace("/dashboard");
    } catch (e: any) {
      setError(e?.message || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="max-w-md w-full space-y-8 p-8">
        <div className="bg-white rounded-2xl shadow-xl p-8">
          <div className="text-center">
            <h2 className="text-3xl font-bold text-gray-900 mb-2">
              Create account
            </h2>
            <p className="text-black mb-8">Join CloudGate SSO</p>
          </div>

          <form className="space-y-4" onSubmit={handleSubmit}>
            <input
              className="w-full p-3 border rounded-lg text-black"
              placeholder="Email"
              type="email"
              value={form.email}
              required
              autoComplete="email"
              onChange={(e) => setForm({ ...form, email: e.target.value })}
            />
            <input
              className="w-full p-3 border rounded-lg text-black"
              placeholder="Username"
              value={form.username}
              required
              autoComplete="username"
              onChange={(e) => setForm({ ...form, username: e.target.value })}
            />
            <div className="grid grid-cols-2 gap-3">
              <input
                className="w-full p-3 border rounded-lg text-black"
                placeholder="First name"
                value={form.first_name}
                required
                onChange={(e) =>
                  setForm({ ...form, first_name: e.target.value })
                }
              />
              <input
                className="w-full p-3 border rounded-lg text-black"
                placeholder="Last name"
                value={form.last_name}
                required
                onChange={(e) =>
                  setForm({ ...form, last_name: e.target.value })
                }
              />
            </div>
            <input
              className="w-full p-3 border rounded-lg text-black"
              placeholder="Password"
              type="password"
              value={form.password}
              required
              minLength={8}
              autoComplete="new-password"
              onChange={(e) => setForm({ ...form, password: e.target.value })}
            />
            {error && <p className="text-red-600 text-sm">{error}</p>}
            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 rounded-lg text-white bg-blue-600 hover:bg-blue-700 cursor-pointer"
            >
              {loading ? "Creating..." : "Create account"}
            </button>
            <div className="text-sm text-gray-600 text-center">
              Have an account?{" "}
              <Link href="/login" className="text-blue-600 hover:underline">
                Login
              </Link>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
