package utils

import "fmt"

// FormatBytes 格式化字节大小
func FormatBytes(bytes float64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes < KB:
		return fmt.Sprintf("%.2f B", bytes)
	case bytes < MB:
		return fmt.Sprintf("%.2f KB", bytes/KB)
	case bytes < GB:
		return fmt.Sprintf("%.2f MB", bytes/MB)
	case bytes < TB:
		return fmt.Sprintf("%.2f GB", bytes/GB)
	default:
		return fmt.Sprintf("%.2f TB", bytes/TB)
	}
}
