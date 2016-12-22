package typetree

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

type itemContainer map[string]*tree

// A Group is a grouping of related named objects which
// can be instantiated and operated on. Group follows
// the zero value principle.
type Group struct {
	lock  sync.RWMutex
	items itemContainer
}

// Register registers the given object type with the
// name.
func (g *Group) Register(name string, tag string, i interface{}) error {

	if i == nil {
		return errors.New("Type argument is nil")
	}
	if name == "" {
		return errors.New("Name is invalid")
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	// zero value principle support code.
	if g.items == nil {
		g.items = make(itemContainer)
	}

	// try to unwrap the type if it is a pointer
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// ensure the given object is a struct
	if t.Kind() != reflect.Struct {
		return errors.Errorf("Type '%v' is not a struct", t)
	}

	// build the tree
	tr, err := buildTree(tag, t)
	if err != nil {
		return errors.Wrapf(err, "Failed to build type '%v'", name)
	}

	g.items[name] = tr

	return nil
}

// New creates a new instance of the given name
func (g *Group) New(name string) (Instance, error) {
	if name == "" {
		return nil, errors.New("Name is invalid")
	}

	g.lock.RLock()
	defer g.lock.RUnlock()

	// zero value principle support code.
	if g.items == nil {
		return nil, errors.Errorf("Given type '%v' not found", name)
	}

	// find the type tree
	tr, ok := g.items[name]
	if !ok {
		return nil, errors.Errorf("Given type '%v' not found", name)
	}

	// build the zero value
	val := reflect.New(tr.Type)

	// build the type tree isntance
	instance := &treeInstance{
		tree: tr,
		val:  val,
	}

	return instance, nil
}
