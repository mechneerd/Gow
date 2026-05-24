# Scaffolding Implementation Roadmap

**Project**: GoW Framework  
**Related Repo**: https://github.com/mechneerd/gow-skeleton.git  
**Date**: 2026-05-23 (Initial Creation)  
**Last Updated**: 2026-05-23 (Skeleton repo created & initial structure pushed)  
**Status**: Phase 1 In Progress — Skeleton Repository Created  
**Goal**: Make `gow new` deliver a Laravel-like scaffolding experience with published, user-owned files.

---

## 1. Current State

- `gow new myapp`, `gow new myapp --api`, and `gow new mysite --minimal` exist.
- Scaffolding logic is currently embedded inside the `gow` CLI binary.
- No interactive wizard.
- **Skeleton repository created**: https://github.com/mechneerd/gow-skeleton.git
- Initial structure with 4 starter kits (minimal, api, web, web-auth) has been pushed.
- Post-generation steps are very limited.

**Problem**: While the skeleton repo now exists, the `gow` CLI does not yet consume it. The main remaining work is implementing the CLI integration.

---

## Progress Log

**2026-05-23**
- Created parallel folder `D:\gow-skeleton`
- Connected to https://github.com/mechneerd/gow-skeleton.git
- Implemented initial structure with 4 starter kits
- Added placeholders, post-install scripts, README, and STRUCTURE.md
- Pushed initial commit (22 files) to GitHub
- Started CLI integration:
  - Created `cmd/gow/scaffold/` package (config, selector, cloner, copier)
  - Added `--minimal`, `--api`, `--auth` flags to `gow new`
  - Refactored `newCmd` to use new scaffold system
  - `gow new` now clones the real skeleton and copies the selected template
  - Implemented full placeholder replacement (`{{ .AppName }}`, `{{ .ModulePath }}`, etc.)
   - Automatic renaming of `go.mod.template` → `go.mod`
   - Added post-install: `go mod tidy` + `.env` creation + nice "Next Steps" output
   - Added `--module` flag to allow custom module path (e.g. github.com/username/myapp)
   - Added `--db` flag (sqlite, mysql, postgres) to choose database driver at creation time
   - Added input validation for `--db` flag
   - Added interactive wizard
   - Added `--force` flag to overwrite existing directories
   - Added `--no-git` flag + automatic `git init` (skipped when flag is used)
   - Polished output messages (cleaner, more professional UX)
   - Added `--skeleton` flag (experimental, custom repository support) — implemented but saved for future release
   - Added `--yes` flag for fully non-interactive mode (defaults to Web + Auth when no other starter kit is specified)

**Current status**: `gow new` is now very close to a production-ready Laravel-like experience. The legacy hardcoded scaffolding function has been fully removed.

---

## 2. Goal

Adopt Laravel’s key principle:

> “Scaffolding is published, not hidden.”

Every file created by `gow new` must land directly in the user’s project tree and be fully editable from day one.

Deliver a professional `gow new` experience that includes:
- Multiple starter kits
- Interactive prompts (optional)
- Automatic post-install steps
- Clear “Next Steps” guidance

---

## 3. Architecture Decision

**Primary Recommendation**: Use a separate published skeleton repository.

- **Skeleton Repository**: `https://github.com/mechneerd/gow-skeleton.git`
- The `gow` CLI will clone (or download) this skeleton and copy it into the target directory.
- This keeps the CLI lightweight and makes the scaffolding easy to version, fork, and improve.

Alternative (fallback): Embedded templates inside the `gow` binary using Go’s `embed` package (less flexible).

**Decision**: Proceed with the external skeleton repository as the main strategy.

---

## 4. Recommended Structure of `gow-skeleton` Repository

The skeleton repo should be organized as follows:

```
gow-skeleton/
├── .github/
├── templates/
│   ├── minimal/
│   │   ├── app/
│   │   ├── routes/
│   │   ├── resources/
│   │   ├── config/
│   │   ├── bootstrap/
│   │   ├── database/
│   │   ├── .env.example
│   │   ├── go.mod.template
│   │   └── main.go
│   ├── api/
│   ├── web/
│   └── web-auth/
├── scripts/
│   ├── post-install.sh
│   └── post-install.ps1
├── README.md
└── .gitignore
```

### Template Placeholders (to be replaced by `gow new`)

Use simple Go template syntax or custom markers:

- `{{ .AppName }}`
- `{{ .ModulePath }}`
- `{{ .DatabaseDriver }}`
- `{{ .IncludeAuth }}`
- `{{ .Year }}`

Example in `go.mod.template`:
```go
module {{ .ModulePath }}

go 1.22
```

---

## 5. Starter Kits to Implement

| Starter Kit   | Folder in Skeleton | Command Flags          | Description                                      | Priority |
|---------------|--------------------|------------------------|--------------------------------------------------|----------|
| Minimal       | `templates/minimal` | `--minimal`            | Basic routing + views                            | High     |
| API           | `templates/api`     | `--api`                | Sanctum + JSON routes + Resources                | High     |
| Web           | `templates/web`     | (default)              | Blade + layouts + components                     | High     |
| Web + Auth    | `templates/web-auth`| `--auth`               | Full session auth (login/register/middleware)    | High     |
| Full          | `templates/full`    | `--full`               | Web + API + Auth + Queue + Mail stubs            | Medium   |

**MVP Target (before public release)**: Deliver **Minimal**, **API**, and **Web + Auth**.

---

## 6. `gow new` Behavior (Target Experience)

### Command Examples
```bash
gow new myapp                    # Interactive or defaults to Web
gow new myapi --api
gow new myblog --auth
gow new mysite --minimal
```

### Interactive Mode (when no flags are passed)
The CLI should ask:
1. Starter kit? (Minimal / API / Web / Web + Auth)
2. Database driver? (SQLite / MySQL / PostgreSQL)
3. Include testing setup? (Y/n)
4. Run migrations after creation? (Y/n)

### Post-Installation Steps (automatic)
1. Replace all placeholders in files
2. Run `go mod tidy`
3. Copy `.env.example` → `.env` (generate basic secret if needed)
4. Initialize git (optional, with `--no-git` flag)
5. Print beautiful “Next Steps” message

---

## 7. Implementation Phases

### Phase 1: Foundation (Immediate – before v1.0)

- [x] Set up proper folder structure in `gow-skeleton` → **Completed** (2026-05-23)
- [ ] Implement basic `gow new` that clones the skeleton repo
- [ ] Add placeholder replacement logic
- [ ] Add post-install steps (`go mod tidy`, `.env` creation)
- [ ] Support `--minimal`, `--api`, `--auth` flags in `gow new`

### Phase 2: Polish & UX (v1.0)

- [ ] Add interactive wizard
- [ ] Improve “Next Steps” output
- [ ] Add `--force` and `--no-git` flags
- [x] Support custom skeleton URL (`--skeleton=...`) → Implemented (experimental, for future release)

### Phase 3: Advanced Features (1.2+)

- [ ] Multiple template variants inside one skeleton
- [ ] Support for third-party scaffolding (similar to Laravel’s `vendor:publish`)
- [ ] Pre-configured Docker / docker-compose options
- [ ] Ability to publish additional stubs from packages

---

## 8. Technical Recommendations

- Use shallow git clone (`--depth=1`) for speed when fetching the skeleton.
- Keep the `gow` binary small — logic for scaffolding should be minimal.
- Use Go’s `text/template` or a simple custom replacer for placeholders.
- Make the skeleton repo the single source of truth for default project structure.
- Add a `gow-skeleton` version file so the CLI can warn if the skeleton is outdated.

---

## 9. Documentation Updates Needed

- Update `docs/COMPLETE_USER_GUIDE.md`
- Update root `README.md`
- Create `docs/guide/scaffolding.md` (optional but recommended)
- Add examples of `gow new` with different starter kits

---

## 10. Open Questions / Decisions

- Should we support `--skeleton` pointing to any Git repo from day one?
- Do we want to embed a fallback skeleton inside the binary for offline use?
- How should we handle versioning of the skeleton repo vs the `gow` CLI?

---

## 11. Next Actions

1. **Team**: Update user-facing documentation (`COMPLETE_USER_GUIDE.md`, root README) to document the powerful new `gow new` command.
2. **Team**: Add support for local skeleton paths (`--skeleton=/path/to/folder`).
3. **Team**: Write proper tests for the `gow new` command and scaffold package.
4. Consider adding more official starter kits in the gow-skeleton repository.

---

**Document Status**: Living document. Update after each implementation phase.

**Last Updated**: 2026-05-24 (--yes flag implemented + legacy scaffold() function removed)
