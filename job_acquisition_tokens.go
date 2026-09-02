package stacksapi

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// IssueJobAcquisitionTokensRequest is the input type for [IssueJobAcquisitionTokens]. The StackKey and JobUUIDs fields are required.
type IssueJobAcquisitionTokensRequest struct {
	StackKey             string   `json:"-"`                               // The key of the stack calling this endpoint. Required.
	JobUUIDs             []string `json:"job_uuids"`                       // The UUIDs of the reserved jobs to issue tokens for. Required.
	TokenLifetimeSeconds int      `json:"token_lifetime_seconds,omitzero"` // The token lifetime in seconds. Optional, defaults to 900 (15 minutes) if not set.
}

// JobAcquisitionToken is a short-lived credential that allows an agent to acquire one specific job.
type JobAcquisitionToken struct {
	JobUUID   string    `json:"job_uuid"`
	Token     string    `json:"job_acquisition_token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IssueJobAcquisitionTokensResponse is the output type for [IssueJobAcquisitionTokens].
type IssueJobAcquisitionTokensResponse struct {
	JobAcquisitionTokens []JobAcquisitionToken `json:"job_acquisition_tokens"`
	NotIssued            []string              `json:"not_issued"`
}

// IssueJobAcquisitionTokens issues short-lived tokens for reserved jobs. A request can partially succeed; callers
// should match tokens to jobs using JobUUID and must not start workloads for jobs listed in NotIssued.
func (c *Client) IssueJobAcquisitionTokens(ctx context.Context, issueReq IssueJobAcquisitionTokensRequest, opts ...RequestOption) (*IssueJobAcquisitionTokensResponse, http.Header, error) {
	path := constructPath("/stacks/%s/job-acquisition-tokens", issueReq.StackKey)
	req, err := c.newRequest(ctx, http.MethodPost, path, issueReq, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	issueResp, header, err := do[IssueJobAcquisitionTokensResponse](ctx, c, req)
	if err != nil {
		return nil, nil, fmt.Errorf("issue job acquisition tokens: %w", err)
	}

	return issueResp, header, nil
}
