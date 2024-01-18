package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
)

// Текущая директория
var dirPath, _ = os.Getwd()

// Имя файла с логом
var fileName = "log.txt"

// Имя папки с логом
var dirName = "logger"

// Debug - выводить ли отладочные сообщения.
var Debug = true

// DebugWriter - дескриптор для вывода отладочных сообщений.
var DebugWriter = os.Stderr

// Проверяем существует ли файл с логом
// true - существует
// false - не существует
func logIsExist() bool {
	path := filepath.Join(dirPath, dirName, fileName)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	} else {
		return true
	}

	return false
}

// Создаем папку для лога и сам файл лога
func createDirAndFile() *os.File {
	err := os.MkdirAll(filepath.Join(dirPath, dirName), 777)
	if err != nil {
		logf("Ошибка при создании папки для лога: %w", err)
		return nil
	}
	file, err := os.OpenFile(filepath.Join(dirPath, dirName, fileName), os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		logf("Ошибка при создании файла лога: %w", err)
		return nil
	}

	return file
}

// Получаем файл для записи
func getFileLog() *os.File {
	if logIsExist() {
		file, err := os.OpenFile(filepath.Join(dirPath, dirName, fileName), os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			logf("Ошибка при открытии файла лога: %w", err)
			return nil
		}

		return file
	} else {
		return createDirAndFile()
	}
}

// Запись лога в файл
func writeToLog(log string) {
	file := getFileLog()
	defer file.Close()
	if file != nil {
		_, err := file.WriteString(log)
		if err != nil {
			logf("Ошибка при записи в лог: %w", err)
			return
		}
	}
}

// logf - выводит отформатированное отладочное сообщение в DebugWriter.
func logf(format string, a ...interface{}) {
	if Debug {
		_, _ = fmt.Fprintf(DebugWriter, format+"\n", a)
		writeToLog(fmt.Sprintf(format+"\n", a))
	}
}
