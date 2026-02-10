import fs from "fs";
import path from "path";
import matter from "gray-matter";

const CONTENT_DIR = path.join(process.cwd(), "src", "content", "docs");

export interface DocMeta {
  title: string;
  description?: string;
  slug: string;
  [key: string]: unknown;
}

export interface TocEntry {
  /** The heading text */
  text: string;
  /** Depth: 2 for h2, 3 for h3, etc. */
  depth: number;
  /** Slug-ified id for linking */
  id: string;
}

/** Extract h2/h3 headings from raw MDX/markdown content */
export function extractHeadings(content: string): TocEntry[] {
  const headings: TocEntry[] = [];
  // Match lines starting with ## or ### (not inside code blocks)
  let inCodeBlock = false;
  for (const line of content.split("\n")) {
    if (line.trim().startsWith("```")) {
      inCodeBlock = !inCodeBlock;
      continue;
    }
    if (inCodeBlock) continue;

    const match = line.match(/^(#{2,3})\s+(.+)$/);
    if (match) {
      const depth = match[1].length;
      const text = match[2]
        .replace(/\*\*([^*]+)\*\*/g, "$1") // strip bold
        .replace(/`([^`]+)`/g, "$1") // strip inline code
        .replace(/\[([^\]]+)\]\([^)]+\)/g, "$1") // strip links
        .trim();
      const id = text
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, "-")
        .replace(/^-|-$/g, "");
      headings.push({ text, depth, id });
    }
  }
  return headings;
}

/** Resolve the filesystem path for a doc slug, or null if not found. */
export function getDocFilePath(slug: string): string | null {
  const possiblePaths = [
    path.join(CONTENT_DIR, `${slug}.mdx`),
    path.join(CONTENT_DIR, slug, "index.mdx"),
  ];
  for (const filePath of possiblePaths) {
    if (fs.existsSync(filePath)) return filePath;
  }
  return null;
}

/** Get the last-modified date of a doc file. */
export function getDocLastModified(slug: string): Date | null {
  const filePath = getDocFilePath(slug);
  if (!filePath) return null;
  const stat = fs.statSync(filePath);
  return stat.mtime;
}

export function getDocBySlug(slug: string): {
  content: string;
  meta: DocMeta;
} | null {
  const filePath = getDocFilePath(slug);
  if (!filePath) return null;

  const raw = fs.readFileSync(filePath, "utf-8");
  const { data, content } = matter(raw);
  return {
    content,
    meta: {
      title: data.title || slug,
      description: data.description,
      slug,
      ...data,
    },
  };
}

export function getAllDocSlugs(): string[] {
  const slugs: string[] = [];

  function walk(dir: string, prefix: string) {
    if (!fs.existsSync(dir)) return;
    const entries = fs.readdirSync(dir, { withFileTypes: true });
    for (const entry of entries) {
      if (entry.isDirectory()) {
        walk(path.join(dir, entry.name), `${prefix}${entry.name}/`);
      } else if (entry.name.endsWith(".mdx")) {
        const name = entry.name.replace(/\.mdx$/, "");
        if (name === "index") {
          slugs.push(prefix.replace(/\/$/, ""));
        } else {
          slugs.push(`${prefix}${name}`);
        }
      }
    }
  }

  walk(CONTENT_DIR, "");
  return slugs.filter(Boolean);
}
