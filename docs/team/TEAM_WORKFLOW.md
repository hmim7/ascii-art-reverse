# Team Workflow & Implementation Playbook: ascii-art-reverse

Welcome to the **ascii-art-reverse** team! This document is our single source of truth for how we collaborate, build, and deliver this project successfully. By adhering to these guidelines, we ensure high code quality, linear history, and predictable delivery.

## Project Overview & Role Assignment
The objective is to build a unified ASCII art utility that supports file persistence, color rendering, custom banner files, text alignment, and reverse engineering reconstruction.

### Team Delivery
- **hmim**: **Member 1 (Lead/CLI Architect)** - Responsible for CLI orchestration (`internal/cli`), `main.go` thin logic, and overall system integration.
- **kchatzian**: **Member 2 (Rendering & Color)** - Focuses on the rendering pipeline, ANSI color injection, and terminal-aware alignment.
- **gtzimoka**: **Member 3 (Banners & I/O)** - Manages the banner loading system, file system validation, and output redirection logic.
- **edamaski**: **Member 4 (Reverse Engineering)** - Dedicated to the greedy-scan algorithm for reconstructing text from existing ASCII art files.

### Suggestions for Role Assignment
- **Architecture & CLI:** Lead (hmim) ensures the `Config` struct and flag parsing handle all audit priority rules.
- **Core Rendering:** (kchatzian) handles the 8-line glyph assembly and ANSI color logic.
- **Infrastructure:** (gtzimoka) ensures banner files are validated and the file-writer handles `O_TRUNC`.
- **Algorithmic:** (edamaski) focuses on the complexity of the greedy-scan reverse logic.

## Project Architecture
We follow a strict KISS and YAGNI approach:
- `/internal`: Core business logic (cli, banner, render, output, reverse).
- `/main.go`: Minimal orchestrator (kept under 60 lines).
- `/banner`: Source .txt font files.

## Technical Specifications

### Core Stack & Constraints
- **Environment:** Go Standard Library ONLY.
- **Dependencies:** No external packages (e.g., no `cobra`, `pflag`, or `color` libraries).

### Code Quality & Standards
- **Formatting:** Code MUST be `gofmt` compliant.
- **Error Handling:** Use centralized usage strings in `internal/cli/errs.go`. Fail fast and exit with status 1.
- **Security:** Validate banner paths to prevent traversal; sanitize all user inputs.

### Testing Requirements
- **Coverage:** All internal logic must be verified against `golden-tests.md`.
- **Pre-commit Rule:** No code is pushed to `dev` if `go test ./...` fails.

## Workflow Implementation

### 1. Understanding the PRD, Edge Cases & Milestones
Before writing any implementation logic, thoroughly review `docs/PRD.md`. You must also consult `docs/edge-cases.md`, `docs/error-cases_universal.md`, `docs/audit-cases.md`, and `docs/golden-tests_universal.md` to ensure your code accounts for required technical constraints and official audit scenarios.

### 2. Working with Tasks
1. Claim an unassigned Task Card in the `.tasks/` directory.
2. Update the `STATUS` to `IN PROGRESS`.
3. Update the card to `DONE` only when criteria are met.
4. **Status Update Rule:** Update checklist items incrementally. Commit task status updates within the same logical commit as the related code change.

### 3. Effective Go & Testing
Code quality is a collective responsibility. To maintain our standard of excellence, every teammate must adhere to the following testing protocols and idiomatic Go practices.
Before marking a task as **DONE** or push to the remote `dev` branch , you **must** verify your implementation against the official good practices checklists in `.docs/.team/checklists/`:  
- [CLI Good Practices Checklist](./checklists/CLI-Good-Practices-Checklist.md)
- [Testing Good Practices Checklist](./checklists/Testing-Good-Practices-Checklist.md)
- [Web Good Practices Checklist](./checklists/Web-Good-Practices-Checklist.md)
- [Rest-API Good Practices Checklist](./checklists/Rest-API-Good-Practices-Checklist.md)

*Recommended as well:*
- [Conventional Commits](./checklists/Conventional-Commits.md)
- [Git Workflow](./checklists/Git-Workflow.md)

We follow the **Table-Driven Tests** pattern as the absolute standard for our Go logic.
- **Race Detector:** Always run tests with `go test -race ./...`.
- **Benchmarks:** Use `testing.B` for performance-critical sections of the reverse algorithm.

### 4. Effective Git Workflow
- **Branching:** Use a `dev` branch for features; merge to `main` only for final release.
- **Conventional Commits:** 
  - `feat(reverse): implement greedy-scan match`
  - `fix(cli): correct alignment usage priority`
- **Push Discipline:** Push regularly. Pull with `git pull --rebase` to maintain a linear history.

### 5. Code Review Process (Pre-Merge)
Peer review is mandatory. Before pushing to `dev`, have at least one teammate verify:
- Standard Library compliance.
- Adherence to `audit_cases.md`.
- No `panic()` calls.

### 6. Final Team Test (Pre-Release)
Before merging the `dev` branch into `main` to finalize the project, the entire team must conduct a joint final test:
1. **Feature Freeze:** Halt all new development on the `dev` branch.
2. **Golden Test Run:** Run the entire `golden-tests_universal.md` suite.
3. **Audit Simulation:** Attempt to break the application using `audit-cases.md` and `edge-cases.md`.

### 7. Personal AI Collaboration Protocol
- Each teammate maintains their own log in `.ai/`.
- Log entries must reference Task IDs.
- Never commit AI code without local verification.

### 8. Team Reminder
- Keep the repository readable.
- Keep commits atomic.
- Keep AI usage transparent.
- Keep task files updated.

## Appendix: Team Log Files
- `.ai/hmim.ai.txt` (Shared Architect Log)
- `.ai/kchatzian.ai.txt` (Rendering & Color)
- `.ai/gtzimoka.ai.txt` (Banners & I/O)
- `.ai/edamaski.ai.txt` (Reverse Engineering)

---
*This project is part of the Zone01 Campus curriculum.*

