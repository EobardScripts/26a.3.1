package pipeline

import (
	"2021/pkg/ringbuffer"
	"fmt"
	"io"
	"time"
)

// Handler - функция-обработчик пайплайна.
// Принимает входной канал предыдущего обработчика.
// Возвращает выходной канал обработанных данных.
// После закрытия входного канала и завершения обработки данных,
// обработчик должен закрыть свой выходной канал.
type Handler func(<-chan int) chan int

// PassMin - возвращает обработчик, который пропускает числа не меньше чем min.
func PassMin(min int) Handler {
	return func(in <-chan int) chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for val := range in {
				if val >= min {
					logf("PassMin: значение %d не меньше %d: передано", val, min)
					out <- val
				} else {
					logf("PassMin: значение %d меньше %d: отфильтровано", val, min)
				}
			}
			logf("PassMin: обработчик завершил работу")
		}()
		return out
	}
}

// PassDivBy - возвращает обработчик, который пропускает числа,
// делящиеся на d без остатка. Нулевые значения не пропускаются.
func PassDivBy(d int) Handler {
	return func(in <-chan int) chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for val := range in {
				if val != 0 && val%d == 0 {
					logf("PassDivBy: значение %d делится на %d: передано", val, d)
					out <- val
				} else {
					logf("PassDivBy: значение %d не делится на %d: отфильтровано", val, d)
				}
			}
			logf("PassDivBy: обработчик завершил работу")
		}()
		return out
	}
}

// RingBuffer - возвращает обработчик кольцевого буфера длиной l с таймаутом t.
// При превышении таймаута отправляет все значения из буфера в выходной канал.
func RingBuffer(l int, t time.Duration) Handler {
	return func(in <-chan int) chan int {
		done := make(chan struct{})
		out := make(chan int)
		b := ringbuffer.New(l)

		// Горутина отправляет входящие значения в буфер
		go func() {
			// Как только закрыт входящий канал, закрываем done
			defer close(done)
			for val := range in {
				b.Push(val)
				logf("RingBuffer: значение %d отправлено в буфер", val)
			}
		}()

		// flush - считывает все значения из буфера и отправляет в out
		flush := func() {
			for i := 0; ; i++ {
				v, ok := b.Pop()
				if !ok {
					logf("RingBuffer: отправлено значений из буфера: %d", i)
					break
				}
				for el := range v {
					out <- el
				}
			}
		}

		// Горутина отправляет значения из буфера в выходной канал с интервалом t
		go func() {
			defer close(out)
			for {
				select {
				// Если истек таймаут, отправляем все значения из буфера в out
				case <-time.After(t):
					flush()
					select {
					// Если входной канал обработчика закрыт, завершаем горутину
					case <-done:
						logf("RingBuffer: обработчик завершил работу")
						return
					default:
						break
					}
				}
			}
		}()

		return out
	}
}

// ToWriter - возвращает обработчик, который отправляет значения
// из пайплайна в указанный io.Writer
func ToWriter(w io.Writer) Handler {
	return func(in <-chan int) chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for val := range in {
				if _, err := fmt.Fprintln(w, val); err != nil {
					logf("ToWriter: ошибка записи: %v", err)
				}
				// Передаем значение дальше в пайплайн
				out <- val
			}
		}()
		return out
	}
}

// ToChan - возвращает обработчик, который отправляет значения
// из пайплайна в указанный канал
func ToChan(c chan<- int) Handler {
	return func(in <-chan int) chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for val := range in {
				// Передаем значение в канал
				c <- val
				logf("ToChan: значение %d передано в канал", val)
				// Передаем значение дальше в пайплайн
				out <- val
			}
		}()
		return out
	}
}
