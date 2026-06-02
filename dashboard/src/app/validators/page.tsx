"use client";

import { useEffect, useState } from "react";
import { formatLALA } from "@/lib/chain";

interface Validator {
  operator_address: string;
  description: { moniker: string };
  status: string;
  tokens: string;
  commission: { commission_rates: { rate: string } };
  jailed: boolean;
}

export default function ValidatorsPage() {
  const [validators, setValidators] = useState<Validator[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const restUrl = process.env.NEXT_PUBLIC_REST_URL || "http://localhost:1317";
        const res = await fetch(`${restUrl}/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED`);
        const data = await res.json();
        setValidators(data.validators || []);
      } catch {
        // Node not available
      } finally {
        setLoading(false);
      }
    };
    load();
    const interval = setInterval(load, 10000);
    return () => clearInterval(interval);
  }, []);

  if (loading) return <div className="text-center py-20 text-gray-500">Loading validators...</div>;

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Validators</h1>
      <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-800 text-gray-400">
            <tr>
              <th className="px-4 py-3 text-left">#</th>
              <th className="px-4 py-3 text-left">Moniker</th>
              <th className="px-4 py-3 text-right">Voting Power</th>
              <th className="px-4 py-3 text-right">Commission</th>
              <th className="px-4 py-3 text-center">Status</th>
            </tr>
          </thead>
          <tbody>
            {validators.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-4 py-8 text-center text-gray-500">
                  No validators found. Make sure the chain is running.
                </td>
              </tr>
            ) : (
              validators.map((v, i) => (
                <tr key={v.operator_address} className="border-t border-gray-800 hover:bg-gray-800/50">
                  <td className="px-4 py-3">{i + 1}</td>
                  <td className="px-4 py-3 font-medium">{v.description.moniker}</td>
                  <td className="px-4 py-3 text-right font-mono">{formatLALA(v.tokens)} LALA</td>
                  <td className="px-4 py-3 text-right">
                    {(parseFloat(v.commission.commission_rates.rate) * 100).toFixed(1)}%
                  </td>
                  <td className="px-4 py-3 text-center">
                    {v.jailed ? (
                      <span className="text-red-400">Jailed</span>
                    ) : (
                      <span className="text-green-400">Active</span>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
