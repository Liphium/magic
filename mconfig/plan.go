package mconfig

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

// Name for a service Docker container
func ContainerName[S ServiceDriver](appName string, profile string, driver S) string {
	appName = EverythingToSnakeCase(appName)
	return fmt.Sprintf("mgc-%s-%s-%s", appName, profile, driver.GetUniqueId())
}

type Plan struct {
	AppName        string                         `json:"app_name"`
	Profile        string                         `json:"profile"`
	Environment    map[string]string              `json:"environment"`
	AllocatedPorts map[uint]uint                  `json:"ports"`
	Containers     map[string]ContainerAllocation `json:"containers"` // Service id -> Container allocation
	Services       map[string]string              `json:"services"`   // Service id -> Data
}

// Name for a service container (get by plan)
func PlannedContainerName[S ServiceDriver](plan *Plan, driver S) string {
	return ContainerName(plan.AppName, plan.Profile, driver)
}

// Turn the plan into printable form
func (p *Plan) ToPrintable() (string, error) {
	// Pretty-print the plan as JSON
	encoded, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

// Convert back to a plan from printable form
func FromPrintable(printable string) (*Plan, error) {
	plan := &Plan{}
	err := json.Unmarshal([]byte(printable), plan)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

// Convert every character except for letters and digits directly to _
func EverythingToSnakeCase(s string) string {
	newString := ""
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			newString += string(unicode.ToLower(r))
		} else if !strings.HasSuffix(newString, "_") {
			newString += "_"
		}
	}
	return newString
}
