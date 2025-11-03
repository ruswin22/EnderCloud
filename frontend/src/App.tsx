import React, { useState } from "react";
import { motion } from "framer-motion";
import { type ServerRequest } from "./types/server";
import { ServerForm } from "./components/ServerForm";
import { StatusBox } from "./components/StatusBox";

export default function App() {
  const [status, setStatus] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  const handleServerCreate = async (data: ServerRequest) => {
    setLoading(true);
    setStatus("");

    try {
      const response = await fetch("http://localhost:8080/api/create-server", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });

      const text = await response.text();
      if (response.ok) {
        setStatus(`Success: ${text}`);
      } else {
        setStatus(`Error: ${text}`);
      }
    } catch (err) {
      setStatus("Server Error: Could not connect to backend.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 flex flex-col items-center justify-center text-white px-4">
      <motion.div
        className="bg-gray-800 p-8 rounded-2xl shadow-lg max-w-md w-full"
        initial={{ opacity: 0, y: -30 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <h1 className="text-3xl font-bold mb-6 text-center text-green-400">
          EnderCloud Server Creator
        </h1>

        <ServerForm onSubmit={handleServerCreate} loading={loading} />
        <StatusBox message={status} />
      </motion.div>

      <p className="mt-6 text-gray-400 text-sm">
        Powered by • React • Go • AWS • Terraform • Ansible
      </p>
    </div>
  );
}
