import { notFound } from "next/navigation";
import { MDXRemote } from "next-mdx-remote/rsc";
import remarkGfm from "remark-gfm";
import { getDocBySlug, getAllDocSlugs, extractHeadings, getDocLastModified } from "@/lib/mdx";
import { getAdjacentPages } from "@/lib/docs-nav";
import { createDocsMetadata, getDocsPageJsonLd } from "@/lib/seo";
import { CodeGroup } from "@/components/CodePanel";
import CopyButton from "@/components/CopyButton";
import TerminalTyping from "@/components/TerminalTyping";
import GitHubStats from "@/components/GitHubStats";
import Aside from "@/components/Aside";
import GoPlayground from "@/components/GoPlayground";
import ReadContractDemo from "@/components/ReadContractDemo";
import { Table, Thead, Tbody, Tr, Th, Td } from "@/components/MdxTable";
import DocsTable from "@/components/DocsTable";
import FlowDiagram from "@/components/FlowDiagram";
import BenchmarkSlider from "@/components/BenchmarkSlider";
import BenchmarkViewer from "@/components/BenchmarkViewer";
import TableOfContents from "@/components/TableOfContents";
import DocsPageFooter from "@/components/DocsPageFooter";

/** Generate a slug id from heading text (matches extractHeadings logic) */
function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "");
}

/** Heading component that auto-generates an id for TOC linking */
function createHeading(level: 2 | 3 | 4) {
  const Tag = `h${level}` as const;
  return function Heading({ children }: { children?: React.ReactNode }) {
    const text =
      typeof children === "string"
        ? children
        : extractText(children);
    const id = slugify(text);
    return (
      <Tag id={id}>
        <a href={`#${id}`} className="heading-anchor group flex items-center gap-2 -ml-7 ">
          <span className="anchor-icon opacity-0 group-hover:opacity-100 text-foreground-muted transition-opacity shrink-0" aria-hidden="true">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
              <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
            </svg>
          </span>
          <span className="heading-text transition-colors duration-150">{children}</span>
        </a>
      </Tag>
    );
  };
}

/** Recursively extract text from React children */
function extractText(node: React.ReactNode): string {
  if (typeof node === "string") return node;
  if (typeof node === "number") return String(node);
  if (Array.isArray(node)) return node.map(extractText).join("");
  if (node && typeof node === "object" && "props" in node) {
    return extractText((node as React.ReactElement<{ children?: React.ReactNode }>).props.children);
  }
  return "";
}

const mdxComponents = {
  CodeGroup,
  CopyButton,
  TerminalTyping,
  GitHubStats,
  Aside,
  GoPlayground,
  ReadContractDemo,
  DocsTable,
  FlowDiagram,
  BenchmarkSlider,
  BenchmarkViewer,
  h2: createHeading(2),
  h3: createHeading(3),
  h4: createHeading(4),
  table: Table,
  thead: Thead,
  tbody: Tbody,
  tr: Tr,
  th: Th,
  td: Td,
};

interface PageProps {
  params: Promise<{ slug: string[] }>;
}

export async function generateStaticParams() {
  const slugs = getAllDocSlugs();
  return slugs.map((slug) => ({
    slug: slug.split("/"),
  }));
}

export async function generateMetadata({ params }: PageProps) {
  const { slug } = await params;
  const slugStr = slug.join("/");
  const doc = getDocBySlug(slugStr);

  if (!doc) return { title: "Not Found" };

  return createDocsMetadata({
    title: doc.meta.title,
    description: doc.meta.description,
    slug: slugStr,
  });
}

export default async function DocPage({ params }: PageProps) {
  const { slug } = await params;
  const slugStr = slug.join("/");
  const doc = getDocBySlug(slugStr);

  if (!doc) {
    notFound();
  }

  const headings = extractHeadings(doc.content);
  const { prev, next } = getAdjacentPages(slugStr);
  const lastModified = getDocLastModified(slugStr);

  // Generate JSON-LD structured data for this documentation page
  const jsonLd = getDocsPageJsonLd({
    title: doc.meta.title,
    description: doc.meta.description,
    slug: slugStr,
    dateModified: lastModified || undefined,
  });

  return (
    <>
      {/* JSON-LD Structured Data for this doc page */}
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify(jsonLd),
        }}
      />
      <div className="flex gap-0">
        <article className="flex-1 min-w-0 max-w-[80ch]">
          <h1 className="heading-1 mb-2">
            {doc.meta.title}
          </h1>
          {doc.meta.description && (
            <p className="text-lead mb-8">
              {doc.meta.description}
            </p>
          )}
          <div className="docs-prose ">
            <MDXRemote
              source={doc.content}
              components={mdxComponents}
              options={{ mdxOptions: { remarkPlugins: [remarkGfm] } }}
            />
          </div>
          <DocsPageFooter
            slug={slugStr}
            prev={prev}
            next={next}
            lastModified={lastModified}
          />
        </article>
        <TableOfContents headings={headings} />
      </div>
    </>
  );
}
