"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { setAuth } from "@/store/auth-store";

export default function AdminSignupPage() {
  const router = useRouter();
  const [form, setForm] = useState({ firstName: "", lastName: "", email: "", password: "" });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const res = await apiFetch<any>("/auth/admin-signup", {
        method: "POST",
        body: JSON.stringify(form),
      });
      setAuth(res.data.token, res.data.user);
      router.push("/dashboard");
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen flex items-center justify-center p-6">
      <form onSubmit={handleSubmit} className="w-full max-w-md rounded-2xl border p-6 shadow-sm space-y-4">
        <h1 className="text-2xl font-bold">Admin Signup</h1>
        <input name="firstName" placeholder="First Name" className="w-full rounded-xl border p-3" onChange={handleChange} />
        <input name="lastName" placeholder="Last Name" className="w-full rounded-xl border p-3" onChange={handleChange} />
        <input name="email" type="email" placeholder="Email" className="w-full rounded-xl border p-3" onChange={handleChange} />
        <input name="password" type="password" placeholder="Password" className="w-full rounded-xl border p-3" onChange={handleChange} />
        {error ? <p className="text-sm text-red-600">{error}</p> : null}
        <button disabled={loading} className="w-full rounded-xl bg-black p-3 text-white">
          {loading ? "Creating..." : "Create Admin"}
        </button>
      </form>
    </main>
  );
}