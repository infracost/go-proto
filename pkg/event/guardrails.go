package event

import (
	"github.com/infracost/go-proto/pkg/rat"
	"github.com/infracost/proto/gen/go/infracost/parser/event"
)

// GuardrailResult is the result of evaluating a single applicable guardrail.
// It is serialized as JSON and sent back to the dashboard API via the addRun
// GraphQL mutation. Every applicable guardrail produces a result, whether or
// not it triggered.
type GuardrailResult struct {
	GuardrailID            string   `json:"guardrailId"`
	Triggered              bool     `json:"triggered"`
	PRComment              bool     `json:"prComment"`
	BlockPR                bool     `json:"blockPr"`
	TriggeringProjectNames []string `json:"triggeringProjectNames"`
	Increase               *rat.Rat `json:"increase"`
	PercentIncrease        *rat.Rat `json:"percentIncrease"`
	TotalMonthlyCost       *rat.Rat `json:"totalMonthlyCost"`
}

// ProjectCostInfo holds the cost data for a single project, used as input
// to guardrail evaluation.
type ProjectCostInfo struct {
	ProjectName          string
	TotalMonthlyCost     *rat.Rat
	PastTotalMonthlyCost *rat.Rat
}

// Guardrails wraps a list of guardrail proto configurations for evaluation.
type Guardrails []*event.Guardrail

// Evaluate checks guardrails against the provided run-level and project-level
// costs and returns a result for each applicable guardrail.
func (gs Guardrails) Evaluate(
	totalMonthlyCost *rat.Rat,
	pastTotalMonthlyCost *rat.Rat,
	projects []ProjectCostInfo,
) []GuardrailResult {
	var results []GuardrailResult

	for _, g := range gs {
		result := GuardrailResult{
			GuardrailID:            g.GetId(),
			PRComment:              g.GetPrComment(),
			BlockPR:                g.GetBlockPr(),
			TriggeringProjectNames: []string{},
		}

		if !hasThresholds(g) {
			results = append(results, result)
			continue
		}

		switch g.GetScope() {
		case event.Guardrail_REPO:
			increase := calcIncrease(totalMonthlyCost, pastTotalMonthlyCost)
			result.Increase = increase.increase
			result.PercentIncrease = increase.percentIncrease
			result.TotalMonthlyCost = totalMonthlyCost
			result.Triggered = thresholdExceeded(g, totalMonthlyCost, pastTotalMonthlyCost)

		case event.Guardrail_PROJECT:
			projectFilter := StringFilterFromProto(g.GetProjectFilter())

			var triggeringNames []string
			var maxIncrease costIncrease

			for _, p := range projects {
				if !projectFilter.Matches(p.ProjectName) {
					continue
				}

				inc := calcIncrease(p.TotalMonthlyCost, p.PastTotalMonthlyCost)
				if thresholdExceeded(g, p.TotalMonthlyCost, p.PastTotalMonthlyCost) {
					triggeringNames = append(triggeringNames, p.ProjectName)
				}

				if maxIncrease.increase == nil || inc.increase.GreaterThan(maxIncrease.increase) {
					maxIncrease = inc
				}
			}

			if maxIncrease.increase != nil {
				result.Increase = maxIncrease.increase
				result.PercentIncrease = maxIncrease.percentIncrease
			} else {
				result.Increase = rat.Zero
				result.PercentIncrease = rat.Zero
			}
			result.TotalMonthlyCost = totalMonthlyCost
			result.Triggered = len(triggeringNames) > 0
			result.TriggeringProjectNames = triggeringNames
		}

		results = append(results, result)
	}

	return results
}

func hasThresholds(g *event.Guardrail) bool {
	return g.IncreaseThreshold != nil || g.IncreasePercentThreshold != nil || g.TotalThreshold != nil
}

func thresholdExceeded(g *event.Guardrail, totalMonthlyCost, pastTotalMonthlyCost *rat.Rat) bool {
	increase := calcIncrease(totalMonthlyCost, pastTotalMonthlyCost)

	// Never trigger on cost decreases or no change.
	if !increase.increase.GreaterThanZero() {
		return false
	}

	// All configured thresholds must be exceeded (AND logic).
	if g.IncreaseThreshold != nil {
		if !increase.increase.GreaterThan(rat.FromProto(g.IncreaseThreshold)) {
			return false
		}
	}

	if g.IncreasePercentThreshold != nil {
		if !increase.percentIncrease.GreaterThan(rat.FromProto(g.IncreasePercentThreshold)) {
			return false
		}
	}

	// Total threshold uses crossing logic: only triggers when crossing from
	// at-or-below to above the threshold.
	if g.TotalThreshold != nil {
		threshold := rat.FromProto(g.TotalThreshold)
		if !pastTotalMonthlyCost.LessThanOrEqual(threshold) || !totalMonthlyCost.GreaterThan(threshold) {
			return false
		}
	}

	return true
}

type costIncrease struct {
	increase        *rat.Rat
	percentIncrease *rat.Rat
}

func calcIncrease(totalMonthlyCost, pastTotalMonthlyCost *rat.Rat) costIncrease {
	increase := totalMonthlyCost.Sub(pastTotalMonthlyCost)

	var percentIncrease *rat.Rat
	if pastTotalMonthlyCost.IsZero() {
		percentIncrease = rat.Zero
	} else {
		// (total / past - 1) * 100, rounded to 0 decimal places
		percentIncrease = totalMonthlyCost.Div(pastTotalMonthlyCost).Sub(rat.New(1)).Mul(rat.New(100)).Round(0)
	}

	return costIncrease{
		increase:        increase,
		percentIncrease: percentIncrease,
	}
}
