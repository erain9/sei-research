package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

const (
	// DBFileName is the name of the encrypted database file
	DBFileName = "sei_accounts.db"
	// DefaultDBPassword is the default password for the encrypted database
	// In production, this should be securely provided, not hardcoded
	DefaultDBPassword = "change-me-in-production"
)

// AccountStore manages secure storage of SEI accounts
type AccountStore struct {
	db     *sql.DB
	dbPath string
	mu     sync.Mutex
}

// NewAccountStore creates a new account store
func NewAccountStore(dbDir string) (*AccountStore, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, DBFileName)
	store := &AccountStore{
		dbPath: dbPath,
	}

	// Initialize the database
	if err := store.openDB(); err != nil {
		return nil, err
	}

	// Create the accounts table if it doesn't exist
	if err := store.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return store, nil
}

// openDB opens the encrypted database
func (s *AccountStore) openDB() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		return nil
	}

	// Determine if the database already exists
	_, err := os.Stat(s.dbPath)
	dbExists := !os.IsNotExist(err)

	// Create connection string with encryption options
	connStr := fmt.Sprintf(
		"%s?_pragma_key=%s&_pragma_cipher_page_size=4096",
		s.dbPath,
		DefaultDBPassword,
	)

	// Open the database connection
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	s.db = db

	// If this is a new database, initialize with some optimization settings
	if !dbExists {
		if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
			return fmt.Errorf("failed to set journal mode: %w", err)
		}
		if _, err := db.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
			return fmt.Errorf("failed to set synchronous mode: %w", err)
		}
	}

	return nil
}

// initSchema creates the necessary tables
func (s *AccountStore) initSchema() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address TEXT UNIQUE NOT NULL,
		mnemonic TEXT NOT NULL,
		public_key TEXT NOT NULL,
		private_key TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_accounts_address ON accounts(address);
	`)
	return err
}

// SaveAccount stores an account in the encrypted database
func (s *AccountStore) SaveAccount(account *Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		return fmt.Errorf("database connection not established")
	}

	// Check if the account already exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM accounts WHERE address = ?", account.Address).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if account exists: %w", err)
	}

	if count > 0 {
		// Account already exists, so we'll skip saving it
		return nil
	}

	// Insert the new account
	_, err = s.db.Exec(
		"INSERT INTO accounts (address, mnemonic, public_key, private_key) VALUES (?, ?, ?, ?)",
		account.Address,
		account.Mnemonic,
		account.PubKey,
		account.PrivateKey,
	)
	if err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}

	return nil
}

// GetAccounts retrieves all stored accounts
func (s *AccountStore) GetAccounts() ([]*Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}

	rows, err := s.db.Query("SELECT address, mnemonic, public_key, private_key FROM accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*Account
	for rows.Next() {
		account := &Account{}
		if err := rows.Scan(&account.Address, &account.Mnemonic, &account.PubKey, &account.PrivateKey); err != nil {
			return nil, fmt.Errorf("failed to scan account row: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating account rows: %w", err)
	}

	return accounts, nil
}

// CountAccounts returns the number of accounts stored in the database
func (s *AccountStore) CountAccounts() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		return 0, fmt.Errorf("database connection not established")
	}

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count accounts: %w", err)
	}

	return count, nil
}

// ExportAccountsJSON exports all accounts to a JSON file (for backup purposes)
func (s *AccountStore) ExportAccountsJSON(filePath string) error {
	accounts, err := s.GetAccounts()
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal accounts to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write accounts to file: %w", err)
	}

	return nil
}

// Close closes the database connection
func (s *AccountStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		err := s.db.Close()
		s.db = nil
		return err
	}
	return nil
}

// DeleteDatabase removes the database file (use with caution)
func (s *AccountStore) DeleteDatabase() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.Printf("Warning: error closing database before deletion: %v", err)
		}
		s.db = nil
	}

	return os.Remove(s.dbPath)
}
