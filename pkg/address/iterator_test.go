package address

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Iterator(t *testing.T) {

	tests := []struct {
		name             string
		address          string
		expectedSegments []any
		expectedTypes    []SegmentType
	}{
		{
			name:             "simple",
			address:          "a.b.c",
			expectedSegments: []any{"a", "b", "c"},
			expectedTypes:    []SegmentType{SegmentTypeString, SegmentTypeString, SegmentTypeString},
		},
		{
			name:             "int index",
			address:          "a[0].b",
			expectedSegments: []any{"a", 0, "b"},
			expectedTypes:    []SegmentType{SegmentTypeString, SegmentTypeInt, SegmentTypeString},
		},
		{
			name:             "int index with string index",
			address:          `a[0].b["something"]`,
			expectedSegments: []any{"a", 0, "b", "something"},
			expectedTypes:    []SegmentType{SegmentTypeString, SegmentTypeInt, SegmentTypeString, SegmentTypeString},
		},
		{
			name:             "list wildcard",
			address:          `a[0].b.#`,
			expectedSegments: []any{"a", 0, "b", "#"},
			expectedTypes:    []SegmentType{SegmentTypeString, SegmentTypeInt, SegmentTypeString, SegmentTypeWildcard},
		},
		{
			name:             "map wildcard",
			address:          `a[0].b.*`,
			expectedSegments: []any{"a", 0, "b", "*"},
			expectedTypes:    []SegmentType{SegmentTypeString, SegmentTypeInt, SegmentTypeString, SegmentTypeWildcard},
		},
		{
			name:             "literal * map key",
			address:          `a[0].b["*"]`,
			expectedSegments: []any{"a", 0, "b", "*"},
			expectedTypes:    []SegmentType{SegmentTypeString, SegmentTypeInt, SegmentTypeString, SegmentTypeString},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parsed, err := Parse(test.address)
			require.NoError(t, err, "invalid address in test case")
			iter := NewIterator(parsed)
			var actualValues []any
			var actualTypes []SegmentType
			for {
				peekedType := iter.PeekType()
				s, i := iter.PopKey()
				if s == nil && i == nil {
					break
				}
				if s != nil {
					actualValues = append(actualValues, *s)
					if i != nil {
						t.Errorf("expected only string or int key, got both")
					}
				}
				if i != nil {
					actualValues = append(actualValues, int(*i))
				}
				actualTypes = append(actualTypes, peekedType)
			}
			require.Equal(t, test.expectedSegments, actualValues)
			require.Equal(t, test.expectedTypes, actualTypes)
		})
	}

}
