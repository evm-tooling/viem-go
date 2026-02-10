import type { MetadataRoute } from "next";
import { siteConfig } from "@/lib/seo";

/**
 * Generates robots.txt for the viem-go documentation site.
 * This controls how search engine crawlers access and index the site.
 *
 * @see https://nextjs.org/docs/app/api-reference/file-conventions/metadata/robots
 */
export default function robots(): MetadataRoute.Robots {
  const baseUrl = siteConfig.url;

  return {
    rules: [
      {
        userAgent: "*",
        allow: "/",
        disallow: [
          "/api/", // Don't index API routes
          "/_next/", // Don't index Next.js internal routes
          "/private/", // Don't index any private content
        ],
      },
      {
        // Specific rules for Googlebot
        userAgent: "Googlebot",
        allow: "/",
        disallow: ["/api/"],
      },
    ],
    sitemap: `${baseUrl}/sitemap.xml`,
    host: baseUrl,
  };
}
