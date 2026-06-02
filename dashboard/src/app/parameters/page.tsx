"use client";

import { useEffect, useState } from "react";

interface MintParams {
  mint_denom: string;
  inflation_rate_change: string;
  inflation_max: string;
  inflation_min: string;
  goal_bonded: string;
  blocks_per_year: string;
}

interface StakingParams {
  unbonding_time: string;
  max_validators: number;
  bond_denom: string;
  min_commission_rate: string;
}

export default function ParametersPage() {
  const [mintParams, setMintParams] = useState<MintParams | null>(null);
  const [stakingParams, setStakingParams] = useState<StakingParams | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const restUrl = process.env.NEXT_PUBLIC_REST_URL || "http://localhost:1317";
        const [mintRes, stakingRes] = await Promise.allSettled([
          fetch(`${restUrl}/cosmos/mint/v1beta1/params`),
          fetch(`${restUrl}/cosmos/staking/v1beta1/params`),
        ]);
        if (mintRes.status === "fulfilled") {
          const d = await mintRes.value.json();
          setMintParams(d.params);
        }
        if (stakingRes.status === "fulfilled") {
          const d = await stakingRes.value.json();
          setStakingParams(d.params);
        }
      } catch {
        // ignore
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  if (loading) return <div className="text-center py-20 text-gray-500">Loading parameters...</div>;

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Chain Parameters</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Mint/Inflation */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
          <h2 className="text-lg font-semibold mb-4 text-purple-400">Mint / Inflation</h2>
          {mintParams ? (
            <div className="space-y-3 text-sm">
              <ParamRow label="Mint Denom" value={mintParams.mint_denom} />
              <ParamRow label="Inflation Rate Change" value={`${(parseFloat(mintParams.inflation_rate_change) * 100).toFixed(1)}%`} />
              <ParamRow label="Inflation Max" value={`${(parseFloat(mintParams.inflation_max) * 100).toFixed(1)}%`} />
              <ParamRow label="Inflation Min" value={`${(parseFloat(mintParams.inflation_min) * 100).toFixed(1)}%`} />
              <ParamRow label="Goal Bonded" value={`${(parseFloat(mintParams.goal_bonded) * 100).toFixed(1)}%`} />
              <ParamRow label="Blocks/Year" value={parseInt(mintParams.blocks_per_year).toLocaleString()} />
            </div>
          ) : (
            <p className="text-gray-500">Not available</p>
          )}
        </div>

        {/* Staking */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
          <h2 className="text-lg font-semibold mb-4 text-purple-400">Staking</h2>
          {stakingParams ? (
            <div className="space-y-3 text-sm">
              <ParamRow label="Bond Denom" value={stakingParams.bond_denom} />
              <ParamRow label="Unbonding Time" value={`${parseInt(stakingParams.unbonding_time) / 86400000000000} days`} />
              <ParamRow label="Max Validators" value={stakingParams.max_validators.toString()} />
              <ParamRow label="Min Commission" value={`${(parseFloat(stakingParams.min_commission_rate) * 100).toFixed(1)}%`} />
            </div>
          ) : (
            <p className="text-gray-500">Not available</p>
          )}
        </div>

        {/* Custom: AI Advisor */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
          <h2 className="text-lg font-semibold mb-4 text-purple-400">AI Advisor Config</h2>
          <div className="space-y-3 text-sm">
            <ParamRow label="Low Util Threshold" value="40%" />
            <ParamRow label="High Util Threshold" value="80%" />
            <ParamRow label="Min Block Gas Limit" value="10,000,000" />
            <ParamRow label="Max Block Gas Limit" value="30,000,000" />
            <ParamRow label="Streak Required (Low)" value="3 epochs" />
            <ParamRow label="Streak Required (High)" value="2 epochs" />
          </div>
        </div>

        {/* Custom: Governance */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
          <h2 className="text-lg font-semibold mb-4 text-purple-400">LalaGov Config</h2>
          <div className="space-y-3 text-sm">
            <ParamRow label="Quorum" value="66%" />
            <ParamRow label="Approval Threshold" value="51%" />
            <ParamRow label="Voting Period" value="1 epoch" />
            <ParamRow label="Max Change/Proposal" value="5%" />
          </div>
        </div>
      </div>
    </div>
  );
}

function ParamRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between items-center py-1 border-b border-gray-800/50">
      <span className="text-gray-400">{label}</span>
      <span className="font-mono text-white">{value}</span>
    </div>
  );
}
