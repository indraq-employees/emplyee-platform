"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { setAuth } from "@/store/auth-store";

export default function AdminSignupPage() {
  const router = useRouter();
  const [step, setStep] = useState<"form" | "otp">("form");
  const [form, setForm] = useState({
    firstName: "",
    lastName: "",
    email: "",
    password: "",
  });
  const [otp, setOtp] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const sendOtp = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setMessage("");
    setLoading(true);

    try {
      await apiFetch<any>("/auth/admin-signup/send-otp", {
        method: "POST",
        body: JSON.stringify(form),
      });
      setStep("otp");
      setMessage(`OTP sent to ${form.email}`);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const verifyOtp = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const res = await apiFetch<any>("/auth/admin-signup/verify-otp", {
        method: "POST",
        body: JSON.stringify({ email: form.email, otp }),
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
      <div className="w-full max-w-md rounded-2xl border p-6 shadow-sm space-y-4">
        <h1 className="text-2xl font-bold">Admin Signup</h1>

        {step === "form" ? (
          <form onSubmit={sendOtp} className="space-y-4">
            <input
              name="firstName"
              value={form.firstName}
              placeholder="First Name"
              className="w-full rounded-xl border p-3"
              onChange={handleChange}
            />
            <input
              name="lastName"
              value={form.lastName}
              placeholder="Last Name"
              className="w-full rounded-xl border p-3"
              onChange={handleChange}
            />
            <input
              name="email"
              type="email"
              value={form.email}
              placeholder="Email"
              className="w-full rounded-xl border p-3"
              onChange={handleChange}
            />
            <input
              name="password"
              type="password"
              value={form.password}
              placeholder="Password"
              className="w-full rounded-xl border p-3"
              onChange={handleChange}
            />
            {error ? <p className="text-sm text-red-600">{error}</p> : null}
            <button disabled={loading} className="w-full rounded-xl bg-black p-3 text-white">
              {loading ? "Sending OTP..." : "Create Admin"}
            </button>
          </form>
        ) : (
          <form onSubmit={verifyOtp} className="space-y-4">
            <p className="text-sm text-gray-600">{message || `Enter the OTP sent to ${form.email}`}</p>
            <input
              value={otp}
              onChange={(e) => setOtp(e.target.value.replace(/\D/g, "").slice(0, 6))}
              inputMode="numeric"
              placeholder="Enter 6-digit OTP"
              className="w-full rounded-xl border p-3"
            />
            {error ? <p className="text-sm text-red-600">{error}</p> : null}
            <button disabled={loading} className="w-full rounded-xl bg-black p-3 text-white">
              {loading ? "Verifying..." : "Verify OTP"}
            </button>
            <button
              type="button"
              disabled={loading}
              className="w-full rounded-xl border p-3"
              onClick={async () => {
                setError("");
                setLoading(true);
                try {
                  await apiFetch<any>("/auth/admin-signup/send-otp", {
                    method: "POST",
                    body: JSON.stringify(form),
                  });
                  setMessage(`New OTP sent to ${form.email}`);
                } catch (err: any) {
                  setError(err.message);
                } finally {
                  setLoading(false);
                }
              }}
            >
              Resend OTP
            </button>
          </form>
        )}
      </div>
    </main>
  );
}