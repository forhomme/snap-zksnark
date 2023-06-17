package main

import "testing"

func BenchmarkCreateEllipticCircuit(b *testing.B) {
	b.Run("Test CreateEllipticCircuit", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			CreateEllipticCircuit()
		}
	})
}
