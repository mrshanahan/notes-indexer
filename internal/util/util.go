package util

func Map[T any, S any](xs []T, f func(T) S) []S {
	ys := make([]S, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}
