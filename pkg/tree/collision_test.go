 package tree

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNoFlattenedTagCollisions walks every struct reachable from the Tree and
// asserts that no two fields contribute the same key to the same serialized
// object.
//
// structToValueObject flattens embedded structs into their parent's entry map
// (so a Deployment's replicas sit alongside the Workload fields it embeds), and
// that flattening writes straight into the map. Two fields tagged the same name
// — one declared on the type, one arriving through an embedded struct — would
// silently overwrite each other, and the survivor depends on declaration order.
// Nothing else in the conversion path reports it, and a round-trip test would
// not catch it either: the value that survives round-trips perfectly well, it is
// just the wrong field.
//
// This is a real hazard for the Kubernetes types in particular, where
// meta.ObjectMeta contributes a flattened "name" and Container separately has
// one of its own. Those two do not collide, because containers are a nested list
// rather than an embedded struct — but a top-level Name added to a kind later
// would collide, and this test is what would report it.
func TestNoFlattenedTagCollisions(t *testing.T) {
	seen := map[reflect.Type]bool{}
	assertNoCollisions(t, reflect.TypeFor[Tree](), seen)
}

// assertNoCollisions checks one struct type's own entry map, then recurses into
// every type it can reach. seen guards against infinite recursion on any
// self-referential type.
func assertNoCollisions(t *testing.T, typ reflect.Type, seen map[reflect.Type]bool) {
	t.Helper()

	typ = elem(typ)
	if typ.Kind() != reflect.Struct || seen[typ] {
		return
	}
	seen[typ] = true

	// owner records which field produced each key, so a failure names both
	// sides of the collision rather than just the key.
	owner := map[string]string{}
	collectKeys(t, typ, typ.String(), owner)

	for i := range typ.NumField() {
		assertNoCollisions(t, typ.Field(i).Type, seen)
	}
}

// collectKeys adds typ's serialized keys to owner, following the same rules as
// structToValueObject: embedded structs (other than the specially-handled
// resource.Resource) flatten into the parent, and everything else contributes
// its tree tag. path is the chain of embedded types the key arrived through,
// used to describe a collision.
func collectKeys(t *testing.T, typ reflect.Type, path string, owner map[string]string) {
	t.Helper()

	for i := range typ.NumField() {
		field := typ.Field(i)

		if field.Anonymous && field.Type.Kind() == reflect.Struct && field.Type != resourceType {
			collectKeys(t, field.Type, path+" -> "+field.Type.String(), owner)
			continue
		}

		key := field.Tag.Get("tree")
		if key == "" || key == "-" {
			continue
		}

		from := path + "." + field.Name
		if prev, ok := owner[key]; ok {
			assert.Failf(t, "duplicate tree tag",
				"tree tag %q is contributed twice to the same object: by %s and by %s. "+
					"Flattening writes both into one entry map, so one silently overwrites the other.",
				key, prev, from)
			continue
		}
		owner[key] = from
	}
}

// elem unwraps the containers the converter recurses through — pointers, and
// the slices/arrays/maps that hold nested structs — so a []Container is checked
// as a Container. Anything that is not a struct falls out of the walk in
// assertNoCollisions.
func elem(typ reflect.Type) reflect.Type {
	for {
		switch typ.Kind() {
		case reflect.Pointer, reflect.Slice, reflect.Array, reflect.Map:
			typ = typ.Elem()
		default:
			return typ
		}
	}
}