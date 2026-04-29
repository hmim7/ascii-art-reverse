# ascii-art-reverse

![Go Version](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MIT](https://img.shields.io/badge/MIT-1BB581?style=for-the-badge&logo=opensourceinitiative&logoColor=white)
[![01](https://img.shields.io/badge/zone01-Athens-916ADE?&labelColor=181717&style=for-the-badge&logo=data:image/svg%2Bxml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IndoaXRlIiBzdHJva2Utd2lkdGg9IjIiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCI+PHBhdGggZD0iTTEyIDJMMiA3bDEwIDUgMTAtNS0xMC01eiIvPjxwYXRoIGQ9Ik0yIDE3bDEwIDUgMTAtNU0yIDEybDEwIDUgMTAtNSIvPjwvc3ZnPg==)](https://github.com/01-edu/public/blob/master/subjects/ascii-art/reverse/README.md)

## Description
A unified ASCII art utility built in Go that integrates six distinct features: core generation, multi-notation color support, terminal-aware alignment, custom banner loading, file output redirection, and a greedy-scan reverse engineering algorithm to reconstruct text from existing ASCII art files.

## Authors & Team Roles
- **hmim** - Lead / CLI Architect (CLI orchestration, flag parsing, and main.go logic).
- **kchatzian** - Rendering & Color (Rendering pipeline, ANSI injection, and terminal-aware alignment).
- **gtzimoka** - Banners & I/O (Banner loading, FS validation, and output redirection).
- **edamaski** - Reverse Engineering (Greedy-scan algorithm for reconstructing text).

## Installation
```bash
# Build the binary
go build -o ascii-art-reverse .
```

## Usage

### Part 1: CLI Usage
```bash
# Standard Usage (defaults to standard banner)
go run . "Hello World"
go run . "Hello World" standard

# Color with Substring Targeting
go run . --color=red "Hello" "Hello World"

# Terminal Alignment
go run . --align=center "Centered Art"

# File Output Redirection
go run . --output=output.txt "Persisted Art"

# Reverse Engineering
go run . --reverse=example.txt
go run . --reverse=example.txt shadow
```

### Error Handling
The CLI enforces strict flag validation and always exits with code 1 on bad input, printing the usage message that matches the first malformed flag in input order:

| Scenario | Message |
|---|---|
| Invalid/missing flag value | Correct `Usage:` for that flag |
| Unknown `--flag` | Hint to use `--` delimiter, then `UsageBasic` |
| Too many positionals | `UsageBasic` (with hint if a known flag option is misplaced) |
| Banner not found / invalid | Grey warning + `UsageBasic` |
| Multiple malformed flags | Single consolidated warning listing all, then first-match usage |

## Implementation Details

- **Architecture:** Modular pipeline using the Go standard library only. `main.go` is a pure orchestrator (≤60 lines); all decision logic lives in the `cli`, `banner`, `render`, `output`, and `reverse` packages.
- **CLI parsing:** Single-pass `ClassifyArgs` routes every token into exactly one bucket (`Positional`, `ColorRules`, `Malformed`, `UnknownFlags`). Validation helpers (`ValidateOrFatal`, `CheckPositionalsOrFatal`, `HandleBannerFallbackOrFatal`) keep `main.go` free of guard logic.
- **Color system:** `BuildColorRules` resolves multi-notation color strings (Named, Hex `#RRGGBB`, `rgb()`, `hsl()`) directly to `[]render.ColorRule` — no intermediate `[]interface{}` bridge.
- **Rendering:** `Mapper` returns 8 empty strings for runes outside `[32, 126]`. `RenderAlignedWithColorRules` handles the all-empty special case so `""` produces no output and `"\n"` produces exactly one blank line.
- **Banner loading:** Files are validated for exactly 855 lines. CRLF is normalized. Leading vs. trailing separator format is detected automatically. All 7 bundled banners are fully loaded and tested.
- **Reverse algorithm:** Greedy column-scan matching art signatures against a pre-built inverse banner map. CRLF normalization applied before splitting. Supports all printable ASCII characters and multi-line art.
- **Gopher easter egg:** `containsNonASCII` treats `\n`, `\t`, and `\` as neutral; any rune outside `[32, 126]` triggers the Peeking Gopher.
- **Testing:** 40 passing unit tests across 5 packages including table-driven tests, golden-file round-trip tests for all 8 example files (+ shadow and multi-line), and targeted edge-case coverage.

## Project Structure
```plaintext
ascii-art-reverse
.
├── .ai/                   # AI Usage Index and Logs
├── .docs/                 # Project Documentation
│   ├── .team/             # Team Workflow and Checklists
│   │   ├── checklists/
│   ├── PRD.md             # Product Requirements Document
│   ├── audit-cases.md     # Official Audit Scenarios
│   ├── error-cases.md     # Error Handling Specifications
│   ├── edge-cases.md      # Boundary Conditions
│   └── golden-tests.md    # Reference I/O Pairs
├── .tasks/                # Agile Task Cards
├── banner/                # Source Font Files (.txt)
│   ├── standard.txt
│   ├── shadow.txt
│   ├── doom.txt
│   ├── block.txt
│   ├── thinkertoy.txt
│   ├── dancing.txt
│   └── greek.txt
├── internal/
│   ├── banner/            # Banner loading & validation
│   ├── cli/               # Flag parsing, validation guards & usage strings
│   ├── output/            # I/O writer abstraction
│   ├── render/            # Rendering, color, alignment & glyph mapping
│   └── reverse/           # Greedy-scan reverse reconstruction
├── main.go                # Pure orchestrator (≤60 lines)
└── README.md
```

## Technical Documentation

- **Data Validation:** Strict `--flag=value` syntax enforced. Banner files validated for exactly 855 lines; path traversal prevented via name sanitization. Invalid banners always fatal — no silent fallback.
- **Warning consolidation:** All malformed flags in one session emit a single consolidated `warning: invalid <cats> flag(s) "<raw…>"` line, followed by individual warnings for duplicate flags, then any double-dash hints.
- **Visual / Output Integrity:** ANSI escape sequences precisely injected to avoid color bleeding. Alignment pads based on visible (non-ANSI) width. Terminal width detected via `COLUMNS` env var → ioctl → fallback 80.
- **Testing:**
  ```bash
  go test ./...
  go test ./... -v
  go vet ./...
  ```

## Project Documentation References:

- PRD
- Golden Tests
- Audit Cases
- Edge Cases

---
*This project is part of the Zone01 Campus curriculum. It is built and maintained according to the guidelines specified in the `.docs/` directory.*
