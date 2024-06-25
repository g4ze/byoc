/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
    webpack: (config, { isServer }) => {
      if (!isServer) {
        config.resolve.fallback = {
          ...config.resolve.fallback,
          pg: false,
          'pg-native': false,
        };
      }
      return config;
    },
};

export default nextConfig;
