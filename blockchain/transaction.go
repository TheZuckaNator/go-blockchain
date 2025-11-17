package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// ============================================================================
// TRANSACTION STRUCTURES
// ============================================================================

// Transaction represents a blockchain transaction that transfers value from
// one party to another. Each transaction has:
// - A unique ID (hash of transaction data)
// - Inputs (proof of funds from sender)
// - Outputs (destination and amount for recipient)
type Transaction struct {
	ID     []byte     // Unique identifier (SHA-256 hash of transaction data)
	Input  []TxInput  // List of transaction inputs (sender information)
	Output []TxOutput // List of transaction outputs (recipient information)
}

// TxInput represents the input side of a transaction.
// It proves that the sender has the authority to spend funds by providing:
// - A digital signature (proves ownership of the private key)
// - The sender's public key (identifies who is sending)
type TxInput struct {
	Signature []byte // ECDSA signature proving sender owns the private key
	PublicKey []byte // Sender's public key (identifies the sender)
}

// TxOutput represents the output side of a transaction.
// It specifies:
// - How much value is being transferred
// - Who the recipient is (identified by their public key)
type TxOutput struct {
	Value     int    // Amount of tokens/coins being transferred
	PublicKey []byte // Recipient's public key (who receives the funds)
}

// ============================================================================
// TRANSACTION CREATION
// ============================================================================

// NewTransaction creates a new transaction and signs it using ECDSA.
// 
// Process:
// 1. Create transaction input (sender info) and output (recipient + amount)
// 2. Generate a unique transaction ID by hashing the transaction data
// 3. Sign the transaction ID with the sender's private key
// 4. Attach the signature to the input to prove authenticity
//
// Parameters:
//   privateKey: sender's ECDSA private key (used to sign the transaction)
//   recipient:  recipient's public key (who receives the funds)
//   amount:     number of tokens/coins to transfer
//
// Returns:
//   *Transaction: pointer to the newly created and signed transaction
func NewTransaction(privateKey ecdsa.PrivateKey, recipient []byte, amount int) *Transaction {
	// Create the transaction input (sender side)
	// Initially empty signature - will be filled after signing
	txIn := TxInput{
		PublicKey: privateKey.PublicKey.X.Bytes(), // Sender's public key
	}

	// Create the transaction output (recipient side)
	txOut := TxOutput{
		Value:     amount,    // Amount to transfer
		PublicKey: recipient, // Recipient's public key
	}

	// Build the transaction structure
	tx := Transaction{
		Input:  []TxInput{txIn},
		Output: []TxOutput{txOut},
	}

	// Generate a unique ID for this transaction by hashing its contents
	tx.ID = tx.hashTransaction()

	// Sign the transaction ID with the sender's private key
	// This proves that the sender authorized this transaction
	// ECDSA signature returns two big integers (r, s)
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, tx.ID)
	if err != nil {
		log.Panic(err)
	}

	// Combine r and s into a single signature byte array
	signature := append(r.Bytes(), s.Bytes()...)
	
	// Attach the signature to the input to prove authenticity
	tx.Input[0].Signature = signature

	return &tx
}

// ============================================================================
// TRANSACTION HASHING
// ============================================================================

// hashTransaction generates a unique identifier for the transaction by hashing
// its core data (sender's public key, recipient's public key, and amount).
// 
// Why hash?
// - Creates a unique, fixed-size identifier for the transaction
// - Any change to transaction data will change the hash
// - Used as the data that gets signed (proves transaction integrity)
//
// Returns:
//   []byte: SHA-256 hash of the transaction data
func (tx *Transaction) hashTransaction() []byte {
	// Combine all transaction data into a single byte array
	// Using: sender's public key + recipient's public key + amount
	combinedData := bytes.Join([][]byte{
		tx.Input[0].PublicKey,                    // Who is sending
		tx.Output[0].PublicKey,                   // Who is receiving
		[]byte(string(rune(tx.Output[0].Value))), // How much is being sent
	}, []byte{})

	// Hash the combined data using SHA-256
	hash := sha256.Sum256(combinedData)
	
	// Return as a byte slice (convert from fixed-size array)
	return hash[:]
}

// ============================================================================
// SERIALIZATION (for storage and network transmission)
// ============================================================================

// Serialize converts the transaction into a byte array for storage or
// transmission over the network.
//
// Why serialize?
// - Transactions need to be stored in the blockchain
// - Transactions need to be sent over the network to other nodes
// - Byte arrays are the universal format for data storage/transmission
//
// Uses Go's gob encoding (efficient binary format)
//
// Returns:
//   []byte: serialized transaction data
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	
	// Create a gob encoder that writes to our buffer
	enc := gob.NewEncoder(&encoded)
	
	// Encode the transaction into the buffer
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	
	// Return the encoded bytes
	return encoded.Bytes()
}

// DeserializeTransaction converts a byte array back into a Transaction object.
//
// Why deserialize?
// - Retrieve transactions from storage
// - Receive transactions from other nodes over the network
// - Reconstruct the original transaction structure from bytes
//
// Parameters:
//   data: serialized transaction bytes (from Serialize())
//
// Returns:
//   *Transaction: reconstructed transaction object
func DeserializeTransaction(data []byte) *Transaction {
	var transaction Transaction
	
	// Create a gob decoder that reads from the byte data
	decoder := gob.NewDecoder(bytes.NewReader(data))
	
	// Decode the bytes back into a Transaction struct
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	
	return &transaction
}