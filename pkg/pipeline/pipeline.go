package pipeline

import (
	"bufio"
	"io"
	"strconv"
)

// Pipeline - пайплайн обработки целых чисел
type Pipeline struct {
	in   chan int      // Входной канал пайплайна
	done chan struct{} // Сигнальный канал завершения работы пайплайна
}

// NewPipe - создает новый пайплайн из цепочки обработчиков
func NewPipe(handlers ...Handler) *Pipeline {
	in := make(chan int)
	done := make(chan struct{})

	// Перебираем обработчики и соединяем вход обработчика с выходом предыдущего обработчика.
	// Вход первого обработчика соединен со входным каналом пайплайна.
	out := in
	for _, h := range handlers {
		out = h(out)
	}

	// Читаем сообщения из выхода последнего обработчика в цепочке,
	// ждем его закрытия и сигнализируем о завершении работы пайплайна
	go func() {
		for range out {
		}
		logf("Все обработчики пайплайна завершили работу")
		close(done)
	}()

	logf("Пайплайн создан. Кол-во обработчиков: %d", len(handlers))
	return &Pipeline{in, done}
}

// Emit - посылает в пайплайн входные значения.
// Возвращает сигнальный канал, который будет закрыт
// после отправки в пайплайн последнего значения.
func (p *Pipeline) Emit(values ...int) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for _, v := range values {
			logf("Emit: получено значение %d", v)
			p.in <- v
		}
		logf("Emit: передача в пайплайн завершена")
	}()
	return done
}

// EmitFromReader - посылает в пайплайн входные значения из io.Reader построчно.
// Если строка не является целым числом, то она игнорируется.
// Возвращает сигнальный канал, который будет закрыт
// после отправки в пайплайн последнего значения.
func (p *Pipeline) EmitFromReader(r io.Reader) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			text := scanner.Text()
			v, err := strconv.Atoi(text)
			if err != nil {
				logf("EmitFromReader: `%s` не является целым числом", text)
			} else {
				logf("EmitFromReader: получено значение %d", v)
				p.in <- v
			}
		}
		if err := scanner.Err(); err != nil {
			logf("EmitFromReader: ошибка чтения: %v", err)
		}
		logf("EmitFromReader: передача в пайплайн завершена")
	}()
	return done
}

// EmitFromChan - посылает в пайплайн входные значения из канала.
// Возвращает сигнальный канал, который будет закрыт
// после отправки в пайплайн последнего значения.
func (p *Pipeline) EmitFromChan(in <-chan int) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for v := range in {
			logf("EmitFromChan: получено значение %d", v)
			p.in <- v
		}
		logf("EmitFromChan: передача в пайплайн завершена")
	}()
	return done
}

// Close - закрывает входной канал пайплайна.
// После этого, все обработчики должны последовательно завершить обработку и
// закрыть свои выходные каналы. После закрытия выходного канала последнего обработчика
// будет закрыт сигнальный канал done.
func (p *Pipeline) Close() {
	logf("Close: пайплайн закрыт")
	close(p.in)
}

// Done - возвращает сигнальный канал окончания работы пайплайна.
// Канал закроется после того, как будет вызван Close
// и последний обработчик в цепочке завершит обработку.
func (p *Pipeline) Done() <-chan struct{} {
	return p.done
}
