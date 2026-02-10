"use client";

import { useCallback, useEffect, useState } from "react";
import SearchModal from "./SearchModal";
import { Button } from "@/components/ui/button";

export default function SearchTrigger({ compact }: { compact?: boolean }) {
  const [open, setOpen] = useState(false);
  const isMac =
    typeof navigator !== "undefined" &&
    navigator.platform.toUpperCase().includes("MAC");

  /* Global Cmd+K / Ctrl+K listener */
  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        setOpen((prev) => !prev);
      }
    }
    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, []);

  /* Lock body scroll when modal is open */
  useEffect(() => {
    if (open) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }
    return () => {
      document.body.style.overflow = "";
    };
  }, [open]);

  const handleClose = useCallback(() => setOpen(false), []);

  if (compact) {
    return (
      <>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={() => setOpen(true)}
          className="w-12 h-9 p-0 rounded-lg text-foreground-muted hover:text-foreground"
          aria-label="Search"
        >
          <svg
            className="w-5 h-5"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={2}
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
            />
          </svg>
        </Button>
        {open ? <SearchModal onClose={handleClose} /> : null}
      </>
    );
  }

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="flex py-3 items-center gap-2 min-w-[300px] h-10 px-3 rounded-[16px] border border-border bg-background-tertiary/60 text-foreground-muted text-sm cursor-pointer transition-all duration-200 hover:border-primary/40 hover:text-foreground-secondary hover:opacity-65"
      >
        <svg
          className="w-4 h-4 shrink-0"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={2}
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
          />
        </svg>
        <span className="flex-1 text-left">Search...</span>
        <kbd className="hidden sm:inline-flex items-center gap-0.5 text-[11px] text-foreground-muted bg-background/60 border border-border rounded px-1.5 py-0.5 font-mono leading-none">
          {isMac ? "⌘" : "Ctrl"}K
        </kbd>
      </button>
      {open ? <SearchModal onClose={handleClose} /> : null}
    </>
  );
}
