package kendra

import (
	"strings"
)

type Kendra struct {
	Indexes     []Index      `tree:"indexes"`
	DataSources []DataSource `tree:"data_sources"`
}

func (k *Kendra) PostProcess() {
	for i := range k.DataSources {
		k.DataSources[i].Relationships = DataSourceRelationships{}
	}

	// Match on index ID, or name when the reference resolves to it.
	// The data source with the lowest address on each index is its first one,
	// so the result does not depend on parse order.
	first := map[*Index]int{}
	for i, ds := range k.DataSources {
		ref := ds.IndexID.Value()
		if ref == "" {
			continue
		}
		for j := range k.Indexes {
			idx := &k.Indexes[j]
			if idx.ID == ref || strings.Contains(ref, idx.ID) || (!idx.Name.IsEmpty() && idx.Name.Value() == ref) {
				k.DataSources[i].Relationships.Index = idx
				if f, ok := first[idx]; !ok || addressOf(&k.DataSources[i]) < addressOf(&k.DataSources[f]) {
					first[idx] = i
				}
				break
			}
		}
	}
	for _, i := range first {
		k.DataSources[i].Relationships.FirstForIndex = true
	}
}

func addressOf(ds *DataSource) string {
	if ds.Definition.Address == nil {
		return ds.ID
	}
	return ds.Definition.Address.String()
}
