package main

import (
	"strings"

	"github.com/fatih/structs"
)

// TODO: clean up the ordering of these or split out into multiple files

// MemStore represents the data structure where metadata will be stored
type MemStore map[string]*metadata

// Storage provides methods to read and write to a store (type MemStore)
type Storage struct {
	store MemStore
}

// toMap converts a value of type metadata to type map[string]interface{} with the same nesting structure
// this facilitates searching each piece of metadata by making it possible to index it by value
func (md metadata) toMap() map[string]interface{} {
	return structs.Map(md)
}

// storage is used throughout application to persist metadata
var storage = Storage{
	store: make(MemStore),
}

/**~ Storage Methods ~**/
func (s Storage) set(k string, v *metadata) {
	s.store[k] = v
}

func (s Storage) get(q map[string]string) []*metadata {
	return searchMetadata(s.store, q)
}

/**~ Helper Methods ~**/
func searchMetadata(ms MemStore, q map[string]string) []*metadata {
	results := []*metadata{}

	// check each piece of saved metadata for each key/value pair of q, returning
	// only those that satisfy all pairs
	for _, metadata := range ms {
		m := metadata.toMap()
		contains := true

		for k, v := range q {
			keys := strings.Split(k, ",")
			contains = containsNestedData(m, keys, v)

			if !contains {
				break
			}
		}

		if contains {
			results = append(results, metadata)
		}
	}

	return results
}

/*
 containsNestedData recursively traverses a generic variable m to determine if it contains a value of type string at
 a given key, whose location is described by the path variable, a slice of string representing the order of keys to access;
 keys that reside in lists, such as Maintainer, can be specified independently of the index of the target data. For example,
 to access any maintainers email, the path would simply be: ["maintainers", "email"]
*/
func containsNestedData(m interface{}, path []string, value string) bool {
	if len(path) == 0 {
		return false
	}

	key := strings.Title(strings.ToLower(path[0]))

	if mapData, isMap := m.(map[string]interface{}); isMap {
		// here we know its a map, but don't know the type of value, so we must check before accessing it
		v := mapData[key]

		stringValue, ok := v.(string)

		// if value is a string, see if it contains the value we're searching for
		if ok {
			return strings.Contains(stringValue, value)
		}

		// if value is anything besides a string, pass it to another function call with the next key in the path
		// TODO: this works for metadata structure specified in instructions (plus nested maps like Spec.Replicas), but breaks when a key points to an
		// array of strings ie:
		// fruits:
		//   - apple
		//   - orange
		return containsNestedData(v, path[1:], value)
	}

	if sliceData, isSlice := m.([]interface{}); isSlice {
		// since we know m is a slice, iterate over each element and check if it contains the current key and desired value
		for _, elem := range sliceData {
			match := containsNestedData(elem, []string{key}, value)

			if match {
				return true
			}
		}
	}

	return false
}
