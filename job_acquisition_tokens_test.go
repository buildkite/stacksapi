package stacksapi

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestIssueJobAcquisitionTokens(t *testing.T) {
	t.Parallel()

	expiresAt := time.Date(2026, time.August, 20, 0, 15, 0, 0, time.UTC)

	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		verifyAuthMethodPath(t, r, "POST", "/stacks/stack-123/job-acquisition-tokens")

		var params IssueJobAcquisitionTokensRequest
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			t.Fatalf("failed to decode issue job acquisition tokens request body: %v", err)
		}

		expectedParams := IssueJobAcquisitionTokensRequest{
			JobUUIDs:             []string{"job-1", "job-2"},
			TokenLifetimeSeconds: 600,
		}
		if diff := cmp.Diff(expectedParams, params); diff != "" {
			t.Errorf("request params mismatch (-want +got):\n%s", diff)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "custom-value")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(IssueJobAcquisitionTokensResponse{
			JobAcquisitionTokens: []JobAcquisitionToken{
				{JobUUID: "job-1", Token: "bkjat_token", ExpiresAt: expiresAt},
			},
			NotIssued: []string{"job-2"},
		}); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	})
	t.Cleanup(server.Close)

	resp, header, err := client.IssueJobAcquisitionTokens(t.Context(), IssueJobAcquisitionTokensRequest{
		StackKey:             "stack-123",
		JobUUIDs:             []string{"job-1", "job-2"},
		TokenLifetimeSeconds: 600,
	})
	if err != nil {
		t.Fatalf("client.IssueJobAcquisitionTokens error = %v, want nil", err)
	}

	expectedResp := &IssueJobAcquisitionTokensResponse{
		JobAcquisitionTokens: []JobAcquisitionToken{
			{JobUUID: "job-1", Token: "bkjat_token", ExpiresAt: expiresAt},
		},
		NotIssued: []string{"job-2"},
	}
	if diff := cmp.Diff(expectedResp, resp); diff != "" {
		t.Errorf("issue job acquisition tokens response mismatch (-want +got):\n%s", diff)
	}
	if got := header.Get("X-Custom-Header"); got != "custom-value" {
		t.Errorf("header X-Custom-Header = %q, want %q", got, "custom-value")
	}
}

func TestIssueJobAcquisitionTokensOmitsDefaultLifetime(t *testing.T) {
	t.Parallel()

	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if _, ok := body["token_lifetime_seconds"]; ok {
			t.Error("token_lifetime_seconds was present, want omitted")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"job_acquisition_tokens":[],"not_issued":[]}`))
	})
	t.Cleanup(server.Close)

	_, _, err := client.IssueJobAcquisitionTokens(t.Context(), IssueJobAcquisitionTokensRequest{
		StackKey: "stack-123",
		JobUUIDs: []string{"job-1"},
	})
	if err != nil {
		t.Fatalf("client.IssueJobAcquisitionTokens error = %v, want nil", err)
	}
}
