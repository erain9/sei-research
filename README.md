# SEI Research

This repository contains tools and research materials for exploring the SEI Network ecosystem. Currently, it includes a tool to generate SEI accounts following the Cosmos standard.

## Account Generator

The account generator creates 10 SEI accounts, each including:

- Mnemonic (seed phrase)
- Address (with SEI prefix)
- Public key
- Private key

## Secure Storage

The account generator now includes secure, encrypted storage functionality:

- Accounts are stored in an encrypted SQLite database (using SQLCipher with AES-256 encryption)
- The database automatically stores generated accounts for later use
- When run, the tool checks for existing accounts before generating new ones
- Database files are excluded from Git using .gitignore

### Security Features

- AES-256 encryption for all stored account data
- Database is password-protected (customize in the code)
- Only stores accounts locally on your machine
- WAL journaling mode for durability and crash resistance
- Thread-safe implementation with mutex protection

### Database Location

By default, the encrypted database is stored in:
```
~/.sei-accounts/sei_accounts.db
```

You can change the storage location by modifying the `DefaultStorageDirectory` constant in the code.

## Understanding Cosmos Accounts

### What is a Cosmos Account?

In the Cosmos ecosystem, an account represents a user's identity on the blockchain. Each account consists of:

1. **Private Key**: The secret key that must remain secure and is used to sign transactions
2. **Public Key**: Derived from the private key and used for verification
3. **Address**: A human-readable string derived from the public key that identifies an account on the network

Cosmos accounts follow the hierarchical deterministic (HD) wallet structure defined in BIP32, with derivation paths following BIP44. This allows a single mnemonic seed to generate multiple accounts across different blockchains.

### Cosmos Account Structure

```
                                 ┌─────────────┐
                                 │  Mnemonic   │ (24-word seed phrase)
                                 └──────┬──────┘
                                        │
                                        ▼
                                 ┌─────────────┐
                                 │ Master Seed │
                                 └──────┬──────┘
                                        │
                                        ▼
┌───────────────────────────────────────────────────────────────────┐
│                       Derivation Path                              │
│      m/44'/118'/0'/0/0  (Cosmos standard for first account)       │
└─────────────────┬─────────────────────────┬───────────────────────┘
                  │                         │
                  ▼                         ▼
           ┌─────────────┐           ┌─────────────┐
           │ Private Key │           │ Public Key  │
           └──────┬──────┘           └──────┬──────┘
                  │                         │
                  │                         ▼
                  │                  ┌─────────────┐
                  │                  │   Address   │ (Bech32 format with prefix)
                  │                  └─────────────┘
                  ▼
       ┌───────────────────┐
       │     Sign          │
       │  Transactions     │
       └───────────────────┘
```

### Address Format in Cosmos

Cosmos uses the Bech32 address format, which consists of:
- A human-readable prefix (network identifier)
- A separator character (1)
- A data part (derived from the public key)
- A checksum

For example, a SEI address looks like: `sei1g7gw7rvn8pnthr6q5f5eagkh9ghwa6g6cyq0nl`
- `sei` is the human-readable prefix specific to the SEI network
- The rest is the encoded public key and checksum

Different Cosmos chains use different prefixes:
- `cosmos` - Cosmos Hub
- `osmo` - Osmosis
- `sei` - SEI Network
- `juno` - Juno Network

### Account Security and Mnemonics

Accounts in Cosmos are secured through BIP39 mnemonics, which are word sequences that can regenerate the private key. This generator uses 24-word mnemonics, providing 256 bits of entropy and maximum security.

The BIP39 process:
1. Generate random entropy (256 bits for 24 words)
2. Add checksum
3. Map to mnemonic words
4. Create seed from mnemonic (+ optional passphrase)
5. Derive master key from seed
6. Derive child keys using the derivation path

## Requirements

- Go 1.20+
- CGO enabled (for SQLCipher compilation)

## Installation

1. Clone this repository

```bash
git clone https://github.com/yourusername/sei-research.git
cd sei-research
```

2. Install dependencies

```bash
go mod tidy
```

## Usage

Run the account generator:

```bash
go run main.go storage.go
```

The program will:
1. Check for existing accounts in the encrypted database
2. If accounts exist, display them
3. If no accounts exist, generate 10 new accounts and store them
4. Display the account details in the terminal

### Output

For each account, the program outputs:
- Address (with `sei` prefix)
- Mnemonic (24 words)
- Public key
- Private key

## Technical Details

The account generator uses the Cosmos SDK to create SEI accounts. Key details:

- Derivation path: `m/44'/118'/0'/0/0` (Cosmos standard)
  - 44' = BIP44 purpose
  - 118' = Cosmos coin type
  - 0' = Account index
  - 0 = Change (external chain)
  - 0 = Address index
- Mnemonic: 24 words (256 bits of entropy)
- Key algorithm: secp256k1 (same as Bitcoin and Ethereum)
- Address format: Bech32 with 'sei' prefix

## Security Warning

The accounts generated by this tool are intended for testing and development purposes only. 
Never use these accounts on mainnet with real funds unless you have secured the private keys and mnemonics appropriately.

## About SEI Network

SEI is a Layer 1 blockchain built on Cosmos SDK that's specifically optimized for trading applications. It features a built-in central limit order book (CLOB) and provides high throughput with low latency, making it ideal for DeFi and trading applications.

For more information, visit [SEI's GitHub](https://github.com/sei-protocol/sei-chain). 