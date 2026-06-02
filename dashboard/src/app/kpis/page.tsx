"use client";

import { useEffect, useState } from "react";

interface KPIData {
  epoch: number;
  avg_block_time_ms: number;
  avg_utilization: number;
  avg_base_fee: number;
  validator_count: number;
  total_staked_ratio: number;
  slashing_events: number;
}

export default function KPIsPage() {
  const [kpis, setKpis] = useState<KPIData[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const restUrl = process.env.NEXT_PUBLIC_REST_URL || "http://localhost:1317";
        const res = await fetch(`${restUrl}/lala/telemetry/v1/kpis`);
        const data = await res.json();
        setKpis(data.kpis || []);
      } catch {
        // Custom query not available yet - show mock
        setKpis([]);
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  if (loading) return <div className="text-center py-20 text-gray-500">Loading KPIs...</div>;

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Epoch KPIs</h1>

      {kpis.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400 mb-2">No KPI data available yet.</p>
          <p className="text-sm text-gray-600">
            KPIs are computed once per epoch (~100 blocks). Wait for the chain to produce data.
          </p>
        </div>
      ) : (
        <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-800 text-gray-400">
              <tr>
                <th className="px-4 py-3 text-left">Epoch</th>
                <th className="px-4 py-3 text-right">Avg Block Time</th>
                <th className="px-4 py-3 text-right">Utilization</th>
                <th className="px-4 py-3 text-right">Avg Base Fee</th>
                <th className="px-4 py-3 text-right">Validators</th>
                <th className="px-4 py-3 text-right">Staked Ratio</th>
                <th className="px-4 py-3 text-right">Slashing Events</th>
              </tr>
            </thead>
            <tbody>
              {kpis.map((kpi) => (
                <tr key={kpi.epoch} className="border-t border-gray-800 hover:bg-gray-800/50">
                  <td className="px-4 py-3 font-mono">{kpi.epoch}</td>
                  <td className="px-4 py-3 text-right font-mono">{kpi.avg_block_time_ms.toFixed(0)}ms</td>
                  <td className="px-4 py-3 text-right font-mono">{(kpi.avg_utilization * 100).toFixed(1)}%</td>
                  <td className="px-4 py-3 text-right font-mono">{(kpi.avg_base_fee / 1e6).toFixed(2)} ulala</td>
                  <td className="px-4 py-3 text-right">{kpi.validator_count}</td>
                  <td className="px-4 py-3 text-right font-mono">{(kpi.total_staked_ratio * 100).toFixed(1)}%</td>
                  <td className="px-4 py-3 text-right">{kpi.slashing_events}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
          <h3 className="text-sm text-gray-400 mb-2">Epoch Duration</h3>
          <p className="text-lg font-bold">100 blocks (~8.3 min)</p>
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
          <h3 className="text-sm text-gray-400 mb-2">Target Block Time</h3>
          <p className="text-lg font-bold">5,000 ms</p>
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
          <h3 className="text-sm text-gray-400 mb-2">Gas Limit</h3>
          <p className="text-lg font-bold">15,000,000</p>
        </div>
      </div>
    </div>
  );
}
