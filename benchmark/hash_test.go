package main

import "testing"

func BenchmarkCreateHashCircuit(b *testing.B) {
	b.Run("Test CreateHashCircuit", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			CreateHashCircuit()
		}
	})
}
