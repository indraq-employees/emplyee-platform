"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { setAuth } from "@/store/auth-store";

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [form, setForm] = useState({ email: "", password: "" });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const email = searchParams.get("email") || "";
    if (email) {
      setForm((prev) => ({ ...prev, email }));
    }
  }, [searchParams]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const res = await apiFetch<any>("/auth/login", {
        method: "POST",
        body: JSON.stringify(form),
      });
      setAuth(res.data.token, res.data.user);
      if (!res.data.user.profileCompleted && res.data.user.role === "employee") {
        router.push("/complete-profile");
      } else {
        router.push("/dashboard");
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen flex items-center justify-center p-6">
      <form onSubmit={handleSubmit} className="w-full max-w-md rounded-2xl border p-6 shadow-sm space-y-4">
        <h1 className="text-2xl font-bold">Login</h1>
        <input
          value={form.email}
          onChange={(e) => setForm((prev) => ({ ...prev, email: e.target.value }))}
          type="email"
          placeholder="Email"
          className="w-full rounded-xl border p-3"
        />
        <input
          value={form.password}
          onChange={(e) => setForm((prev) => ({ ...prev, password: e.target.value }))}
          type="password"
          placeholder="Password"
          className="w-full rounded-xl border p-3"
        />
        {error ? <p className="text-sm text-red-600">{error}</p> : null}
        <button disabled={loading} className="w-full rounded-xl bg-black p-3 text-white">
          {loading ? "Logging in..." : "Login"}
        </button>
      </form>
    </main>
  );
}