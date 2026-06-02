import { StargateClient } from "@cosmjs/stargate";

const RPC_URL = process.env.NEXT_PUBLIC_RPC_URL || "http://localhost:26657";
const REST_URL = process.env.NEXT_PUBLIC_REST_URL || "http://localhost:1317";

export async function getStargateClient(): Promise<StargateClient> {
  return StargateClient.connect(RPC_URL);
}

export async function fetchREST<T>(path: string): Promise<T> {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 10000);
  try {
    const res = await fetch(`${REST_URL}${path}`, {
      next: { revalidate: 5 },
      signal: controller.signal,
    });
    if (!res.ok) throw new Error(`REST error: ${res.status}`);
    return res.json();
  } finally {
    clearTimeout(timeout);
  }
}

export function formatLALA(amount: string | number): string {
  const num = typeof amount === "string" ? parseInt(amount, 10) : amount;
  return (num / 1_000_000).toLocaleString(undefined, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 6,
  });
}

export const CHAIN_CONFIG = {
  chainId: process.env.NEXT_PUBLIC_CHAIN_ID || "lalachain-testnet-1",
  chainName: "LalaChain Testnet",
  rpc: RPC_URL,
  rest: REST_URL,
  bip44: { coinType: 118 },
  bech32Config: {
    bech32PrefixAccAddr: "lala",
    bech32PrefixAccPub: "lalapub",
    bech32PrefixValAddr: "lalavaloper",
    bech32PrefixValPub: "lalavaloperpub",
    bech32PrefixConsAddr: "lalavalcons",
    bech32PrefixConsPub: "lalavalconspub",
  },
  currencies: [
    { coinDenom: "LALA", coinMinimalDenom: "ulala", coinDecimals: 6 },
  ],
  feeCurrencies: [
    { coinDenom: "LALA", coinMinimalDenom: "ulala", coinDecimals: 6, gasPriceStep: { low: 0, average: 0.025, high: 0.04 } },
  ],
  stakeCurrency: { coinDenom: "LALA", coinMinimalDenom: "ulala", coinDecimals: 6 },
};
