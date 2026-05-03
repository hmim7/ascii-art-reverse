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

### CLI Usage
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
go run . --reverse=example00.txt
go run . --reverse=colored_hello.txt standard

# Reverse preserves alignment — right-aligned art returns text with leading spaces
go run . --output=out.txt --align=right "RGB" doom
go run . --reverse=out.txt doom   # returns "              RGB"
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
- **Alignment system:** Terminal width is detected via `ioctl` (Unix) or `COLUMNS` env var, falling back to 80 for file output. Implements specific padding algorithms for `center`, `right`, and gap-distributed `justify` modes.
- **Output system:** `internal/output` provides a stream-agnostic `io.Writer` abstraction, resolving `os.Stdout` or a named file with `O_TRUNC` and `0644` permissions.
- **Rendering:** `Mapper` returns 8 empty strings for runes outside `[32, 126]`. `RenderAlignedWithColorRules` handles the all-empty special case so `""` produces no output and `"\n"` produces exactly one blank line.
- **Banner loading:** Files are validated for exactly 855 lines. CRLF is normalized. Leading vs. trailing separator format is detected automatically. All 7 bundled banners are fully loaded and tested.
- **Reverse algorithm:** Five-phase pipeline:
  1. **Preprocessing** — CRLF normalized to `\n`; trailing empty lines stripped; `render.StripANSI` applied per line so color codes never affect glyph matching.
  2. **Reverse map** — `buildReverseMap` inverts `bannerMap` into a `signature → rune` lookup. Each glyph is padded to its maximum width before hashing. Runes are iterated 126 → 32 so that on a visual collision the higher ASCII value (lowercase) wins.
  3. **Greedy scanner** (`matchSegmentFrom`) — glyph widths deduplicated and sorted descending; scanner tries widest candidate first at each column, ensuring wider glyphs are never shadowed by a narrower sub-match.
  4. **Alignment preservation** (`matchSegment`) — `col=0` is tried first; when alignment padding is divisible by the banner's space glyph width (`spaceW`), leading columns decode as space glyphs and the full offset is preserved. When not divisible, `floor(indent / spaceW)` spaces are prepended and matching resumes from `indent` (best-effort). An all-blank block always routes to `col=0` to decode a single space glyph.
  5. **Multi-line support** — a blank line in the art file passes through as `""`, reconstructing the original `\n` separator in the output.

  Consistent alignment-preserving output verified across all six bundled banner fonts: standard, shadow, thinkertoy, doom, dancing, greek.
- **Gopher easter egg:** `containsNonASCII` treats `\n`, `\t`, and `\` as neutral; any rune outside `[32, 126]` triggers the Peeking Gopher.
- **Testing:** 100% table-driven testing suite across all packages. Includes golden-file round-trip integration tests, comprehensive fuzzing for every entry point, and performance benchmarks for core algorithms.

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
│   │   ├── loader.go
│   │   └── loader_test.go
│   ├── cli/               # Flag parsing, validation guards & usage strings
│   │   ├── errs.go
│   │   ├── flags.go
│   │   └── cli_test.go
│   ├── output/            # I/O writer abstraction
│   │   ├── output.go
│   │   └── output_test.go
│   ├── render/            # Rendering, color, alignment & glyph mapping
│   │   ├── alignment.go
│   │   ├── alignment_unix.go
│   │   ├── alignment_windows.go
│   │   ├── colorizer.go
│   │   ├── gopher.go
│   │   ├── mapper.go
│   │   ├── parser.go
│   │   ├── renderer.go
│   │   └── render_test.go
│   └── reverse/           # Greedy-scan reverse reconstruction
│       ├── reverser.go
│       └── reverser_test.go
├── main.go                # Pure orchestrator (≤60 lines)
└── README.md
```

## Technical Documentation

- **Data Validation:** Strict `--flag=value` syntax enforced. Banner files validated for exactly 855 lines; path traversal prevented via name sanitization. Invalid banners always fatal — no silent fallback.
- **Warning consolidation:** All malformed flags in one session emit a single consolidated `warning: invalid <cats> flag(s) "<raw…>"` line, followed by individual warnings for duplicate flags, then any double-dash hints.
- **Memory & Quality:** Structs are optimized for memory alignment to minimize padding and pointer scanning. The codebase is strictly verified against `shadow`, `fieldalignment`, and `errcheck` analyzers.
- **Verification Suite:**
  ```bash
  go test -race -v ./...        # Run 100% table-driven suite with race detector
  go test -fuzz=Fuzz ./...      # Execute fuzzing targets
  go test -bench=. ./...        # Run performance benchmarks
  go test -cover ./...          # Check code coverage
  go tool cover -func=cover.out # Show total coverage percentage
  go test ./... -coverprofile=cover.out
  ```

## Project Documentation References:

- [PRD](.docs/PRD.md)
- [Error Cases](.docs/error-cases.md)
- [Audit Cases](.docs/audit-cases.md)
- [Edge Cases](.docs/edge-cases.md)
- [Golden Tests](.docs/golden-tests.md)

---
*This project is part of the Zone01 Campus curriculum. It is built and maintained according to the guidelines specified in the `.docs/` directory.*
