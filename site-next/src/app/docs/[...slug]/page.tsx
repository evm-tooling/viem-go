import { notFound } from "next/navigation";
import { MDXRemote } from "next-mdx-remote/rsc";
import { getDocBySlug, getAllDocSlugs, extractHeadings } from "@/lib/mdx";
import { CodeGroup } from "@/components/CodePanel";
import CopyButton from "@/components/CopyButton";
import TerminalTyping from "@/components/TerminalTyping";
import GitHubStats from "@/components/GitHubStats";
import Aside from "@/components/Aside";
import TableOfContents from "@/components/TableOfContents";

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
        <a href={`#${id}`} className="no-underline text-inherit hover:text-inherit">
          {children}
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

/* Table components with inline styles to defeat Tailwind preflight */
function Table({ children }: { children?: React.ReactNode }) {
  return (
    <div style={{ overflowX: "auto", margin: "1.5rem 0" }}>
      <table
        style={{
          display: "table",
          width: "100%",
          borderCollapse: "collapse",
          border: "1px solid rgba(57,145,205,0.15)",
          borderRadius: "0.5rem",
          overflow: "hidden",
          fontSize: "0.9375rem",
        }}
      >
        {children}
      </table>
    </div>
  );
}
function Thead({ children }: { children?: React.ReactNode }) {
  return <thead style={{ display: "table-header-group" }}>{children}</thead>;
}
function Tbody({ children }: { children?: React.ReactNode }) {
  return <tbody style={{ display: "table-row-group" }}>{children}</tbody>;
}
function Tr({ children }: { children?: React.ReactNode }) {
  return (
    <tr style={{ display: "table-row", borderBottom: "1px solid rgba(57,145,205,0.1)" }}>
      {children}
    </tr>
  );
}
function Th({ children }: { children?: React.ReactNode }) {
  return (
    <th
      style={{
        display: "table-cell",
        padding: "0.75rem 1rem",
        textAlign: "left",
        background: "#252d3a",
        fontWeight: 600,
        color: "#ffffff",
        borderBottom: "1px solid rgba(57,145,205,0.25)",
      }}
    >
      {children}
    </th>
  );
}
function Td({ children }: { children?: React.ReactNode }) {
  return (
    <td
      style={{
        display: "table-cell",
        padding: "0.75rem 1rem",
        textAlign: "left",
        color: "#b8c5d4",
        background: "rgba(37,45,58,0.3)",
        borderBottom: "1px solid rgba(57,145,205,0.08)",
      }}
    >
      {children}
    </td>
  );
}

const mdxComponents = {
  CodeGroup,
  CopyButton,
  TerminalTyping,
  GitHubStats,
  Aside,
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

  return {
    title: `${doc.meta.title} - viem-go`,
    description: doc.meta.description,
  };
}

export default async function DocPage({ params }: PageProps) {
  const { slug } = await params;
  const slugStr = slug.join("/");
  const doc = getDocBySlug(slugStr);

  if (!doc) {
    notFound();
  }

  const headings = extractHeadings(doc.content);

  return (
    <div className="flex gap-0">
      <article className="flex-1 min-w-0 max-w-[75ch]">
        <h1 className="text-3xl font-bold text-white mb-2">
          {doc.meta.title}
        </h1>
        {doc.meta.description && (
          <p className="text-lg text-gray-3 mb-8 leading-relaxed">
            {doc.meta.description}
          </p>
        )}
        <div className="docs-prose pr-12">
          <MDXRemote source={doc.content} components={mdxComponents} />
        </div>
      </article>
      <TableOfContents headings={headings} />
    </div>
  );
}
