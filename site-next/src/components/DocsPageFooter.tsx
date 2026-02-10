import Link from "next/link";
import type { FlatNavPage } from "@/lib/docs-nav";


interface DocsPageFooterProps {
  slug: string;
  prev: FlatNavPage | null;
  next: FlatNavPage | null;
  lastModified: Date | null;
}

function formatDate(date: Date): string {
  return date.toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  }) + ", " + date.toLocaleTimeString("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
  });
}

export default function DocsPageFooter({
  slug,
  prev,
  next,
  lastModified,
}: DocsPageFooterProps) {
  const editUrl = `https://github.com/ChefBingbong/viem-go/edit/main/site-next/src/content/docs/${slug}.mdx`;

  return (
    <footer className="mt-10 max-w-[80ch]">
       <div className="flex items-center justify-between text-sm text-foreground-muted mb-6 pb-4 border-b-2">
        <a
          href={editUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center gap-1.5 hover:text-foreground transition-colors"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M12 20h9" />
            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
          </svg>
          Suggest changes to this page
        </a>
        {lastModified && (
          <span>Last updated: {formatDate(lastModified)}</span>
        )}
      </div>
      {/* Prev / Next navigation */}
      <div className="flex items-start justify-between">
        {prev ? (
          <Link
            href={`/docs/${prev.slug}`}
            className="group flex flex-col items-start gap-2"
          >
            <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 px-3 py-1 text-md text-foreground-muted group-hover:text-foreground group-hover:border-foreground/30 transition-colors">
              &larr; Previous
            </span>
            <span className="text-lg font-semibold text-foreground group-hover:text-primary transition-colors">
              {prev.label}
            </span>
          </Link>
        ) : (
          <div />
        )}

        {next ? (
          <Link
            href={`/docs/${next.slug}`}
            className="group flex flex-col items-end gap-2"
          >
            <span className="inline-flex items-center gap-1.5 rounded-full border border-border/60 px-3 py-1 text-md text-foreground-muted group-hover:text-foreground group-hover:border-foreground/30 transition-colors">
              Next &rarr;
            </span>
            <span className="text-lg font-semibold text-foreground group-hover:text-primary transition-colors">
              {next.label}
            </span>
          </Link>
        ) : (
          <div />
        )}
      </div>

      {/* Divider */}
      <div className="mt-8 border-t border-border/30" />

      {/* Copyright + social icons */}
      <div className="mt-5 flex items-center justify-between pb-8">
        <span className="text-xs text-foreground-muted">
          &copy; Copyright {new Date().getFullYear()} viem-go. All rights reserved.
        </span>
        <div className="flex items-center gap-3">
          {/* Twitter / X */}
          <a
            href="https://twitter.com"
            target="_blank"
            rel="noopener noreferrer"
            className="text-foreground-muted hover:text-foreground transition-colors"
            title="Twitter"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
              <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
            </svg>
          </a>
          {/* GitHub */}
          <a
            href="https://github.com/ChefBingbong/viem-go"
            target="_blank"
            rel="noopener noreferrer"
            className="text-foreground-muted hover:text-foreground transition-colors"
            title="GitHub"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12" />
            </svg>
          </a>
          {/* Discord */}
          <a
            href="https://discord.gg"
            target="_blank"
            rel="noopener noreferrer"
            className="text-foreground-muted hover:text-foreground transition-colors"
            title="Discord"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
              <path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028 14.09 14.09 0 0 0 1.226-1.994.076.076 0 0 0-.041-.106 13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z" />
            </svg>
          </a>
        </div>
      </div>
    </footer>
  );
}
