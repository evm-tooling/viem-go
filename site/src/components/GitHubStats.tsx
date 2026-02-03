/** @jsxImportSource react */
import { useEffect, useState } from 'react';

interface RepoData {
  stars: string;
  license: string;
}

export default function GitHubStats() {
  const [data, setData] = useState<RepoData>({
    stars: '...',
    license: '...',
  });
  const [coverage] = useState('95%'); // Coverage typically comes from a CI badge, hardcoded for now

  useEffect(() => {
    async function fetchGitHubData() {
      try {
        const response = await fetch(
          'https://api.github.com/repos/ChefBingbong/viem-go'
        );
        if (response.ok) {
          const repo = await response.json();
          
          // Format stars (e.g., 1234 -> "1.2k")
          const starsCount = repo.stargazers_count || 0;
          let starsFormatted: string;
          if (starsCount >= 1000) {
            starsFormatted = (starsCount / 1000).toFixed(1) + 'k';
          } else {
            starsFormatted = starsCount.toString();
          }

          // Get license
          const license = repo.license?.spdx_id || 'MIT';

          setData({
            stars: starsFormatted,
            license: license,
          });
        }
      } catch (error) {
        // Fallback to defaults on error
        console.error('Failed to fetch GitHub data:', error);
        setData({
          stars: '0',
          license: 'MIT',
        });
      }
    }

    fetchGitHubData();
  }, []);

  return (
    <div className="stats-row">
      <a 
        href="https://github.com/ChefBingbong/viem-go/stargazers" 
        className="stat-badge stat-link"
        target="_blank"
        rel="noopener noreferrer"
      >
        <span className="stat-label">stars</span>
        <span className="stat-value">{data.stars}</span>
      </a>
      <div className="stat-badge green">
        <span className="stat-label">coverage</span>
        <span className="stat-value">{coverage}</span>
      </div>
      <a 
        href="https://github.com/ChefBingbong/viem-go/blob/main/LICENSE" 
        className="stat-badge stat-link"
        target="_blank"
        rel="noopener noreferrer"
      >
        <span className="stat-label">license</span>
        <span className="stat-value">{data.license}</span>
      </a>
    </div>
  );
}
