package work

// Неоптимальный рекурсивный Фибоначчи — демонстрация CPU-нагрузки.
func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

// Оптимизированная итеративная версия.
// Используется для "после оптимизации".
func FibFast(n int) int {
	if n < 2 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
