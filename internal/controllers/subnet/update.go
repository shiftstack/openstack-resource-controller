/*
Copyright 2024 The ORC Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package subnet

import "slices"

func updateFunc[T1, T2 any](currentValue T1, newValue T2, equalFunc func(T1, T2) bool) *T2 {
	if !equalFunc(currentValue, newValue) {
		return &newValue
	}
	return nil
}

func update[T comparable](currentValue T, newValue T) *T {
	return updateFunc(currentValue, newValue, func(v1, v2 T) bool { return v1 == v2 })
}

// updateSetFunc returns a pointer to newValue if it is semantically different than currentValue.
// Use updateVectorFunc when comparing two slices whose order is relevant.
func updateVectorFunc[T1, T2 any](currentValue []T1, newValue []T2, equalFunc func(T1, T2) bool) *[]T2 {
	if len(currentValue) != len(newValue) {
		return nil
	}

	sliceEqualFunc := func(slice1 []T1, slice2 []T2) bool {
		for i := range slice1 {
			if !equalFunc(slice1[i], slice2[i]) {
				return false
			}
		}
		return true
	}

	return updateFunc(currentValue, newValue, sliceEqualFunc)
}

func updateVector[T comparable](currentValue, newValue []T) *[]T {
	equalFunc := func(v1, v2 T) bool { return v1 == v2 }
	return updateVectorFunc(currentValue, newValue, equalFunc)
}

// updateSetFunc returns a pointer to newValue if it is semantically different than currentValue.
// Use updateSetFunc when comparing two slices whose order is irrelevant, and whose duplicate values are irrelevant.
func updateSetFunc[T1, T2 any](currentValue []T1, newValue []T2, equalFunc func(T1, T2) bool) *[]T2 {
	sliceEqualFunc := func(slice1 []T1, slice2 []T2) bool {
		for i := range slice1 {
			if !slices.ContainsFunc(slice2, func(e T2) bool { return equalFunc(slice1[i], e) }) {
				return false
			}
		}
		for i := range slice2 {
			if !slices.ContainsFunc(slice1, func(e T1) bool { return equalFunc(e, slice2[i]) }) {
				return false
			}
		}
		return true
	}

	return updateFunc(currentValue, newValue, sliceEqualFunc)
}
