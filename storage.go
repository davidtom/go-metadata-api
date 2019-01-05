package main

// MemStore represents the data structure where metadata will be stored
type MemStore map[string]*metadata

// Storage provides methods to read and write to a store (type MemStore)
type Storage struct {
	store MemStore
}

var storage = Storage{
	make(MemStore),
}

func (s Storage) set(k string, v *metadata) {
	s.store[k] = v
}

func (s Storage) get() []*metadata {
	results := []*metadata{}

	for _, v := range s.store {
		results = append(results, v)
	}

	return results
}
