package commands

func wrap(line string, limit int) string {
	if limit <= 3 {
		return line
	}
	wrapped := []rune(line)
	if len(wrapped) > limit {
		wrapped = append(wrapped[:limit-3], []rune("...")...)
	}
	return string(wrapped)
}
