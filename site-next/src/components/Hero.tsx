import Link from "next/link";
import CopyButton from "./CopyButton";
import TerminalTyping from "./TerminalTyping";
import GitHubStats from "./GitHubStats";
import Image from "next/image";

export default function Hero() {
  return (
    <section className="relative pt-28 pb-6">
      {/* Background image overlay */}
      <div
        className="absolute top-2 left-[45%] -translate-x-1/2 w-[110vw] bottom-0 bg-no-repeat bg-[center_top] bg-cover opacity-25 z-0 pointer-events-none scale-110"
        style={{
          backgroundImage: "url('/svg/hero-bg.svg')",
          maskImage:
            "linear-gradient(to bottom, black 0%, black 30%, transparent 100%)",
          WebkitMaskImage:
            "linear-gradient(to bottom, black 0%, black 30%, transparent 100%)",
        }}
      />

      <div className="relative z-10 flex justify-between items-stretch gap-12 mb-20 mt-6 max-lg:flex-col max-lg:items-center max-lg:text-center">
        {/* Left - text content */}
        <div className="max-w-[420px] flex flex-col gap-6">
          <Image
            className="items-left w-auto max-sm:h-[38px]"
            width={100}
            height={80}
            src="/svg/golem-logo-text-light.svg"
            alt="viem-go logo"
          />
          <p className="text-[1.25rem] leading-relaxed text-gray-2">
            Build reliable blockchain apps & libraries with{" "}
            <strong className="text-white font-semibold">idiomatic Go</strong>,{" "}
            <strong className="text-white font-semibold">type-safe</strong>, and{" "}
            <strong className="text-white font-semibold">composable</strong> modules that
            interface with Ethereum — inspired by{" "}
            <a
              href="https://viem.sh"
              className="text-accent hover:underline"
            >
              viem
            </a>
          </p>
          <div className="flex gap-2 flex-wrap max-lg:justify-center">
            <Link
              href="/docs/getting-started"
              className="px-4 py-2.5 rounded-lg font-medium text-[1rem] bg-accent text-white hover:bg-accent-high transition-all no-underline"
            >
              Get started
            </Link>
            <Link
              href="/docs/introduction"
              className="px-4 py-2.5 rounded-lg font-medium text-[1rem] bg-transparent border border-gray-4 text-white hover:border-accent hover:text-accent transition-all no-underline"
            >
              Why viem-go?
            </Link>
            <a
              href="https://github.com/ChefBingbong/viem-go"
              className="px-4 py-2.5 rounded-lg font-medium text-[1rem] bg-transparent border border-gray-4 text-white hover:border-accent hover:text-accent transition-all no-underline"
              target="_blank"
              rel="noopener noreferrer"
            >
              GitHub
            </a>
          </div>
        </div>

        {/* Right - terminal + stats */}
        <div className="w-[520px] shrink-0 flex flex-col justify-start gap-4 max-lg:hidden">
          {/* Install terminal */}
          <div className="flex flex-col rounded-lg overflow-hidden border border-accent/20 bg-gray-6/80 min-h-[180px]">
            <div className="flex justify-between items-center bg-dark-deep/60 border-b border-accent/15 pr-2">
              <div className="flex items-center gap-2 px-4 py-2.5 text-[0.875rem] font-medium text-white bg-white/5 border-b-2 border-accent" style={{ fontFamily: "'SF Mono', Menlo, Monaco, 'Courier New', monospace" }}>
                <svg
                  className="w-[18px] h-[18px] text-[#00ADD8]"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path d="M1.811 10.231c-.047 0-.058-.023-.035-.059l.246-.315c.023-.035.081-.058.128-.058h4.172c.046 0 .058.035.035.07l-.199.303c-.023.036-.082.07-.117.07zM.047 11.306c-.047 0-.059-.023-.035-.058l.245-.316c.023-.035.082-.058.129-.058h5.328c.047 0 .07.035.058.07l-.093.28c-.012.047-.058.07-.105.07zm2.828 1.075c-.047 0-.059-.035-.035-.07l.163-.292c.023-.035.07-.07.117-.07h2.337c.047 0 .07.035.07.082l-.023.28c0 .047-.047.082-.082.082zm12.129-2.36c-.736.187-1.239.327-1.963.514-.176.046-.187.058-.34-.117-.174-.199-.303-.327-.548-.444-.737-.362-1.45-.257-2.115.175-.795.514-1.204 1.274-1.192 2.22.011.935.654 1.706 1.577 1.835.795.105 1.46-.175 1.987-.77.105-.13.198-.27.315-.434H10.47c-.245 0-.304-.152-.222-.35.152-.362.432-.97.596-1.274a.315.315 0 01.292-.187h4.253c-.023.316-.023.631-.07.947a4.983 4.983 0 01-.958 2.29c-.841 1.11-1.94 1.8-3.33 1.986-1.145.152-2.209-.07-3.143-.77-.865-.655-1.356-1.52-1.484-2.595-.152-1.274.222-2.419.993-3.424.83-1.086 1.928-1.776 3.272-2.02 1.098-.2 2.15-.07 3.096.571.62.41 1.063.97 1.356 1.648.07.105.023.164-.117.2m3.868 6.461c-1.064-.024-2.034-.328-2.852-1.029a3.665 3.665 0 01-1.262-2.255c-.21-1.32.152-2.489.947-3.529.853-1.122 1.881-1.706 3.272-1.95 1.192-.21 2.314-.095 3.33.595.923.63 1.496 1.484 1.648 2.605.198 1.578-.257 2.863-1.344 3.962-.771.783-1.718 1.273-2.805 1.495-.315.06-.63.07-.934.106zm2.78-4.72c-.011-.153-.011-.27-.034-.387-.21-1.157-1.274-1.81-2.384-1.554-1.087.245-1.788.935-2.045 2.033-.21.912.234 1.835 1.075 2.21.643.28 1.285.244 1.905-.07.923-.48 1.425-1.228 1.484-2.233z" />
                </svg>
                <span>golang</span>
              </div>
              <CopyButton text="go get github.com/ChefBingbong/viem-go" />
            </div>
            <div className="flex-1 flex items-start py-2 px-5 text-[1rem]" style={{ fontFamily: "'SF Mono', Menlo, Monaco, 'Courier New', monospace" }}>
              <TerminalTyping />
            </div>
          </div>
          <GitHubStats />
        </div>
      </div>
    </section>
  );
}
