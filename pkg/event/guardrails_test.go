package event

import (
	"testing"

	"github.com/infracost/go-proto/pkg/rat"
	"github.com/infracost/proto/gen/go/infracost/parser/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuardrails_NoThresholds(t *testing.T) {
	gs := Guardrails{
		{Id: "g1", Name: "no thresholds", Scope: event.Guardrail_REPO, PrComment: true},
	}

	results := gs.Evaluate(rat.New(1000), rat.New(500), nil)
	require.Len(t, results, 1)
	assert.Equal(t, "g1", results[0].GuardrailID)
	assert.False(t, results[0].Triggered)
	assert.True(t, results[0].PRComment)
}

func TestGuardrails_CostDecrease(t *testing.T) {
	gs := Guardrails{
		{
			Id:                "g1",
			Scope:             event.Guardrail_REPO,
			IncreaseThreshold: rat.New(10).Proto(),
		},
	}

	results := gs.Evaluate(rat.New(500), rat.New(1000), nil)
	require.Len(t, results, 1)
	assert.False(t, results[0].Triggered)
	// Cost data is still populated
	assert.True(t, results[0].Increase.Equals(rat.New(-500)))
}

func TestGuardrails_NoCostChange(t *testing.T) {
	gs := Guardrails{
		{
			Id:                "g1",
			Scope:             event.Guardrail_REPO,
			IncreaseThreshold: rat.New(10).Proto(),
		},
	}

	results := gs.Evaluate(rat.New(1000), rat.New(1000), nil)
	require.Len(t, results, 1)
	assert.False(t, results[0].Triggered)
}

func TestGuardrails_RepoScope_IncreaseThreshold(t *testing.T) {
	gs := Guardrails{
		{
			Id:                "g1",
			Name:              "increase check",
			Scope:             event.Guardrail_REPO,
			IncreaseThreshold: rat.New(250).Proto(),
			BlockPr:           true,
		},
	}

	t.Run("exceeded", func(t *testing.T) {
		results := gs.Evaluate(rat.New(1500), rat.New(1200), nil)
		require.Len(t, results, 1)
		assert.Equal(t, "g1", results[0].GuardrailID)
		assert.True(t, results[0].Triggered)
		assert.True(t, results[0].BlockPR)
		assert.Empty(t, results[0].TriggeringProjectNames)
		assert.True(t, results[0].Increase.Equals(rat.New(300)))
	})

	t.Run("not exceeded", func(t *testing.T) {
		results := gs.Evaluate(rat.New(1400), rat.New(1200), nil)
		require.Len(t, results, 1)
		assert.False(t, results[0].Triggered)
		// Cost data still present
		assert.True(t, results[0].Increase.Equals(rat.New(200)))
	})
}

func TestGuardrails_RepoScope_PercentThreshold(t *testing.T) {
	gs := Guardrails{
		{
			Id:                       "g1",
			Scope:                    event.Guardrail_REPO,
			IncreasePercentThreshold: rat.New(10).Proto(),
		},
	}

	t.Run("exceeded", func(t *testing.T) {
		// 1200 -> 1500 = 25% increase
		results := gs.Evaluate(rat.New(1500), rat.New(1200), nil)
		require.Len(t, results, 1)
		assert.True(t, results[0].Triggered)
		assert.True(t, results[0].PercentIncrease.Equals(rat.New(25)))
	})

	t.Run("not exceeded", func(t *testing.T) {
		// 1200 -> 1300 = ~8% increase
		results := gs.Evaluate(rat.New(1300), rat.New(1200), nil)
		require.Len(t, results, 1)
		assert.False(t, results[0].Triggered)
	})
}

func TestGuardrails_RepoScope_TotalThreshold_Crossing(t *testing.T) {
	gs := Guardrails{
		{
			Id:             "g1",
			Scope:          event.Guardrail_REPO,
			TotalThreshold: rat.New(1000).Proto(),
		},
	}

	t.Run("crosses from below to above", func(t *testing.T) {
		results := gs.Evaluate(rat.New(1050), rat.New(950), nil)
		require.Len(t, results, 1)
		assert.True(t, results[0].Triggered)
	})

	t.Run("already above", func(t *testing.T) {
		results := gs.Evaluate(rat.New(1100), rat.New(1050), nil)
		require.Len(t, results, 1)
		assert.False(t, results[0].Triggered)
	})

	t.Run("still below", func(t *testing.T) {
		results := gs.Evaluate(rat.New(950), rat.New(900), nil)
		require.Len(t, results, 1)
		assert.False(t, results[0].Triggered)
	})

	t.Run("crosses from at threshold to above", func(t *testing.T) {
		results := gs.Evaluate(rat.New(1050), rat.New(1000), nil)
		require.Len(t, results, 1)
		assert.True(t, results[0].Triggered)
	})
}

func TestGuardrails_RepoScope_ANDLogic(t *testing.T) {
	gs := Guardrails{
		{
			Id:                       "g1",
			Scope:                    event.Guardrail_REPO,
			IncreaseThreshold:        rat.New(100).Proto(),
			IncreasePercentThreshold: rat.New(20).Proto(),
		},
	}

	t.Run("both exceeded", func(t *testing.T) {
		// 1000 -> 1300 = $300 increase, 30%
		results := gs.Evaluate(rat.New(1300), rat.New(1000), nil)
		require.Len(t, results, 1)
		assert.True(t, results[0].Triggered)
	})

	t.Run("increase exceeded but not percent", func(t *testing.T) {
		// 1000 -> 1100 = $100 increase, 10%
		results := gs.Evaluate(rat.New(1100), rat.New(1000), nil)
		require.Len(t, results, 1)
		assert.False(t, results[0].Triggered)
	})
}

func TestGuardrails_ProjectScope(t *testing.T) {
	gs := Guardrails{
		{
			Id:                "g1",
			Scope:             event.Guardrail_PROJECT,
			IncreaseThreshold: rat.New(100).Proto(),
		},
	}

	projects := []ProjectCostInfo{
		{ProjectName: "proj-a", TotalMonthlyCost: rat.New(1200), PastTotalMonthlyCost: rat.New(1000)},
		{ProjectName: "proj-b", TotalMonthlyCost: rat.New(500), PastTotalMonthlyCost: rat.New(450)},
	}

	results := gs.Evaluate(rat.New(1700), rat.New(1450), projects)
	require.Len(t, results, 1)
	assert.True(t, results[0].Triggered)
	assert.Equal(t, []string{"proj-a"}, results[0].TriggeringProjectNames)
	assert.True(t, results[0].Increase.Equals(rat.New(200)))
}

func TestGuardrails_ProjectScope_NoneTriggered(t *testing.T) {
	gs := Guardrails{
		{
			Id:                "g1",
			Scope:             event.Guardrail_PROJECT,
			IncreaseThreshold: rat.New(500).Proto(),
		},
	}

	projects := []ProjectCostInfo{
		{ProjectName: "proj-a", TotalMonthlyCost: rat.New(1200), PastTotalMonthlyCost: rat.New(1000)},
	}

	results := gs.Evaluate(rat.New(1200), rat.New(1000), projects)
	require.Len(t, results, 1)
	assert.False(t, results[0].Triggered)
	// Max increase is still populated
	assert.True(t, results[0].Increase.Equals(rat.New(200)))
}

func TestGuardrails_ProjectScope_WithFilter(t *testing.T) {
	gs := Guardrails{
		{
			Id:                "g1",
			Scope:             event.Guardrail_PROJECT,
			IncreaseThreshold: rat.New(50).Proto(),
			ProjectFilter: &event.StringFilter{
				Include: []string{"prod-*"},
			},
		},
	}

	projects := []ProjectCostInfo{
		{ProjectName: "prod-api", TotalMonthlyCost: rat.New(600), PastTotalMonthlyCost: rat.New(500)},
		{ProjectName: "staging-api", TotalMonthlyCost: rat.New(600), PastTotalMonthlyCost: rat.New(500)},
	}

	results := gs.Evaluate(rat.New(1200), rat.New(1000), projects)
	require.Len(t, results, 1)
	assert.True(t, results[0].Triggered)
	assert.Equal(t, []string{"prod-api"}, results[0].TriggeringProjectNames)
}

func TestGuardrails_ProjectScope_MaxIncrease(t *testing.T) {
	gs := Guardrails{
		{
			Id:                "g1",
			Scope:             event.Guardrail_PROJECT,
			IncreaseThreshold: rat.New(50).Proto(),
		},
	}

	projects := []ProjectCostInfo{
		{ProjectName: "proj-a", TotalMonthlyCost: rat.New(600), PastTotalMonthlyCost: rat.New(500)},
		{ProjectName: "proj-b", TotalMonthlyCost: rat.New(800), PastTotalMonthlyCost: rat.New(500)},
	}

	results := gs.Evaluate(rat.New(1400), rat.New(1000), projects)
	require.Len(t, results, 1)
	assert.True(t, results[0].Triggered)
	// Should use the max increase ($300 from proj-b, not $100 from proj-a)
	assert.True(t, results[0].Increase.Equals(rat.New(300)))
	assert.Len(t, results[0].TriggeringProjectNames, 2)
}

func TestGuardrails_MultipleGuardrails(t *testing.T) {
	gs := Guardrails{
		{
			Id:                "g1",
			Scope:             event.Guardrail_REPO,
			IncreaseThreshold: rat.New(100).Proto(),
			PrComment:         true,
		},
		{
			Id:                "g2",
			Scope:             event.Guardrail_REPO,
			IncreaseThreshold: rat.New(500).Proto(),
		},
	}

	// $300 increase triggers g1 but not g2; both are returned
	results := gs.Evaluate(rat.New(1300), rat.New(1000), nil)
	require.Len(t, results, 2)

	assert.Equal(t, "g1", results[0].GuardrailID)
	assert.True(t, results[0].Triggered)
	assert.True(t, results[0].PRComment)

	assert.Equal(t, "g2", results[1].GuardrailID)
	assert.False(t, results[1].Triggered)
	assert.False(t, results[1].PRComment)
}

func TestGuardrails_PercentIncrease_FromZero(t *testing.T) {
	gs := Guardrails{
		{
			Id:                       "g1",
			Scope:                    event.Guardrail_REPO,
			IncreasePercentThreshold: rat.New(10).Proto(),
		},
	}

	// From $0 to $100: percent increase should be 0 (not infinity), so percent
	// threshold is not exceeded.
	results := gs.Evaluate(rat.New(100), rat.New(0), nil)
	require.Len(t, results, 1)
	assert.False(t, results[0].Triggered)
}