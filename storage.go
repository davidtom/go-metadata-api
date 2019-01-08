package main

import (
	"strings"

	"github.com/fatih/structs"
)

// MemStore represents the data structure where metadata will be stored
type MemStore map[string]*metadata

// Storage provides methods to read and write to a store (type MemStore)
type Storage struct {
	store MemStore
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

// searchMetadata searches all metadata in a MemStore for any that match all query parameters supplied
func searchMetadata(ms MemStore, q map[string]string) []*metadata {
	results := []*metadata{}

	// check each piece of saved metadata for each key/value pair of q, returning
	// only those that satisfy all pairs
	for _, metadata := range ms {
		m := metadata.toMap()
		contains := true

		for k, v := range q {
			keys := strings.Split(k, ",")
			contains = metadataContainsValue(m, keys, v)

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
func metadataContainsValue(m interface{}, path []string, value string) bool {
	if len(path) == 0 {
		return false
	}

	key := strings.Title(strings.ToLower(path[0]))

	if mapData, isMap := m.(map[string]interface{}); isMap {
		// here we know its a map, but don't know the type of value, so we must check before accessing it
		v := mapData[key]

		// we will handle both strings and slice of strings here, so create a variable to use in both cases
		tempSlice := []string{}

		if sliceValue, isSliceString := v.([]string); isSliceString {
			tempSlice = sliceValue
		} else if stringValue, isString := v.(string); isString {
			tempSlice = []string{stringValue}
		}

		for _, val := range tempSlice {
			match := strings.Contains(strings.ToLower(val), strings.ToLower(value))

			if match {
				return true
			}
		}

		// if value is anything besides a string or slice of string, pass it to another function call with the next key in the path
		return metadataContainsValue(v, path[1:], value)
	}

	// if m is not a map, it must be a slice; pass each value in it back to this function with the current key and check return values
	if sliceData, isSlice := m.([]interface{}); isSlice {
		for _, elem := range sliceData {
			match := metadataContainsValue(elem, []string{key}, value)

			if match {
				return true
			}
		}
	}

	return false
}

// toMap converts a value of type metadata to type map[string]interface{} with the same nesting structure
// this facilitates searching each piece of metadata by making it possible to index it by value
func (md metadata) toMap() map[string]interface{} {
	return structs.Map(md)
}
