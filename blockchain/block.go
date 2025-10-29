package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

// Block represents an immutable unit in the blockchain ledger.
// Each block contains a collection of validated transactions, cryptographic proof
// of its predecessor, and metadata required for consensus verification.
type Block struct {
	Timestamp     time.Time      // Unix timestamp of block creation for temporal ordering
	Transactions  []*Transaction // Merkle tree of state transitions within this block
	PrevBlockHash []byte         // Cryptographic link to parent block, ensuring chain integrity
	Hash          []byte         // SHA-256 digest serving as this block's unique identifier
	Validator     []byte         // Public key of the validator node that proposed this block
	Nonce         int            // Proof-of-work solution or consensus-specific challenge value
}

// NewBlock constructs and initializes a new Block instance with validated transactions.
// The block is cryptographically linked to its predecessor via prevBlockHash, establishing
// an immutable chain. The validator parameter identifies the consensus participant responsible
// for block proposal, enabling accountability in the network.
//
// Returns a pointer to the newly created Block with its hash computed.
func NewBlock(transactions []*Transaction, prevBlockHash []byte, validator []byte) *Block {
	block := &Block{
		Timestamp:     time.Now(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Validator:     validator,
	}
	// Compute deterministic hash immediately to maintain referential integrity
	block.Hash = block.calculateHash()
	return block
}

// calculateHash generates a SHA-256 digest of the block's canonical representation.
// This cryptographic commitment includes all transactions (via their hashes), the previous
// block hash, and temporal data. The resulting hash serves as both a unique identifier
// and tamper-evident seal, as any modification would produce a different hash value.
//
// Implementation uses SHA-256 for its collision resistance and preimage security properties.
func (b *Block) calculateHash() []byte {
	// Aggregate transaction hashes to create a compact cryptographic commitment
	var txHashes []byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Hash()...)
	}

	// Concatenate all block components for deterministic hashing
	hash := sha256.Sum256(bytes.Join([][]byte{
		b.PrevBlockHash,
		txHashes,
		[]byte(b.Timestamp.String()),
	}, []byte{}))

	return hash[:]
}

// Serialize encodes the Block into a byte slice using gob encoding.
// This enables efficient persistence to disk or transmission over the network
// while preserving the complete block structure including all nested data.
//
// Panics if encoding fails, as this indicates a critical system error rather
// than a recoverable condition (e.g., corrupted memory).
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		// Serialization failure is non-recoverable; panic to prevent data corruption
		panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock reconstructs a Block from its serialized byte representation.
// This is the inverse operation of Serialize(), used when loading blocks from
// persistent storage or receiving them from network peers.
//
// Parameters:
//   - data: gob-encoded byte slice representing a serialized Block
//
// Returns a pointer to the reconstructed Block instance.
// Panics if deserialization fails due to malformed data.
func DeserializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	if err != nil {
		// Corrupted block data threatens chain integrity; fail fast
		panic(err)
	}

	return &block
}
