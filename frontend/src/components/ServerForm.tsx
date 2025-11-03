import React, { useState } from "react";
import { type ServerRequest } from "../types/server";

interface Props {
  onSubmit: (data: ServerRequest) => Promise<void>;
  loading: boolean;
}

export const ServerForm: React.FC<Props> = ({ onSubmit, loading }) => {
  const [formData, setFormData] = useState<ServerRequest>({
    serverName: "",
    seed: "",
    version: "1.21.11",
    ram: "1G",
  });

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="serverName" className="block mb-1 after:content-['*'] after:text-red-500 after:ml-1">
          Server Name
        </label>
        <input
          type="text"
          id="serverName"
          name="serverName"
          value={formData.serverName}
          onChange={handleChange}
          className="w-full p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:ring-2 focus:ring-green-500"
          required
        />
      </div>

      <div>
        <label htmlFor="seed" className="block mb-1">
          Seed <span className="text-gray-400 text-sm">(optional)</span>
        </label>
        <input
          type="text"
          id="seed"
          name="seed"
          value={formData.seed}
          onChange={handleChange}
          placeholder="Leave empty for random"
          className="w-full p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:ring-2 focus:ring-green-500"
        />
      </div>

      <div>
        <label htmlFor="version" className="block mb-1 after:content-['*'] after:text-red-500 after:ml-1">
          Version
        </label>
        <select
          id="version"
          name="version"
          value={formData.version}
          onChange={handleChange}
          className="w-full p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:ring-2 focus:ring-green-500"
          required
        >
          <option value="1.21.11">1.21.11</option>
          <option value="1.21.10">1.21.10</option>
          <option value="1.21.9">1.21.9</option>
        </select>
      </div>

      <div>
        <label htmlFor="ram" className="block mb-1 after:content-['*'] after:text-red-500 after:ml-1">
          RAM
        </label>
        <select
          id="ram"
          name="ram"
          value={formData.ram}
          onChange={handleChange}
          className="w-full p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:ring-2 focus:ring-green-500"
          required
        >
          <option value="1G">1 GB</option>
          <option value="2G">2 GB</option>
          <option value="4G">4 GB</option>
        </select>
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full py-2 bg-green-500 hover:bg-green-600 rounded font-semibold transition disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-green-400"
      >
        {loading ? "Creating..." : "Create Server"}
      </button>
    </form>
  );
};