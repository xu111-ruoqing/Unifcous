/** @type {import('next').NextConfig} */
const isDev = process.env.NODE_ENV === 'development';
const devPort = process.env.PORT || '3000';

const nextConfig = {
  reactStrictMode: true,
  distDir: isDev ? `.next-dev-${devPort}` : '.next',
  env: {
    API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:8080',
  },
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.API_BASE_URL || 'http://localhost:8080'}/api/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;

