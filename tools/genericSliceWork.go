package tools

// Equal represents a type that can call the method Equal/1 on some other type T
type Equal[T any] interface {
	Equal(T) bool
}

// GetUniqueFromSlice returns a list of elements from candidates that do not appear in incumbents.
func GetUniqueFromSlice[T any, X Equal[T]](incumbents []T, candidates []X) []X {
	var uniqueCandidates []X

	if len(incumbents) > 0 {
		for _, candidate := range candidates {
			found := false
			for _, incumbent := range incumbents {
				if candidate.Equal(incumbent) {
					found = true
					break
				}
			}

			if !found {
				uniqueCandidates = append(uniqueCandidates, candidate)
			}
		}

		return uniqueCandidates
	}

	return candidates
}

// GetDuplicatesFromSlice returns a list of elements from candidates ⊆ incumbents.
func GetDuplicatesFromSlice[T any, X Equal[T]](incumbents []T, candidates []X) []X {
	var duplicates []X

	if len(incumbents) > 0 {
		for _, candidate := range candidates {
			found := false
			for _, incumbent := range incumbents {
				if candidate.Equal(incumbent) {
					found = true
					break
				}
			}

			if found {
				duplicates = append(duplicates, candidate)
			}
		}
	}

	return duplicates
}

// SliceContains returns true if elem ⊆ slice
func SliceContains[T any, X Equal[T]](slice []T, elem X) bool {
	for _, sliceElem := range slice {
		if elem.Equal(sliceElem) {
			return true
		}
	}

	return false
}

// SliceContainsSet returns true if set ⊆ slice
func SliceContainsSet[T any, X Equal[T]](slice []T, set []X) bool {
	for _, elem := range set {
		if !SliceContains(slice, elem) {
			return false
		}
	}

	return true
}

// SlicesEqual returns true if a == b
func SlicesEqual[T any, X Equal[T]](a []T, b []X) bool {
	// For slices to be equal, they have to have the same number of elements
	if len(a) != len(b) {
		return false
	}

	return SliceContainsSet(a, b)
}
