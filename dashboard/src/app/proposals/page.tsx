"use client";

import { useEffect, useState } from "react";

interface Proposal {
  proposal_id: string;
  parameter: string;
  current_value: string;
  proposed_value: string;
  rationale: string;
  votes_approve: string;
  votes_reject: string;
  outcome: string;
}

export default function ProposalsPage() {
  const [proposals, setProposals] = useState<Proposal[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const restUrl = process.env.NEXT_PUBLIC_REST_URL || "http://localhost:1317";
        const res = await fetch(`${restUrl}/lala/lalagov/v1/history`);
        const data = await res.json();
        setProposals(data.proposals || []);
      } catch {
        setProposals([]);
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  if (loading) return <div className="text-center py-20 text-gray-500">Loading proposals...</div>;

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Governance Proposals</h1>
      <p className="text-gray-400 mb-6">
        AI-originated parameter change proposals voted on by validators.
      </p>

      {proposals.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400 mb-2">No proposals yet.</p>
          <p className="text-sm text-gray-600">
            The AI advisor generates proposals based on epoch telemetry analysis.
          </p>
        </div>
      ) : (
        <div className="space-y-4">
          {proposals.map((p) => (
            <div
              key={p.proposal_id}
              className="bg-gray-900 border border-gray-800 rounded-xl p-6"
            >
              <div className="flex items-center justify-between mb-3">
                <span className="font-mono text-sm text-gray-400">
                  Proposal #{p.proposal_id}
                </span>
                <span
                  className={`px-3 py-1 rounded-full text-xs font-medium ${
                    p.outcome === "passed"
                      ? "bg-green-900/50 text-green-400"
                      : "bg-red-900/50 text-red-400"
                  }`}
                >
                  {p.outcome.toUpperCase()}
                </span>
              </div>
              <h3 className="text-lg font-semibold mb-2">
                Change <span className="text-purple-400">{p.parameter}</span>
              </h3>
              <p className="text-sm text-gray-400 mb-3">{p.rationale}</p>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                <div>
                  <span className="text-gray-500">Current:</span>{" "}
                  <span className="font-mono">{p.current_value}</span>
                </div>
                <div>
                  <span className="text-gray-500">Proposed:</span>{" "}
                  <span className="font-mono">{p.proposed_value}</span>
                </div>
                <div>
                  <span className="text-gray-500">Approve:</span>{" "}
                  <span className="font-mono text-green-400">{p.votes_approve}</span>
                </div>
                <div>
                  <span className="text-gray-500">Reject:</span>{" "}
                  <span className="font-mono text-red-400">{p.votes_reject}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
