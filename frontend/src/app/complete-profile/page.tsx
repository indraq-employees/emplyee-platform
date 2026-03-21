"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";

export default function CompleteProfilePage() {
  const router = useRouter();
  const [form, setForm] = useState({
    nickName: "",
    department: "",
    location: "",
    designation: "",
    jobRole: "",
    employmentType: "",
    employeeStatus: "Active",
    sourceOfHire: "Web",
    dateOfJoining: "",
    dateOfBirth: "",
    maritalStatus: "",
    aboutMe: "",
    expertise: "",
    uan: "",
    pan: "",
    workPhoneNumber: "",
    personalMobileNumber: "",
    extension: "",
    personalEmailAddress: "",
    seatingLocation: "",
    presentAddress: "",
    permanentAddress: "",
    tags: [] as string[],
  });
  const [tagsInput, setTagsInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await apiFetch("/profile/me", {
        method: "PUT",
        body: JSON.stringify({
          ...form,
          tags: tagsInput
            .split(",")
            .map((tag) => tag.trim())
            .filter(Boolean),
        }),
      });
      router.push("/dashboard");
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen p-6">
      <div className="mx-auto max-w-4xl rounded-2xl border p-6 shadow-sm">
        <h1 className="text-2xl font-bold">Complete Profile</h1>
        <form onSubmit={handleSubmit} className="mt-6 grid gap-4 md:grid-cols-2">
          <input name="nickName" placeholder="Nick Name" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="department" placeholder="Department" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="location" placeholder="Location" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="designation" placeholder="Designation" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="jobRole" placeholder="Job Role" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="employmentType" placeholder="Employment Type" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="employeeStatus" placeholder="Employee Status" className="rounded-xl border p-3" onChange={handleChange} value={form.employeeStatus} />
          <input name="sourceOfHire" placeholder="Source Of Hire" className="rounded-xl border p-3" onChange={handleChange} value={form.sourceOfHire} />
          <input name="dateOfJoining" type="date" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="dateOfBirth" type="date" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="maritalStatus" placeholder="Marital Status" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="expertise" placeholder="Expertise" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="uan" placeholder="UAN" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="pan" placeholder="PAN" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="workPhoneNumber" placeholder="Work Phone Number" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="personalMobileNumber" placeholder="Personal Mobile Number" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="extension" placeholder="Extension" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="personalEmailAddress" placeholder="Personal Email Address" className="rounded-xl border p-3" onChange={handleChange} />
          <input name="seatingLocation" placeholder="Seating Location" className="rounded-xl border p-3" onChange={handleChange} />
          <input placeholder="Tags comma separated" className="rounded-xl border p-3" value={tagsInput} onChange={(e) => setTagsInput(e.target.value)} />
          <textarea name="aboutMe" placeholder="About Me" className="rounded-xl border p-3 md:col-span-2" onChange={handleChange} rows={4} />
          <textarea name="presentAddress" placeholder="Present Address" className="rounded-xl border p-3 md:col-span-2" onChange={handleChange} rows={3} />
          <textarea name="permanentAddress" placeholder="Permanent Address" className="rounded-xl border p-3 md:col-span-2" onChange={handleChange} rows={3} />
          {error ? <p className="text-sm text-red-600 md:col-span-2">{error}</p> : null}
          <button disabled={loading} className="rounded-xl bg-black p-3 text-white md:col-span-2">
            {loading ? "Saving..." : "Save Profile"}
          </button>
        </form>
      </div>
    </main>
  );
}