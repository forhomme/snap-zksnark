package models

import "github.com/consensys/gnark/frontend"

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
