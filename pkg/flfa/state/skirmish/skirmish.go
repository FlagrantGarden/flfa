package skirmish

import (
	"fmt"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state"
	"github.com/google/uuid"
)

type Skirmish struct {
	Scenario  string
	Attackers []string
	Defenders []string
	Companies []data.Company
	Updates   string
}

func (skirmish Skirmish) Initialize() *Skirmish {
	if skirmish.Scenario != "" {
		return &skirmish
	}
	return &Skirmish{}
}

type Updates struct {
	Company       string
	Type          string
	Subtype       string
	Actor         uuid.UUID
	Target        uuid.UUID
	StartLocation Location `mapstructure:"start_location"`
	EndLocation   Location `mapstructure:"end_location"`
	Result        string
}

type Location struct {
	X int
	Y int
}

type Result struct {
	Activation    string
	InflictHits   int    `mapstructure:"inflict_hits"`
	ReceiveHits   int    `mapstructure:"receive_hits"`
	ActorResolve  string `mapstructure:"actor_resolve"`
	TargetResolve string `mapstructure:"target_resolve"`
}

type TestResult string

const (
	Pass TestResult = "pass"
	Fail TestResult = "fail"
)

func ValidTestResults() []TestResult {
	return []TestResult{Pass, Fail}
}

func (result TestResult) Passed() (passed bool, err error) {
	var validResultList []string
	for _, validResult := range ValidTestResults() {
		validResultList = append(validResultList, string(validResult))
	}

	if result == Pass {
		return true, nil
	} else if result == Fail {
		return false, nil
	}

	return false, fmt.Errorf("unexpected test result value '%s', should be one of: %s", string(result), strings.Join(validResultList, ", "))
}

func Kind() *state.Kind {
	return &state.Kind{
		Name:       "skirmish",
		FolderName: "skirmishes",
	}
}
