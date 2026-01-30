import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Improve dev server stability in containerized environments
  devIndicators: false,

  // Use experimental webpackMemoryOptimizations to reduce memory pressure
  experimental: {
    webpackMemoryOptimizations: true,
  },

  // Allow images from MinIO storage
  images: {
    remotePatterns: [
      {
        protocol: "http",
        hostname: "localhost",
        port: "9000",
        pathname: "/credfolio/**",
      },
    ],
  },

  // Proxy API requests through Next.js to avoid port forwarding issues in devcontainer
  async rewrites() {
    return [
      {
        source: "/api/graphql",
        destination: "http://localhost:8080/graphql",
      },
    ];
  },
};

export default nextConfig;
