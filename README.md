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

## Desktop GUI (Ubuntu)

DevMem ships an optional desktop app, `devmem-gui`, built with
[Fyne](https://fyne.io). It reads the same `~/.devmem.db` database as the CLI and
puts search, at-risk work, the timeline and scanning behind a tabbed window:

- **Search** — type a term and press Enter to find matching commits.
- **At risk** — checkouts with uncommitted, unpushed or stashed work, filterable
  by repo or branch.
- **Timeline** — the most recent commits across every indexed workspace.
- **Scan** — pick a workspace root and index it in the background.

Selecting any row lets you copy its path or open the folder in Files.

### Install

Grab the `devmem-gui` archive for your platform from the
[releases page](https://github.com/rushikeshg25/devmem/releases) (Linux/amd64),
unpack it, and install the binary, icon and launcher entry:

```bash
install -Dm755 devmem-gui ~/.local/bin/devmem-gui
install -Dm644 packaging/icon.png ~/.local/share/icons/hicolor/256x256/apps/devmem.png
install -Dm644 packaging/devmem-gui.desktop ~/.local/share/applications/devmem-gui.desktop
```

DevMem then appears in the GNOME activities launcher, or run `devmem-gui` directly.

### Build from source

The GUI needs CGo and the OpenGL/X11 development headers:

```bash
sudo apt-get install -y gcc libgl1-mesa-dev xorg-dev   # Ubuntu/Debian
make build-gui    # produces ./devmem-gui
make run-gui      # build and launch
```

> The `devmem` CLI itself stays pure Go and needs none of these — only the GUI
> binary requires CGo, which is why it is a separate build.

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
make build      # compile the CLI binary
make build-gui  # compile the desktop GUI (needs CGo + GL/X11 headers)
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
The desktop `devmem-gui` is built for Linux/amd64 (CGo) and published as its own
archive bundled with the `.desktop` entry and icon.

To cut a release:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

## License

MIT. See [LICENSE](LICENSE).
