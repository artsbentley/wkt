# wkt: Git Worktree Management CLI

## Overview

wkt is a Go-based CLI tool that simplifies Git worktree workflows by providing convenient commands
for cloning, adding, removing, and managing worktrees, as well as creating branches from GitHub issues.

## Features

- Clone repositories as bare with an adjustable directory structure
- Add new worktrees for local branches
- Add worktrees for remote branches via fzf selection
- Remove and clean up worktrees with branch deletion
- Create new branches from GitHub issues using GitHub CLI and fzf
- (Experimental) Merge or rebase worktrees against a target branch

## Prerequisites

- Go (>= 1.16)
- Git
- fzf (https://github.com/junegunn/fzf)
- GitHub CLI (`gh`) for the `issue` command

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/artsbentley/wkt.git
   cd wkt
   ```

2. Build the binary:
   ```bash
   just build    # or go build -o bin/wkt ./cmd/wkt/main.go
   ```

3. (Optional) Copy to a directory in your PATH:
   ```bash
   cp ./bin/wkt /usr/local/bin/wkt
   ```

## Usage

Run `wkt --help` to see available commands:

```bash
wkt --help
```

### Commands

- `wkt clone <repo-url> [--tree]`  
  Clone a repository as a bare repo under a `.bare` directory.  
  `--tree, -t`: create worktrees for all branches.

- `wkt add <branch-name> [--base <base>] [--upstream]`  
  Create a new worktree sibling directory with a branch based on `<base>` (default `main`).  
  `--upstream, -u`: set upstream and push new branch.

- `wkt add-remote`  
  Select a remote branch via `fzf` to add as a worktree.

- `wkt remove`  
  Select an existing worktree via `fzf`, remove it, delete remote branch, and prune.

- `wkt issue`  
  Create a new worktree from a GitHub issue: choose an issue via `fzf`, branch off `main`, push.

- `wkt merge [--target <branch>] [--rebase] [--push]`  
  (Experimental) Merge or rebase the current worktree onto `<target>` (default `main`).  
  `--rebase, -r`: use rebase.  
  `--push, -p`: push after merge/rebase.

## Examples

```bash
# Clone repo and all branches as worktrees
wkt clone https://github.com/user/repo.git --tree

# Create a feature worktree from main
wkt add feature/cool-feature -b main

# Add a remote branch via fzf
wkt add-remote

# Remove a worktree and its branch
wkt remove

# Create a worktree from GitHub issue 42
wkt issue

# Merge current worktree with main and push
wkt merge --target main --push
```

## Contributing

Contributions are welcome! Feel free to open issues or pull requests.
