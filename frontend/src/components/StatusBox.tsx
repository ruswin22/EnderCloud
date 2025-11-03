import React from "react";

interface Props {
  message: string;
}

export const StatusBox: React.FC<Props> = ({ message }) => {
  if (!message) return null;

  return (
    <div className="mt-4 text-sm p-3 bg-gray-700 rounded-lg text-center">
      {message}
    </div>
  );
};
