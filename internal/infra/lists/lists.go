package lists

type List[T any] []T

type MapFunc[T, V any] func(item T) V

func Map[L ~[]T, R List[V], T, V any](l L, f MapFunc[T, V]) R {
	if len(l) == 0 {
		return nil
	}

	r := make(R, len(l))
	for i, t := range l {
		r[i] = f(t)
	}

	return r
}
