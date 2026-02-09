"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useState, useEffect, useRef } from "react";
import { docsNav, type NavItem } from "@/lib/docs-nav";
import { useSidebar } from "./SidebarContext";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/cn";
import { LayoutGroup, AnimatePresence, motion } from "framer-motion";

/* ── helpers ── */

function isChildActive(item: NavItem, pathname: string): boolean {
  if (item.slug && pathname === `/docs/${item.slug}`) return true;
  if (item.items) return item.items.some((c) => isChildActive(c, pathname));
  return false;
}

/** Recursively find the slug of the first leaf in a group */
function findFirstLeafSlug(item: NavItem): string | null {
  if (item.slug) return item.slug;
  if (item.items) {
    for (const child of item.items) {
      const slug = findFirstLeafSlug(child);
      if (slug) return slug;
    }
  }
  return null;
}

type HoverState = {
  hoveredId: string | null;
  setHoveredId: (id: string | null) => void;
};

type AccordionState = {
  openGroupId: string | null;
  setOpenGroupId: React.Dispatch<React.SetStateAction<string | null>>;
};

/** Find the group ID that contains the current active page */
function findActiveGroupId(
  sections: NavItem[],
  pathname: string,
  depth: number = 0,
  parentPath: string = "",
): string | null {
  for (const item of sections) {
    if (!item.items) continue;

    // depth 0 = section headers (Introduction, Clients, etc.) — not collapsible
    if (depth === 0) {
      const found = findActiveGroupId(item.items, pathname, 1, item.label);
      if (found) return found;
      continue;
    }

    const itemPath = `${parentPath}/${item.label}`;

    // depth 1+ = collapsible groups
    if (isChildActive(item, pathname)) {
      return `group:${itemPath}`;
    }
  }
  return null;
}

/* ── Shared highlight ── */

function Highlight() {
  return (
    <motion.span
      layoutId="docsSidebarHover"
      className="absolute inset-0 rounded-md bg-gradient-to-r from-primary/15 via-primary/5 to-transparent ring-1 ring-primary/10 will-change-transform pointer-events-none"
      transition={{ type: "spring", stiffness: 700, damping: 45 }}
    />
  );
}

/* ── Leaf: a simple page link ── */

function NavLeaf({
  item,
  depth,
  pathname,
  hover,
  onToggle
}: {
  item: NavItem;
  depth: number;
  pathname: string;
  hover: HoverState;
  onToggle: () => void
}) {
  const isActive = pathname === `/docs/${item.slug}`;
  const id = `leaf:${item.slug}`;
  const showHighlight = hover.hoveredId === id;

  return (
    <div
      className={cn(
        "relative",
        depth === 1 && "border-l-2 transition-colors duration-150",
        depth === 1 && (isActive ? "!border-brand-blue-light" : "!border-gray-5/30"),
      )}
      onClick={onToggle}
      onMouseEnter={() => hover.setHoveredId(id)}
      onFocus={() => hover.setHoveredId(id)}
    >
      <Button
        asChild
        variant="ghost"
        size="sm"

        className={cn(
          "relative w-full justify-start rounded-md h-auto py-1 pr-3 overflow-hidden",
          depth > 1 ? "pl-4" : "pl-4",
          isActive
            ? "text-accent font-medium hover:text-accent"
            : "text-gray-3"
        )}
      >
        <Link href={`/docs/${item.slug}`} className="relative">
          {showHighlight ? <Highlight /> : null}
          <span className="relative z-10">{item.label}</span>
        </Link>
      </Button>
    </div>
  );
}

/* ── Collapsible group ── */

function NavGroupItem({
  item,
  depth,
  pathname,
  hover,
  isOpen,
  onToggle,
  accordion,
  parentPath,
}: {
  item: NavItem;
  depth: number;
  pathname: string;
  hover: HoverState;
  isOpen: boolean;
  onToggle: () => void;
  accordion: AccordionState;
  parentPath: string;
}) {
  const router = useRouter();
  const childActive = isChildActive(item, pathname);
  const itemPath = `${parentPath}/${item.label}`;
  const id = `group:${itemPath}`;
  const showHighlight =
    hover.hoveredId === id || (!hover.hoveredId && childActive);

  const handleClick = () => {
    onToggle();
    // When expanding (not collapsing), navigate to the first leaf
    if (!isOpen) {
      const slug = findFirstLeafSlug(item);
      if (slug) router.push(`/docs/${slug}`);
    }
  };

  return (
    <div
      className={cn(
        "w-full",
        depth === 1 && "border-l-2 transition-colors duration-150",
        depth === 1 && (childActive ? "!border-brand-blue-light" : "!border-gray-5/30"),
      )}
      onClick={handleClick}
    >
      <div
        className={cn(
          "w-full rounded-r-md transition-colors duration-150",
          childActive && isOpen ? "bg-primary/[0.04]" : "",
          isOpen ? "bg-primary/[0.06]" : ""
        )}
      >
      <Button
        type="button"
        variant="ghost"
        size="sm"
        
        onMouseEnter={() => hover.setHoveredId(id)}
        onFocus={() => hover.setHoveredId(id)}
        className={cn(
          "relative w-full justify-between rounded-md h-auto py-1 pr-3 overflow-hidden",
          depth > 1 ? "pl-4" : "pl-4",
          childActive ? "text-gray-1 font-semibold" : "text-gray-3"
        )}
      >
        {showHighlight ? <Highlight /> : null}
        <span className="truncate relative z-10">{item.label}</span>
        <motion.svg
          animate={{ rotate: isOpen ? 90 : 0 }}
          transition={{ type: "spring", stiffness: 400, damping: 25 }}
          xmlns="http://www.w3.org/2000/svg"
          width="12"
          height="12"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2.5"
          strokeLinecap="round"
          strokeLinejoin="round"
          className="shrink-0 text-gray-4 relative z-10"
        >
          <polyline points="9 18 15 12 9 6" />
        </motion.svg>
      </Button>

      <AnimatePresence initial={false}>
        {isOpen && (
          <motion.div
            key={`collapse-${item.label}`}
            initial={{ height: 0, opacity: 0 }}
            animate={{
              height: "auto",
              opacity: 1,
              transition: {
                height: { duration: 0.3, ease: [0.16, 1, 0.3, 1] },
                opacity: { duration: 0.2, delay: 0.06 },
              },
            }}
            exit={{
              height: 0,
              opacity: 0,
              transition: {
                height: { duration: 0.25, ease: [0.16, 1, 0.3, 1] },
                opacity: { duration: 0.12 },
              },
            }}
            className="overflow-hidden"
          >
            <div className="relative ml-3  pl-0.5">
              <NavChildren
                items={item.items!}
                depth={depth + 1}
                pathname={pathname}
                hover={hover}
                accordion={accordion}
                parentPath={itemPath}
              />
            </div>
          </motion.div>
        )}
      </AnimatePresence>
      </div>
    </div>
  );
}

/* ── Container for sibling nav items ── */

function NavChildren({
  items,
  depth,
  pathname,
  hover,
  accordion,
  parentPath,
}: {
  items: NavItem[];
  depth: number;
  pathname: string;
  hover: HoverState;
  accordion: AccordionState;
  parentPath: string;
}) {
  return (
    <div className="flex flex-col gap-0">
      {items.map((child) => {
              const groupId = `group:${parentPath}/${child.label}`;
              const isOpen = accordion.openGroupId === groupId;
        if (child.items) {

          return (
            <NavGroupItem
              key={child.label}
              item={child}
              depth={depth}
              pathname={pathname}
              hover={hover}
              isOpen={isOpen}
              onToggle={() => {
                accordion.setOpenGroupId((prev) =>
                  prev === groupId ? null : groupId
                );
              }}
              accordion={accordion}
              parentPath={parentPath}
            />
          );
        }

        return (
          <NavLeaf
            key={child.slug!}
            item={child}
            depth={depth}
            pathname={pathname}
            hover={hover}
            onToggle={() => {
              accordion.setOpenGroupId((prev) =>
                prev === groupId ? null : groupId
              );
            }}
          />
        );
      })}
    </div>
  );
}

/* ── Top-level section header ── */

function NavSection({
  item,
  pathname,
  hover,
  accordion,
}: {
  item: NavItem;
  pathname: string;
  hover: HoverState;
  accordion: AccordionState;
}) {
  return (
    <div className="mt-1 first:mt-0">
      <span className="flex items-center text-[0.8125rem] font-bold text-gray-1 uppercase tracking-wider mb-2 py-1 px-3">
        {item.label}
      </span>
      {item.items && (
        <div className="ml-3 pl-0.5">
          <NavChildren
            items={item.items}
            depth={1}
            pathname={pathname}
            hover={hover}
            accordion={accordion}
            parentPath={item.label}
          />
        </div>
      )}
    </div>
  );
}

/* ── Social footer ── */

function SocialFooter() {
  return (
    <div className="flex items-center gap-3 px-4 py-4 border-t border-gray-5/50 mt-auto">
      <Button
        asChild
        variant="ghost"
        size="sm"
        className="h-9 w-9 p-0 rounded-md text-gray-3 hover:text-gray-1"
      >
        <a
          href="https://github.com/ChefBingbong/viem-go"
          target="_blank"
          rel="noopener noreferrer"
          title="GitHub"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
          </svg>
        </a>
      </Button>
      <Button
        asChild
        variant="ghost"
        size="sm"
        className="h-9 w-9 p-0 rounded-md text-gray-3 hover:text-gray-1"
      >
        <a
          href="https://twitter.com/ChefBingbong"
          target="_blank"
          rel="noopener noreferrer"
          title="Twitter"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
          </svg>
        </a>
      </Button>
      <Button
        asChild
        variant="ghost"
        size="sm"
        className="h-9 w-9 p-0 rounded-md text-gray-3 hover:text-gray-1"
      >
        <a
          href="https://discord.gg/your-discord"
          target="_blank"
          rel="noopener noreferrer"
          title="Discord"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <path d="M20.317 4.3698a19.7913 19.7913 0 00-4.8851-1.5152.0741.0741 0 00-.0785.0371c-.211.3753-.4447.8648-.6083 1.2495-1.8447-.2762-3.68-.2762-5.4868 0-.1636-.3933-.4058-.8742-.6177-1.2495a.077.077 0 00-.0785-.037 19.7363 19.7363 0 00-4.8852 1.515.0699.0699 0 00-.0321.0277C.5334 9.0458-.319 13.5799.0992 18.0578a.0824.0824 0 00.0312.0561c2.0528 1.5076 4.0413 2.4228 5.9929 3.0294a.0777.0777 0 00.0842-.0276c.4616-.6304.8731-1.2952 1.226-1.9942a.076.076 0 00-.0416-.1057c-.6528-.2476-1.2743-.5495-1.8722-.8923a.077.077 0 01-.0076-.1277c.1258-.0943.2517-.1923.3718-.2914a.0743.0743 0 01.0776-.0105c3.9278 1.7933 8.18 1.7933 12.0614 0a.0739.0739 0 01.0785.0095c.1202.099.246.1981.3728.2924a.077.077 0 01-.0066.1276 12.2986 12.2986 0 01-1.873.8914.0766.0766 0 00-.0407.1067c.3604.698.7719 1.3628 1.225 1.9932a.076.076 0 00.0842.0286c1.961-.6067 3.9495-1.5219 6.0023-3.0294a.077.077 0 00.0313-.0552c.5004-5.177-.8382-9.6739-3.5485-13.6604a.061.061 0 00-.0312-.0286zM8.02 15.3312c-1.1825 0-2.1569-1.0857-2.1569-2.419 0-1.3332.9555-2.4189 2.157-2.4189 1.2108 0 2.1757 1.0952 2.1568 2.419 0 1.3332-.9555 2.4189-2.1569 2.4189zm7.9748 0c-1.1825 0-2.1569-1.0857-2.1569-2.419 0-1.3332.9554-2.4189 2.1569-2.4189 1.2108 0 2.1757 1.0952 2.1568 2.419 0 1.3332-.946 2.4189-2.1568 2.4189z" />
          </svg>
        </a>
      </Button>
    </div>
  );
}

/* ── Sidebar content ── */

function SidebarContent({ pathname }: { pathname: string }) {
  const [hoveredId, setHoveredId] = useState<string | null>(null);
  const [openGroupId, setOpenGroupId] = useState<string | null>(
    () => findActiveGroupId(docsNav, pathname),
  );
  const hover: HoverState = { hoveredId, setHoveredId };
  const accordion: AccordionState = { openGroupId, setOpenGroupId };
  const scrollRef = useRef<HTMLElement>(null);

  /* Sync the open accordion group when the pathname changes */
  useEffect(() => {
    const activeGroup = findActiveGroupId(docsNav, pathname);
    if (activeGroup) {
      setOpenGroupId(activeGroup);
    }
  }, [pathname]);

  return (
    <>
      <LayoutGroup id="docsSidebar">
        <motion.nav
          ref={scrollRef}
          layoutScroll
          className="flex-1 overflow-y-auto py-4 px-3 hide-scrollbar"
          onMouseLeave={() => setHoveredId(null)}
          onBlurCapture={() => setHoveredId(null)}
        >
          {docsNav.map((section) => (
            <NavSection
              key={section.label}
              item={section}
              pathname={pathname}
              hover={hover}
              accordion={accordion}
            />
          ))}
        </motion.nav>
      </LayoutGroup>
      <SocialFooter />
    </>
  );
}

/* ── Main sidebar component ── */

export default function DocsSidebar() {
  const pathname = usePathname();
  const { open, close } = useSidebar();

  return (
    <>
      {/* Desktop sidebar */}
      <aside className="hidden lg:flex w-[260px] xl:w-[320px] 2xl:w-[380px] pl-4 xl:pl-12 2xl:pl-20 shrink-0 h-full flex-col  bg-dark-bg overflow-y-auto">
        <SidebarContent pathname={pathname} />
      </aside>

      {/* Mobile sidebar overlay */}
      {open && (
        <>
          <div
            className="fixed inset-0 z-40 bg-dark-deep/60 backdrop-blur-sm lg:hidden"
            onClick={close}
          />
          <aside className="fixed top-12 left-0 z-50 w-[300px] max-w-[85vw] h-[calc(100vh-3rem)] bg-dark-bg border-r border-accent/15 flex flex-col lg:hidden shadow-2xl animate-slide-in">
            <SidebarContent pathname={pathname} />
          </aside>
        </>
      )}
    </>
  );
}
