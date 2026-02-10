import type { Metadata } from "next";

/**
 * Site-wide SEO configuration for viem-go documentation
 */

export const siteConfig = {
  name: "viem-go",
  title: "viem-go · Go Interface for Ethereum",
  description:
    "Build reliable blockchain apps & libraries with idiomatic Go, type-safe and composable modules that interface with Ethereum — inspired by viem",
  url: "https://www.viemgolem.sh",
  ogImage: "/og-image.png", // Default OG image
  creator: "@viemgo", // Update to your Twitter handle
  keywords: [
    "Go",
    "Golem",
    "viem-go",
    "Golang",
    "Ethereum",
    "blockchain",
    "viem",
    "web3",
    "smart contracts",
    "EVM",
    "cryptocurrency",
    "decentralized",
    "RPC",
    "JSON-RPC",
    "Ethereum client",
    "Go Ethereum library",
    "type-safe",
    "composable",
  ] as string[],
  authors: [
    {
      name: "viem-go",
      url: "https://www.viemgolem.sh",
    },
  ] as Array<{ name: string; url: string }>,
  links: {
    github: "https://github.com/evm-tooling/viem-go", // Update to your repo
  },
};

/**
 * Base metadata shared across all pages
 */
export function createBaseMetadata(): Metadata {
  return {
    metadataBase: new URL(siteConfig.url),
    title: {
      default: siteConfig.title,
      template: `%s - ${siteConfig.name}`,
    },
    description: siteConfig.description,
    keywords: siteConfig.keywords,
    authors: siteConfig.authors,
    creator: siteConfig.creator,
    publisher: siteConfig.name,
    formatDetection: {
      email: true,
      address: false,
      telephone: false,
    },
    openGraph: {
      type: "website",
      locale: "en_US",
      url: siteConfig.url,
      title: siteConfig.title,
      description: siteConfig.description,
      siteName: siteConfig.name,
      images: [
        {
          url: siteConfig.ogImage,
          width: 1200,
          height: 630,
          alt: `${siteConfig.name} - Go Interface for Ethereum`,
        },
      ],
    },
    twitter: {
      card: "summary_large_image",
      title: siteConfig.title,
      description: siteConfig.description,
      images: [siteConfig.ogImage],
      creator: siteConfig.creator,
    },
    robots: {
      index: true,
      follow: true,
      nocache: false,
      googleBot: {
        index: true,
        follow: true,
        noimageindex: false,
        "max-video-preview": -1,
        "max-image-preview": "large",
        "max-snippet": -1,
      },
    },
    icons: {
      icon: [
        { url: "/favicons/golem-icon-only-light.svg", media: "(prefers-color-scheme: light)" },
        { url: "/favicons/golem-icon-only-dark.svg", media: "(prefers-color-scheme: dark)" },
      ],
      apple: "/favicons/golem-icon-only-dark.svg",
    },
    manifest: "/manifest.json",
    alternates: {
      canonical: siteConfig.url,
    },
    category: "technology",
  };
}

/**
 * Generate metadata for documentation pages
 */
export function createDocsMetadata({
  title,
  description,
  slug,
}: {
  title: string;
  description?: string;
  slug: string;
}): Metadata {
  const pageTitle = title;
  const pageDescription = description || siteConfig.description;
  const pageUrl = `${siteConfig.url}/docs/${slug}`;

  return {
    title: pageTitle,
    description: pageDescription,
    openGraph: {
      type: "article",
      title: `${pageTitle} - ${siteConfig.name}`,
      description: pageDescription,
      url: pageUrl,
      siteName: siteConfig.name,
      images: [
        {
          url: siteConfig.ogImage,
          width: 1200,
          height: 630,
          alt: `${pageTitle} - ${siteConfig.name}`,
        },
      ],
    },
    twitter: {
      card: "summary_large_image",
      title: `${pageTitle} - ${siteConfig.name}`,
      description: pageDescription,
      images: [siteConfig.ogImage],
      creator: siteConfig.creator,
    },
    alternates: {
      canonical: pageUrl,
    },
  };
}

/**
 * JSON-LD structured data for the homepage
 */
export function getHomePageJsonLd() {
  return {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: siteConfig.name,
    applicationCategory: "DeveloperApplication",
    operatingSystem: "Cross-platform",
    description: siteConfig.description,
    url: siteConfig.url,
    author: {
      "@type": "Organization",
      name: siteConfig.name,
      url: siteConfig.url,
    },
    offers: {
      "@type": "Offer",
      price: "0",
      priceCurrency: "USD",
    },
    programmingLanguage: "Go",
  };
}

/**
 * JSON-LD structured data for documentation pages
 */
export function getDocsPageJsonLd({
  title,
  description,
  slug,
  dateModified,
}: {
  title: string;
  description?: string;
  slug: string;
  dateModified?: Date;
}) {
  return {
    "@context": "https://schema.org",
    "@type": "TechArticle",
    headline: title,
    description: description || siteConfig.description,
    url: `${siteConfig.url}/docs/${slug}`,
    author: {
      "@type": "Organization",
      name: siteConfig.name,
      url: siteConfig.url,
    },
    publisher: {
      "@type": "Organization",
      name: siteConfig.name,
      url: siteConfig.url,
    },
    ...(dateModified && {
      dateModified: dateModified.toISOString(),
    }),
    inLanguage: "en-US",
    isAccessibleForFree: true,
    about: {
      "@type": "ComputerLanguage",
      name: "Go",
    },
  };
}

/**
 * JSON-LD for the documentation website
 */
export function getWebsiteJsonLd() {
  return {
    "@context": "https://schema.org",
    "@type": "WebSite",
    name: siteConfig.name,
    url: siteConfig.url,
    description: siteConfig.description,
    potentialAction: {
      "@type": "SearchAction",
      target: {
        "@type": "EntryPoint",
        urlTemplate: `${siteConfig.url}/docs?search={search_term_string}`,
      },
      "query-input": "required name=search_term_string",
    },
  };
}
