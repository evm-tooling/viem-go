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
    label: "ABI",
    items: [
      { label: "Parse", slug: "abi/parse" },
      { label: "Encode Parameters", slug: "abi/encode-parameters" },
      { label: "Decode Parameters", slug: "abi/decode-parameters" },
      { label: "Encode Packed", slug: "abi/encode-packed" },
      { label: "Get Abi Item", slug: "abi/get-abi-item" },
    ],
  },
  {
    label: "Contract",
    items: [
      { label: "Reading Contracts", slug: "contract/read-contract" },
      { label: "Writing Contracts", slug: "contract/write-contract" },
      { label: "Get Code", slug: "contract/get-code" },
      { label: "Get Storage At", slug: "contract/get-storage-at" },
      { label: "Multicall", slug: "contract/multicall" },
      { label: "Contract Events", slug: "contract/contract-events" },
      { label: "Deploy Contract", slug: "contract/deploy-contract" },
      { label: "Simulate Contract", slug: "contract/simulate-contract" },
      { label: "Estimate Contract Gas", slug: "contract/estimate-contract-gas" },
      { label: "Contract Utilities", slug: "contract/contract-utilities" },
      {
        label: "ABI",
        items: [
          { label: "Introduction", slug: "contract/abi-introduction" },
          { label: "Encoding", slug: "contract/abi-encoding" },
          { label: "Decoding", slug: "contract/abi-decoding" },
          { label: "Types", slug: "contract/abi-types" },
          { label: "Selectors & Items", slug: "contract/abi-selectors" },
        ],
      },
    ],
  },
  {
    label: "Utilities",
    items: [
      { label: "Units", slug: "utilities/units" },
      { label: "Hex", slug: "utilities/hex" },
      { label: "Bytes", slug: "utilities/bytes" },
      {
        label: "Data",
        items: [
          { label: "Concat", slug: "utilities/data/concat" },
          { label: "Pad", slug: "utilities/data/pad" },
          { label: "Slice", slug: "utilities/data/slice" },
          { label: "Trim", slug: "utilities/data/trim" },
          { label: "Size", slug: "utilities/data/size" },
          { label: "Is Hex", slug: "utilities/data/is-hex" },
          { label: "Is Bytes", slug: "utilities/data/is-bytes" },
        ],
      },
      {
        label: "Encoding",
        items: [
          { label: "To Hex", slug: "utilities/encoding/to-hex" },
          { label: "From Hex", slug: "utilities/encoding/from-hex" },
          { label: "To Bytes", slug: "utilities/encoding/to-bytes" },
          { label: "From Bytes", slug: "utilities/encoding/from-bytes" },
          { label: "RLP", slug: "utilities/encoding/rlp" },
        ],
      },
      {
        label: "Hash",
        items: [
          { label: "Keccak256", slug: "utilities/hash/keccak256" },
          { label: "SHA-256", slug: "utilities/hash/sha256" },
          { label: "RIPEMD-160", slug: "utilities/hash/ripemd160" },
          { label: "To Function Selector", slug: "utilities/hash/to-function-selector" },
          { label: "To Event Selector", slug: "utilities/hash/to-event-selector" },
          { label: "To Function Hash", slug: "utilities/hash/to-function-hash" },
          { label: "To Event Hash", slug: "utilities/hash/to-event-hash" },
          { label: "To Signature", slug: "utilities/hash/to-signature" },
          { label: "Normalize Signature", slug: "utilities/hash/normalize-signature" },
          { label: "Is Hash", slug: "utilities/hash/is-hash" },
        ],
      },
      {
        label: "Signature",
        items: [
          { label: "Hash Message", slug: "utilities/signature/hash-message" },
          { label: "Hash Typed Data", slug: "utilities/signature/hash-typed-data" },
          { label: "Verify Message", slug: "utilities/signature/verify-message" },
          { label: "Verify Typed Data", slug: "utilities/signature/verify-typed-data" },
          { label: "Recover Message Address", slug: "utilities/signature/recover-message-address" },
          { label: "Recover Typed Data Address", slug: "utilities/signature/recover-typed-data-address" },
          { label: "Recover Address", slug: "utilities/signature/recover-address" },
          { label: "Recover Public Key", slug: "utilities/signature/recover-public-key" },
          { label: "Parse Signature", slug: "utilities/signature/parse-signature" },
          { label: "Serialize Signature", slug: "utilities/signature/serialize-signature" },
          { label: "Parse Compact Signature", slug: "utilities/signature/parse-compact-signature" },
          { label: "Serialize Compact Signature", slug: "utilities/signature/serialize-compact-signature" },
          { label: "Signature To Compact Signature", slug: "utilities/signature/signature-to-compact-signature" },
          { label: "Compact Signature To Signature", slug: "utilities/signature/compact-signature-to-signature" },
          { label: "Parse ERC-6492 Signature", slug: "utilities/signature/parse-erc6492-signature" },
          { label: "Serialize ERC-6492 Signature", slug: "utilities/signature/serialize-erc6492-signature" },
          { label: "Is ERC-6492 Signature", slug: "utilities/signature/is-erc6492-signature" },
        ],
      },
      {
        label: "Transaction",
        items: [
          { label: "Parse Transaction", slug: "utilities/transaction/parse-transaction" },
          { label: "Serialize Transaction", slug: "utilities/transaction/serialize-transaction" },
          { label: "Get Transaction Type", slug: "utilities/transaction/get-transaction-type" },
          { label: "Get Serialized Transaction Type", slug: "utilities/transaction/get-serialized-transaction-type" },
          { label: "Assert Transaction", slug: "utilities/transaction/assert-transaction" },
          { label: "Assert Request", slug: "utilities/transaction/assert-request" },
          { label: "Serialize Access List", slug: "utilities/transaction/serialize-access-list" },
        ],
      },
      {
        label: "Blob",
        items: [
          { label: "To Blobs", slug: "utilities/blob/to-blobs" },
          { label: "From Blobs", slug: "utilities/blob/from-blobs" },
          { label: "To Blob Sidecars", slug: "utilities/blob/to-blob-sidecars" },
          { label: "Blobs To Commitments", slug: "utilities/blob/blobs-to-commitments" },
          { label: "Blobs To Proofs", slug: "utilities/blob/blobs-to-proofs" },
          { label: "Commitment To Versioned Hash", slug: "utilities/blob/commitment-to-versioned-hash" },
          { label: "Commitments To Versioned Hashes", slug: "utilities/blob/commitments-to-versioned-hashes" },
          { label: "Sidecars To Versioned Hashes", slug: "utilities/blob/sidecars-to-versioned-hashes" },
        ],
      },
      {
        label: "ENS",
        items: [
          { label: "Namehash", slug: "utilities/ens/namehash" },
          { label: "Labelhash", slug: "utilities/ens/labelhash" },
          { label: "Normalize", slug: "utilities/ens/normalize" },
          { label: "Encode Labelhash", slug: "utilities/ens/encode-labelhash" },
          { label: "Encoded Label To Labelhash", slug: "utilities/ens/encoded-label-to-labelhash" },
          { label: "To Coin Type", slug: "utilities/ens/to-coin-type" },
          { label: "Packet To Bytes", slug: "utilities/ens/packet-to-bytes" },
        ],
      },
      {
        label: "KZG",
        items: [
          { label: "Setup KZG", slug: "utilities/kzg/setup-kzg" },
          { label: "Define KZG", slug: "utilities/kzg/define-kzg" },
        ],
      },
      {
        label: "Authorization",
        items: [
          { label: "Hash Authorization", slug: "utilities/authorization/hash-authorization" },
        ],
      },
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
