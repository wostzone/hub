package kvmem

import "sort"

// Map2SortedKeys returns the map keys as a list sorted by key
func Map2SortedKeys[V any](mapInput map[string]V) []string {
	keys := make([]string, 0, len(mapInput))
	for key := range mapInput {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Map2SortedValues returns the map values as a list sorted by key
//func Map2SortedValues[V any](mapInput map[string]V) []V {
//	keys := make([]string, 0, len(mapInput))
//	for key := range mapInput {
//		keys = append(keys, key)
//	}
//	sort.Strings(keys)
//	listOutput := make([]V, 0, len(keys))
//	for _, key := range keys {
//		listOutput = append(listOutput, mapInput[key])
//	}
//	return listOutput
//}
