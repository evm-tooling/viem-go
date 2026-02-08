"use client";

import { useEffect, useState } from "react";

interface RepoData {
  stars: string;
  license: string;
}

export default function GitHubStats() {
  const [data, setData] = useState<RepoData>({
    stars: "...",
    license: "...",
  });
  const [coverage] = useState("95%");

  useEffect(() => {
    async function fetchGitHubData() {
      try {
        const response = await fetch(
          "https://api.github.com/repos/ChefBingbong/viem-go"
        );
        if (response.ok) {
          const repo = await response.json();
          const starsCount = repo.stargazers_count || 0;
          let starsFormatted: string;
          if (starsCount >= 1000) {
            starsFormatted = (starsCount / 1000).toFixed(1) + "k";
          } else {
            starsFormatted = starsCount.toString();
          }
          const license = repo.license?.spdx_id || "MIT";
          setData({ stars: starsFormatted, license });
        }
      } catch (error) {
        console.error("Failed to fetch GitHub data:", error);
        setData({ stars: "0", license: "MIT" });
      }
    }
    fetchGitHubData();
  }, []);

  return (
    <div className="flex justify-center items-center gap-1.5 shrink-0 mt-auto text-center">
      {/* Stars badge */}
      <a
        href="https://github.com/ChefBingbong/viem-go/stargazers"
        className="flex-1 h-12 p-1.5 flex items-center justify-center gap-2 rounded-lg border border-accent/20 bg-gray-5/60 text-center no-underline transition-all duration-200 hover:border-accent/50 hover:bg-gray-5/90 hover:-translate-y-px"
        target="_blank"
        rel="noopener noreferrer"
      >
        <span className="text-[15px] font-medium opacity-70 bg-dark-deep/60 px-3 h-full flex justify-center text-center items-center rounded-md text-gray-1 leading-[15.5px]">
          stars
        </span>
        <span className="text-[15px] font-medium flex-1 text-center text-gray-1 hover:text-accent">
          {data.stars}
        </span>
      </a>

      {/* Coverage badge */}
      <div className="flex-1 h-12 p-1.5 flex items-center justify-center gap-2 rounded-lg border border-success/30 bg-success/10 text-center">
        <span className="text-[15px] font-medium opacity-80 bg-dark-deep px-3 h-full flex justify-center text-center items-center rounded-md text-gray-1 leading-[15.5px]">
          coverage
        </span>
        <span className="text-[15px] font-medium flex-1 text-center text-success">
          {coverage}
        </span>
      </div>

      {/* License badge */}
      <a
        href="https://github.com/ChefBingbong/viem-go/blob/main/LICENSE"
        className="flex-1 h-12 p-1.5 flex items-center justify-center gap-2 rounded-lg border border-accent/20 bg-gray-5/60 text-center no-underline transition-all duration-200 hover:border-accent/50 hover:bg-gray-5/90 hover:-translate-y-px max-lg:hidden"
        target="_blank"
        rel="noopener noreferrer"
      >
        <span className="text-[15px] font-medium opacity-70 bg-dark-deep/60 px-3 h-full flex justify-center text-center items-center rounded-md text-gray-1 leading-[15.5px]">
          license
        </span>
        <span className="text-[15px] font-medium flex-1 text-center text-gray-1 hover:text-accent">
          {data.license}
        </span>
      </a>
    </div>
  );
}
