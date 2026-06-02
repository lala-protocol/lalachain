import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "LalaChain Dashboard",
  description: "Monitor LalaChain testnet: validators, KPIs, governance proposals, and AI advisor state",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-gray-950 text-gray-100">
        <nav className="border-b border-gray-800 bg-gray-900">
          <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <span className="text-xl font-bold text-purple-400">LalaChain</span>
              <span className="text-sm text-gray-500">Testnet Dashboard</span>
            </div>
            <div className="flex gap-6 text-sm">
              <a href="/" className="hover:text-purple-400 transition">Dashboard</a>
              <a href="/validators" className="hover:text-purple-400 transition">Validators</a>
              <a href="/kpis" className="hover:text-purple-400 transition">KPIs</a>
              <a href="/proposals" className="hover:text-purple-400 transition">Proposals</a>
              <a href="/parameters" className="hover:text-purple-400 transition">Parameters</a>
            </div>
          </div>
        </nav>
        <main className="max-w-7xl mx-auto px-4 py-8">
          {children}
        </main>
      </body>
    </html>
  );
}
