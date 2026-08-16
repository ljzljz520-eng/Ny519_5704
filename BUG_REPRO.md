# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	example.com/sine3d/cmd/sine3d	[no test files]
ok  	example.com/sine3d/internal/curve	0.001s
--- FAIL: TestCoordinateWorkflowIncludesSamplingIntervalEnd (0.00s)
    server_test.go:32: last X = 0, want 6.283185307179586
FAIL
FAIL	example.com/sine3d/internal/httpapi	0.002s
ok  	example.com/sine3d/internal/preset	0.001s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/sine3d): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/sine3d): exit `0`
