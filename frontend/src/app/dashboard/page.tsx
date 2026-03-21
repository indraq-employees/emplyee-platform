"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { getUser, logout } from "@/store/auth-store";

export default function DashboardPage() {
  const router = useRouter();
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const user = getUser();
    if (!user) {
      router.push("/login");
      return;
    }

    apiFetch<any>("/auth/me")
      .then((res) => setData(res.data))
      .catch(() => router.push("/login"))
      .finally(() => setLoading(false));
  }, [router]);

  if (loading) return <div className="p-6">Loading...</div>;
  if (!data) return null;

  const user = data.user;

  return (
    <main className="min-h-screen p-6">
      <div className="mx-auto max-w-5xl space-y-6">
        <div className="flex items-center justify-between rounded-2xl border p-6">
          <div>
            <h1 className="text-3xl font-bold">Welcome, {user.firstName}</h1>
            <p className="mt-2 text-sm text-gray-600">Role: {user.role}</p>
          </div>
          <button
            className="rounded-xl border px-4 py-2"
            onClick={() => {
              logout();
              router.push("/login");
            }}
          >
            Logout
          </button>
        </div>

        <div className="rounded-2xl border p-6">
          <h2 className="text-xl font-semibold">Account Summary</h2>
          <div className="mt-4 space-y-2 text-sm">
            <p>Email: {user.email}</p>
            <p>Profile Completed: {user.profileCompleted ? "Yes" : "No"}</p>
          </div>
        </div>

        {user.role !== "employee" ? (
          <div className="rounded-2xl border p-6">
            <h2 className="text-xl font-semibold">Admin Actions</h2>
            <div className="mt-4 flex gap-3">
              <Link className="rounded-xl bg-black px-4 py-2 text-white" href="/dashboard/employees/create">
                Create Employee
              </Link>
              <Link className="rounded-xl border px-4 py-2" href="/dashboard/employees">
                View Employees
              </Link>
            </div>
          </div>
        ) : (
          <div className="rounded-2xl border p-6">
            <h2 className="text-xl font-semibold">Employee Actions</h2>
            <div className="mt-4 flex gap-3">
              <Link className="rounded-xl bg-black px-4 py-2 text-white" href="/complete-profile">
                Complete Profile
              </Link>
            </div>
          </div>
        )}
      </div>
    </main>
  );
}