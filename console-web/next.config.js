/** @type {import('next').NextConfig} */
const nextConfig = {
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
