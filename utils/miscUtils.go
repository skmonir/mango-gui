package utils

import "fmt"

func ConvertMemoryInMb(memory uint64) uint64 {
	return memory / 1024 / 1024
}

func ParseMemoryInMb(memory uint64) string {
	return fmt.Sprintf("%v MB", memory/1024/1024)
}

func ParseMemoryInKb(memory uint64) string {
	return fmt.Sprintf("%v KB", memory/1024)
}
