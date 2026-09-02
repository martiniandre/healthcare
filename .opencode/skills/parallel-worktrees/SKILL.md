---
name: parallel-worktrees
description: Coordinates multiple simultaneous agent sessions working in isolated git worktrees (../healthcare-worktrees/) so one feature never affects another. Claims worktree slugs/branches and reserves migration numbers/identifiers in a shared registry via node scripts/worktrees.mjs, enforces the additive-anchored edit protocol for shared files (policy.go, main.go, routes.tsx, AppSidebar, i18n), rebase-before-PR with renumber-on-collision, and releases identifiers after merge. Use when any feature/fix/refactor task begins, when running multiple features in parallel, or whenever the AGENTS.md worktree isolation rule applies.
---

# Parallel Worktrees — Feature Isolation Protocol

Every feature, fix or refactor runs in a **dedicated worktree** at `../healthcare-worktrees/{slug}`. When multiple features may be in flight at the same time, this skill guarantees they never affect each other.

## Why isolation is real — and where it is not

- Worktrees give **filesystem isolation**: an agent in worktree A physically cannot write into worktree B's files. This is already bulletproof.
- The only thing that crosses over is at **merge time**, when two branches touching the same files both land on `main`. Two failure modes:
  1. **Identifier collisions** (breaking): two features both create `020_*.up.sql`, or the same module folder, route path, branch or slug. After both merge, `main` has duplicates — golang-migrate breaks at runtime.
  2. **Adjacent-line insertions** (ordinary conflicts): both features append to the same spot of a shared map/file. Git flags a conflict; a careless resolve can silently drop an entry.
- Both are prevented here: a **shared registry** for identifiers, and an **additive-anchored edit protocol** for shared files, plus **rebase-before-PR**.

## The shared registry

`../healthcare-worktrees/registry.json` — a sibling of the main repo, outside every worktree. All worktrees read the **same** file, it is never committed, and it carries no merge semantics. It records every active worktree and the identifiers (migration numbers, etc.) it has reserved.

Never edit the JSON by hand — use the helper:

| Command | Purpose |
| --- | --- |
| `node scripts/worktrees.mjs list [--json]` | Show active worktrees + reserved identifiers |
| `node scripts/worktrees.mjs claim --slug <slug> --type <tipo> [--branch <branch>] [--note <text>]` | Register a new worktree (dir must exist first) |
| `node scripts/worktrees.mjs reserve-migration --slug <slug> [--name <snake_case>]` | Reserve the next unused migration number |
| `node scripts/worktrees.mjs release --slug <slug>` | Free the worktree's identifiers after merge |

The helper locks the registry for every mutation, so concurrent agents claim unique slugs and reserve unique numbers even when they run at the same instant.

## Step 0 — Activate a worktree (claim)

1. `git fetch origin main`
2. `node scripts/worktrees.mjs list` — loop over the rows. If a slug for the same feature, or a name/reserved identifier you need, is already taken, STOP and pick another slug (add a numeric or noun suffix). Never guess silently.
3. Create the worktree **from `origin/main`**:
   `git worktree add ../healthcare-worktrees/{slug} -b {tipo}/{slug} origin/main` (run from the main repo directory; `{tipo}` is `feat/`, `fix/`, `chore/`, `refactor/`, `docs/`).
4. Claim it: `node scripts/worktrees.mjs claim --slug {slug} --type {tipo} --note "<one-line description>"`. If the claim fails with a collision, go back to step 2.
5. All work — code, tests, validations, commits, push — happens **inside the worktree**. The main repo directory stays on `main`. Never `git checkout -b` in the main repo directory.

If you run the claimed worktree's steps from a fresh agent session, the first command reads `list` again — the registry is how a later session learns what earlier ones reserved.

## Step 1 — Reserve unique identifiers BEFORE creating files

Anything that must be unique project-wide is claimed before any file is written:

| Identifier | Command / rule |
| --- | --- |
| Migration number | `node scripts/worktrees.mjs reserve-migration --slug <slug> --name <snake_case_name>` — never compute "local max + 1"; sibling agents compute the same number. |
| Module folder / Go package | Derive from the domain; confirm via `list` `--note` values that no other feature is building the same module. |
| Route base path / i18n key prefix | Confirm against `list`; if the same path is claimed, differ the slug, not the path. |

## Step 2 — Shared-file protocol (files every feature touches)

`backend/internal/app/policy/policy.go`, `backend/cmd/api/main.go`, `frontend/src/app/routes.tsx`, `shared/components/AppSidebar`, and i18n files are edited by every feature. Rules:

- **Additive only.** You may add entries and lines. You may NOT delete, rename, reformat, reindent or restructure existing lines — that is the conflict surface other agents depend on.
- **Anchored inserts.** Insert where git diffs cleanly:
  - `policy.go` map keys: insert at the **sorted key position** of the map (never append before the map's closing brace).
  - `main.go`: new import in **alphabetical position**; `Xxx.Register(...)`, HTTP handler and router argument near the related module block.
  - `routes.tsx`: `lazy()` import and `<Route>` in **path order**.
  - `AppSidebar`: menu item after the alphabetical/role-ordered neighbour.
  - i18n: add `{domain}` keys only; never rename or move existing keys.
- **Own only your leaf files.** Module subfolders (`internal/modules/{domain}/`, `frontend/src/modules/{domain}/`, `frontend/e2e/{domain}.spec.ts`) are private to the worktree — edit them freely.
- If the shared edit you need is not additive (a refactor of a shared file), STOP and warn the user: it will conflict with every in-flight feature.

## Step 3 — Rebase before opening the PR

Before `git push` and `gh pr create`:

1. `git fetch origin main`
2. If `origin/main` advanced since the worktree was created: `git rebase origin/main`
3. Resolve conflicts keeping both intents (see `resolving-merge-conflicts`). For additive maps/files a parallel feature's entries must survive alongside yours.
4. If the rebase introduced a duplicate migration number (another feature's same-number migration now exists), renumber your pair: delete the pair, run `node scripts/worktrees.mjs release --slug <slug>` + `reserve-migration` again, and recreate with the new number.
5. Re-run the validations (backend `go build` + `go vet` + `go test`, RBAC gate, frontend `npm run lint` + `npm run test` + `npm run build` + `npm run test:e2e`). Cheap — a clobbered shared-file entry shows up immediately.

## Step 4 — Collision handling

- **Another PR merged while you worked on the same files**: fetch, rebase again, resolve, revalidate (RBAC gate: `go test -run TestEveryCompiledGRPCRPCIsRegisteredInPolicy ./internal/app/policy/...`), then push.
- **Regained a reserved number after merge**: the released `migrations[]` in the registry no longer counts; reserve the next free number and rename your pair.
- **Never** `git push --force` to a shared branch.

## Step 5 — Release

After the PR is merged (or abandoned):

1. From the main repo directory: `git worktree remove ../healthcare-worktrees/{slug}`
2. `git branch -d {tipo}/{slug}`
3. `node scripts/worktrees.mjs release --slug {slug}` — frees the slug, branch and reserved numbers for the next feature.

## Rules of thumb

- Claim before you write; release after you merge.
- Check `list` before every identifier decision.
- Write registry mutations only through the helper (it locks).
- Never touch files outside your module subtree except through the additive-anchored protocol.
- A feature that cannot be scoped to additive shared-file edits must be flagged to the user before proceeding.