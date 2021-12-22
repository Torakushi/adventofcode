package day22

import "testing"

func BenchmarkFirstPart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		process(true)
	}
}
func BenchmarkSecondPart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		process(false)
	}
}
