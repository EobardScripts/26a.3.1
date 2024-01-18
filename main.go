package main

import (
	"2021/pkg/pipeline"
	"os"
	"time"
)

func main() {
	// Создаем пайплайн
	p := pipeline.NewPipe(
		pipeline.PassMin(0),                   // Фильтруем отрицательные значения
		pipeline.PassDivBy(3),                 // Фильтруем значения не кратные 3
		pipeline.RingBuffer(2, time.Second*5), // Кольцевой буфер на 2 элемента с таймаутом 5 сек
		pipeline.ToWriter(os.Stdout),          // Выводим в stdout
	)

	// Читаем из stdin и ожидаем завершения чтения
	<-p.EmitFromReader(os.Stdin)

	// Тк мы больше ничего не собираемся передавать в пайплайн, то закрываем его
	p.Close()

	// Дожидаемся завершения всех обработчиков пайплайна
	<-p.Done()
}
