package convert

// Slice converts a slice of T's to a slice of U's.
func Slice[T, U any](v []T, fn func(T) U) []U {
	if len(v) == 0 {
		return []U{}
	}

	out := make([]U, 0, len(v))
	for i := range v {
		out = append(out, fn(v[i]))
	}

	return out
}
