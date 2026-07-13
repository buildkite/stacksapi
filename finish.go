package stacksapi

import (
	"context"
	"fmt"
	"net/http"
	"unicode/utf8"
)

type FinishJobRequest struct {
	StackKey   string `json:"-"`
	JobUUID    string `json:"-"`
	ExitStatus int    `json:"exit_status,omitempty"`
	Detail     string `json:"detail"`
	// SignalReason optionally records why the stack terminated the job (for
	// example "stack_error"), giving the job a machine-readable failure
	// reason alongside the human-readable Detail. Known values match the
	// signal reasons documented for retry rules:
	// https://buildkite.com/docs/pipelines/configure/retry
	SignalReason string `json:"signal_reason,omitempty"`
}

func (c *Client) FinishJob(ctx context.Context, finishJobReq FinishJobRequest, opts ...RequestOption) (http.Header, error) {
	const croppedMessage = "… (Message was cropped because it exceeded the max size)"
	const maxDetailSize = (4 * 1024) - len(croppedMessage)

	if len(finishJobReq.Detail) > maxDetailSize {
		c.logger.Warn("Detail exceeds 4KB limit, cropping", "original_size", len(finishJobReq.Detail), "cropped_size", maxDetailSize)
		// Walk back to the nearest rune boundary so we don't slice a multi-byte
		// UTF-8 character in half; partial runes get encoded as U+FFFD by
		// encoding/json, inflating the payload past the server's 4KB cap.
		cut := maxDetailSize
		for cut > 0 && !utf8.RuneStart(finishJobReq.Detail[cut]) {
			cut--
		}
		finishJobReq.Detail = finishJobReq.Detail[:cut] + croppedMessage
	}

	path := constructPath("/stacks/%s/jobs/%s/finish", finishJobReq.StackKey, finishJobReq.JobUUID)
	req, err := c.newRequest(ctx, http.MethodPost, path, finishJobReq, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	_, header, err := do[struct{}](ctx, c, req)
	if err != nil {
		return nil, fmt.Errorf("finish job: %w", err)
	}

	return header, nil
}
