interface DocsTableProps {
  headers: string[];
  rows: (string | React.ReactNode)[][];
}

export default function DocsTable({ headers, rows }: DocsTableProps) {
  return (
    <div className="docs-feature-table">
      <div className="docs-feature-table__header">
        {headers.map((h, i) => (
          <span key={i} style={{ flex: 1 }}>{h}</span>
        ))}
      </div>
      {rows.map((row, ri) => (
        <div key={ri} className="docs-feature-table__row">
          {row.map((cell, ci) => (
            <span key={ci} className={ci === 0 ? "docs-feature-table__name" : "docs-feature-table__status"}>
              {cell}
            </span>
          ))}
        </div>
      ))}
    </div>
  );
}
