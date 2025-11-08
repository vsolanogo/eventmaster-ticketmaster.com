import React from "react";

export const Footer = () => {
  return (
    <footer className="bg-gray-100 py-6">
      <div className="container mx-auto px-5 text-center text-sm text-gray-500">
        © {new Date().getFullYear()} EventsMaster — <a href="https://github.com/vsolanogo/" target="_blank" rel="noopener noreferrer" className="text-gray-600 hover:text-gray-900">@Vitalii</a>
      </div>
    </footer>
  );
};