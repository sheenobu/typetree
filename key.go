package typetree

// A Key is a location in a tree of types.
type Key []string

// Keys converts the list of strings to a Key
func Keys(names ...string) Key {
	return Key(names)
}
