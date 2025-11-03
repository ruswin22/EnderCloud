import React, { useState } from "react";
import { ServerForm } from "./components/ServerForm";
import { StatusBox } from "./components/StatusBox";
import { type ServerRequest } from "./types/server";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

export default function App() {
  const [status, setStatus] = useState("");
  const [loading, setLoading] = useState(false);
  const [serverIP, setServerIP] = useState<string | null>(null);

  const handleCreate = async (data: ServerRequest) => {
    try {
      setLoading(true);
      setStatus("🚀 Creating Minecraft server... please wait (may take a few minutes)");

      const res = await fetch(`${BACKEND_URL}/api/create-server`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });

      const result = await res.json();

      if (res.ok) {
        setServerIP(result.public_ip);
        setStatus(`✅ Server ready at ${result.public_ip}`);
      } else {
        setStatus(`❌ Failed: ${result.error || "Unknown error"}`);
      }
    } catch (err) {
      console.error(err);
      setStatus("❌ Could not connect to backend");
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    try {
      setLoading(true);
      setStatus("🧹 Destroying server... please wait");

      const res = await fetch(`${BACKEND_URL}/api/delete-server`, {
        method: "DELETE",
      });

      const result = await res.json();

      if (res.ok) {
        setServerIP(null);
        setStatus("🧨 Server destroyed successfully!");
      } else {
        setStatus(`❌ Delete failed: ${result.error}`);
      }
    } catch (err) {
      console.error(err);
      setStatus("❌ Could not reach backend");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white flex items-center justify-center">
      <div className="bg-gray-800 p-8 rounded-2xl shadow-2xl w-full max-w-md">
        <h1 className="text-3xl font-bold text-center mb-6">🟢 EnderCloud</h1>

        <ServerForm onSubmit={handleCreate} loading={loading} />

        {serverIP && (
          <div className="mt-4 text-center">
            <p className="text-green-400 text-sm mb-2">
              Server IP: <span className="font-mono">{serverIP}:25565</span>
            </p>
            <button
              onClick={handleDelete}
              disabled={loading}
              className="w-full py-2 bg-red-600 hover:bg-red-700 rounded font-semibold transition disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-red-400"
            >
              {loading ? "Deleting..." : "Delete Server"}
            </button>
          </div>
        )}

        <StatusBox message={status} />
      </div>
    </div>
  );
}
