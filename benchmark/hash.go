package main

import (
	"bytes"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	witness2 "github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	r1cs2 "github.com/consensys/gnark/frontend/cs/r1cs"
	mimc2 "github.com/consensys/gnark/std/hash/mimc"
	"log"
	"smart-contract-service/internal"
	models2 "smart-contract-service/models/circuit"
	"strings"
)

type Circuit struct {
	Secret frontend.Variable
	Hash   frontend.Variable `gnark:",public"`
}

type Data struct {
	PreImage string
}

func (circuit *Circuit) Define(api frontend.API) error {
	// hash function
	mimc, err := mimc2.NewMiMC(api)
	if err != nil {
		return err
	}
	mimc.Write(circuit.Secret)
	api.AssertIsEqual(circuit.Hash, mimc.Sum())
	return nil
}

func CreateHashCircuit() (err error) {
	var circuit Circuit

	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs2.NewBuilder, &circuit)
	if err != nil {
		log.Fatal(err)
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		log.Fatal(err)
	}

	internal.Serialize(r1cs, internal.R1csPath)
	internal.Serialize(pk, internal.PkPath)
	internal.Serialize(vk, internal.VkPath)

	// read R1CS, proving key and verifying keys
	ccs := groth16.NewCS(ecc.BN254)
	pk = groth16.NewProvingKey(ecc.BN254)
	vk = groth16.NewVerifyingKey(ecc.BN254)
	internal.Deserialize(ccs, internal.R1csPath)
	internal.Deserialize(pk, internal.PkPath)
	internal.Deserialize(vk, internal.VkPath)

	assignment := &models2.Circuit{}
	b := make([]byte, 32)
	preImage := &Data{PreImage: "data"}
	preImageByte := []byte(fmt.Sprintf("%v", preImage))
	copy(b, preImageByte)
	hash := internal.MimcHash(b)

	assignment.Secret = frontend.Variable(b)
	assignment.Hash = frontend.Variable(hash)

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

func main() {
	CreateHashCircuit()
}
