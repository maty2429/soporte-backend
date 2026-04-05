package handlers

// ptrOrDefault returns the dereferenced value of p if non-nil, otherwise def.
func ptrOrDefault[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

// ptrOrStr returns the dereferenced value of p if non-nil, otherwise an empty string.
func ptrOrStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
