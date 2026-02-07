import Link from "next/link";

export default function Footer() {
  return (
    <footer className="w-full bg-dark-deep border-t border-accent/15 pt-16 pb-8 px-8 mt-8">
      <div className="max-w-[1120px] mx-auto flex justify-between gap-12 flex-wrap">
        {/* Brand */}
        <div className="max-w-[280px]">
          <img
            src="/svg/golem-logo-text-light.svg"
            alt="viem-go"
            className="h-8 mb-4"
          />
          <p className="text-gray-2 text-[1rem] leading-relaxed">
            Go Interface for Ethereum
          </p>
        </div>

        {/* Links */}
        <div className="flex gap-16 flex-wrap">
          <div className="flex flex-col gap-3">
            <h4 className="text-[1rem] font-semibold text-gray-3 mb-2 uppercase tracking-wider">
              Resources
            </h4>
            <Link
              href="/docs/introduction"
              className="text-gray-2 no-underline text-[0.9rem] hover:text-white transition-colors"
            >
              Documentation
            </Link>
            <Link
              href="/docs/getting-started"
              className="text-gray-2 no-underline text-[0.9rem] hover:text-white transition-colors"
            >
              Getting Started
            </Link>
            <Link
              href="/docs/examples"
              className="text-gray-2 no-underline text-[0.9rem] hover:text-white transition-colors"
            >
              Examples
            </Link>
          </div>
          <div className="flex flex-col gap-3">
            <h4 className="text-[1rem] font-semibold text-gray-3 mb-2 uppercase tracking-wider">
              Community
            </h4>
            <a
              href="https://github.com/ChefBingbong/viem-go"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-2 no-underline text-[0.9rem] hover:text-white transition-colors"
            >
              GitHub
            </a>
            <a
              href="https://github.com/ChefBingbong/viem-go/discussions"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-2 no-underline text-[0.9rem] hover:text-white transition-colors"
            >
              Discussions
            </a>
            <a
              href="https://twitter.com/ChefBingbong"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-2 no-underline text-[0.9rem] hover:text-white transition-colors"
            >
              Twitter
            </a>
          </div>
          <div className="flex flex-col gap-3">
            <h4 className="text-[1rem] font-semibold text-gray-3 mb-2 uppercase tracking-wider">
              More
            </h4>
            <a
              href="https://viem.sh"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-2 no-underline text-[0.9rem] hover:text-white transition-colors"
            >
              viem (TypeScript)
            </a>
            <a
              href="https://github.com/ethereum/go-ethereum"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-2 no-underline text-[0.9rem] hover:text-white transition-colors"
            >
              go-ethereum
            </a>
            <a
              href="https://ethereum.org"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-2 no-underline text-[0.9rem] hover:text-white transition-colors"
            >
              Ethereum
            </a>
          </div>
        </div>
      </div>

      {/* Bottom bar */}
      <div className="max-w-[1120px] mx-auto mt-12 pt-8 border-t border-accent/10 flex justify-between items-center flex-wrap gap-4">
        <p className="text-gray-4 text-md">Released under the MIT License.</p>
        <p className="text-gray-4 text-md">
          Inspired by{" "}
          <a
            href="https://viem.sh"
            target="_blank"
            rel="noopener noreferrer"
            className="text-accent no-underline hover:underline"
          >
            viem
          </a>
        </p>
      </div>
    </footer>
  );
}
