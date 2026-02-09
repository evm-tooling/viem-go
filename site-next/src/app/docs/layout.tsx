import Header from "@/components/Header";
import DocsSidebar from "@/components/DocsSidebar";
import { SidebarProvider } from "@/components/SidebarContext";

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SidebarProvider>
      <Header />
      <div className="flex h-[calc(100vh-3rem)] bg-dark-bg">
        <DocsSidebar />
        <main className="flex-1 min-w-0 overflow-y-auto bg-background-tertiary/60 relative z-10 lg:rounded-tl-3xl lg:border-l lg:border-t lg:border-[1px] lg:border-accent/10 lg:-ml-px">
          <div className="py-8 px-6 lg:px-12 lg:pr-16">
            {children}
          </div>
        </main>
      </div>
    </SidebarProvider>
  );
}
