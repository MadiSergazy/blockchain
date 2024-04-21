package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	err := sign()

	if err != nil {
		log.Fatalln(err)
	}

}

func sign() error {
	// privateKey, err := crypto.GenerateKey()
	// if err != nil {
	// 	return err
	// }

	// Need to load the private key file for the configured beneficiary so the
	// account can get credited with fees and tips.
	path := fmt.Sprintf("%s%s.ecdsa", "zblock/accounts/", "kennedy")
	privateKey, err := crypto.LoadECDSA(path)
	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	fmt.Println("Address: ", address)

	// type tx struct {
	// 	To    string
	// 	Value uint64
	// }

	// v := tx{
	// 	To:    "8dc79feefd3b86e2f9991def0e5ccd9a5128e104682407b308594bc1032ac7f0",
	// 	Value: 100,
	// }

	// type txSign struct {
	// 	tx
	// 	Sig []byte
	// }

	v := struct {
		Name string
	}{
		Name: "Bill",
	}

	data, err := stamp(v)
	if err != nil {
		return fmt.Errorf("stamp: %w", err)
	}

	sig, err := crypto.Sign(data, privateKey)
	if err != nil {
		return err
	}

	fmt.Printf("SIG: 0x%s\n", hex.EncodeToString(sig))
	// v2 := tx{
	// 	To:    "8dc79feefd3b86e2f9991def0e5ccd9a5128e104682407b308594bc1032ac7f0",
	// 	Value: 100,
	// }
	v2 := struct {
		Name string
	}{
		Name: "Bill",
	}

	data2, err := stamp(v2)
	if err != nil {
		return fmt.Errorf("stamp: %w", err)
	}
	// data2, err := json.Marshal(v2)
	// if err != nil {
	// 	return err
	// }

	// //hash data into 32 byte
	// txHash := crypto.Keccak256(data)

	// Sign the hash with the private key to produce a signature.
	// sig, err := crypto.Sign(txHash, privateKey)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("Signature: ", sig)

	// txHash2 := crypto.Keccak256(data2)
	//Node side =================================================================
	sigPublicKey, err := crypto.Ecrecover(data2, sig)
	if err != nil {
		return err
	}

	fmt.Println("PKLEN:", len(sigPublicKey))

	rs := sig[:crypto.RecoveryIDOffset]

	if !crypto.VerifySignature(sigPublicKey, data2, rs) {
		return errors.New("invalid signature")
	}

	x, y := elliptic.Unmarshal(crypto.S256(), sigPublicKey)
	pubKey := ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}

	address = crypto.PubkeyToAddress(pubKey).String()
	fmt.Println("Address returned: ", address)

	return nil
}

// stamp returns a hash of 32 bytes that represents this data with
// the Ardan stamp embedded into the final hash.
func stamp(value any) ([]byte, error) {

	// Marshal the data.
	v, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// This stamp is used so signatures we produce when signing data
	// are always unique to the Ardan blockchain.
	stamp := []byte(fmt.Sprintf("\x19Madi Signed Message:\n%d", len(v)))

	// Hash the stamp and txHash together in a final 32 byte array
	// that represents the data.
	data := crypto.Keccak256(stamp, v)

	return data, nil
}
