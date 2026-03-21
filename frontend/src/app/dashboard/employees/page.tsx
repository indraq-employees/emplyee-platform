"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";

export default function EmployeesPage() {
  const [employees, setEmployees] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiFetch<any>("/employees")
      .then((res) => setEmployees(res.data || []))
      .finally(() => setLoading(false));
  }, []);

  return (
    <main className="min-h-screen p-6">
      <div className="mx-auto max-w-5xl rounded-2xl border p-6">
        <h1 className="text-2xl font-bold">Employees</h1>
        {loading ? (
          <p className="mt-4">Loading...</p>
        ) : (
          <div className="mt-6 overflow-x-auto">
            <table className="w-full border-collapse text-sm">
              <thead>
                <tr className="border-b text-left">
                  <th className="p-3">Name</th>
                  <th className="p-3">Email</th>
                  <th className="p-3">Role</th>
                  <th className="p-3">Profile Completed</th>
                </tr>
              </thead>
              <tbody>
                {employees.map((employee) => (
                  <tr key={employee.id} className="border-b">
                    <td className="p-3">{employee.firstName} {employee.lastName}</td>
                    <td className="p-3">{employee.email}</td>
                    <td className="p-3">{employee.role}</td>
                    <td className="p-3">{employee.profileCompleted ? "Yes" : "No"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </main>
  );
}