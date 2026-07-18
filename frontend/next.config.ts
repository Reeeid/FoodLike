import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // standalone は Docker等の自前ホスト用(Dockerfileが .next/standalone を使う)。
  // Vercel はルーティングを自前で組むため standalone だと全ルート404になる。
  // Vercelビルド時は VERCEL=1 が立つので、そのときだけ無効化する。
  output: process.env.VERCEL ? undefined : "standalone",
};

export default nextConfig;
