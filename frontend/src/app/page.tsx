import Link from "next/link";

export default function HomePage() {
  return (
    <main className="min-h-screen flex items-center justify-center p-6">
      <div className="w-full max-w-xl rounded-2xl border p-8 shadow-sm">
        <h1 className="text-3xl font-bold">Employee Platform</h1>
        <p className="mt-3 text-sm text-gray-600">
          Admin can sign up and create employees. Employees can only log in..
        </p>
        <div className="mt-6 flex gap-3">
          <Link className="rounded-xl bg-black px-4 py-2 text-white" href="/login">
            Login
          </Link>
          <Link className="rounded-xl border px-4 py-2" href="/admin-signup">
            Admin Signup
          </Link>
        </div>
      </div>
    </main>
  );
}