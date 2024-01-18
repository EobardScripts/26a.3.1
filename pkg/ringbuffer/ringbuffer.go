package ringbuffer

import (
	"sync"
)

// RingBuffer Кольцевой буфер
type RingBuffer struct {
	data []*Data
	// размер буфера
	size int
	// Последний элемент в буфере
	lastInsert int
	// Следующий элемент для чтения из буфера
	nextRead int
	// Мьютекс для потокобезопасного чтения/записи
	mutex sync.RWMutex
}

// Data Отдельная структура для хранения данных (для удобства)
type Data struct {
	Value int
}

func New(size int) *RingBuffer {
	//страховка на случай, если размер буфера будет указан меньше 1
	if size < 1 {
		size = 1
	}

	return &RingBuffer{
		data:       make([]*Data, size),
		size:       size,
		lastInsert: -1,
		mutex:      sync.RWMutex{},
	}
}

// Push Добавление элемента в буфер
func (r *RingBuffer) Push(input int) {
	newData := &Data{Value: input}
	r.mutex.Lock()
	r.lastInsert = (r.lastInsert + 1) % r.size
	r.data[r.lastInsert] = newData

	if r.nextRead == r.lastInsert {
		r.nextRead = (r.nextRead + 1) % r.size
	}
	r.mutex.Unlock()
}

// Pop Извлечение элементов из буфера
func (r *RingBuffer) Pop() ([]int, bool) {
	defer r.mutex.Unlock()
	var output []int
	mark := true
	for {
		r.mutex.Lock()
		if r.data[r.nextRead] != nil {
			output = append(output, r.data[r.nextRead].Value)
			r.data[r.nextRead] = nil
		}
		if r.nextRead == r.lastInsert || r.lastInsert == -1 {
			break
		}
		r.nextRead = (r.nextRead + 1) % r.size
		r.mutex.Unlock()
	}

	if len(output) == 0 {
		mark = false
	}
	return output, mark
}
