package utils

// Truncate обрезает строку до указанной длины
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// FormatFileSize форматирует размер файла в человекочитаемый вид
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return string(rune(bytes)) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return string(rune(bytes/div)) + string("KMGTPE"[exp]) + "B"
}
