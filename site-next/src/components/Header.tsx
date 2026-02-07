"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import SearchTrigger from "./SearchTrigger";
import SidebarToggle from "./SidebarToggle";
import Image from "next/image";

function NavLink({
  href,
  external,
  children,
}: {
  href: string;
  external?: boolean;
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const isActive = !external && pathname.startsWith(href);

  const base =
    "hidden sm:flex items-center text-sm font-medium no-underline px-3 py-1 transition-colors relative";
  const active = isActive
    ? "text-accent after:absolute after:bottom-[-11px] after:left-0 after:right-0 after:h-[2px] after:bg-accent after:rounded-full"
    : "text-gray-2 hover:text-white";

  if (external) {
    return (
      <a
        href={href}
        target="_blank"
        rel="noopener noreferrer"
        className={`${base} ${active}`}
      >
        {children}
      </a>
    );
  }
  return (
    <Link href={href} className={`${base} ${active}`}>
      {children}
    </Link>
  );
}

export default function Header() {
  const pathname = usePathname();
  return (
    <header className="sticky top-0 z-50 w-full  bg-dark-deep/80 py-[5px]">
      <div className="px-4 sm:px-6 h-12 flex items-center justify-between gap-3">
        {/* Left: logo area (matches sidebar width) + search (aligns with main content) */}
        <div className="flex items-center min-w-0">
          <div className={`flex items-center gap-3 sm:gap-4 shrink-0 ${pathname === "/" ? "" : "lg:w-[260px] xl:w-[320px] 2xl:w-[355px] xl:pl-14 2xl:pl-20"}`}>
            <SidebarToggle />
            <Link href="/" className="flex items-center gap-2 shrink-0">
              <Image
              height={90}
              width={90}
                src="/svg/golem-logo-full-light.svg"
                alt="viem-go"
              />
            </Link>
          </div>
          { pathname !== "/" && <div className="hidden sm:block ">
            <SearchTrigger />
          </div>}{}
        </div>

        {/* Right: nav links */}
        <nav className="flex items-center gap-2 shrink-0">
          <NavLink href="/docs">Docs</NavLink>
          <NavLink href="https://github.com/ChefBingbong/viem-go" external>
            GitHub
          </NavLink>
          {/* Mobile search (icon only) */}
          { pathname !== "/" && <div className="sm:hidden">
            <SearchTrigger compact />
          </div>}
          {/* Version dropdown */}
          <VersionDropdown />
        </nav>
      </div>
    </header>
  );
}

function VersionDropdown() {
  return (
    <div className="relative group">
      <button className="flex items-center gap-1 text-sm font-medium text-gray-2 bg-transparent border border-gray-5 px-3 py-1.5 rounded-md cursor-pointer hover:text-white hover:border-gray-4 hover:bg-white/5 transition-colors">
        v0.1.0
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="12"
          height="12"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          className="transition-transform group-hover:rotate-180"
        >
          <polyline points="6 9 12 15 18 9" />
        </svg>
      </button>
      <div className="absolute top-full right-0 mt-2 min-w-[160px] bg-gray-6 border border-gray-5 rounded-lg p-1 opacity-0 invisible -translate-y-1 transition-all group-hover:opacity-100 group-hover:visible group-hover:translate-y-0 z-50">
        <a
          href="https://github.com/ChefBingbong/viem-go/releases"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center justify-between gap-2 text-gray-2 no-underline text-sm px-3 py-2 rounded hover:text-white hover:bg-white/[0.08] transition-colors"
        >
          Releases
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="opacity-50"
          >
            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
            <polyline points="15 3 21 3 21 9" />
            <line x1="10" y1="14" x2="21" y2="3" />
          </svg>
        </a>
        <a
          href="https://github.com/ChefBingbong/viem-go/tree/main/examples"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center justify-between gap-2 text-gray-2 no-underline text-sm px-3 py-2 rounded hover:text-white hover:bg-white/[0.08] transition-colors"
        >
          Examples
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="opacity-50"
          >
            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
            <polyline points="15 3 21 3 21 9" />
            <line x1="10" y1="14" x2="21" y2="3" />
          </svg>
        </a>
        <a
          href="https://github.com/ChefBingbong/viem-go/blob/main/.github/CONTRIBUTING.md"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center justify-between gap-2 text-gray-2 no-underline text-sm px-3 py-2 rounded hover:text-white hover:bg-white/[0.08] transition-colors"
        >
          Contributing
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="opacity-50"
          >
            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
            <polyline points="15 3 21 3 21 9" />
            <line x1="10" y1="14" x2="21" y2="3" />
          </svg>
        </a>
      </div>
    </div>
  );
}
