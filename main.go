package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
)

// Account structure is unchanged, just renamed fields to be more consistent
type Account struct {
	Mnemonic   string
	Address    string
	PubKey     string
	PrivateKey string
}

// Default configuration
const (
	DefaultAccountCount     = 10
	DefaultStorageDirectory = ".sei-accounts"
)

func init() {
	// Set up Sei network configuration
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("sei", "seipub")
	config.SetBech32PrefixForValidator("seivaloper", "seivaloperpub")
	config.SetBech32PrefixForConsensusNode("seivalcons", "seivalconspub")
	config.Seal()
}

// generateAccount creates a new account with mnemonic
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
	// Create a home directory for storing accounts if not specified
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Create storage directory
	storageDir := filepath.Join(homeDir, DefaultStorageDirectory)

	// Initialize account store for secure storage
	accountStore, err := NewAccountStore(storageDir)
	if err != nil {
		fmt.Printf("Error initializing account store: %v\n", err)
		os.Exit(1)
	}
	defer accountStore.Close()

	// Check if we already have accounts
	count, err := accountStore.CountAccounts()
	if err != nil {
		fmt.Printf("Error counting accounts: %v\n", err)
		os.Exit(1)
	}

	// If we already have accounts, retrieve and display them
	if count >= DefaultAccountCount {
		fmt.Println("Using existing SEI accounts from secure storage")
		printStoredAccounts(accountStore)
		return
	}

	// We need to generate new accounts
	fmt.Printf("Generating %d SEI Accounts\n", DefaultAccountCount)
	fmt.Println("=======================")

	// Generate and store accounts
	for i := count + 1; i <= DefaultAccountCount; i++ {
		// Generate new account
		account, err := generateAccount()
		if err != nil {
			fmt.Printf("Error generating account %d: %v\n", i, err)
			os.Exit(1)
		}

		// Save account to secure storage
		if err := accountStore.SaveAccount(account); err != nil {
			fmt.Printf("Error saving account %d: %v\n", i, err)
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

	fmt.Println("All accounts have been securely stored on disk.")
	fmt.Printf("You can find them in: %s\n", storageDir)
}

// printStoredAccounts displays all accounts from secure storage
func printStoredAccounts(store *AccountStore) {
	accounts, err := store.GetAccounts()
	if err != nil {
		fmt.Printf("Error retrieving accounts: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=======================")
	for i, account := range accounts {
		fmt.Printf("Account #%d\n", i+1)
		fmt.Printf("Address: %s\n", account.Address)
		fmt.Printf("Mnemonic: %s\n", account.Mnemonic)
		fmt.Printf("Public Key: %s\n", account.PubKey)
		fmt.Printf("Private Key: %s\n", account.PrivateKey)
		fmt.Println("=======================")
	}
}
