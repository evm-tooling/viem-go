const communityLinks = [
  {
    title: "Discussions",
    description:
      "Join the discussions on GitHub to ask questions and share ideas",
    href: "https://github.com/ChefBingbong/viem-go/discussions",
  },
  {
    title: "Issues",
    description: "Report issues or request features on GitHub Issues",
    href: "https://github.com/ChefBingbong/viem-go/issues",
  },
  {
    title: "viem Docs",
    description:
      "Check out the original viem documentation for conceptual understanding",
    href: "https://viem.sh",
  },
];

export default function CommunitySection() {
  return (
    <section className="w-full py-8 px-8 bg-gradient-to-b from-dark-deep/30 to-transparent">
      <div className="max-w-[1120px] mx-auto mb-10 text-center">
        <h2 className="text-[2.5rem] font-semibold text-white mb-2">
          Community
        </h2>
        <p className="text-[1.qrem] text-gray-2 max-w-[600px] mx-auto leading-relaxed">
          Check out the following places for more viem-go content
        </p>
      </div>
      <div className="grid grid-cols-3 gap-5 max-w-[1120px] mx-auto max-lg:grid-cols-2 max-sm:grid-cols-1">
        {communityLinks.map((link) => (
          <a
            key={link.title}
            href={link.href}
            target="_blank"
            rel="noopener noreferrer"
            className="block bg-gray-6/50 border border-accent/20 rounded-xl p-6 transition-all duration-200 hover:border-accent/40 hover:bg-gray-6/70 hover:-translate-y-0.5 no-underline group"
          >
            <h3 className="text-[1.25rem] font-semibold text-white mb-1.5 group-hover:text-accent transition-colors">
              {link.title}
            </h3>
            <p className="text-[1rem] text-gray-2 leading-relaxed">
              {link.description}
            </p>
          </a>
        ))}
      </div>
    </section>
  );
}
