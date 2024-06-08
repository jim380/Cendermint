package rest_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
)

func TestGetConspubMonikerMap(t *testing.T) {
	data, err := os.ReadFile("../testutil/json/validators.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	// Override the RESTAddr
	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	// Call the function
	result := rest.GetConspubMonikerMap()

	// Expected result
	expected := map[string]string{
		"uEUR1gpesU4bnSWL2TOXOf3org2mCYhQHMYkiCJyMD4=": "Ubik Capital",
		"Qajjf1kiAJ0M1UcH1TSUYLP13kgE128Av1XmGQO711c=": "GAME",
		"XiGz/D6eg3KdjaFB0uYIJwkOTW5xZcFRxJmHcQYB3zg=": "WeStaking",
	}

	// Compare the result with the expected map
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetConspubMonikerMap() = %v, want %v", result, expected)
	}
}
