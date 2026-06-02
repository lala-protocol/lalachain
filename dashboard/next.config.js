/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    NEXT_PUBLIC_RPC_URL: process.env.NEXT_PUBLIC_RPC_URL || "http://localhost:26657",
    NEXT_PUBLIC_REST_URL: process.env.NEXT_PUBLIC_REST_URL || "http://localhost:1317",
    NEXT_PUBLIC_CHAIN_ID: process.env.NEXT_PUBLIC_CHAIN_ID || "lalachain-testnet-1",
    NEXT_PUBLIC_CHAIN_NAME: "LalaChain Testnet",
    NEXT_PUBLIC_DENOM: "ulala",
    NEXT_PUBLIC_DISPLAY_DENOM: "LALA",
    NEXT_PUBLIC_DENOM_EXPONENT: "6",
  },
};

module.exports = nextConfig;
