/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    unoptimized: true,
  },
  async redirects() {
    return [
      {
        source: "/",
        destination: "/c",
        permanent: true,
      }
    ]
  }
}

module.exports = nextConfig
