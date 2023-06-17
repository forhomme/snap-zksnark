package models

import (
	"github.com/consensys/gnark/frontend"
	mimc2 "github.com/consensys/gnark/std/hash/mimc"
)

type Circuit struct {
	Secret frontend.Variable
	Hash   frontend.Variable `gnark:",public"`
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
