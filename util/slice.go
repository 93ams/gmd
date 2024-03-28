package util

func Apply[T any](s T, opts []func(T)) T {
	for _, f := range opts {
		f(s)
	}
	return s
}
