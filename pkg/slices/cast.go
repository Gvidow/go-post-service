package slices

func ToPointers[S ~[]E, P []*E, E any](s S) P {
	res := make([]*E, len(s))
	for ind, elem := range s {
		res[ind] = &elem
	}
	return res
}
