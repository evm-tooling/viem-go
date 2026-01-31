import { defineConfig } from 'vocs'

export default defineConfig({
  title: 'viem-go',
  titleTemplate: '%s | viem-go',
  description: 'Go Interface for Ethereum',
  rootDir: '.',
  topNav: [
    { text: 'Docs', link: '/docs/getting-started' },
    { text: 'Examples', link: '/docs/examples' },
  ],
  socials: [
    {
      icon: 'github',
      link: 'https://github.com/ChefBingbong/viem-go',
    },
  ],
  sidebar: {
    '/docs/': [
      {
        text: 'Introduction',
        items: [
          { text: 'Why viem-go', link: '/docs/introduction' },
          { text: 'Getting Started', link: '/docs/getting-started' },
          { text: 'Examples', link: '/docs/examples' },
        ],
      },
      {
        text: 'Clients',
        items: [
          { text: 'Introduction', link: '/docs/clients/intro' },
          { text: 'Public Client', link: '/docs/clients/public' },
          { text: 'Wallet Client', link: '/docs/clients/wallet' },
          {
            text: 'Transports',
            items: [
              { text: 'HTTP', link: '/docs/clients/transports/http' },
              { text: 'WebSocket', link: '/docs/clients/transports/websocket' },
            ],
          },
        ],
      },
      {
        text: 'Accounts',
        items: [
          { text: 'Overview', link: '/docs/accounts/overview' },
          { text: 'Private Key', link: '/docs/accounts/private-key' },
          { text: 'Mnemonic', link: '/docs/accounts/mnemonic' },
          { text: 'HD Wallet', link: '/docs/accounts/hd-wallet' },
        ],
      },
      {
        text: 'Contract',
        items: [
          { text: 'Reading Contracts', link: '/docs/contract/read-contract' },
          { text: 'Writing Contracts', link: '/docs/contract/write-contract' },
          { text: 'ABI Encoding', link: '/docs/contract/abi-encoding' },
        ],
      },
      {
        text: 'Utilities',
        items: [
          { text: 'Units', link: '/docs/utilities/units' },
          { text: 'Hashing', link: '/docs/utilities/hashing' },
          { text: 'Signatures', link: '/docs/utilities/signatures' },
          { text: 'Addresses', link: '/docs/utilities/addresses' },
        ],
      },
    ],
  },
})
