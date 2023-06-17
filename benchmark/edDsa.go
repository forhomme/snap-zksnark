package main

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	eddsa2 "github.com/consensys/gnark-crypto/signature/eddsa"
	"github.com/consensys/gnark/backend/groth16"
	witness2 "github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	r1cs2 "github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	mimc2 "github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
	"log"
	"smart-contract-service/internal"
	"strings"
)

type EddsaCircuit struct {
	PublicKey eddsa.PublicKey   `gnark:",public"`
	Signature eddsa.Signature   `gnark:",public"`
	Message   frontend.Variable `gnark:",public"`
}

func (circuit *EddsaCircuit) Define(api frontend.API) error {
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	mimc, err := mimc2.NewMiMC(api)
	if err != nil {
		return err
	}

	// verify the signature in the cs
	return eddsa.Verify(curve, circuit.Signature, circuit.Message, circuit.PublicKey, &mimc)
}

func CreateEddsaCircuit() (err error) {
	var circuit EddsaCircuit

	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs2.NewBuilder, &circuit)
	if err != nil {
		log.Fatal(err)
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		log.Fatal(err)
	}

	internal.Serialize(r1cs, internal.R1csEddsaPath)
	internal.Serialize(pk, internal.PkEddsaPath)
	internal.Serialize(vk, internal.VkEddsaPath)

	// read R1CS, proving key and verifying keys
	ccs := groth16.NewCS(ecc.BN254)
	pk = groth16.NewProvingKey(ecc.BN254)
	vk = groth16.NewVerifyingKey(ecc.BN254)
	internal.Deserialize(ccs, internal.R1csEddsaPath)
	internal.Deserialize(pk, internal.PkEddsaPath)
	internal.Deserialize(vk, internal.VkEddsaPath)

	// instantiate hash function
	f := bn254.NewMiMC()

	// create a eddsa key pair
	privateKey, err := eddsa2.New(tedwards.BN254, rand.Reader)
	publicKey := privateKey.Public()

	// note that the message is on 4 bytes
	//msg := []byte{0xde, 0xad, 0xf0, 0x0d}
	preImage := &Data{PreImage: "data"}
	preImageByte := []byte(fmt.Sprintf("%v", preImage))
	b := make([]byte, 32)
	copy(b, preImageByte)

	// sign the message
	signature, err := privateKey.Sign(b, f)

	// verifies signature
	isValid, err := publicKey.Verify(signature, b, f)
	if !isValid {
		return errors.New("not valid")
	}
	// declare the witness
	assignment := &EddsaCircuit{}

	// assign message value
	assignment.Message = b

	// public key bytes
	_publicKey := publicKey.Bytes()

	// assign public key values
	assignment.PublicKey.Assign(tedwards.BN254, _publicKey[:32])

	// assign signature values
	assignment.Signature.Assign(tedwards.BN254, signature)

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		return
	}

	var proofBuf bytes.Buffer
	proof.WriteTo(&proofBuf)

	publicWitness, _ := witness.Public()
	dataBin, _ := publicWitness.MarshalBinary()

	valReader := strings.NewReader(proofBuf.String())
	proof = groth16.NewProof(ecc.BN254)
	proof.ReadFrom(valReader)

	witness, _ = witness2.New(ecc.BN254.ScalarField())
	err = witness.UnmarshalBinary(dataBin)
	if err != nil {
		return err
	}
	// verify the proof using witness
	err = groth16.Verify(proof, vk, witness)
	if err != nil {
		return err
	}
	return nil
}
