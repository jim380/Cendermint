package utils_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/jim380/Cendermint/utils"
)

type ConsensusState struct {
	Result struct {
		RoundState struct {
			Height string `json:"height"`
			Round  int    `json:"round"`
			Step   int    `json:"step"`
		} `json:"round_state"`
	} `json:"result"`
}

func TestParseConsensusOutput(t *testing.T) {
	file, err := os.Open("../testutil/json/dump_consensus_state.json")
	if err != nil {
		t.Fatalf("Failed to open JSON file: %v", err)
	}
	defer file.Close()

	var state ConsensusState
	if err := json.NewDecoder(file).Decode(&state); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	tests := []struct {
		name     string
		target   string
		reg      string
		matchGr  int
		expected float64
	}{
		{
			name:     "Extract height",
			target:   state.Result.RoundState.Height,
			reg:      `(\d+)`,
			matchGr:  1,
			expected: 16545953,
		},
		{
			name:     "Extract round",
			target:   fmt.Sprint(state.Result.RoundState.Round),
			reg:      `(\d+)`,
			matchGr:  1,
			expected: 0,
		},
		{
			name:     "Extract step",
			target:   fmt.Sprint(state.Result.RoundState.Step),
			reg:      `(\d+)`,
			matchGr:  1,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ParseConsensusOutput(tt.target, tt.reg, tt.matchGr)
			if result != tt.expected {
				t.Errorf("ParseConsensusOutput() = %v, want %v", result, tt.expected)
			}
		})
	}
}
