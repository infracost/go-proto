package address

import (
	"github.com/hashicorp/hcl/v2"
	parserpb "github.com/infracost/proto/gen/go/infracost/parser"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// FromHCLTraversals creates an Address from one or more HCL traversals.
func FromHCLTraversals(traversals ...hcl.Traversal) *Address {
	segments := make([]*parserpb.Segment, 0, len(traversals))
	for _, t := range traversals {
		for _, p := range t {
			switch part := p.(type) {
			case hcl.TraverseRoot:
				segments = append(segments, &parserpb.Segment{
					Value: part.Name,
				})
			case hcl.TraverseAttr:
				segments = append(segments, &parserpb.Segment{
					Value: part.Name,
				})
			case hcl.TraverseIndex:
				strX, intX := getHCLIndexValue(part)
				if len(segments) == 0 || segments[len(segments)-1].IndexInt != nil || segments[len(segments)-1].IndexString != nil {
					segments = append(segments, &parserpb.Segment{
						IndexInt:    intX,
						IndexString: strX,
					})
				} else {
					segments[len(segments)-1].IndexInt = intX
					segments[len(segments)-1].IndexString = strX
				}
			}
		}
	}
	return &Address{proto: &parserpb.Address{Segments: segments}}
}

func getHCLIndexValue(part hcl.TraverseIndex) (*string, *int64) {
	switch part.Key.Type() {
	case cty.String:
		s := part.Key.AsString()
		return &s, nil
	case cty.Number:
		var intVal int64
		_ = gocty.FromCtyValue(part.Key, &intVal)
		return nil, &intVal
	default:
		return nil, nil
	}
}
