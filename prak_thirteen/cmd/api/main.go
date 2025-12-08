package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // регистрирует /debug/pprof/* на DefaultServeMux
	"runtime"

	"github.com/CyberGeo335/prak_thirteen/internal/work"
)

func main() {
	// Включаем профили блокировок/мьютексов (доп. практика).
	enableLocks()

	// Регистрируем наш "тяжёлый" эндпоинт на DefaultServeMux.
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		n := 38 // достаточно тяжело для CPU

		// Ручное измерение времени выполнения.
		defer work.TimeIt("FibFast(38)")()

		res := work.FibFast(n)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(w, "%d\n", res)
	})

	log.Println("Server on :8080; pprof on /debug/pprof/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// enableLocks включает профили блокировок и мьютексов.
func enableLocks() {
	runtime.SetBlockProfileRate(1)     // включить Block profile
	runtime.SetMutexProfileFraction(1) // включить Mutex profile
}
