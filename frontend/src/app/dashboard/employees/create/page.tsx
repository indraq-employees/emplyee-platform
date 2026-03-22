"use client";

import { useState } from "react";
import { apiFetch } from "@/lib/api";

export default function CreateEmployeePage() {
  const [form, setForm] = useState({
    firstName: "",
    lastName: "",
    email: "",
    password: "",
    role: "employee",
  });
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setResult(null);
    setLoading(true);

    try {
      const res = await apiFetch<any>("/employees", {
        method: "POST",
        body: JSON.stringify(form),
      });
      setResult(res.data);
      setForm({
        firstName: "",
        lastName: "",
        email: "",
        password: "",
        role: "employee",
      });
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen p-6">
      <div className="mx-auto max-w-2xl rounded-2xl border p-6 shadow-sm">
        <h1 className="text-2xl font-bold">Create Employee</h1>
        <p className="mt-2 text-sm text-gray-600">
          Password is optional. If you leave it empty, a temporary password will be generated and emailed.
        </p>

        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <input
            name="firstName"
            value={form.firstName}
            onChange={handleChange}
            placeholder="First Name"
            className="w-full rounded-xl border p-3"
          />
          <input
            name="lastName"
            value={form.lastName}
            onChange={handleChange}
            placeholder="Last Name"
            className="w-full rounded-xl border p-3"
          />
          <input
            name="email"
            value={form.email}
            onChange={handleChange}
            placeholder="Email"
            className="w-full rounded-xl border p-3"
          />
          <input
            name="password"
            value={form.password}
            onChange={handleChange}
            type="password"
            placeholder="Temporary Password (optional)"
            className="w-full rounded-xl border p-3"
          />
          <select name="role" value={form.role} onChange={handleChange} className="w-full rounded-xl border p-3">
            <option value="employee">Employee</option>
            <option value="manager">Manager</option>
          </select>

          {error ? <p className="text-sm text-red-600">{error}</p> : null}

          <button className="rounded-xl bg-black px-4 py-3 text-white" disabled={loading}>
            {loading ? "Creating..." : "Create Employee"}
          </button>
        </form>

        {result ? (
          <div className="mt-6 rounded-xl border p-4 text-sm space-y-1">
            <p><strong>Employee created:</strong> {result.user.firstName} {result.user.lastName}</p>
            <p><strong>Email:</strong> {result.user.email}</p>
            <p><strong>Employee ID:</strong> {result.profile.employeeId}</p>
            <p><strong>Invite email sent:</strong> Yes</p>
            <p><strong>Login link:</strong> {result.loginLink}</p>
          </div>
        ) : null}
      </div>
    </main>
  );
}