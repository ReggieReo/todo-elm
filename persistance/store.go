package persistence

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v4"
	"golang.org/x/crypto/bcrypt"
)

// ErrUserExists is returned when trying to create a user that already exists.
var ErrUserExists = errors.New("user already exists")

// Store manages the BadgerDB connection and operations.
type Store struct {
	db *badger.DB
}

// NewStore initializes and returns a new Store instance.
// It creates the database directory if it doesn't exist.
func NewStore(baseDir string) (*Store, error) {
	dbDir := filepath.Join(baseDir, "badger") // Store DB in a subdir
	if err := os.MkdirAll(dbDir, 0777); err != nil {
		return nil, fmt.Errorf("failed to create BadgerDB directory %s: %w", dbDir, err)
	}
	log.Printf("BadgerDB directory: %s\n", dbDir)

	opts := badger.DefaultOptions(dbDir)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the underlying BadgerDB database.
func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

// userKey generates the database key for a user.
func userKey(username string) []byte {
	return []byte("user:" + username)
}

// CreateUser attempts to create a new user in the database.
// It hashes the password before storing.
// Returns ErrUserExists if the username is already taken.
func (s *Store) CreateUser(username, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	key := userKey(username)

	err = s.db.Update(func(txn *badger.Txn) error {
		// 1. Check if user already exists
		_, err = txn.Get(key)
		if err == nil {
			// Key exists, username is taken
			return ErrUserExists // Use the custom error
		}
		if err != badger.ErrKeyNotFound {
			// Different error occurred during Get
			return fmt.Errorf("failed checking username: %w", err)
		}

		// 2. Key not found, safe to set the new user data
		err = txn.Set(key, hashedPassword)
		if err != nil {
			return fmt.Errorf("failed saving user: %w", err)
		}
		return nil // Commit transaction
	})

	return err // Return the result of the transaction (nil, ErrUserExists, or other error)
}

// AuthenticateUser checks if the username exists and the password is correct.
// Returns the username on success, or an error otherwise.
func (s *Store) AuthenticateUser(username, password string) (string, error) {
	key := userKey(username)
	var hashedPassword []byte

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return errors.New("invalid username or password") // Generic error for security
			}
			return fmt.Errorf("failed retrieving user: %w", err)
		}

		// Retrieve the hashed password
		// Value() requires a function to process the value; copy it out.
		err = item.Value(func(val []byte) error {
			hashedPassword = append([]byte{}, val...) // Copy value since it's only valid during tx
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed reading password hash: %w", err)
		}
		return nil
	})

	if err != nil {
		return "", err // Return error from View transaction (e.g., user not found)
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		// Passwords don't match (bcrypt.ErrMismatchedHashAndPassword) or other bcrypt error
		return "", errors.New("invalid username or password") // Generic error for security
	}

	// Authentication successful
	return username, nil
}
