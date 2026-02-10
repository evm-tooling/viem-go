import type { MetadataRoute } from "next";
import { getAllDocSlugs, getDocLastModified } from "@/lib/mdx";
import { siteConfig } from "@/lib/seo";

/**
 * Generates a dynamic sitemap for the viem-go documentation site.
 * This helps search engines discover and index all pages.
 *
 * @see https://nextjs.org/docs/app/api-reference/file-conventions/metadata/sitemap
 */
export default function sitemap(): MetadataRoute.Sitemap {
  const baseUrl = siteConfig.url;

  // Static pages
  const staticPages: MetadataRoute.Sitemap = [
    {
      url: baseUrl,
      lastModified: new Date(),
      changeFrequency: "weekly",
      priority: 1.0,
    },
  ];

  // Documentation pages - dynamically generated from MDX files
  const docSlugs = getAllDocSlugs();
  const docPages: MetadataRoute.Sitemap = docSlugs.map((slug) => {
    const lastModified = getDocLastModified(slug);

    // Determine priority based on page depth and type
    let priority = 0.8;
    const depth = slug.split("/").length;

    // Introduction pages get higher priority
    if (slug.endsWith("introduction") || slug === "getting-started") {
      priority = 0.9;
    } else if (depth === 1) {
      priority = 0.85;
    } else if (depth === 2) {
      priority = 0.75;
    } else {
      priority = 0.7;
    }

    return {
      url: `${baseUrl}/docs/${slug}`,
      lastModified: lastModified || new Date(),
      changeFrequency: "weekly" as const,
      priority,
    };
  });

  return [...staticPages, ...docPages];
}
