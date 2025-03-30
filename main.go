package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
)

func init() {
	// Set up Sei network configuration
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("sei", "seipub")
	config.SetBech32PrefixForValidator("seivaloper", "seivaloperpub")
	config.SetBech32PrefixForConsensusNode("seivalcons", "seivalconspub")
	config.Seal()
}

// Account represents a Sei account with its components
type Account struct {
	Mnemonic   string
	Address    string
	PubKey     string
	PrivateKey string
}

// GenerateAccount creates a new account with mnemonic
func generateAccount() (*Account, error) {
	// Generate a random mnemonic
	entropySizeInBits := 256 // 24 words
	entropy, err := bip39.NewEntropy(entropySizeInBits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	// Derive key from mnemonic using BIP44 HD path for Sei
	// Cosmos coin type is 118, Sei uses the same standard
	derivationPath := "m/44'/118'/0'/0/0"

	// Derive private key from mnemonic
	seed := bip39.NewSeed(mnemonic, "")
	master, ch := hd.ComputeMastersFromSeed(seed)

	// Get private key from derivation path
	derivedPrivateKey, err := hd.DerivePrivateKeyForPath(master, ch, derivationPath)
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %w", err)
	}

	// Create private key object
	privKey := &secp256k1.PrivKey{Key: derivedPrivateKey}

	// Get public key
	pubKey := privKey.PubKey()

	// Get address from public key
	addr := sdk.AccAddress(pubKey.Address())

	// Format the public key
	pubKeyHex := hex.EncodeToString(pubKey.Bytes())

	return &Account{
		Mnemonic:   mnemonic,
		Address:    addr.String(),
		PubKey:     pubKeyHex,
		PrivateKey: hex.EncodeToString(privKey.Key),
	}, nil
}

func main() {
	fmt.Println("Generating 10 SEI Accounts")
	fmt.Println("=======================")

	for i := 1; i <= 10; i++ {
		// Generate new account
		account, err := generateAccount()
		if err != nil {
			fmt.Printf("Error generating account %d: %v\n", i, err)
			os.Exit(1)
		}

		// Print account details
		fmt.Printf("Account #%d\n", i)
		fmt.Printf("Address: %s\n", account.Address)
		fmt.Printf("Mnemonic: %s\n", account.Mnemonic)
		fmt.Printf("Public Key: %s\n", account.PubKey)
		fmt.Printf("Private Key: %s\n", account.PrivateKey)
		fmt.Println("=======================")
	}
}
