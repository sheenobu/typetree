# typetree

typetree is a library for operating on arbitrary type trees. typetree
aims to hide the reflection magic.

## Usage

Some initial structures, `SubType` and `MyType`:

```go
type SubType struct {
    Count int64 `mytag:"count"`
    Items []string `mytag:"items"`
}

type MyType struct {
    Name    string `mytag:"name"`
    SubType *SubType `mytag:"sub_type"`
}
```

Build the typetree group and register our type. We
do not register `SubType` as the group will descend into `MyType`
to inspect it.

```go
var group typetree.Group
group.Register("myTypeName", "mytag", &myTypeName{})
```

We create an instance using the predetermined name. What we get back
is a `typetree.Instance`, which wraps our region of memory into 
a set of operations:

```go
var instance typetree.Instance
instance, _ = group.New("myTypeName")
```

If we were to go ahead and cast to `MyType`, we would get an empty, but
valid structure:

```go 
// IMPORTANT: Both instance and myType are sharing a region of memory 
myType, _ = instance.Interface().(*MyType)
```

We can set and get the field values via the `instance`. We
can find the fields via `typetree.Keys`:

```go
// set the name field
err := instance.Set(typetree.Keys("Name"), "value") 

// get the name
var value typetree.Value
value, err = instance.Get(typetree.Keys("name"))
valueString, err := value.Interface().(string)
```

Nested element operations work the same way:

```go
// Set a nested element. The SubType pointer will be constructed
err := instance.Set(typetree.Keys("sub_type", "count"), 12)
```

Additionally, our `myType` variable from above gets populated as 
it points to the same region of memory:

```go
// this should be true 
myType.SubType.Count == 12
```

We can also append to arrays. The boolean return value will
be false if appending failed due to a missing or wrongly typed field:

```go
// Append to an array
value, err = instance.Get(typetree.Keys("sub_type", "items"))
err := value.Append("Hello")
```