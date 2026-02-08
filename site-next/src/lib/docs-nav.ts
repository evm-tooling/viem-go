export interface NavItem {
  label: string;
  slug?: string;
  items?: NavItem[];
}

export const docsNav: NavItem[] = [
  {
    label: "Introduction",
    items: [
      { label: "Why viem-go", slug: "introduction" },
      { label: "Getting Started", slug: "getting-started" },
      { label: "Examples", slug: "examples" },
    ],
  },
  {
    label: "Clients",
    items: [
      { label: "Introduction", slug: "clients/intro" },
      { label: "Public Client", slug: "clients/public" },
      { label: "Wallet Client", slug: "clients/wallet" },
      {
        label: "Transports",
        items: [
          { label: "HTTP", slug: "clients/transports/http" },
          { label: "WebSocket", slug: "clients/transports/websocket" },
        ],
      },
    ],
  },
  {
    label: "Chains",
    items: [
      { label: "Introduction", slug: "chains/introduction" },
      { label: "Mainnet", slug: "chains/mainnet" },
      { label: "Sepolia", slug: "chains/sepolia" },
      { label: "Custom Chains", slug: "chains/custom" },
    ],
  },
  {
    label: "Accounts",
    items: [
      { label: "Overview", slug: "accounts/overview" },
      {
        label: "Key Types",
        items: [
          { label: "Private Key", slug: "accounts/private-key" },
          { label: "Mnemonic", slug: "accounts/mnemonic" },
          { label: "HD Wallet", slug: "accounts/hd-wallet" },
          { label: "Custom & address-only", slug: "accounts/custom" },
        ],
      },
    ],
  },
  {
    label: "Contract",
    items: [
      { label: "Reading Contracts", slug: "contract/read-contract" },
      { label: "Writing Contracts", slug: "contract/write-contract" },
      {
        label: "ABI",
        items: [
          { label: "Encoding", slug: "contract/abi-encoding" },
          { label: "Decoding", slug: "contract/abi-decoding" },
          { label: "Types", slug: "contract/abi-types" },
        ],
      },
    ],
  },
  {
    label: "Utilities",
    items: [
      { label: "Units", slug: "utilities/units" },
      { label: "Hashing", slug: "utilities/hashing" },
      {
        label: "Crypto",
        items: [
          { label: "Signatures", slug: "utilities/signatures" },
          { label: "Addresses", slug: "utilities/addresses" },
        ],
      },
    ],
  },
];
