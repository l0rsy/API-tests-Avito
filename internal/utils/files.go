package utils

import (
	"runtime"
)

// GetProjectRoot возвращает абсолютный путь до корня проекта
func GetProjectRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("Не удалось определить путь до файла")
	}
	return findGoModDir(file)
}
