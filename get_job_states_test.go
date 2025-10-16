package stacksapi

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetJobStates(t *testing.T) {
	t.Parallel()

	t.Run("successful get job states", func(t *testing.T) {
		t.Parallel()

		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			verifyAuthMethodPath(t, r, "POST", "/stacks/stack-123/jobs/get-states")

			var params GetJobStatesRequest
			if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
				t.Fatalf("json.NewDecoder(r.Body).Decode(&params) error: got %v, want nil", err)
			}

			expectedParams := GetJobStatesRequest{
				JobUUIDs: []string{"job-1", "job-2", "job-3"},
			}

			if diff := cmp.Diff(expectedParams, params); diff != "" {
				t.Errorf("request params mismatch (-want +got):\n%s", diff)
			}

			response := &GetJobStatesResponse{
				States: map[string]string{
					"job-1": "running",
					"job-2": "finished",
					"job-3": "failed",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Fatalf("json.NewEncoder(w).Encode error: got %v, want nil", err)
			}
		})
		t.Cleanup(func() { server.Close() })

		req := GetJobStatesRequest{
			StackKey: "stack-123",
			JobUUIDs: []string{"job-1", "job-2", "job-3"},
		}

		response, _, err := client.GetJobStates(t.Context(), req)
		if err != nil {
			t.Fatalf("client.GetJobStates error: got %v, want nil", err)
		}

		expectedResponse := &GetJobStatesResponse{
			States: map[string]string{
				"job-1": "running",
				"job-2": "finished",
				"job-3": "failed",
			},
		}

		if diff := cmp.Diff(expectedResponse, response); diff != "" {
			t.Errorf("response mismatch (-want +got):\n%s", diff)
		}
	})
}
