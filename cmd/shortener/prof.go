package main

import (
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func runProf() {
	// Записываем 30 секунд
	timer := time.NewTimer(time.Second * 30)
	go func() {
		<-timer.C

		// создаём файл журнала профилирования памяти
		fmem, err := os.Create(`profiles/result.pprof`)
		if err != nil {
			panic(err)
		}
		defer fmem.Close()
		runtime.GC() // получаем статистику по использованию памяти
		if err := pprof.WriteHeapProfile(fmem); err != nil {
			panic(err)
		}
	}()
}
