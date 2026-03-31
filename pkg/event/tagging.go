package event

import (
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/infracost/go-proto/pkg/address"
	"github.com/infracost/proto/gen/go/infracost/parser/event"
	"github.com/infracost/proto/gen/go/infracost/provider"
)

type TaggingPolicyResult struct {
	Name                   string                           `json:"name"`
	TagPolicyID            string                           `json:"tagPolicyId"`
	Message                string                           `json:"message"`
	FailingResources       []TagPolicyResultResource        `json:"resources"`
	PassingResources       []TagPolicyResultPassingResource `json:"passingResources"`
	PastFailingResources   []TagPolicyResultResource        `json:"pastResources"`
	TotalDetectedResources int                              `json:"totalDetectedResources"`
	TotalTaggableResources int                              `json:"totalTaggableResources"`
	BlockPR                bool                             `json:"blockPr"`
	PRComment              bool                             `json:"prComment"`
}

type TagPolicyResultResource struct {
	Address                  string               `json:"address"`
	ResourceType             string               `json:"resourceType"`
	Path                     string               `json:"path"`
	ModulePath               string               `json:"modulePath"`
	ModuleCallPath           string               `json:"moduleCallPath"`
	ProviderLink             string               `json:"providerLink"`
	InvalidTags              []InvalidTag         `json:"invalidTags"`
	ProjectNames             []string             `json:"projectNames"`
	MissingMandatoryTags     []string             `json:"missingMandatoryTags"`
	PropagationProblems      []PropagationProblem `json:"propagationProblems"`
	Line                     int                  `json:"line"`
	ModuleCallLine           int                  `json:"moduleCallLine"`
	SupportsDefaultTags      bool                 `json:"supportsDefaultTags"`
	HasDefaultTags           bool                 `json:"hasDefaultTags"`
	DefaultTagsNotPropagated bool                 `json:"defaultTagsNotPropagated"`
}

type InvalidTag struct {
	Key                  string   `json:"key"`
	Value                string   `json:"value"`
	ValidRegex           string   `json:"validRegex"`
	Suggestion           string   `json:"suggestion"`
	Message              string   `json:"message"`
	ValidValues          []string `json:"validValues"`
	ValidValueCount      int      `json:"validValueCount"`
	ValidValuesTruncated bool     `json:"validValuesTruncated"`
	FromDefaultTags      bool     `json:"fromDefaultTags"`
	MissingMandatory     bool     `json:"missingMandatory"`
}

type PropagationProblem struct {
	Attribute    string   `json:"attribute"`
	From         string   `json:"from"`
	To           string   `json:"to"`
	ValidSources []string `json:"validSources"`
	AffectedTags []string `json:"affectedTags"`
}

type TagPolicyResultPassingResource struct {
	Address      string   `json:"address"`
	ProjectNames []string `json:"projectNames"`
}

// hasDefaultTags returns true if any tag in the list is a default tag.
func hasDefaultTags(tags []*provider.Tag) bool {
	for _, tag := range tags {
		if tag.IsDefault {
			return true
		}
	}
	return false
}

type TagPolicies []*event.TagPolicy

func (t TagPolicies) EvaluateAgainstResources(resources []*provider.Resource, projectInfo *provider.ProjectInfo) []TaggingPolicyResult {

	// filter out policies for other repos/branches etc.
	var filteredPolicies []*event.TagPolicy
	for _, policy := range t {
		if StringFilterFromProto(policy.GetProjectFilter()).Matches(projectInfo.Name) &&
			StringFilterFromProto(policy.GetBranchFilter()).Matches(projectInfo.BranchName) {
			filteredPolicies = append(filteredPolicies, policy)
		}
	}

	// grab resources which support tags
	compatibleResources := make([]*provider.Resource, 0, len(resources))
	for _, resource := range resources {
		if resource.Tagging != nil && resource.Tagging.SupportsTags {
			compatibleResources = append(compatibleResources, resource)
		}
	}

	// preallocate the results
	results := make([]TaggingPolicyResult, 0, len(filteredPolicies))

	for _, policy := range filteredPolicies {

		var filteredResources []*provider.Resource
		for _, resource := range compatibleResources {
			if StringFilterFromProto(policy.GetResourceFilter()).Matches(resource.Type) {
				filteredResources = append(filteredResources, resource)
			}
		}

		result := TaggingPolicyResult{
			Name:                   policy.Name,
			TagPolicyID:            policy.Id,
			Message:                policy.Message,
			BlockPR:                policy.BlockPr,
			PRComment:              policy.PrComment,
			FailingResources:       make([]TagPolicyResultResource, 0),
			PassingResources:       make([]TagPolicyResultPassingResource, 0),
			PastFailingResources:   make([]TagPolicyResultResource, 0),
			TotalDetectedResources: 0,
			TotalTaggableResources: 0,
		}

		for _, resource := range filteredResources {
			evaluateResourceAgainstPolicy(&result, resource, policy, projectInfo.Name)
		}

		if len(result.PassingResources) > 0 || len(result.FailingResources) > 0 {
			results = append(results, result)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].TagPolicyID < results[j].TagPolicyID
	})

	return results
}

type tagCandidate struct {
	address string
	tagging *provider.Tagging
}

// evaluateResourceAgainstPolicy evaluates a resource and its taggable children against a tag policy.
func evaluateResourceAgainstPolicy(result *TaggingPolicyResult, resource *provider.Resource, policy *event.TagPolicy, projectName string) {
	// collect candidates: the resource itself plus any taggable child resources
	candidates := []tagCandidate{{
		address: resource.Name,
		tagging: resource.Tagging,
	}}
	for _, sub := range resource.ChildResources {
		if sub.Tagging != nil && sub.Tagging.SupportsTags {
			candidates = append(candidates, tagCandidate{
				address: strings.Join([]string{resource.Name, sub.Name}, "."),
				tagging: sub.Tagging,
			})
		}
	}

	baseHasDefaultTags := resource.Tagging != nil && hasDefaultTags(resource.Tagging.Tags)

	for _, candidate := range candidates {
		if !MapFilterFromProto(policy.GetTagFilter()).Matches(candidate.tagging.Tags) {
			continue
		}

		result.TotalDetectedResources++
		result.TotalTaggableResources++

		resourceResult := evaluateTagPolicyOnResource(candidate, policy.GetRequirements())
		if len(resourceResult.invalidTags) == 0 && len(resourceResult.missingMandatoryTags) == 0 && len(resourceResult.propagationProblems) == 0 {
			result.PassingResources = append(result.PassingResources, TagPolicyResultPassingResource{
				Address:      candidate.address,
				ProjectNames: []string{projectName},
			})
			continue
		}

		failingResource := buildFailingResource(resource, candidate, resourceResult, baseHasDefaultTags, projectName)
		result.FailingResources = append(result.FailingResources, failingResource)
	}
}

// buildFailingResource constructs a TagPolicyResultResource for a resource that failed policy evaluation.
func buildFailingResource(resource *provider.Resource, candidate tagCandidate, evalResult resourceTagPolicyResult, baseHasDefaultTags bool, projectName string) TagPolicyResultResource {
	candidateHasDefaultTags := hasDefaultTags(candidate.tagging.Tags)

	failingResource := TagPolicyResultResource{
		MissingMandatoryTags:     evalResult.missingMandatoryTags,
		InvalidTags:              evalResult.invalidTags,
		PropagationProblems:      evalResult.propagationProblems,
		DefaultTagsNotPropagated: !candidateHasDefaultTags && baseHasDefaultTags,
		Address:                  candidate.address,
		ResourceType:             resource.Type,
		ProviderLink:             resource.ProviderLink,
		ProjectNames:             []string{projectName},
		SupportsDefaultTags:      resource.Tagging.SupportsDefaultTags,
		HasDefaultTags:           baseHasDefaultTags,
	}

	if meta := resource.Metadata; meta != nil {
		failingResource.Path = meta.Filename
		failingResource.Line = int(meta.StartLine)
	}

	// only add module data if we're at least 1 module "deep"
	if stack := resource.CallStack; stack != nil && len(stack.Frames) > 1 {
		failingResource.ModulePath = stack.Frames[1].Source
		failingResource.ModuleCallPath = address.FromProto(stack.Frames[0].Address).String()
		failingResource.ModuleCallLine = int(stack.Frames[0].SourceRange.StartLine)
	}

	return failingResource
}

type resourceTagPolicyResult struct {
	missingMandatoryTags []string
	invalidTags          []InvalidTag
	propagationProblems  []PropagationProblem
}

func evaluateTagPolicyOnResource(resource tagCandidate, requirements []*event.TagPolicyRequirement) resourceTagPolicyResult {

	result := resourceTagPolicyResult{
		missingMandatoryTags: make([]string, 0),
		invalidTags:          make([]InvalidTag, 0),
		propagationProblems:  make([]PropagationProblem, 0),
	}

	for _, problem := range resource.tagging.PropagationProblems {

		// ensure non-nil slices for api
		validValues := problem.ValidValues
		if len(validValues) == 0 {
			validValues = []string{}
		}
		affectedTags := problem.AffectedTags
		if len(affectedTags) == 0 {
			affectedTags = []string{}
		}

		result.propagationProblems = append(result.propagationProblems, PropagationProblem{
			Attribute:    problem.Attribute,
			From:         problem.ActualValue,
			To:           problem.TagRecipient,
			ValidSources: validValues,
			AffectedTags: affectedTags,
		})
	}

	var hasSyntheticKeys bool
	for _, tag := range resource.tagging.Tags {
		if tag.IsKeySynthetic {
			hasSyntheticKeys = true
			break
		}
	}

	for _, requirement := range requirements {

		var tagExists bool
		for _, tag := range resource.tagging.Tags {
			if tag.Key == requirement.GetKey() {
				tagExists = true

				// ignore synthetic tag values
				if tag.IsValueSynthetic {
					continue
				}

				switch requirement.GetType() {
				case event.TagPolicyRequirement_ANY:
					// nothing to do here
				case event.TagPolicyRequirement_LIST:
					//  check if the value is in the allowed values
					if !slices.Contains(requirement.GetAllowedValues(), tag.Value) {
						result.invalidTags = append(result.invalidTags, createSimplifiedInvalidTagForList(requirement.GetKey(), tag.Value, requirement.GetAllowedValues(), tag.IsDefault, requirement.GetMessage()))
					}
				case event.TagPolicyRequirement_REGEX:
					compiled, err := convertJSRegexToGo(requirement.GetValueRegex())
					if err == nil && !compiled.MatchString(tag.Value) {
						result.invalidTags = append(result.invalidTags, InvalidTag{
							Key:             requirement.GetKey(),
							Value:           tag.Value,
							ValidRegex:      requirement.GetValueRegex(),
							Message:         requirement.GetMessage(),
							FromDefaultTags: tag.IsDefault,
						})
					}
				}

				break
			}
		}

		// only report missing mandatory keys if there are no synthetic keys
		if requirement.GetMandatory() && !tagExists && !hasSyntheticKeys {
			if requirement.GetType() == event.TagPolicyRequirement_LIST {
				invalid := createSimplifiedInvalidTagForList(requirement.GetKey(), "", requirement.GetAllowedValues(), false, requirement.GetMessage())
				invalid.MissingMandatory = true
				result.invalidTags = append(result.invalidTags, invalid)
			} else {
				result.missingMandatoryTags = append(result.missingMandatoryTags, requirement.GetKey())
			}
		}
	}

	return result
}

const maxStoredValidTagValues = 5
const maxValidTagValuesForSuggestions = 250

type suggestionMatch struct {
	Value    string
	Distance int
}

func createSimplifiedInvalidTagForList(key string, value string, validValues []string, fromDefault bool, message string) InvalidTag {
	result := InvalidTag{
		Key:             key,
		Value:           value,
		ValidValues:     validValues,
		ValidValueCount: len(validValues),
		Message:         message,
		FromDefaultTags: fromDefault,
	}

	// if the number of valid values is too long to provide suggestions, truncate it to the maximum we can store and skip suggesting
	if len(result.ValidValues) > maxValidTagValuesForSuggestions {
		result.ValidValuesTruncated = true
		result.ValidValues = result.ValidValues[:maxStoredValidTagValues]
		return result
	}

	if value == "" {
		if len(result.ValidValues) > maxStoredValidTagValues {
			result.ValidValuesTruncated = true
			result.ValidValues = result.ValidValues[:maxStoredValidTagValues]
		}
		return result
	}

	var closeMatches []suggestionMatch
	allMatches := make([]suggestionMatch, 0, len(result.ValidValues))
	for _, candidate := range result.ValidValues {
		distance := levenshtein.ComputeDistance(strings.ToLower(value), strings.ToLower(candidate))
		if distance < 2 {
			closeMatches = append(closeMatches, suggestionMatch{
				Value:    candidate,
				Distance: distance,
			})
		}
		allMatches = append(allMatches, suggestionMatch{
			Value:    candidate,
			Distance: distance,
		})
	}

	// if we didn't find a close match yet, look for prefix matches e.g. prod == production
	if len(closeMatches) == 0 {
		for _, candidate := range result.ValidValues {
			if strings.HasPrefix(strings.ToLower(candidate), strings.ToLower(value)) || strings.HasPrefix(strings.ToLower(value), strings.ToLower(candidate)) {
				closeMatches = append(closeMatches, suggestionMatch{
					Value:    candidate,
					Distance: 0, // set these up as perfect matches, as if these are the only match type, they're all considered equal anyway
				})
			}
		}
	}

	// if there is one standout match, use it as the suggestion and return early
	if len(closeMatches) == 1 {
		result.Suggestion = closeMatches[0].Value
		result.ValidValues = []string{closeMatches[0].Value}
		result.ValidValuesTruncated = len(validValues) > 1
		return result
	}

	sort.Slice(allMatches, func(i, j int) bool {
		if allMatches[i].Distance == allMatches[j].Distance {
			return allMatches[i].Value < allMatches[j].Value
		}
		return allMatches[i].Distance < allMatches[j].Distance
	})

	if len(allMatches) > maxStoredValidTagValues {
		allMatches = allMatches[:maxStoredValidTagValues]
	}

	result.ValidValues = make([]string, len(allMatches))
	for i, m := range allMatches {
		result.ValidValues[i] = m.Value
	}
	result.ValidValuesTruncated = len(allMatches) < len(validValues)

	return result
}

func convertJSRegexToGo(jsRegex string) (*regexp.Regexp, error) {
	// Check if it looks like a JS regex (starts and ends with /)
	if !strings.HasPrefix(jsRegex, "/") {
		// Not a JS regex pattern with delimiters, return as is
		return regexp.Compile(jsRegex)
	}

	// Find the last slash that's not escaped
	lastSlashIndex := -1
	for i := len(jsRegex) - 1; i > 0; i-- {
		if jsRegex[i] == '/' && (i == 0 || jsRegex[i-1] != '\\') {
			lastSlashIndex = i
			break
		}
	}

	if lastSlashIndex == -1 || lastSlashIndex == 0 {
		return regexp.Compile(jsRegex)
	}

	// Extract the pattern and flags
	pattern := jsRegex[1:lastSlashIndex]
	flags := ""
	if lastSlashIndex < len(jsRegex)-1 {
		flags = jsRegex[lastSlashIndex+1:]
	}

	// Process the flags
	flagStr := ""
	for _, flag := range flags {
		switch flag {
		case 'i', 'm', 's', 'U':
			// Case-insensitive
			flagStr += string(flag)
		}
	}

	if flagStr != "" {
		flagStr = "(?" + flagStr + ")"
	}

	converted, err := regexp.Compile(flagStr + pattern)
	if err != nil {
		// fall back to the original regex if we converted to something invalid
		return regexp.Compile(jsRegex)
	}

	return converted, nil
}
