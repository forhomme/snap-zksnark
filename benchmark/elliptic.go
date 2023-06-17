package main

import (
	"bytes"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	witness2 "github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	r1cs2 "github.com/consensys/gnark/frontend/cs/r1cs"
	"log"
	"smart-contract-service/internal"
	"strings"
)

type EllipticCurve struct {
	X frontend.Variable
	Y frontend.Variable `gnark:",public"`
}

func (circuit *EllipticCurve) Define(api frontend.API) error {
	// compute x**3 and store it in the local variable x3.
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)

	// compute x**3 + x + 5 and store it in the local variable res
	res := api.Add(x3, circuit.X, 5)

	// assert that the statement x**3 + x + 5 == y is true.
	api.AssertIsEqual(circuit.Y, res)
	return nil
}

func CreateEllipticCircuit() (err error) {
	var circuit EllipticCurve

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

	assignment := &EllipticCurve{
		X: 3,
		Y: 35,
	}

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
