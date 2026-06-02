"use client";

import { useEffect, useState } from "react";
import { formatLALA } from "@/lib/chain";

interface NetworkStats {
  latestBlock: number;
  blockTime: string;
  totalValidators: number;
  bondedTokens: string;
  inflation: string;
  communityPool: string;
}

async function fetchNetworkStats(): Promise<NetworkStats> {
  const restUrl = process.env.NEXT_PUBLIC_REST_URL || "http://localhost:1317";

  try {
    const [blockRes, stakingRes, mintRes, poolRes] = await Promise.allSettled([
      fetch(`${restUrl}/cosmos/base/tendermint/v1beta1/blocks/latest`),
      fetch(`${restUrl}/cosmos/staking/v1beta1/pool`),
      fetch(`${restUrl}/cosmos/mint/v1beta1/inflation`),
      fetch(`${restUrl}/cosmos/distribution/v1beta1/community_pool`),
    ]);

    const block = blockRes.status === "fulfilled" ? await blockRes.value.json() : null;
    const staking = stakingRes.status === "fulfilled" ? await stakingRes.value.json() : null;
    const minting = mintRes.status === "fulfilled" ? await mintRes.value.json() : null;
    const pool = poolRes.status === "fulfilled" ? await poolRes.value.json() : null;

    return {
      latestBlock: block?.block?.header?.height ? parseInt(block.block.header.height) : 0,
      blockTime: block?.block?.header?.time || "N/A",
      totalValidators: staking?.pool ? 1 : 0,
      bondedTokens: staking?.pool?.bonded_tokens || "0",
      inflation: minting?.inflation || "0",
      communityPool: pool?.pool?.[0]?.amount || "0",
    };
  } catch {
    return {
      latestBlock: 0,
      blockTime: "N/A",
      totalValidators: 0,
      bondedTokens: "0",
      inflation: "0",
      communityPool: "0",
    };
  }
}

function StatCard({ label, value, unit }: { label: string; value: string; unit?: string }) {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
      <p className="text-sm text-gray-400 mb-1">{label}</p>
      <p className="text-2xl font-bold text-white">
        {value}
        {unit && <span className="text-sm text-gray-500 ml-1">{unit}</span>}
      </p>
    </div>
  );
}

export default function DashboardPage() {
  const [stats, setStats] = useState<NetworkStats | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const load = async () => {
      try {
        const data = await fetchNetworkStats();
        setStats(data);
      } catch (e: unknown) {
        setError(e instanceof Error ? e.message : "Failed to fetch");
      }
    };
    load();
    const interval = setInterval(load, 5000);
    return () => clearInterval(interval);
  }, []);

  if (error) {
    return (
      <div className="text-center py-20">
        <h1 className="text-2xl font-bold mb-4">LalaChain Dashboard</h1>
        <p className="text-yellow-400">Waiting for node connection...</p>
        <p className="text-sm text-gray-500 mt-2">{error}</p>
        <p className="text-sm text-gray-600 mt-4">
          Make sure lalachaind is running with API enabled at{" "}
          {process.env.NEXT_PUBLIC_REST_URL || "http://localhost:1317"}
        </p>
      </div>
    );
  }

  if (!stats) {
    return <div className="text-center py-20 text-gray-500">Loading...</div>;
  }

  const inflationPct = (parseFloat(stats.inflation) * 100).toFixed(2);

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Network Overview</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <StatCard label="Latest Block" value={stats.latestBlock.toLocaleString()} />
        <StatCard label="Bonded Tokens" value={formatLALA(stats.bondedTokens)} unit="LALA" />
        <StatCard label="Inflation" value={`${inflationPct}%`} />
        <StatCard label="Community Pool" value={formatLALA(stats.communityPool)} unit="LALA" />
      </div>

      <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
        <h2 className="text-lg font-semibold mb-4">Chain Info</h2>
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-400">Chain ID:</span>{" "}
            <span className="font-mono">{process.env.NEXT_PUBLIC_CHAIN_ID}</span>
          </div>
          <div>
            <span className="text-gray-400">Denom:</span>{" "}
            <span className="font-mono">ulala (LALA)</span>
          </div>
          <div>
            <span className="text-gray-400">Block Time:</span>{" "}
            <span className="font-mono">~5s</span>
          </div>
          <div>
            <span className="text-gray-400">Last Block Time:</span>{" "}
            <span className="font-mono text-xs">{stats.blockTime}</span>
          </div>
        </div>
      </div>
    </div>
  );
}
