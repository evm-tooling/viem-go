"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { usePathname } from "next/navigation";

interface SidebarContextValue {
  open: boolean;
  toggle: () => void;
  close: () => void;
  /** true when wrapped in SidebarProvider (docs pages) */
  hasSidebar: boolean;
}

const SidebarContext = createContext<SidebarContextValue>({
  open: false,
  toggle: () => {},
  close: () => {},
  hasSidebar: false,
});

export function SidebarProvider({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const [open, setOpen] = useState(false);

  const toggle = useCallback(() => setOpen((prev) => !prev), []);
  const close = useCallback(() => setOpen(false), []);

  /* Close mobile sidebar on navigation */
  useEffect(() => {
    setOpen(false);
  }, [pathname]);

  /* Lock body scroll when sidebar overlay is open */
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

  return (
    <SidebarContext.Provider value={{ open, toggle, close, hasSidebar: true }}>
      {children}
    </SidebarContext.Provider>
  );
}

export function useSidebar() {
  return useContext(SidebarContext);
}
