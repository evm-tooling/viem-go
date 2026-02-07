"use client";

import { useCallback, useEffect, useState } from "react";
import SearchModal from "./SearchModal";

export default function SearchTrigger({ compact }: { compact?: boolean }) {
  const [open, setOpen] = useState(false);
  const [isMac, setIsMac] = useState(true);

  useEffect(() => {
    setIsMac(navigator.platform.toUpperCase().indexOf("MAC") >= 0);
  }, []);

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
        <button
          onClick={() => setOpen(true)}
          className="flex items-center justify-center w-12 h-9 rounded-lg text-gray-3 hover:text-white hover:bg-white/[0.08] transition-colors"
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
        </button>
        <SearchModal open={open} onClose={handleClose} />
      </>
    );
  }

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="flex items-center gap-2 min-w-[300px] h-10 px-3 rounded-[12px] border border-gray-5 bg-dark-bg text-gray-4 text-sm cursor-pointer transition-all duration-150 hover:border-gray-4 hover:text-gray-3"
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
        <kbd className="hidden sm:inline-flex items-center gap-0.5 text-[11px] text-gray-4 bg-dark-bg/60 border border-gray-5 rounded px-1.5 py-0.5 font-mono leading-none">
          {isMac ? "⌘" : "Ctrl"}K
        </kbd>
      </button>
      <SearchModal open={open} onClose={handleClose} />
    </>
  );
}
