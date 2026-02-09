export function Table({ children }: { children?: React.ReactNode }) {
  return (
    <div className="overflow-x-auto my-6 rounded-lg border border-card-border overflow-hidden">
      <table className="w-full border-collapse text-[0.875rem]">
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
    <tr className="border-b border-accent/10 last:border-b-0">{children}</tr>
  );
}

export function Th({ children }: { children?: React.ReactNode }) {
  return (
    <th className="px-4 py-2.5 text-left text-[0.8125rem] font-semibold text-foreground bg-background-tertiary/70 border-b border-accent/20">
      {children}
    </th>
  );
}

export function Td({ children }: { children?: React.ReactNode }) {
  return (
    <td className="px-4 py-2.5 text-left text-[0.875rem] text-foreground-secondary bg-background-secondary/60">
      {children}
    </td>
  );
}
