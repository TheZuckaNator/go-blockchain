package blockchain

// TODO: Implement transaction logic
import (
	"crypto/ecdsa"
)

// A block consistents of timestamp, txns, prev block, and validators pub key

// A Transaction: contains an input/output. \
// The input holds senders pub key and sig
// Output recipeints public key and value.

// NewTransactions: Create a new transaction, sign it with ECDSA + assogm OD


// Serialize and Deserialize it for storage and retrieval

// struct

// Transaction: represents blockchain transaction
type Transaction struct {
	ID []byte
	Input []TxInput
	Output []TxOutput
}

type TxInput struct {
	Signiture []byte
	PublicKey []byte
}

type TxOutput struct {
	Value int
	PublicKey []byte
}

// Create transaction and sign it with ECDSA
func NewTransactions(privateKey ecdsa.PrivateKey, recipeint []byte, amount int) *Transaction {
	TxInput :TxInput{}
	TxOut :TxOutput{Value: amount, PublicKey: recipeint}
	tx := Transaction{
		Input: []TxInput{txIn},
		Output: []TxOutput{TxOut},
	}
	tx.ID = tx.hashTransaction();

	//SIgn the transaction w/ senders privatekey
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, tx.ID)

	if err != nil {
		log.(Panic(err))
	}

	signiture := append(r.Bytes(), s.Bytes()...)
	txIn.Signiture = signiture
	return &tx
}

// hashTransaction: hash the transaction data to create a unique id
func (tx *Transaction) hashTransaction() []byte{
	var hash [32]byte
	hash = sha256.Sum256(bytes.Join([][]byte{
		tx.Input[0].PublicKey,
		tx.Output[0].PublicKey,
		[]byte(string(tx.Output[0].Value))
	},[]byte{}))
	return hash[:]

}

// Serialize and Deserialize
// Serialize the transaction into a byte array to be saved
func(tx *Transaction) Serialize []byte {
	var encoded bytes.Buffer
	enc gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// Deserialize: Deserialize a txn from a byte array
// points to transaction
func DeserializeTransaction(data []byte) *Transaction {
	var transaction Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return &transaction
}
