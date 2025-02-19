package helper

// GenericBuilder defines a builder for a given type.
type GenericBuilder[T any] struct {
	supplier func() T
	steps    []func(*T)
}

// NewGenericBuilder initializes a new builder with a supplier function.
func NewGenericBuilder[T any](supplier func() T) *GenericBuilder[T] {
	return &GenericBuilder[T]{supplier: supplier}
}

// With adds a configuration step that modifies the built object.
func (b *GenericBuilder[T]) With(consumer func(*T)) *GenericBuilder[T] {
	b.steps = append(b.steps, consumer)
	return b
}

// Build creates the object and applies all configuration steps.
func (b *GenericBuilder[T]) Build() T {
	obj := b.supplier()
	for _, step := range b.steps {
		step(&obj)
	}
	return obj
}

/*
// Example usage:
type Npc struct {
	Name string
}

func main() {

	npcBuilder := NewGenericBuilder(func() Npc {
		return Npc{}
	}).With(func(n *Npc) {
		n.Name = "John Doe"
	})

	npc := npcBuilder.Build()
	fmt.Println(npc.Name)
}
	// usage of updating an existing npc
	npc := Npc{Name: "Jane Doe"}
	npcBuilder := NewGenericBuilder(func() Npc {
		return npc
	}).With(func(n *Npc) {
		n.Name = "John Doe"
	})

	npc = npcBuilder.Build()
	fmt.Println(npc.Name)
}
*/

// EOF
