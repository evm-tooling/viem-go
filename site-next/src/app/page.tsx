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
      <main>
        <div className="max-w-[1120px] mx-auto px-6">
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
