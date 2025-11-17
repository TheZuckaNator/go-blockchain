package blockchain

// TODO: Implement consensus mechanism

// proof of stake function
import (
    "log"
    "math/rand"
    "time"
)
/*
1. POSValidaor stuct- is a validator in the POS system recognized by the public key
2. Proof of Stake(): selects a validator based on thier stake. more stake a validator has the higher the probabilty og being chosen to validate a block
*/

//1. POSValidaor stuct- is a validator in the POS system recognized by the public key
type POSValidaor struct {
	PublicKey []byte
	Stake int
}
//Proof of stake the higher stake the higher probability of getting chosen to validate a block
func ProofOfStake(validators map[string]*POSValidaor) string {
	totalStake := 0
	for _, validator := range validators{
		totalStake += validator.Stake
	}

	// select a validator randomly
	rand.Seed(time.Now().UnixNano())
	random = rand.InTn(totalStake)

	// loop to select a validator based on stake
	for _, validator := range validators {
		random -= validator.Stake
		if random <= 0 {
		return string(validator.PublicKey)
		}
	}
	log.Panic("Unable to find a validator")
	return ""
}