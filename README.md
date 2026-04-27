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
# Standard Usage
go run . "Hello World" standard

# Color with Substring Targeting
go run . --color=red "Hello" "Hello World"

# Terminal Alignment
go run . --align=center "Centered Art"

# File Output Redirection
go run . --output=output.txt "Persisted Art"

# Reverse Engineering
go run . --reverse=example.txt standard
```

## Implementation Details

- **Architecture:** Refactored Modular Pipeline using the Go standard library only. Logic is decoupled into `cli`, `banner`, `render`, `output`, and `reverse` packages.
- **Core Algorithm:** Uses a glyph-mapping system for generation and a greedy-scanning algorithm for the reverse feature. The reverse algorithm scans art columns to match signatures against a pre-loaded banner map.
- **Special Features:** Supports multi-notation color (Named, Hex, RGB, HSL), intelligent banner separator detection, and a "Peeking Gopher" easter egg for non-ASCII inputs.
- **Testing:** Comprehensive test suite including Table-Driven unit tests, Golden File integration tests, and race detection.
- **Performance Benchmarks:** Optimized O(1) glyph lookup after initial banner load and efficient column-scanning in the reverse algorithm.

## Error Handling

### CLI Fallbacks
- **Missing/Invalid Input:** Gracefully prints a specific usage message based on flag priority (Reverse > Output > Align > Color) and exits with status code 1.

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
├── internal/
│   ├── banner/            # Banner loading logic
│   ├── cli/               # Flag parsing & Usage strings
│   ├── output/            # I/O redirection abstraction
│   ├── render/            # Rendering & Color pipeline
│   └── reverse/           # Reverse reconstruction logic
├── main.go                # Orchestrator
└── README.md
```

## Project Documentation References:

- PRD
- Golden Tests
- Audit Cases
- Edge Cases

## Technical Documentation

- **Data Validation:** Strict flag validation enforcing `--flag=value` syntax. Banner files are validated for exactly 855 lines and path traversal is prevented via sanitization.
- **Visual / Output Integrity:** ANSI escape sequences are precisely injected to avoid "color bleeding." Alignment logic detects terminal width via the standard library.
- **Testing:** The project includes a robust testing suite. Run the tests using:
  ```bash
  go test ./... -v -race
  ```

---
*This project is part of the Zone01 Campus curriculum. It is built and maintained according to the guidelines specified in the `.docs/` directory.*
