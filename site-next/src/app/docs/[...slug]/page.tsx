import { notFound } from "next/navigation";
import { MDXRemote } from "next-mdx-remote/rsc";
import remarkGfm from "remark-gfm";
import { getDocBySlug, getAllDocSlugs, extractHeadings } from "@/lib/mdx";
import { CodeGroup } from "@/components/CodePanel";
import CopyButton from "@/components/CopyButton";
import TerminalTyping from "@/components/TerminalTyping";
import GitHubStats from "@/components/GitHubStats";
import Aside from "@/components/Aside";
import GoPlayground from "@/components/GoPlayground";
import ReadContractDemo from "@/components/ReadContractDemo";
import { Table, Thead, Tbody, Tr, Th, Td } from "@/components/MdxTable";
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
        <a href={`#${id}`} className="heading-anchor">
          {children}
          <span className="anchor-icon" aria-hidden="true">#</span>
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
      <article className="flex-1 min-w-0 max-w-[80ch]">
        <h1 className="heading-1 mb-2">
          {doc.meta.title}
        </h1>
        {doc.meta.description && (
          <p className="text-lead mb-8">
            {doc.meta.description}
          </p>
        )}
        <div className="docs-prose">
          <MDXRemote
            source={doc.content}
            components={mdxComponents}
            options={{ mdxOptions: { remarkPlugins: [remarkGfm] } }}
          />
        </div>
      </article>
      <TableOfContents headings={headings} />
    </div>
  );
}
