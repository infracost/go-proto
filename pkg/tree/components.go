package tree

import "reflect"

// relationshipsField is the conventional name of the struct field, present on
// resource types that link to other resources, that holds those links (e.g.
// ec2.Instance.Relationships). It is the only place cross-resource pointers
// live, so it is the only field Components walks. The field is tagged `tree:"-"`
// so it is never serialized.
const relationshipsField = "Relationships"

// Component is a set of resources whose costs are coupled: each resource is
// reachable from the others by following Relationships links (in either
// direction). Because a resource's cost function can only observe other
// resources through its Relationships, modifying one resource can only change
// the cost of resources within the same component. This lets callers (e.g. the
// savings calculation) recompute costs over a single component instead of the
// whole tree.
type Component []Resource

// Addresses returns the string form of every resource address in the component,
// in the component's resource order. Resources with no address are skipped.
func (c Component) Addresses() []string {
	out := make([]string, 0, len(c))
	for _, r := range c {
		if addr := r.GetBase().Definition.Address; addr != nil {
			out = append(out, addr.String())
		}
	}
	return out
}

// Components partitions the tree's supported resources into connected
// components based on their Relationships. Two resources share a component when
// one references the other — directly or transitively — through a Relationships
// field. A resource with no relationships forms a component of one.
//
// Relationships must already be linked, so call PostProcess before Components
// (FromProto callers get this via the standard round-trip). Resources are
// matched by Definition.Address, so relationship targets held by value (e.g.
// []T fields) are treated identically to pointer targets. Resources without an
// address cannot be matched and are returned as singleton components.
//
// The result is deterministic: components appear in the order their first
// resource appears in ToResources, and resources within a component keep that
// same order.
func (t *Tree) Components() []Component {
	resources := t.ToResources(false)
	n := len(resources)

	uf := newUnionFind(n)

	// Index resources by address so relationship links (which may be value
	// copies, not pointers into the canonical slices) can be resolved back to
	// the resource they name.
	indexByAddress := make(map[string]int, n)
	for i, r := range resources {
		if addr := r.GetBase().Definition.Address; addr != nil {
			indexByAddress[addr.String()] = i
		}
	}

	for i, r := range resources {
		rel := reflect.ValueOf(r).Elem().FieldByName(relationshipsField)
		if !rel.IsValid() {
			continue
		}
		forEachRelatedAddress(rel, func(addr string) {
			if j, ok := indexByAddress[addr]; ok {
				uf.union(i, j)
			}
		})
	}

	// Group resources by their component root, preserving first-seen order.
	groups := make(map[int]Component, n)
	order := make([]int, 0, n)
	for i, r := range resources {
		root := uf.find(i)
		if _, seen := groups[root]; !seen {
			order = append(order, root)
		}
		groups[root] = append(groups[root], r)
	}

	out := make([]Component, 0, len(order))
	for _, root := range order {
		out = append(out, groups[root])
	}
	return out
}

// forEachRelatedAddress walks a Relationships struct value and calls visit with
// the address of every resource it links to, across *T, []*T and []T fields.
func forEachRelatedAddress(rel reflect.Value, visit func(string)) {
	for i := range rel.NumField() {
		if !rel.Type().Field(i).IsExported() {
			continue
		}
		collectResourceAddresses(rel.Field(i), visit)
	}
}

func collectResourceAddresses(v reflect.Value, visit func(string)) {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return
		}
		if addr, ok := addressOf(v); ok {
			visit(addr)
		}
	case reflect.Slice:
		for i := range v.Len() {
			collectResourceAddresses(v.Index(i), visit)
		}
	case reflect.Struct:
		// A value (non-pointer) relationship target — GetBase has a pointer
		// receiver, so take the element's address to satisfy the interface.
		if v.CanAddr() {
			if addr, ok := addressOf(v.Addr()); ok {
				visit(addr)
			}
		}
	}
}

// addressOf returns the address string of v if it is a resource with a
// non-nil address.
func addressOf(v reflect.Value) (string, bool) {
	if !v.CanInterface() {
		return "", false
	}
	r, ok := v.Interface().(Resource)
	if !ok {
		return "", false
	}
	base := r.GetBase()
	if base == nil || base.Definition.Address == nil {
		return "", false
	}
	return base.Definition.Address.String(), true
}

// unionFind is a disjoint-set structure with path halving and union by size.
type unionFind struct {
	parent []int
	size   []int
}

func newUnionFind(n int) *unionFind {
	uf := &unionFind{parent: make([]int, n), size: make([]int, n)}
	for i := range uf.parent {
		uf.parent[i] = i
		uf.size[i] = 1
	}
	return uf
}

func (uf *unionFind) find(x int) int {
	for uf.parent[x] != x {
		uf.parent[x] = uf.parent[uf.parent[x]] // path halving
		x = uf.parent[x]
	}
	return x
}

func (uf *unionFind) union(a, b int) {
	ra, rb := uf.find(a), uf.find(b)
	if ra == rb {
		return
	}
	if uf.size[ra] < uf.size[rb] {
		ra, rb = rb, ra
	}
	uf.parent[rb] = ra
	uf.size[ra] += uf.size[rb]
}
