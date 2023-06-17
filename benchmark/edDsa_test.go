package main

import "testing"

func BenchmarkCreateEddsaCircuit(b *testing.B) {
	b.Run("Test CreateEDDSACircuit", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			CreateEddsaCircuit()
		}
	})
}
