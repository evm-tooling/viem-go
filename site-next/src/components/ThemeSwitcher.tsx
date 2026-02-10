"use client";

import { Sun, Moon, Monitor } from "lucide-react";
import { useTheme, type Theme } from "./ThemeContext";

const options: { value: Theme; icon: typeof Sun; label: string }[] = [
  { value: "light", icon: Sun, label: "Light" },
  { value: "dark", icon: Moon, label: "Dark" },
  { value: "auto", icon: Monitor, label: "System" },
];

const ThemeSwitcher = () => {
  const { theme, setTheme } = useTheme();

  const active = options.find((o) => o.value === theme)!;
  const ActiveIcon = active.icon;

  return (
    <div className="relative group">
      <button
        className="flex items-center justify-center h-8 w-8 rounded-md border border-border bg-secondary text-foreground-secondary hover:text-foreground hover:bg-background-tertiary transition-all duration-150"
        aria-label={`Theme: ${active.label}`}
      >
        <ActiveIcon className="h-4 w-4" />
      </button>

      <div className="absolute top-full right-0 mt-2 min-w-[140px] rounded-lg border border-card-border bg-card p-1 shadow-xl shadow-black/30 z-50 opacity-0 invisible -translate-y-1 transition-all group-hover:opacity-100 group-hover:visible group-hover:translate-y-0">
        {options.map(({ value, icon: Icon, label }) => (
          <button
            key={value}
            onClick={() => setTheme(value)}
            className={`flex w-full items-center gap-2.5 rounded-md px-3 py-2 text-sm transition-colors duration-150 ${
              theme === value
                ? "text-primary bg-accent/40"
                : "text-foreground-secondary hover:text-foreground hover:bg-secondary"
            }`}
          >
            <Icon className="h-4 w-4 shrink-0" />
            {label}
          </button>
        ))}
      </div>
    </div>
  );
};

export default ThemeSwitcher;
