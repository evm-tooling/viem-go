/**
 * MdxTable — semantic table components for MDX content
 * 
 * Uses tokens: card-border, border, accent, foreground, foreground-secondary,
 * background-tertiary, background-secondary
 */

export function Table({ children }: { children?: React.ReactNode }) {
  return (
    <div className="overflow-x-auto my-6 rounded-lg border border-card-border overflow-hidden">
      <table className="w-full border-collapse text-[var(--fs-small)]">
        {children}
      </table>
    </div>
  );
}

export function Thead({ children }: { children?: React.ReactNode }) {
  return <thead>{children}</thead>;
}

export function Tbody({ children }: { children?: React.ReactNode }) {
  return <tbody>{children}</tbody>;
}

export function Tr({ children }: { children?: React.ReactNode }) {
  return (
    <tr className="border-b border-border/30 last:border-b-0">{children}</tr>
  );
}

export function Th({ children }: { children?: React.ReactNode }) {
  return (
    <th className="px-4 py-2.5 text-left text-[var(--fs-fine)] font-semibold text-foreground bg-background-elevated border-b border-border">
      {children}
    </th>
  );
}

export function Td({ children }: { children?: React.ReactNode }) {
  return (
    <td className="px-4 py-2.5 text-left text-[var(--fs-small)] text-foreground-secondary bg-background-secondary/60">
      {children}
    </td>
  );
}
