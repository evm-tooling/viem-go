import Header from "@/components/Header";
import Hero from "@/components/Hero";
import FeaturesGrid from "@/components/FeaturesGrid";
import CodeComparison from "@/components/CodeComparison";
import WhyGoSection from "@/components/WhyGoSection";
import FeaturesSection from "@/components/FeaturesSection";
import ComparisonSection from "@/components/ComparisonSection";
import CommunitySection from "@/components/CommunitySection";
import SupportSection from "@/components/SupportSection";
import Footer from "@/components/Footer";

export default function Home() {
  return (
    <>
      <Header />
      <main className="relative overflow-x-clip">
        {/* ── Decorative glow orbs ── */}

        {/* 1 · Hero — large top-right */}
        <div
          aria-hidden
          className="pointer-events-none absolute -top-82 -z-50 right-[20%] h-[600px] w-[1220px] rounded-full opacity-50 blur-[120px]"
          style={{ background: "radial-gradient(circle, hsl(215 85% 55% / 0.45), transparent 70%)" }}
        />

        {/* 2 · Hero — top-left accent */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-24 left-[-8%] h-[420px] w-[420px] rounded-full opacity-20 blur-[100px]"
          style={{ background: "radial-gradient(circle, hsl(215 85% 60% / 0.4), transparent 70%)" }}
        />

        {/* 3 · Hero — small center-right, lower */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-[280px] right-[8%] h-[320px] w-[320px] rounded-full opacity-15 blur-[90px]"
          style={{ background: "radial-gradient(circle, hsl(215 90% 60% / 0.35), transparent 70%)" }}
        />

        {/* 4 · FeaturesGrid carousel — centered glow underneath */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-[18%] left-[30%] h-[350px] w-[550px] rounded-full opacity-12 blur-[110px]"
          style={{ background: "radial-gradient(ellipse, hsl(215 85% 55% / 0.3), transparent 70%)" }}
        />

        {/* 5 · Comparison section — right side */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-[30%] right-[-5%] h-[450px] w-[450px] rounded-full opacity-18 blur-[110px]"
          style={{ background: "radial-gradient(circle, hsl(215 85% 55% / 0.4), transparent 70%)" }}
        />

        {/* 6 · Comparison/WhyGo transition — left side */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-[40%] left-[-6%] h-[500px] w-[500px] rounded-full opacity-20 blur-[110px]"
          style={{ background: "radial-gradient(circle, hsl(215 85% 55% / 0.4), transparent 70%)" }}
        />

        {/* 7 · WhyGo section — small right accent */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-[48%] right-[5%] h-[300px] w-[300px] rounded-full opacity-14 blur-[80px]"
          style={{ background: "radial-gradient(circle, hsl(215 90% 58% / 0.35), transparent 70%)" }}
        />

        {/* 8 · WhyGo/Features transition — center-left */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-[55%] left-[10%] h-[380px] w-[380px] rounded-full opacity-16 blur-[100px]"
          style={{ background: "radial-gradient(circle, hsl(215 85% 55% / 0.35), transparent 70%)" }}
        />

        {/* 9 · Features section — right side */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-[63%] right-[-10%] h-[520px] w-[520px] rounded-full opacity-18 blur-[120px]"
          style={{ background: "radial-gradient(circle, hsl(215 85% 55% / 0.45), transparent 70%)" }}
        />

        {/* 10 · Features/Community gap — small left */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-[72%] left-[-4%] h-[340px] w-[340px] rounded-full opacity-14 blur-[90px]"
          style={{ background: "radial-gradient(circle, hsl(215 90% 58% / 0.3), transparent 70%)" }}
        />

        {/* 11 · Community section — center */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-[78%] left-[35%] h-[450px] w-[450px] rounded-full opacity-16 blur-[110px]"
          style={{ background: "radial-gradient(circle, hsl(215 85% 55% / 0.35), transparent 70%)" }}
        />

        {/* 12 · Community — small right accent */}
        <div
          aria-hidden
          className="pointer-events-none absolute top-[82%] right-[2%] h-[280px] w-[280px] rounded-full opacity-12 blur-[80px]"
          style={{ background: "radial-gradient(circle, hsl(215 90% 60% / 0.3), transparent 70%)" }}
        />

        {/* 13 · Support section — left */}
        <div
          aria-hidden
          className="pointer-events-none absolute bottom-[12%] left-[5%] h-[400px] w-[400px] rounded-full opacity-16 blur-[100px]"
          style={{ background: "radial-gradient(circle, hsl(215 85% 55% / 0.35), transparent 70%)" }}
        />

        {/* 14 · Support/Footer — right */}
        <div
          aria-hidden
          className="pointer-events-none absolute bottom-[4%] right-[-6%] h-[380px] w-[380px] rounded-full opacity-14 blur-[100px]"
          style={{ background: "radial-gradient(circle, hsl(215 85% 58% / 0.3), transparent 70%)" }}
        />

        <div className="relative max-w-[1120px] mx-auto px-6">
          <Hero />
          <FeaturesGrid />
        </div>
        {/* <CodeComparison /> */}
        <ComparisonSection />
        <WhyGoSection />
        <FeaturesSection />
        <CommunitySection />
        <SupportSection />
      </main>
      <Footer />
    </>
  );
}
