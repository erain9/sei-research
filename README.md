# SEI Account Generator

This tool generates 10 SEI accounts following the Cosmos standard. Each account includes:

- Mnemonic (seed phrase)
- Address (with SEI prefix)
- Public key
- Private key

## Requirements

- Go 1.20+

## Installation

1. Clone this repository

```bash
git clone https://github.com/yourusername/sei-account-generator.git
cd sei-account-generator
```

2. Install dependencies

```bash
go mod tidy
```

## Usage

Run the generator:

```bash
go run main.go
```

The program will generate 10 SEI accounts and display their details.

## Output

For each account, the program outputs:
- Address (with `sei` prefix)
- Mnemonic (24 words)
- Public key
- Private key

## Technical Details

This generator uses the Cosmos SDK to create SEI accounts. Key details:

- Derivation path: `m/44'/118'/0'/0/0` (Cosmos standard)
- Mnemonic: 24 words (256 bits of entropy)
- Key algorithm: secp256k1

## Security Warning

The accounts generated by this tool are intended for testing and development purposes only. 
Never use these accounts on mainnet with real funds unless you have secured the private keys and mnemonics appropriately.

## SEI Network

SEI is a Layer 1 blockchain built on Cosmos SDK that's specifically optimized for trading applications. For more information, visit [SEI's GitHub](https://github.com/sei-protocol/sei-chain). 