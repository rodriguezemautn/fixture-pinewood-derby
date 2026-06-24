# Skill Registry — fixture

Last updated: 2026-06-24

## Project Context

- **Stack**: Go 1.26.1 + Svelte 5 (PWA) + SQLite
- **Architecture**: Hexagonal/Clean Architecture (Go backend)
- **Testing**: Go testing (built-in), strict TDD mode enabled
- **Mode**: hybrid (openspec + engram)

## Skills

### SDD Pipeline (system — always available)
| Skill | Phase | Purpose |
|-------|-------|---------|
| sdd-explore | Explore | Investigate codebase, think through problems |
| sdd-propose | Propose | Create change proposals with scope and approach |
| sdd-spec | Spec | Write GWT specifications from proposals |
| sdd-design | Design | Technical design with architecture decisions |
| sdd-tasks | Tasks | Break down specs into implementation checklist |
| sdd-apply | Apply | Implement task definitions |
| sdd-verify | Verify | Validate implementation against specs |
| sdd-archive | Archive | Close out completed changes |

### Domain Skills (project-level)
| Skill | Location | Purpose |
|-------|----------|---------|
| **ui-ux-pro-max** | `.opencode/skills/ui-ux-pro-max/` | UI/UX design intelligence: 67 styles, 161 color palettes, 57 font pairings, design system generator |
| go-testing | system | Go testing patterns for Gentleman.Dots |
| skill-creator | system | New AI skill creation |

### User Skills (system-wide)
| Skill | Trigger | Purpose |
|-------|---------|---------|
| go-testing | Go tests, Bubbletea TUI testing | Go testing patterns |
| skill-creator | Creating AI skills | New skill creation |
| customize-opencode | Editing opencode config | AI configuration |

## Conventions (project-level)
See `openspec/config.yaml` for project-specific SDD rules.

## Conventions (system-level)
See `~/.config/opencode/AGENTS.md` for system-wide agent rules.
