package utils_test

import (
	"reflect"
	"testing"

	"github.com/jim380/Cendermint/utils"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		mapValue map[string][]string
		index    int
		want     map[string][]string
	}{
		{
			name: "Basic test",
			mapValue: map[string][]string{
				"a": {"value1", "2"},
				"b": {"value2", "3"},
				"c": {"value3", "1"},
			},
			index: 1,
			want: map[string][]string{
				"b": {"value2", "3"},
				"a": {"value1", "2"},
				"c": {"value3", "1"},
			},
		},
		{
			name:     "Empty map",
			mapValue: map[string][]string{},
			index:    1,
			want:     map[string][]string{},
		},
		{
			name: "Single element",
			mapValue: map[string][]string{
				"a": {"value1", "2"},
			},
			index: 1,
			want: map[string][]string{
				"a": {"value1", "2"},
			},
		},
		{
			name: "Non-numeric values",
			mapValue: map[string][]string{
				"a": {"value1", "a"},
				"b": {"value2", "b"},
				"c": {"value3", "c"},
			},
			index: 1,
			want: map[string][]string{
				"a": {"value1", "a"},
				"b": {"value2", "b"},
				"c": {"value3", "c"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.Sort(tt.mapValue, tt.index)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}
