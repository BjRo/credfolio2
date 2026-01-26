import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Improve dev server stability in containerized environments
  devIndicators: false,

  // Use experimental webpackMemoryOptimizations to reduce memory pressure
  experimental: {
    webpackMemoryOptimizations: true,
  },
};

export default nextConfig;
