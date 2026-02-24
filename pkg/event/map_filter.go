package event

import (
	eventpb "github.com/infracost/proto/gen/go/infracost/parser/event"
	"github.com/infracost/proto/gen/go/infracost/provider"
)

type MapFilter struct {
	*eventpb.MapFilter
}

func MapFilterFromProto(pb *eventpb.MapFilter) *MapFilter {
	return &MapFilter{pb}
}

func (t *MapFilter) Matches(tags []*provider.Tag) bool {

	if t == nil || t.MapFilter == nil || (len(t.Include) == 0 && len(t.Exclude) == 0) {
		return true
	}

	tagMap := make(map[string]string, len(tags))
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}

	// 1. if there are exclude rules and all excluded tag key/values match, return false
	// 2. if there are any include rules, return true if they ALL match
	// 3. if there are no include rules, return true

	if len(t.GetExclude()) > 0 {
		if func() bool {
			for excludeKey, excludeValue := range t.GetExclude() {
				if actualValue, ok := tagMap[excludeKey]; !ok || !matchWildcard(actualValue, excludeValue) {
					return false
				}
			}
			return true
		}() {
			return false
		}
	}

	if len(t.GetInclude()) > 0 {
		for includeKey, includeValue := range t.GetInclude() {
			if actualValue, ok := tagMap[includeKey]; !ok || !matchWildcard(actualValue, includeValue) {
				return false
			}
		}
		return true
	}

	return true
}
