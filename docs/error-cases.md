# Error Cases: ascii-art-reverse

Validation rules and expected behavior when the program receives invalid input,
malformed flags, missing files, or out-of-contract arguments. All error cases
print to **stderr** and exit with **status 1** unless noted otherwise.

---

## Summary Table

| ID | Category | Trigger | Expected Behavior |
|:--:|----------|---------|-------------------|
| E01 | CLI Core | No arguments | Usage (basic) + exit 1 |
| E02 | CLI Core | Too many positional args | Usage (basic) + exit 1 |
| E03 | CLI Core | Empty string `""` | No output; exit 0 |
| E04 | CLI Core | Unknown flag | Usage (basic) + exit 1 |
| E05 | CLI Core | Banner name as only arg | Usage (basic) + exit 1 |
| E06 | FS | Non-existent banner | Warning to stderr; fallback to standard |
| E07 | FS | Corrupted banner (≠855 lines) | Warning to stderr; fallback to standard |
| E08 | FS | Path traversal in banner name | Sanitized; fallback to standard |
| E09 | FS | Non-ASCII input only | Peeking Gopher rendered; exit 0 |
| E10 | FS | Mixed ASCII + non-ASCII | Non-ASCII silently dropped; render valid chars |
| E11 | Color | `--color` without `=` | Usage (color) + exit 1 |
| E12 | Color | Missing string after color | Usage (color) + exit 1 |
| E13 | Color | Invalid color value | Warning to stderr; render without color |
| E14 | Color | All color flags invalid | Warnings; render without color |
| E15 | Color | Invalid + valid colors mixed | Warning for invalid; valid colors applied |
| E16 | Output | `--output` without `=` | Usage (output) + exit 1 |
| E17 | Output | `--output` value not `.txt` | Usage (output) + exit 1 |
| E18 | Output | `--output` value empty | Usage (output) + exit 1 |
| E19 | Output | Duplicate `--output` flags | Warning; last value wins |
| E20 | Output | Invalid + valid output flags | Warning for invalid; valid output used |
| E21 | Alignment | `--align` without `=` | Usage (align) + exit 1 |
| E22 | Alignment | `--align` unknown type | Warning; fallback to left |
| E23 | Alignment | Duplicate `--align` flags | Warning; last value wins |
| E24 | Alignment | Invalid + valid align flags | Warning for invalid; valid align used |
| E25 | Reverse | `--reverse` without `=` | Usage (reverse) + exit 1 |
| E26 | Reverse | File not found | Error to stderr; exit 1 |
| E27 | Reverse | File is not valid ASCII art | Error to stderr; exit 1 |

---

## Category 1: CLI Core

### E01 — No Arguments
**Input:** `go run .`
**Expected (stderr):**
```
Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```
**Rules:** Exit status 1. No art produced.

---

### E02 — Too Many Positional Arguments
**Input:** `go run . "hello" standard "there"`
**Expected (stderr):**
```
Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```
**Rules:** After flag extraction, more than 2 positional tokens is an error; exit 1.

---

### E03 — Empty String Input
**Input:** `go run . ""`
**Expected:** No output; exit 0 (success).
**Rules:** Empty string is not an error — the program exits cleanly with no art.

---

### E04 — Unknown Flag
**Input:** `go run . --foo=bar "hello"`
**Expected (stderr):** Usage (basic); exit 1.
**Rules:** Any `--` token that does not match a known prefix (`--output=`, `--align=`, `--color=`, `--reverse=`, `--stdin`) is an unknown flag and causes immediate rejection.

---

### E05 — Banner Name as Only Argument
**Input:** `go run . standard`
**Expected (stderr):**
```
Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```
**Rules:** A single token that matches a known banner name with no preceding string is an error; exit 1.

---

## Category 2: File System / Banner

### E06 — Non-Existent Banner
**Input:** `go run . "hello" ghost`
**Expected (stderr):**
```
warning: banner "ghost" not found, default banner "standard" applied
```
**Expected (stdout):** ASCII art for "hello" using standard banner.
**Rules:** Non-fatal; program continues with standard fallback.

---

### E07 — Corrupted Banner File (≠ 855 Lines)
**Input:** `go run . "hello" broken` (where broken.txt has wrong line count)
**Expected (stderr):**
```
warning: banner "broken" invalid (expected 855 lines), default banner "standard" applied
```
**Expected (stdout):** ASCII art for "hello" using standard banner.

---

### E08 — Path Traversal in Banner Name
**Input:** `go run . "hello" ../../../etc/passwd`
**Expected (stderr):** Warning about invalid/not-found banner; fallback to standard.
**Rules:** `filepath.Base()` strips path components; only `[a-z0-9\-_]` characters allowed in banner name after sanitization. Path traversal attempts resolve to an invalid name and fall back to standard.

---

### E09 — Non-ASCII Input Only (Peeking Gopher)
**Input:** `go run . "こんにちは"`
**Expected (stdout):** Peeking Gopher ANSI art; exit 0.
**Rules:** Trigger condition: `text != ""` AND `containsNonASCII(text)` AND every segment is blank after `filterASCII`. Not an error — exit 0.

---

### E10 — Mixed ASCII and Non-ASCII
**Input:** `go run . "helloこんにちは"`
**Expected (stdout):** ASCII art for "hello" only; non-ASCII runes silently dropped.
**Rules:** `filterASCII` removes runes outside [32, 126]; gopher does NOT trigger because at least one segment is non-blank.

---

## Category 3: Color

### E11 — `--color` Without `=`
**Input:** `go run . --color red "banana"`
**Expected (stderr):**
```
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```
**Rules:** Exit 1. `--color red` is a malformed flag (space instead of `=`).

---

### E12 — Missing String After Color Flag
**Input:** `go run . --color=red`
**Expected (stderr):** Usage (color); exit 1.
**Rules:** No positional text argument provided; program cannot render.

---

### E13 — Invalid Color Value
**Input:** `go run . --color=notacolor "hello"`
**Expected (stderr):**
```
warning: invalid color "notacolor", rendering without color
```
**Expected (stdout):** ASCII art for "hello" in default color.
**Rules:** Non-fatal; ANSI code is empty string; `WrapWithColor` short-circuits.

---

### E14 — All Color Flags Invalid
**Input:** `go run . --color=badcolor "he" --color=#12 "ll" "hello"`
**Expected (stderr):** Two warning lines (one per invalid color).
**Expected (stdout):** ASCII art for "hello" with no color applied.

---

### E15 — Mixed Valid and Invalid Colors
**Input:** `go run . --color=red "he" --color=notAColor "ll" --color=blue "o" "hello"`
**Expected (stderr):**
```
warning: invalid color "notAColor", rendering without color
```
**Expected (stdout):** "he" in red, "o" in blue, "ll" in default color.

---

## Category 4: Output

### E16 — `--output` Without `=`
**Input:** `go run . --output test00.txt banana standard`
**Expected (stderr):**
```
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --output=<fileName.txt> something standard
```
**Rules:** Exit 1. Token `--output test00.txt` is malformed.

---

### E17 — Output Value Not `.txt`
**Input:** `go run . --output=file.csv "hello"`
**Expected (stderr):** Usage (output); exit 1.
**Rules:** `--output` value must end with `.txt`; any other extension is rejected.

---

### E18 — Output Value Empty
**Input:** `go run . --output= "hello"` (empty value after `=`)
**Expected (stderr):** Usage (output); exit 1.

---

### E19 — Duplicate `--output` Flags
**Input:** `go run . --output=first.txt --output=second.txt "hello" standard`
**Expected (stderr):**
```
warning: output redirected to "second.txt"; previous flag "--output=first.txt" ignored
```
**Expected:** `second.txt` created; `first.txt` not created (or unchanged).

---

### E20 — Invalid Output Alongside Valid Flags
**Input:** `go run . --output test.pdf --color=blue "hello" standard`
**Expected (stderr):**
```
warning: invalid output flag "--output test.pdf" ignored
```
**Expected (stdout):** "hello" in blue to terminal (no file written).

---

## Category 5: Alignment

### E21 — `--align` Without `=`
**Input:** `go run . --align right "something" standard`
**Expected (stderr):**
```
Usage: go run . [OPTION] [STRING] [BANNER]

Example: go run . --align=right something standard
```
**Rules:** Exit 1.

---

### E22 — `--align` Unknown Type
**Input:** `go run . --output=test.txt --align=middle "hello" standard`
**Expected (stderr):**
```
warning: invalid align flag "--align=middle" ignored
```
**Expected:** `test.txt` created; art written with default left alignment.

---

### E23 — Duplicate `--align` Flags
**Input:** `go run . --align=center --align=right "hello" standard`
**Expected (stderr):**
```
warning: previous align flag "--align=center" overridden by "--align=right"
```
**Expected (stdout):** "hello" right-aligned.

---

### E24 — Invalid Align Alongside Valid Flags
**Input:** `go run . --color=red --align=diagonal "hello" standard`
**Expected (stderr):**
```
warning: invalid align flag "--align=diagonal" ignored
```
**Expected (stdout):** "hello" in red with default left alignment.

---

## Category 6: Reverse

### E25 — `--reverse` Without `=`
**Input:** `go run . --reverse example00.txt`
**Expected (stderr):**
```
Usage: go run . [OPTION]

EX: go run . --reverse=<fileName>
```
**Rules:** Exit 1. Space instead of `=` is malformed.

---

### E26 — Reverse File Not Found
**Input:** `go run . --reverse=nonexistent.txt`
**Expected (stderr):** File I/O error message; exit 1.
**Rules:** Hard error (no fallback); the reverse pipeline cannot proceed without the input file.

---

### E27 — Reverse File Not Valid ASCII Art
**Input:** `go run . --reverse=notart.txt` (file contains random bytes)
**Expected (stderr):** Error or empty output; exit 1.
**Rules:** The reverse pipeline must not panic; if no glyph matches are found, report an error.
