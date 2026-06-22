# DevMem

A searchable memory layer for your git workspaces.

If you work out of many date-based workspace folders (`20-june/`, `20-june2/`, ...)
each holding the same repos cloned and worked on independently, it gets hard to
recall where you did a piece of work, especially once a workspace is deleted.
DevMem indexes those checkouts into SQLite so you can search commits, branches and
repositories, and recover the uncommitted or unpushed work that would otherwise vanish.

## Requirements

- Go 1.24 or newer
- A `git` binary on `PATH`

## Install

Install the latest release binary (Linux/macOS) and add it to your PATH:

```bash
curl -fsSL https://raw.githubusercontent.com/rushikeshg25/devmem/master/install.sh | sh
```

This drops `devmem` into `$HOME/.local/bin` (override with `DEVMEM_BIN_DIR`) and
tells you if that directory isn't on your PATH. Prebuilt archives for every
platform are on the [releases page](https://github.com/rushikeshg25/devmem/releases).

Or build directly with Go:

```bash
go install github.com/rushikeshg25/devmem@latest
```

Or from a checkout:

```bash
make build      # produces ./devmem in the project directory
make install    # installs the binary into $GOBIN
```

## Quick start

```bash
devmem scan ~/workspaces      # index everything once
devmem search timezone        # find where you worked on it later
```

## Commands

All commands accept the global `--db <path>` flag (default `~/.devmem.db`).

### scan

```bash
devmem scan <root>
```

Discovers every git checkout under `<root>` and indexes its repositories, branches,
commits and working-tree state. Re-running is safe and incremental: only new commits
are added. Linked worktrees are indexed as separate checkouts with their own branch.

```text
Indexed 3 repos across 2 workspaces (5 new commits)
```

### search

```bash
devmem search <term> [--limit N] [--wip]
```

Substring search across commit messages, branch and ref names, and repository names.

| Flag      | Default | Description                                                  |
| --------- | ------- | ------------------------------------------------------------ |
| `--limit` | 50      | Maximum number of results.                                   |
| `--wip`   | false   | Show only checkouts with uncommitted, unpushed or stashed work matching the term, instead of commits. |

```text
2026-06-20  erpai-report  [feat/report-timezone]
    161054b  fix: timezone handling in exports
    /home/you/workspaces/20-june
```

### repo

```bash
devmem repo <name>
```

Shows every indexed checkout of a repository (across all workspaces), each with its
current branch, working-tree status and most recent commits.

```text
erpai-report  [feat/report-timezone] dirty ahead+1
    /home/you/workspaces/20-june2/erpai-report
    2026-06-20  e1ebc33  wip: local only change
    2026-06-20  161054b  add kafka retries
```

### timeline

```bash
devmem timeline [--limit N]
```

Recent commit activity across all indexed repositories, newest first.

| Flag      | Default | Description                          |
| --------- | ------- | ------------------------------------ |
| `--limit` | 30      | Maximum number of commits to show.   |

### prune

```bash
devmem prune
```

Removes index entries for checkouts whose working directory no longer exists on disk,
for example after deleting a workspace folder, and drops any workspaces left empty.

```text
removed /home/you/workspaces/20-june2/erpai-report
Pruned 1 repos and 1 empty workspaces
```

### Global flag

| Flag   | Default        | Description                          |
| ------ | -------------- | ------------------------------------ |
| `--db` | `~/.devmem.db` | Path to the devmem SQLite database.  |

## How it works

- Discovery walks the root and records each git checkout, treating a `.git`
  file as a linked worktree (each indexed with its own branch and working-tree state).
- Metadata comes from shelling out to native `git` (`log --all`, `status`,
  `rev-list`, `stash list`).
- Storage is SQLite. The same commit hash can appear in many workspaces, so commits
  are keyed by `(repo_id, hash)` and every workspace's copy is kept distinct.
- Search is plain `LIKE` matching, which is fast and dependency-free at this scale.

## Development

```bash
make build      # compile the binary
make test       # run the test suite
make cover      # run tests with a coverage summary
make check      # fmt, vet and test
make help       # list all targets
```

## Releasing

Releases are built and published automatically by GitHub Actions whenever a
version tag is pushed. [GoReleaser](https://goreleaser.com) cross-compiles the
binaries (linux/darwin/windows × amd64/arm64), stamps the version into
`--version`, and attaches the archives plus `checksums.txt` to a GitHub Release.

To cut a release:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

## License

MIT. See [LICENSE](LICENSE).
