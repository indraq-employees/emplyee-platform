"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { setAuth } from "@/store/auth-store";

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [step, setStep] = useState<"credentials" | "otp">("credentials");
  const [form, setForm] = useState({ email: "", password: "" });
  const [otp, setOtp] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const email = searchParams.get("email") || "";
    const invited = searchParams.get("invited") === "1";

    if (email) {
      setForm((prev) => ({ ...prev, email }));
    }

    if (invited) {
      setMessage("Your account is ready. Use the password shared in your email, then verify with OTP.");
    }
  }, [searchParams]);

  const sendOtp = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await apiFetch<any>("/auth/login/send-otp", {
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
      const res = await apiFetch<any>("/auth/login/verify-otp", {
        method: "POST",
        body: JSON.stringify({ email: form.email, otp }),
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
      <div className="w-full max-w-md rounded-2xl border p-6 shadow-sm space-y-4">
        <h1 className="text-2xl font-bold">Login</h1>
        {message ? <p className="text-sm text-gray-600">{message}</p> : null}

        {step === "credentials" ? (
          <form onSubmit={sendOtp} className="space-y-4">
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
              {loading ? "Sending OTP..." : "Continue"}
            </button>
          </form>
        ) : (
          <form onSubmit={verifyOtp} className="space-y-4">
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
                  await apiFetch<any>("/auth/login/send-otp", {
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