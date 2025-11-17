package blockchain

import (
    "crypto/rand"  // Cryptographically secure random (GOOD for blockchain!)
    "log"
    "math/big"
)

// POSValidator represents a validator in the Proof of Stake consensus system.
// Each validator is identified by their public key and has a stake (amount of tokens locked).
type POSValidator struct {
    PublicKey []byte  // Unique identifier for the validator
    Stake     int     // Amount of tokens staked (higher stake = higher chance of selection)
}

// ProofOfStake selects a validator to create the next block using weighted random selection.
// The selection probability is proportional to each validator's stake.
// 
// Algorithm:
// 1. Calculate the total stake across all validators
// 2. Generate a CRYPTOGRAPHICALLY SECURE random number between 0 and totalStake
// 3. Iterate through validators, subtracting their stake from the random number
// 4. When the random number reaches 0 or below, that validator is selected
//
// Example: If Validator A has 70 tokens and Validator B has 30 tokens:
// - Total stake = 100
// - Random number between 0-99
// - If random = 45, Validator A is selected (45 < 70)
// - If random = 85, Validator B is selected (85-70=15, 15 < 30)
//
// Security Note:
//   Uses crypto/rand instead of math/rand to prevent validator selection manipulation.
//   Math/rand is predictable and can be exploited in blockchain systems.
//
// Parameters:
//   validators: map of validator addresses to their POSValidator structs
//
// Returns:
//   string: the public key of the selected validator
func ProofOfStake(validators map[string]*POSValidator) string {
    // Step 1: Calculate the total stake of all validators
    // This gives us the range for our weighted random selection
    totalStake := 0
    for _, validator := range validators {
        totalStake += validator.Stake
    }

    // Step 2: Generate a cryptographically secure random number in the range [0, totalStake)
    // crypto/rand provides unpredictable randomness (essential for fair validator selection)
    randomBig, err := rand.Int(rand.Reader, big.NewInt(int64(totalStake)))
    if err != nil {
        log.Panic(err)
    }
    
    // Convert big.Int to int64 for our selection algorithm
    random := randomBig.Int64()
    
    // Step 3: Weighted selection - iterate through validators
    // Subtract each validator's stake from our random number
    // The first validator that brings random to 0 or below wins
    for _, validator := range validators {
        random -= int64(validator.Stake)  // Cast to int64 for proper subtraction
        
        // If random is 0 or negative, this validator is selected
        // This gives validators with higher stakes a proportionally higher chance
        if random <= 0 {
            return string(validator.PublicKey)
        }
    }
    
    // This should never happen if totalStake > 0 and validators exist
    // Panic indicates a critical error in the selection logic
    log.Panic("Unable to find a validator")
    return ""
}