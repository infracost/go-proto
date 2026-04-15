package event

import (
	"testing"
	"time"

	"github.com/infracost/go-proto/pkg/rat"
	"github.com/infracost/proto/gen/go/infracost/parser/event"
	"github.com/infracost/proto/gen/go/infracost/rational"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ratProto(v int64) *rational.Rat {
	return rat.New(v).Proto()
}

func mkBudget(id string, tags []*event.BudgetTag, amount, currentCost int64) *event.Budget {
	return &event.Budget{
		Id:          id,
		Tags:        tags,
		Amount:      ratProto(amount),
		CurrentCost: ratProto(currentCost),
		PrComment:   true,
		StartedAt:   timestamppb.New(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
		EndedAt:     timestamppb.New(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)),
	}
}

func mkCostInfo(tags map[string]string, cost int64) ResourceCostInfo {
	return ResourceCostInfo{
		Tags:        tags,
		MonthlyCost: rat.New(cost),
	}
}

func TestBudgets_Empty(t *testing.T) {
	bs := Budgets{}
	results := bs.Evaluate(nil)
	assert.Empty(t, results)
}

func TestBudgets_NoResources_FiltersOut(t *testing.T) {
	bs := Budgets{
		mkBudget("b1", []*event.BudgetTag{{Key: "env", Value: "prod"}}, 1000, 500),
	}

	results := bs.Evaluate(nil)
	assert.Empty(t, results, "budget should be filtered out when no resources match")
}

func TestBudgets_MatchingResources_IncludesBudget(t *testing.T) {
	bs := Budgets{
		mkBudget("b1", []*event.BudgetTag{{Key: "env", Value: "prod"}}, 1000, 500),
	}

	resources := []ResourceCostInfo{
		mkCostInfo(map[string]string{"env": "prod"}, 300),
	}

	results := bs.Evaluate(resources)
	require.Len(t, results, 1)
	assert.Equal(t, "b1", results[0].BudgetID)
	// CurrentCost comes from the proto, not computed from resources.
	assert.True(t, results[0].CurrentCost.Equals(rat.New(500)))
	assert.True(t, results[0].Amount.Equals(rat.New(1000)))
}

func TestBudgets_NonMatchingResources_FiltersOut(t *testing.T) {
	bs := Budgets{
		mkBudget("b1", []*event.BudgetTag{{Key: "env", Value: "prod"}}, 1000, 500),
	}

	resources := []ResourceCostInfo{
		mkCostInfo(map[string]string{"env": "staging"}, 300),
	}

	results := bs.Evaluate(resources)
	assert.Empty(t, results)
}

func TestBudgets_PartialTagMatch_FiltersOut(t *testing.T) {
	// Budget requires both env=prod AND team=backend.
	// Resource only has env=prod — should not match.
	bs := Budgets{
		mkBudget("b1", []*event.BudgetTag{
			{Key: "env", Value: "prod"},
			{Key: "team", Value: "backend"},
		}, 1000, 500),
	}

	resources := []ResourceCostInfo{
		mkCostInfo(map[string]string{"env": "prod"}, 300),
	}

	results := bs.Evaluate(resources)
	assert.Empty(t, results)
}

func TestBudgets_AllTagsMatch(t *testing.T) {
	bs := Budgets{
		mkBudget("b1", []*event.BudgetTag{
			{Key: "env", Value: "prod"},
			{Key: "team", Value: "backend"},
		}, 1000, 800),
	}

	resources := []ResourceCostInfo{
		mkCostInfo(map[string]string{"env": "prod", "team": "backend", "owner": "alice"}, 300),
	}

	results := bs.Evaluate(resources)
	require.Len(t, results, 1)
	assert.True(t, results[0].CurrentCost.Equals(rat.New(800)))
}

func TestBudgets_MultipleBudgets_FiltersByTags(t *testing.T) {
	bs := Budgets{
		mkBudget("b1", []*event.BudgetTag{{Key: "env", Value: "prod"}}, 1000, 500),
		mkBudget("b2", []*event.BudgetTag{{Key: "team", Value: "frontend"}}, 500, 400),
		mkBudget("b3", []*event.BudgetTag{{Key: "team", Value: "infra"}}, 2000, 100),
	}

	resources := []ResourceCostInfo{
		mkCostInfo(map[string]string{"env": "prod", "team": "frontend"}, 300),
	}

	results := bs.Evaluate(resources)
	require.Len(t, results, 2, "b3 should be filtered out (team=infra not in resources)")

	assert.Equal(t, "b1", results[0].BudgetID)
	assert.True(t, results[0].CurrentCost.Equals(rat.New(500)))

	assert.Equal(t, "b2", results[1].BudgetID)
	assert.True(t, results[1].CurrentCost.Equals(rat.New(400)))
}

func TestBudgets_ResourceWithoutTags_NoMatch(t *testing.T) {
	bs := Budgets{
		mkBudget("b1", []*event.BudgetTag{{Key: "env", Value: "prod"}}, 1000, 500),
	}

	resources := []ResourceCostInfo{
		{MonthlyCost: rat.New(500)},
	}

	results := bs.Evaluate(resources)
	assert.Empty(t, results)
}

func TestBudgets_CustomOverrunMessage(t *testing.T) {
	bs := Budgets{
		{
			Id:                   "b1",
			Tags:                 []*event.BudgetTag{{Key: "env", Value: "prod"}},
			Amount:               ratProto(1000),
			CurrentCost:          ratProto(1500),
			PrComment:            true,
			CustomOverrunMessage: "Contact FinOps team",
			StartedAt:            timestamppb.New(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
			EndedAt:              timestamppb.New(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)),
		},
	}

	resources := []ResourceCostInfo{
		mkCostInfo(map[string]string{"env": "prod"}, 100),
	}

	results := bs.Evaluate(resources)
	require.Len(t, results, 1)
	assert.Equal(t, "Contact FinOps team", results[0].CustomOverrunMessage)
	assert.True(t, results[0].CurrentCost.Equals(rat.New(1500)))
}

func TestBudgets_DatesPropagated(t *testing.T) {
	start := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)

	bs := Budgets{
		{
			Id:        "b1",
			Tags:      []*event.BudgetTag{{Key: "env", Value: "prod"}},
			Amount:    ratProto(1000),
			StartedAt: timestamppb.New(start),
			EndedAt:   timestamppb.New(end),
		},
	}

	resources := []ResourceCostInfo{
		mkCostInfo(map[string]string{"env": "prod"}, 100),
	}

	results := bs.Evaluate(resources)
	require.Len(t, results, 1)
	assert.Equal(t, start, results[0].StartDate)
	assert.Equal(t, end, results[0].EndDate)
}

func TestBudgets_BudgetName(t *testing.T) {
	bs := Budgets{
		{
			Id:          "b1",
			Name:        "Production budget",
			Tags:        []*event.BudgetTag{{Key: "env", Value: "prod"}},
			Amount:      ratProto(1000),
			CurrentCost: ratProto(500),
		},
	}

	resources := []ResourceCostInfo{
		mkCostInfo(map[string]string{"env": "prod"}, 100),
	}

	results := bs.Evaluate(resources)
	require.Len(t, results, 1)
	assert.Equal(t, "Production budget", results[0].BudgetName)
}

func TestBudgets_NilCurrentCost_DefaultsToZero(t *testing.T) {
	bs := Budgets{
		{
			Id:     "b1",
			Tags:   []*event.BudgetTag{{Key: "env", Value: "prod"}},
			Amount: ratProto(1000),
			// CurrentCost deliberately nil
		},
	}

	resources := []ResourceCostInfo{
		mkCostInfo(map[string]string{"env": "prod"}, 100),
	}

	results := bs.Evaluate(resources)
	require.Len(t, results, 1)
	assert.True(t, results[0].CurrentCost.IsZero())
}

func TestBudgets_TagsFromMultipleResources(t *testing.T) {
	// Budget needs env=prod AND team=backend.
	// No single resource has both, but the set of resources together does.
	bs := Budgets{
		mkBudget("b1", []*event.BudgetTag{
			{Key: "env", Value: "prod"},
			{Key: "team", Value: "backend"},
		}, 1000, 500),
	}

	resources := []ResourceCostInfo{
		mkCostInfo(map[string]string{"env": "prod"}, 100),
		mkCostInfo(map[string]string{"team": "backend"}, 200),
	}

	results := bs.Evaluate(resources)
	require.Len(t, results, 1, "tags from different resources should combine for matching")
}