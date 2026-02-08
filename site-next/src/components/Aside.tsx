interface AsideProps {
  type?: "note" | "tip" | "caution" | "danger";
  title?: string;
  children: React.ReactNode;
}

const icons: Record<string, React.ReactNode> = {
  note: (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="10" />
      <line x1="12" y1="16" x2="12" y2="12" />
      <line x1="12" y1="8" x2="12.01" y2="8" />
    </svg>
  ),
  tip: (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M9 18h6M10 22h4M12 2a7 7 0 0 1 4 12.9V17H8v-2.1A7 7 0 0 1 12 2z" />
    </svg>
  ),
  caution: (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" />
      <line x1="12" y1="9" x2="12" y2="13" /><line x1="12" y1="17" x2="12.01" y2="17" />
    </svg>
  ),
  danger: (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="10" /><line x1="15" y1="9" x2="9" y2="15" /><line x1="9" y1="9" x2="15" y2="15" />
    </svg>
  ),
};

const styles: Record<
  string,
  { border: string; bg: string; accent: string; content?: string }
> = {
  note: {
    border: "border-[#c678dd]/20",
    bg: "bg-[#c678dd]/[0.35]",
    accent: "text-[#c678dd]",
    content: "text-[#d4b0e8]",
  },
  tip: {
    border: "border-[#c678dd]/25",
    bg: "bg-[#c678dd]/10",
    accent: "text-[#c678dd]",
    content: "[&_:not(pre)>code]:text-[#c678dd]",
  },
  caution: {
    border: "border-warning/25",
    bg: "bg-warning/10",
    accent: "text-warning",
  },
  danger: {
    border: "border-destructive/25",
    bg: "bg-destructive/10",
    accent: "text-destructive",
  },
};

export default function Aside({
  type = "note",
  title,
  children,
}: AsideProps) {
  const style = styles[type] || styles.note;
  const icon = icons[type] || icons.note;

  return (
    <div
      className={`my-6 flex items-start gap-3 rounded-xl border ${style.border} ${style.bg} px-4 py-4`}
    >
      <span className={`mt-0.5 shrink-0 ${style.accent}`}>{icon}</span>
      <div
        className={`text-sm leading-relaxed [&>p]:mb-0 ${style.content || "text-gray-2"}`}
      >
        {children}
      </div>
    </div>
  );
}
