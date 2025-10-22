# stacksapi

Go client library for the Buildkite Stacks API. This package enables custom stack implementations to communicate with Buildkite's agent infrastructure, handling job scheduling, lifecycle management, and stack registration.

## Installation

```bash
go get github.com/buildkite/stacksapi
```

## Authentication

Requires a **Buildkite Cluster Token** (not a REST or GraphQL API token, or an unclustered Agent Registration Token). The token is passed to `NewClient()`.

## Usage

```go
import "github.com/buildkite/stacksapi"

client, err := stacksapi.NewClient(
    os.Getenv("BUILDKITE_CLUSTER_TOKEN"),
    stacksapi.WithLogger(logger),
)
```

## Client Options

- `WithLogger(logger *slog.Logger)` - Configure structured logging. To integrate with non-`slog` loggers, we recommend using log adapters like [`slog-zap`](https://github.com/samber/slog-zap), [`slog-zerolog`](https://github.com/samber/slog-zerolog) and friends.
- `WithBaseURL(url *url.URL)` - Override the default API endpoint
- `WithHTTPClient(client *http.Client)` - Use a custom HTTP client
- `WithRetrierOptions(...roko.RetrierOpt)` - Configure client-wide retry behavior. The default retry behaviour is an exponential backoff with jitter, retrying on 429 and 5xx responses. Retry behaviour can also be customized per-request using `WithRetrier` and `WithNoRetry` request options.
- `LogHTTPPayloads()` - Enable request/response payload logging (may log sensitive data)
- `PrependToUserAgent(prefix string)` - Customize the User-Agent header. Note that when using this library, the user agent will always end with `stacksapi v$SOME_VERSION`

### Register a Stack

```go
stack, _, err := client.RegisterStack(ctx, stacksapi.RegisterStackRequest{
    Key:      "my-stack",
    Type:     stacksapi.StackTypeCustom,
    QueueKey: stacksapi.DefaultQueue,
    Metadata: map[string]string{"version": "1.0"},
})
```

### Fetch Scheduled Jobs

```go
jobs, _, err := client.ScheduledJobs(ctx, stacksapi.ScheduledJobsRequest{
    StackKey: "my-stack",
})
```

### Get Job Details

```go
job, _, err := client.GetJob(ctx, stacksapi.GetJobRequest{
    StackKey: "my-stack",
    JobUUID:  jobUUID,
})
```

### Finish a Job

```go
_, err = client.FinishJob(ctx, stacksapi.FinishJobRequest{
    StackKey:   "my-stack",
    JobUUID:    jobUUID,
    ExitStatus: 0,
    Detail:     "Job completed successfully",
})
```

### Deregister a Stack

```go
_, err = client.DeregisterStack(ctx, "my-stack")
```

## Request Options

Per-request options can be passed to any API method:

- `WithRetrier(retrier *roko.Retrier)` - Override the default retry behavior
- `WithNoRetry()` - Disable retries for the request

## Stack Types

- `StackTypeCustom` - For third-party stack implementations (recommended)
- `StackTypeKubernetes` - Reserved for first-party `agent-stack-k8s`
- `StackTypeElastic` - Reserved for first-party `elastic-stack-for-aws`

## Examples

See [example/main.go](example/main.go) for a complete working example.

## License

See [LICENSE.txt](LICENSE.txt).
