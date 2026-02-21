# Branching Strategy — chess-engine
**Feature**: chess-engine | **Date**: 2026-02-21 | **Status**: Approved
**Author**: Apex (nw-platform-architect)

---

## 1. Model: Trunk-Based Development

The project uses **Trunk-Based Development** (TBD). The single integration branch is `main`. All code reaches `main` quickly (within one working day at most). There are no long-lived feature branches, no develop/release branches.

**Why TBD for this project**:
- Solo/small team: merge conflicts are rare; coordination overhead of branching models like GitFlow is pure waste
- "main is always releasable" is a first-class requirement; TBD is the model designed to satisfy it
- Continuous integration on every commit to main provides immediate feedback
- Release is via tags, not branches; no release branch management required

---

## 2. Branch Rules

### 2.1 `main` Branch (Trunk)

| Rule | Enforcement |
|---|---|
| Direct push prohibited | GitHub branch protection: "Require a pull request before merging" |
| CI must pass before merge | GitHub branch protection: "Require status checks to pass before merging" |
| 1 approval required on PRs | GitHub branch protection: "Require approvals: 1" (waived for solo developer with self-review) |
| No force push | GitHub branch protection: "Allow force pushes: disabled" |
| No branch deletion | GitHub branch protection: "Allow deletion: disabled" |
| Linear history | GitHub branch protection: "Require linear history" (squash or rebase merge only) |

**Required status checks** (must pass before PR merge):
- `format`
- `vet`
- `lint`
- `test-coverage`
- `perft`
- `acceptance`
- `benchmarks`
- `build`

### 2.2 Feature Branches

**Naming convention**:

| Type | Pattern | Example |
|---|---|---|
| Feature | `feat/{short-description}` | `feat/move-generator` |
| Bug fix | `fix/{short-description}` | `fix/castling-through-check` |
| Chore | `chore/{short-description}` | `chore/update-golangci-lint` |
| Documentation | `docs/{short-description}` | `docs/benchmark-guide` |
| Test | `test/{short-description}` | `test/perft-kiwipete` |

**Lifetime**: Feature branches live < 1 day. A branch open longer than 24 hours is a signal that scope is too large. Break the work into smaller increments.

**Deletion**: Branches are deleted after PR merge (GitHub setting: "Automatically delete head branches").

**No commits directly to feature branches from multiple authors**: Each feature branch belongs to one developer. No shared feature branches.

---

## 3. Commit Conventions

All commits follow the **Conventional Commits** specification (https://www.conventionalcommits.org/).

### 3.1 Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### 3.2 Types

| Type | Purpose | Triggers Changelog Section |
|---|---|---|
| `feat` | New feature | "Features" |
| `fix` | Bug fix | "Bug Fixes" |
| `perf` | Performance improvement | "Performance" |
| `test` | Adding or correcting tests | not in changelog |
| `docs` | Documentation only | not in changelog |
| `chore` | Build, deps, tooling | not in changelog |
| `refactor` | Code restructure (no behaviour change) | not in changelog |
| `ci` | CI workflow changes | not in changelog |

### 3.3 Scope (Optional)

Scope identifies the affected package: `chess`, `engine`, `tui`, `web`, `cmd`.

Examples:
```
feat(chess): implement castling move generation
fix(engine): respect movetime deadline in quiescence search
perf(chess): reduce allocations in LegalMoves
test(chess): add perft depth-4 for Kiwipete position
chore(ci): pin golangci-lint to v1.58.0
docs: add benchmark interpretation guide
```

### 3.4 Breaking Changes

Breaking changes (not expected in v0.x) are signalled by `BREAKING CHANGE:` in the footer or `!` after the type:

```
feat(chess)!: change Apply() signature to return (Game, bool) instead of (Game, error)

BREAKING CHANGE: callers must update error handling at Apply() call sites
```

### 3.5 Commit Validation

A Git commit-msg hook validates the Conventional Commits format. Install via:
```bash
# .git/hooks/commit-msg (future: add to onboarding docs)
# Validates: type, optional scope, colon, space, description
```

CI does not enforce commit message format automatically (no commitlint in pipeline) to avoid friction on a solo project. Convention is by team agreement.

---

## 4. Release Workflow

Releases are created by tagging `main`. There is no release branch.

### 4.1 Versioning

**Semantic Versioning** (SemVer): `vMAJOR.MINOR.PATCH`

| Increment | Trigger |
|---|---|
| PATCH (`v0.1.1`) | Bug fix commits merged to main |
| MINOR (`v0.2.0`) | Feature commits merged to main |
| MAJOR (`v1.0.0`) | Breaking change or product stability milestone |

During `v0.x`, MINOR increments are used for both features and breaking changes per SemVer spec.

### 4.2 Tagging Procedure

```bash
# Ensure you are on main and CI has passed
git checkout main
git pull origin main

# Create annotated tag (annotated tags are required; they carry metadata)
git tag -a v0.2.0 -m "chess package: complete move generator with perft validation"

# Push tag — this triggers the release workflow
git push origin v0.2.0
```

**Lightweight tags are not used**. Annotated tags only; `git tag -a` is enforced by convention.

### 4.3 Release Checklist

Before tagging:
- [ ] All acceptance criteria for the milestone are passing in CI
- [ ] Perft at depths 1-4 passes on main
- [ ] Benchmark NPS >= 100k in CI
- [ ] `CHANGELOG.md` reviewed (auto-generated by release workflow)
- [ ] `go.mod` version (`go 1.22` or higher) is accurate

### 4.4 Pre-releases

For milestone previews: `v0.1.0-alpha.1`, `v0.1.0-beta.1`. The release workflow automatically marks these as pre-release in GitHub Releases based on the presence of `-alpha` or `-beta` in the tag name.

---

## 5. GitHub Repository Settings

The following settings must be configured in the GitHub repository UI (Settings > Branches):

```
Branch protection rule for: main
  [x] Require a pull request before merging
      Required approvals: 1
      [x] Dismiss stale pull request approvals when new commits are pushed
  [x] Require status checks to pass before merging
      [x] Require branches to be up to date before merging
      Status checks required:
        - format
        - vet
        - lint
        - test-coverage
        - perft
        - acceptance
        - benchmarks
        - build
  [x] Require linear history
  [ ] Allow force pushes (disabled)
  [ ] Allow deletions (disabled)

General settings:
  [x] Automatically delete head branches
  Merge options:
    [ ] Allow merge commits
    [x] Allow squash merging  (default: "Pull request title and commit details")
    [x] Allow rebase merging
```

**Squash merge is the default** for PRs. Short-lived feature branches produce a single, clean commit on main. The squash commit title becomes the Conventional Commits message for the changelog.

---

## 6. Workflow: Typical Development Cycle

```
1. Pull latest main
   git checkout main && git pull origin main

2. Create feature branch
   git checkout -b feat/alpha-beta-search

3. Work in small commits (internal to branch, any format)
   git commit -m "wip: scaffold search.go"
   git commit -m "wip: add negamax frame"
   git commit -m "wip: add alpha-beta pruning"

4. Push branch and open PR
   git push -u origin feat/alpha-beta-search
   gh pr create --title "feat(engine): implement alpha-beta search with iterative deepening"

5. CI runs all gates on the PR branch

6. Squash merge when green
   # GitHub UI or:
   gh pr merge --squash --delete-branch

7. The squash commit on main has the Conventional Commits title

8. Tag when milestone is complete
   git tag -a v0.2.0 -m "engine: alpha-beta search complete"
   git push origin v0.2.0
```

---

## 7. Hotfix Procedure

Since main is always releasable, hotfixes follow the same flow as regular features:

```
1. Create fix branch from main (NOT from a release tag)
   git checkout main && git pull
   git checkout -b fix/castling-through-check

2. Fix, test locally, push PR

3. CI gates run as normal

4. Merge to main

5. If the fix needs a patch release, tag immediately
   git tag -a v0.1.1 -m "fix(chess): correct castling through check detection"
   git push origin v0.1.1
```

There is no separate hotfix branch or cherry-pick from release branches. Trunk-Based Development eliminates the need for them.
