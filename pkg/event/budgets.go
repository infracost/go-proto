package event

import (
	"time"

	"github.com/infracost/go-proto/pkg/rat"
	"github.com/infracost/proto/gen/go/infracost/parser/event"
)

// BudgetResult represents a budget that is relevant to the current PR.
// It is serialized as JSON and sent back to the dashboard API via the
// addRun GraphQL mutation.
type BudgetResult struct {
	BudgetID             string      `json:"budgetId"`
	BudgetName           string      `json:"budgetName"`
	Tags                 []BudgetTag `json:"tags"`
	StartDate            time.Time   `json:"startDate"`
	EndDate              time.Time   `json:"endDate"`
	Amount               *rat.Rat    `json:"amount"`
	CurrentCost          *rat.Rat    `json:"currentCost"`
	CustomOverrunMessage string      `json:"customOverrunMessage,omitempty"`
}

// BudgetTag is a key-value pair defining a budget's tag scope.
type BudgetTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ResourceCostInfo holds pre-computed cost and tag data for a single resource.
// Callers build these from their own resource types (proto, breakdown, etc.)
// before passing to evaluation functions that need per-resource cost+tag data.
type ResourceCostInfo struct {
	Tags        map[string]string
	MonthlyCost *rat.Rat
}

// Budgets wraps a list of budget proto configurations.
type Budgets []*event.Budget

// Evaluate filters budgets to those whose tags appear on at least one of the
// provided resources, and returns the matching budgets as results. The caller
// passes all scanned resources (not just changed ones) — further narrowing to
// diff-only resources is done by the VCS comment layer. The current_cost field
// is pre-resolved by the dashboard from cloud billing data — no client-side
// cost computation is performed.
func (bs Budgets) Evaluate(resources []ResourceCostInfo) []BudgetResult {
	var results []BudgetResult

	// Collect all tag key/value pairs from resources into a set for matching.
	resourceTags := make(map[string]struct{})
	for _, r := range resources {
		for k, v := range r.Tags {
			resourceTags[k+"\x00"+v] = struct{}{}
		}
	}

	for _, b := range bs {
		tags := make([]BudgetTag, 0, len(b.GetTags()))
		for _, t := range b.GetTags() {
			tags = append(tags, BudgetTag{Key: t.GetKey(), Value: t.GetValue()})
		}

		// Skip budgets whose tags don't appear on any scanned resources.
		// The VCS comment layer may further narrow this to diff-only resources.
		if !allTagsPresent(tags, resourceTags) {
			continue
		}

		currentCost := rat.FromProto(b.GetCurrentCost())
		if currentCost == nil {
			currentCost = rat.Zero
		}

		result := BudgetResult{
			BudgetID:             b.GetId(),
			BudgetName:           b.GetName(),
			Tags:                 tags,
			Amount:               rat.FromProto(b.GetAmount()),
			CurrentCost:          currentCost,
			CustomOverrunMessage: b.GetCustomOverrunMessage(),
		}

		if b.GetStartedAt() != nil {
			result.StartDate = b.GetStartedAt().AsTime()
		}
		if b.GetEndedAt() != nil {
			result.EndDate = b.GetEndedAt().AsTime()
		}

		results = append(results, result)
	}

	return results
}

// allTagsPresent returns true if every budget tag appears in the resource tag set.
func allTagsPresent(budgetTags []BudgetTag, resourceTags map[string]struct{}) bool {
	for _, t := range budgetTags {
		if _, ok := resourceTags[t.Key+"\x00"+t.Value]; !ok {
			return false
		}
	}
	return true
}