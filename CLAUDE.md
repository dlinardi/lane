# CLAUDE.md

> Context file for Claude Code working on the `lane` project.

## Project Overview

**lane** is a CLI tool for managing parallel development lanes using git worktrees and tmux. It enables developers to run multiple AI coding sessions (Claude Code, Aider, Codex, etc.) simultaneously on different tickets/branches without context switching confusion.

**Core workflow:**
```bash
lane init                           # Initialize in a repo
lane new air-800 track-referrals    # Create worktree + branch + tmux window + launch agent
lane go air-801                     # Switch to another lane's tmux window
lane ls                             # List active lanes
lane done air-800                   # Clean up when merged
```

**Repository:** `github.com/dlinardi/lane`  
**Linear Project:** https://linear.app/lanecli

---

## Tech Stack

| Component | Technology | Notes |
|-----------|------------|-------|
| Language | Go 1.22+ | Single binary distribution |
| CLI Framework | [Cobra](https://github.com/spf13/cobra) | Commands, flags, help text |
| Config | [Viper](https://github.com/spf13/viper) | TOML config files |
| TUI/Styling | [Charm](https://charm.sh) | Lip Gloss, Bubble Tea, Huh, Log |
| Git Operations | `exec.Command("git", ...)` | NOT go-git (incomplete worktree support) |
| Testing | Standard `testing` package | + testify for assertions |

---

## Project Structure

```
lane/
├── cmd/                      # Cobra commands (one file per command)
│   ├── root.go               # Root command, global flags, version
│   ├── init.go               # lane init
│   ├── new.go                # lane new
│   ├── ls.go                 # lane ls
│   ├── go.go                 # lane go
│   ├── status.go             # lane status
│   ├── done.go               # lane done
│   ├── shell.go              # lane shell
│   └── config.go             # lane config
├── internal/                 # Private packages
│   ├── config/               # Configuration loading/saving/validation
│   │   ├── types.go          # Config structs
│   │   ├── defaults.go       # Default values
│   │   ├── loader.go         # Load from TOML
│   │   ├── validate.go       # Validation logic
│   │   └── context.go        # LaneContext detection
│   ├── git/                  # Git operations (via exec)
│   │   ├── exec.go           # Git command executor wrapper
│   │   ├── repo.go           # Repository detection
│   │   ├── worktree.go       # Worktree create/list/remove
│   │   ├── branch.go         # Branch operations
│   │   ├── status.go         # Working directory status
│   │   └── testutil/         # Test helpers (temp repos)
│   ├── tmux/                 # tmux operations (via exec)
│   │   ├── detect.go         # Environment detection
│   │   ├── exec.go           # tmux command executor
│   │   └── window.go         # Window create/switch/close
│   ├── agent/                # AI agent launching
│   │   └── launch.go         # Launch configured agent
│   ├── lane/                 # Core lane logic
│   │   ├── types.go          # Lane, LaneStatus structs
│   │   ├── validate.go       # Input validation
│   │   ├── branch.go         # Branch name generation
│   │   ├── create.go         # Lane creation
│   │   ├── list.go           # Lane listing
│   │   ├── status.go         # Lane status
│   │   └── remove.go         # Lane removal
│   └── ui/                   # Terminal UI components
│       ├── theme.go          # Color theme
│       ├── table.go          # Table rendering
│       ├── status.go         # Status box rendering
│       ├── prompt.go         # Confirmation prompts
│       ├── log.go            # Structured logging
│       └── error.go          # Error display
├── main.go                   # Entry point: calls cmd.Execute()
├── go.mod
├── go.sum
├── Makefile                  # build, test, lint, install targets
├── .goreleaser.yml           # Release automation
├── .golangci.yml             # Linter config
└── CLAUDE.md                 # This file
```

---

## Coding Conventions

### Go Style

- Follow standard Go conventions (gofmt, go vet)
- Use `golangci-lint` for linting
- Error messages: lowercase, no punctuation (`"failed to create worktree"`)
- Wrap errors with context: `fmt.Errorf("create worktree: %w", err)`
- Define custom error types in each package: `var ErrNotGitRepo = errors.New("not a git repository")`

### Package Design

- `cmd/` — Thin command handlers, delegate to `internal/` packages
- `internal/` — All business logic, testable in isolation
- Each `internal/` package should be independently testable
- Avoid circular dependencies between internal packages

### Naming

- Commands: `lane <verb>` (init, new, ls, go, done, status, config)
- Packages: singular (`config`, `git`, `lane`, `tmux`)
- Files: lowercase, descriptive (`worktree.go`, `branch.go`)
- Test files: `*_test.go` in same package

### Testing

- Every package needs tests
- Use table-driven tests where appropriate
- Git tests should use `internal/git/testutil` for temp repos
- tmux tests: unit test with mocks, integration tests skippable in CI
- Target: 80%+ coverage on core packages (`git`, `config`, `lane`)

```go
// Example test structure
func TestCreateWorktree_Success(t *testing.T) {
    // Arrange
    repoPath := testutil.CreateTempRepo(t)
    testutil.CreateCommit(t, repoPath, "initial")
    
    // Act
    err := CreateWorktree(repoPath, "/tmp/wt", "feat/test", "main")
    
    // Assert
    assert.NoError(t, err)
    assert.DirExists(t, "/tmp/wt")
}
```

---

## Git Workflow

### Branch Naming

All branches follow this pattern:
```
<type>/<ticket-id>/<description>
```

Examples:
- `feat/lan-5/init-setup`
- `fix/lan-42/null-pointer`
- `docs/lan-75/readme-polish`

Valid types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

### Commit Messages

Format:
```
<type>: <description>

[LAN-XX]
```

Examples:
- `feat: implement worktree creation`
- `fix: handle missing base branch`
- `test: add git executor tests`

The Linear issue ID should be included (husky hook auto-prepends from branch name).

### PR Flow

1. Create branch from `main`: `git checkout -b feat/lan-XX/description`
2. Implement changes with tests
3. Push and create PR to `main`
4. CI must pass (lint, test, build)
5. Squash merge

---

## Key Architectural Decisions

### 1. exec.Command for Git (not go-git)

**Decision:** Use `exec.Command("git", ...)` instead of go-git library.

**Rationale:**
- go-git has incomplete worktree support
- Users already have git installed
- Easier debugging (run same commands manually)
- More predictable behavior

**Implementation:** `internal/git/exec.go` provides a `GitExecutor` interface for testability.

### 2. Worktree Location Options

**Decision:** Support both nested (`.lane/worktrees/`) and sibling (`../<repo>--lanes/`) worktree locations.

**Rationale:**
- Nested is simpler (everything in one place)
- Sibling avoids nested repo concerns
- User choice during `lane init`

**Config:** `worktree_location = "nested" | "sibling"`

### 3. tmux Required for Full Experience

**Decision:** tmux integration is core but graceful degradation when not in tmux.

**Rationale:**
- tmux enables the parallel workflow that's the main value prop
- But basic lane management should work without it

**Implementation:** Check `$TMUX` env var; print helpful message if not in tmux.

### 4. Configuration in Repo (not global)

**Decision:** Config lives in `.lane/config.toml` within each repo.

**Rationale:**
- Different repos may need different settings
- Config travels with the repo
- No global state to manage

**Future:** May add `~/.config/lane/config.toml` for global defaults.

---

## Common Tasks

### Adding a New Command

1. Create `cmd/<command>.go`
2. Define Cobra command with `Use`, `Short`, `Long`, `Example`
3. Add flags with defaults from config
4. Implement `RunE` function that delegates to `internal/` packages
5. Register in `cmd/root.go` init: `rootCmd.AddCommand(newCmd)`
6. Add tests in `cmd/<command>_test.go`

### Adding a Git Operation

1. Add function to appropriate file in `internal/git/`
2. Use `GitExecutor` interface for running commands
3. Return typed errors (`ErrBranchExists`, etc.)
4. Add tests using `testutil.CreateTempRepo()`

### Adding UI Components

1. Add to `internal/ui/`
2. Use Lip Gloss for styling
3. Respect `NO_COLOR` env var and `ui.color` config
4. Keep it simple — we're a CLI, not a TUI app

---

## Environment & Config

### Config File Structure

```toml
# .lane/config.toml

default_base = "develop"
default_type = "feat"
branch_pattern = "{type}/{ticket}/{description}"
worktree_location = "nested"  # or "sibling"
worktree_dir = ".lane/worktrees"

[agent]
enabled = true
command = "claude"

[tmux]
auto_session = true
session_name = "lanes"
window_pattern = "{ticket}"

[ui]
color = true
status_format = "detailed"
```

### Environment Variables

- `TMUX` — Detected to check if inside tmux
- `EDITOR` — Used by `lane config edit`
- `NO_COLOR` — Disables colored output

---

## Linear Integration

This project uses Linear for issue tracking:
- Workspace: `lanecli`
- Project: `lane`
- Issue prefix: `LAN-`

Issues are organized into sprints (Sprint 0–8), each with specific goals.

When working on a task:
1. Check the Linear issue for acceptance criteria
2. Create branch: `<type>/lan-<number>/<description>`
3. Reference issue in commits and PR

---

## Quick Reference

### Make Targets

```bash
make build      # Build to ./bin/lane
make test       # Run all tests
make lint       # Run golangci-lint
make install    # Install to $GOPATH/bin
make clean      # Remove build artifacts
```

### Useful Commands

```bash
# Run a specific test
go test -v -run TestCreateWorktree ./internal/git/

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Build and test locally
make build && ./bin/lane --version

# Format code
gofmt -w .
```

---

## Current Sprint Focus

Check Linear for the current sprint, but the general progression is:

- **Sprint 0:** Project setup, CI/CD, CLI skeleton
- **Sprint 1:** Core git operations (worktree, branch)
- **Sprint 2:** Configuration system, `lane init`
- **Sprint 3:** Basic commands (new, ls, done) without tmux
- **Sprint 4:** tmux integration
- **Sprint 5:** UI polish with Charm
- **Sprint 6:** Config command, shell completions
- **Sprint 7:** Docs, testing, release prep
- **Sprint 8:** Launch

---

## Getting Help

- **Spec document:** `docs/lane-spec.md` (full product spec)
- **Sprint plan:** `docs/lane-sprint-plan.md` (all tasks)
- **Linear:** https://linear.app/lanecli

When in doubt about implementation details, check the spec or ask for clarification.
