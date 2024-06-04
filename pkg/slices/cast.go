package slices

func ToPointers[S ~[]E, P []*E, E any](s S) P {
	res := make([]*E, len(s))
	for ind := range s {
		res[ind] = &s[ind]
	}
	return res
}
