package work

import "testing"

// Бенчмарк для медленной рекурсивной функции.
func BenchmarkFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Fib(30) // 30, чтобы бенч не жил вечность
	}
}

// Бенчмарк для быстрой итеративной версии.
func BenchmarkFibFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FibFast(30)
	}
}
