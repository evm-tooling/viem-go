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
import BenchmarkBigNumber from "@/components/BenchMarkSection";

export default function Home() {
  return (
    <>
      <Header />
      <main className="relative overflow-x-clip">
        {/* ── Decorative glow orbs — all use primary token ── */}

        {/* 1 · Hero — large top-right */}
        <div
          aria-hidden
          className="glow-orb -top-82 -z-50 right-[20%] h-[600px] w-[1220px] opacity-50 blur-[120px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.45), transparent 70%)" }}
        />

        {/* 2 · Hero — top-left accent */}
        <div
          aria-hidden
          className="glow-orb top-24 left-[-8%] h-[420px] w-[420px] opacity-20 blur-[100px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.4), transparent 70%)" }}
        />

        {/* 3 · Hero — small center-right */}
        <div
          aria-hidden
          className="glow-orb top-[280px] right-[8%] h-[320px] w-[320px] opacity-15 blur-[90px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.35), transparent 70%)" }}
        />

        {/* 4 · FeaturesGrid — centered */}
        <div
          aria-hidden
          className="glow-orb top-[18%] left-[30%] h-[350px] w-[550px] opacity-12 blur-[110px]"
          style={{ background: "radial-gradient(ellipse, hsl(var(--primary) / 0.3), transparent 70%)" }}
        />

        {/* 5 · Comparison — right */}
        <div
          aria-hidden
          className="glow-orb top-[30%] right-[-5%] h-[450px] w-[450px] opacity-18 blur-[110px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.4), transparent 70%)" }}
        />

        {/* 6 · Comparison/WhyGo — left */}
        <div
          aria-hidden
          className="glow-orb top-[40%] left-[-6%] h-[500px] w-[500px] opacity-20 blur-[110px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.4), transparent 70%)" }}
        />

        {/* 7 · WhyGo — small right */}
        <div
          aria-hidden
          className="glow-orb top-[48%] right-[5%] h-[300px] w-[300px] opacity-14 blur-[80px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.35), transparent 70%)" }}
        />

        {/* 8 · WhyGo/Features — center-left */}
        <div
          aria-hidden
          className="glow-orb top-[55%] left-[10%] h-[380px] w-[380px] opacity-16 blur-[100px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.35), transparent 70%)" }}
        />

        {/* 9 · Features — right */}
        <div
          aria-hidden
          className="glow-orb top-[63%] right-[-10%] h-[520px] w-[520px] opacity-18 blur-[120px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.45), transparent 70%)" }}
        />

        {/* 10 · Features/Community — left */}
        <div
          aria-hidden
          className="glow-orb top-[72%] left-[-4%] h-[340px] w-[340px] opacity-14 blur-[90px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.3), transparent 70%)" }}
        />

        {/* 11 · Community — center */}
        <div
          aria-hidden
          className="glow-orb top-[78%] left-[35%] h-[450px] w-[450px] opacity-16 blur-[110px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.35), transparent 70%)" }}
        />

        {/* 12 · Community — small right */}
        <div
          aria-hidden
          className="glow-orb top-[82%] right-[2%] h-[280px] w-[280px] opacity-12 blur-[80px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.3), transparent 70%)" }}
        />

        {/* 13 · Support — left */}
        <div
          aria-hidden
          className="glow-orb bottom-[12%] left-[5%] h-[400px] w-[400px] opacity-16 blur-[100px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.35), transparent 70%)" }}
        />

        {/* 14 · Support/Footer — right */}
        <div
          aria-hidden
          className="glow-orb bottom-[4%] right-[-6%] h-[380px] w-[380px] opacity-14 blur-[100px]"
          style={{ background: "radial-gradient(circle, hsl(var(--primary) / 0.3), transparent 70%)" }}
        />

        <div className="relative max-w-[1120px] mx-auto px-6">
          <Hero />
          <FeaturesGrid />
        </div>

        <BenchmarkBigNumber />
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
