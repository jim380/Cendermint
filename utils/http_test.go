package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jim380/Cendermint/utils"
)

func TestHttpQuery(t *testing.T) {
	tests := []struct {
		name       string
		mockStatus int
		mockBody   string
		want       string
		wantErr    bool
	}{
		{
			name:       "Valid response",
			mockStatus: http.StatusOK,
			mockBody:   `{"message": "success"}`,
			want:       `{"message": "success"}`,
			wantErr:    false,
		},
		{
			name:       "Not found",
			mockStatus: http.StatusNotFound,
			mockBody:   `{"error": "not found"}`,
			want:       `{"error": "not found"}`,
			wantErr:    false,
		},
		{
			name:       "Internal server error",
			mockStatus: http.StatusInternalServerError,
			mockBody:   `{"error": "internal server error"}`,
			want:       `{"error": "internal server error"}`,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockBody))
			}))
			defer mockServer.Close()

			got, err := utils.HttpQuery(mockServer.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("HttpQuery() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
