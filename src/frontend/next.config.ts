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
      // Proxy MinIO storage requests - needed because browser can't reach MinIO directly in devcontainer
      {
        source: "/storage/:path*",
        destination: `${process.env.MINIO_INTERNAL_URL || "http://credfolio2-minio:9000"}/credfolio/:path*`,
      },
    ];
  },
};

export default nextConfig;
