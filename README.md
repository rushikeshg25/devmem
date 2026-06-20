# DevMem

A searchable memory layer for your git workspaces.

If you work out of many date-based workspace folders (`20-june/`, `20-june2/`, …)
each holding the same repos cloned and worked on independently, it gets hard to
recall *where* you did a piece of work — especially once a workspace is deleted.
DevMem indexes those checkouts into SQLite so you can search commits, branches and
repositories, and recover the uncommitted/unpushed work that would otherwise vanish.

## Install

```bash
make build      # produces ./devmem
make install    # installs into $GOBIN
```

Requires Go 1.24+ and a `git` binary on `PATH`.

## Usage

```bash
# Index every repo under a workspace root (re-run anytime — it's incremental)
devmem scan ~/workspaces

# Search commit messages, branches and repo names
devmem search timezone

# Only surface dirty / unpushed / stashed work matching a term
devmem search --wip report

# Show every indexed checkout of a repo and its recent commits
devmem repo erpai-report

# Recent activity across all repos
devmem timeline
```

The database lives at `~/.devmem.db` by default; override with `--db <path>`.

## How it works

- **Discovery** walks the root and records each git checkout, treating a `.git`
  *file* as a linked worktree (each indexed with its own branch and WIP state).
- **Metadata** comes from shelling out to native `git` (`log --all`, `status`,
  `rev-list`, `stash list`).
- **Storage** is SQLite. The same commit hash can appear in many workspaces, so
  commits are keyed by `(repo_id, hash)` — every workspace's copy is kept distinct.
- **Search** is plain `LIKE` matching; fast and dependency-free at this scale.

## Development

```bash
make check      # fmt + vet + test
make cover      # tests with coverage
```
