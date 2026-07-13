package stacksapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
)

func TestFinishJob(t *testing.T) {
	t.Parallel()

	t.Run("successful finish job", func(t *testing.T) {
		t.Parallel()

		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			verifyAuthMethodPath(t, r, "POST", "/stacks/stack-123/jobs/456/finish")

			var params FinishJobRequest
			if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
				t.Fatalf("json.NewDecoder(r.Body).Decode(&params) error: got %v, want nil", err)
			}

			expectedParams := FinishJobRequest{
				ExitStatus:   -10,
				Detail:       "show me the money",
				SignalReason: "stack_error",
			}

			if diff := cmp.Diff(expectedParams, params); diff != "" {
				t.Errorf("request params mismatch (-want +got):\n%s", diff)
			}

			w.Header().Set("X-Custom-Header", "custom-value")
			w.WriteHeader(http.StatusNoContent)
		})
		t.Cleanup(func() { server.Close() })

		req := FinishJobRequest{
			StackKey:     "stack-123",
			JobUUID:      "456",
			ExitStatus:   -10,
			Detail:       "show me the money",
			SignalReason: "stack_error",
		}

		header, err := client.FinishJob(t.Context(), req)
		if err != nil {
			t.Fatalf("client.FinishJob returned an error: %v", err)
		}

		want, got := "custom-value", header.Get("X-Custom-Header")
		if want != got {
			t.Errorf("header X-Custom-Header: want %q, got %q", want, got)
		}
	})

	t.Run("crops detail when exceeds 4KB", func(t *testing.T) {
		t.Parallel()

		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			verifyAuthMethodPath(t, r, "POST", "/stacks/stack-123/jobs/456/finish")

			var params FinishJobRequest
			if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
				t.Fatalf("json.NewDecoder(r.Body).Decode(&params) error: got %v, want nil", err)
			}

			if got, want := len(params.Detail), 4*1024; got != want {
				t.Errorf("len(params.Detail) = %d, want %d", got, want)
			}

			w.WriteHeader(http.StatusNoContent)
		})
		t.Cleanup(func() { server.Close() })

		largeDetail := strings.Repeat("x", 5*1024)

		req := FinishJobRequest{
			StackKey:   "stack-123",
			JobUUID:    "456",
			ExitStatus: -10,
			Detail:     largeDetail,
		}

		_, err := client.FinishJob(t.Context(), req)
		if err != nil {
			t.Fatalf("client.FinishJob error: got %v, want nil", err)
		}
	})

	t.Run("crops multi-byte UTF-8 detail on rune boundary", func(t *testing.T) {
		t.Parallel()

		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			verifyAuthMethodPath(t, r, "POST", "/stacks/stack-123/jobs/456/finish")

			var params FinishJobRequest
			if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
				t.Fatalf("json.NewDecoder(r.Body).Decode(&params) error: got %v, want nil", err)
			}

			if got, max := len(params.Detail), 4*1024; got > max {
				t.Errorf("len(params.Detail) = %d, want <= %d", got, max)
			}

			if !utf8.ValidString(params.Detail) {
				t.Errorf("params.Detail is not valid UTF-8")
			}

			if strings.ContainsRune(params.Detail, utf8.RuneError) {
				t.Errorf("params.Detail contains the Unicode replacement character (U+FFFD), indicating a partial rune was sliced")
			}

			w.WriteHeader(http.StatusNoContent)
		})
		t.Cleanup(func() { server.Close() })

		// "╭" is a 3-byte UTF-8 rune. Repeating it 2000 times produces 6000 bytes,
		// which exceeds the 4KB cap and forces cropping in the middle of the run.
		largeDetail := strings.Repeat("╭", 2000)

		req := FinishJobRequest{
			StackKey:   "stack-123",
			JobUUID:    "456",
			ExitStatus: -10,
			Detail:     largeDetail,
		}

		_, err := client.FinishJob(t.Context(), req)
		if err != nil {
			t.Fatalf("client.FinishJob error: got %v, want nil", err)
		}
	})
}
