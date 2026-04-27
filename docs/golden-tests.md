# Golden Test Suite: ascii-art-reverse

Mandatory expected-I/O pairs for all feature categories. Each case locks the
exact command and its required output. Use `cat -e` for line-ending validation
(`$` marks each line end). Cases are grouped by feature category.

---

## Category 1: CLI Core

### Summary

| ID | Input Command | Expected Output |
|----|--------------|-----------------|
| G01 | `go run . ""` | *(nothing)* |
| G02 | `go run . "\n"` | one blank line |
| G03 | `go run . "Hello\n"` | 8-line "Hello" + blank line |
| G04 | `go run . "hello"` | 8-line standard art |
| G05 | `go run . "HeLlO"` | mixed glyph mapping |
| G06 | `go run . "Hello There"` | horizontal concatenation |
| G07 | `go run . "1Hello 2There"` | numeric + alpha glyphs |
| G08 | `go run . "{Hello There}"` | brace glyphs |
| G09 | `go run . "Hello\nThere"` | two 8-line blocks |
| G10 | `go run . "Hello\n\nThere"` | two blocks + blank line between |

### G04 — `go run . "hello" | cat -e`
```
 _              _   _          $
| |            | | | |         $
| |__     ___  | | | |   ___   $
|  _ \   / _ \ | | | |  / _ \  $
| | | | |  __/ | | | | | (_) | $
|_| |_|  \___| |_| |_|  \___/  $
                               $
                               $
```

### G09 — `go run . "Hello\nThere" | cat -e`
```
 _    _          _   _          $
| |  | |        | | | |         $
| |__| |   ___  | | | |   ___   $
|  __  |  / _ \ | | | |  / _ \  $
| |  | | |  __/ | | | | | (_) | $
|_|  |_|  \___| |_| |_|  \___/  $
                                $
                                $
 _______   _                           $
|__   __| | |                          $
   | |    | |__     ___   _ __    ___  $
   | |    |  _ \   / _ \ | '__|  / _ \ $
   | |    | | | | |  __/ | |    |  __/ $
   |_|    |_| |_|  \___| |_|     \___| $
                                       $
                                       $
```

### G10 — `go run . "Hello\n\nThere" | cat -e`
```
 _    _          _   _          $
| |  | |        | | | |         $
| |__| |   ___  | | | |   ___   $
|  __  |  / _ \ | | | |  / _ \  $
| |  | | |  __/ | | | | | (_) | $
|_|  |_|  \___| |_| |_|  \___/  $
                                $
                                $
$
 _______   _                           $
|__   __| | |                          $
   | |    | |__     ___   _ __    ___  $
   | |    |  _ \   / _ \ | '__|  / _ \ $
   | |    | | | | |  __/ | |    |  __/ $
   |_|    |_| |_|  \___| |_|     \___| $
                                       $
                                       $
```

---

## Category 2: Color

### Summary

| ID | Command | Expected |
|----|---------|----------|
| G11 | `--color red "banana"` | Usage (color); exit 1 |
| G12 | `--color=red "hello world"` | All glyphs in red ANSI |
| G13 | `--color=green "1 + 1 = 2"` | All glyphs in green ANSI |
| G14 | `--color=yellow "(%&) ??"` | All glyphs in yellow ANSI |
| G15 | `--color=red "ello" "hello"` | Only "ello" in red |
| G16 | `--color=red "e" "hello"` | Only "e" glyph in red |
| G17 | `--color=red "ll" "hello"` | Only "ll" glyphs in red |
| G18 | `--color=orange GuYs "HeY GuYs"` | Only "GuYs" in orange |
| G19 | `--color=blue B "RGB()"` | Only "B" in blue |
| G20 | `--color=red "he" --color=blue "lo" "hello"` | "he"=red, "lo"=blue |
| G21 | `--color=red "ll" --color=green "ll" "hello"` | "ll" in green (last wins) |
| G22 | Multi-color + invalid color | Warning; valid colors applied |
| G23 | All invalid colors | All warnings; no color |

### G11 — Usage Message (color)
**Input:** `go run . --color red "banana"`
**Expected (stderr):**
```
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```

### G18 — Case-Sensitive Substring Color
**Input:** `go run . --color=orange GuYs "HeY GuYs"`
**Expected:** ASCII art for "HeY GuYs"; only the glyphs for "GuYs" (columns 4–7) wrapped in `\x1b[38;5;208m`…`\x1b[0m`; "HeY " remains default.

### G20 — Multi-Color Two Targets
**Input:** `go run . --color=red "he" --color=blue "lo" "hello"`
**Expected:** 
- Glyph columns for "h","e" → wrapped in `\x1b[31m`…`\x1b[0m`
- Glyph columns for "l" (3rd char) → default (no rule hits "l" alone)
- Glyph columns for "l","o" → wrapped in `\x1b[34m`…`\x1b[0m`

### G22 — Multi-Color with Invalid Color
**Input:** `go run . --color=red "he" --color=notAColor "ll" --color=blue "o" "hello"`
**Expected (stderr):**
```
warning: invalid color "notAColor", rendering without color
```
**Expected (stdout):** "he" in red, "o" in blue, "ll" in default color.

### G23 — All Colors Invalid
**Input:** `go run . --color=badcolor "he" --color=#12 "ll" "hello"`
**Expected (stderr):**
```
warning: invalid color "badcolor", rendering without color
warning: invalid color "#12", rendering without color
```
**Expected (stdout):** "hello" art with no ANSI codes.

### Color Notation Equivalence
All four commands below must produce visually identical output (red "hello"):
```bash
go run . --color=red "hello"
go run . --color=#ff0000 "hello"
go run . --color=rgb(255,0,0) "hello"
go run . --color=hsl(0,100%,50%) "hello"
```

---

## Category 3: File System / Banner Selection

### Summary

| ID | Command | Expected |
|----|---------|----------|
| G24 | `"banana" standard abc` | Usage (basic); exit 1 |
| G25 | `"hello" standard` | Standard 8-line "hello" |
| G26 | `"hello world" shadow` | Shadow-style "hello world" |
| G27 | `"nice 2 meet you" thinkertoy` | Thinkertoy number+word render |
| G28 | `"you & me" standard` | Ampersand glyph |
| G29 | `"123" shadow` | Shadow numeric glyphs |
| G30 | `"/(\")" thinkertoy` | Slash/paren/quote glyphs |
| G31 | `"ABCDEFGHIJKLMNOPQRSTUVWXYZ" shadow` | All 26 uppercase shadow glyphs |
| G32 | `"It's Working" thinkertoy` | Apostrophe + mixed case |
| G33 | `"hello" ghost` | Warning + standard fallback |

### G24 — Too Many Arguments
**Input:** `go run . "banana" standard abc`
**Expected (stderr):**
```
Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```

### G25 — `go run . "hello" standard | cat -e`
```
 _              _   _          $
| |            | | | |         $
| |__     ___  | | | |   ___   $
|  _ \   / _ \ | | | |  / _ \  $
| | | | |  __/ | | | | | (_) | $
|_| |_|  \___| |_| |_|  \___/  $
                               $
                               $
```

### G26 — `go run . "hello world" shadow | cat -e`
```
                                                                                        $
_|                _| _|                                                     _|       _| $
_|_|_|     _|_|   _| _| _|    _|       _|      _|      _| _|    _| _|_|     _| _|_|_| $
_|    _| _|_|_|_| _| _|_|_|_| _|       _|      _|      _| _|    _| _|    _| _| _|_|_| $
_|    _| _|       _| _| _|    _|         _|  _|  _|  _|   _|    _| _|_|_|   _| _|_|_| $
_|    _|   _|_|_| _| _|   _|_|             _|      _|       _|_|   _|       _|   _|_| $
                                                                                       $
                                                                                       $
```

### G33 — Non-Existent Banner Fallback
**Input:** `go run . "hello" ghost`
**Expected (stderr):**
```
warning: banner "ghost" not found, default banner "standard" applied
```
**Expected (stdout):** Same as G25 (standard "hello").

---

## Category 4: Output (`--output` flag)

### Summary

| ID | Command | Expected |
|----|---------|----------|
| G34 | `--output test00.txt banana standard` | Usage (output); exit 1 |
| G35 | `--output=test00.txt "First\nTest" shadow` | File: two 8-line shadow blocks |
| G36 | `--output=test01.txt "hello" standard` | File: 8-line standard "hello" |
| G37 | `--output=test02.txt "123 -> #$%" standard` | File: numbers+symbols |
| G38 | `--output=test03.txt "432 -> #$%&@" shadow` | File: shadow symbols |
| G39 | `--output=test04.txt "There" shadow` | File: shadow "There" |
| G40 | `--output=test05.txt "123 -> \"#$%@" thinkertoy` | File: thinkertoy quotes |
| G41 | `--output=test06.txt "2 you" thinkertoy` | File: thinkertoy two-word |
| G42 | `--output=test07.txt "Testing long output!" standard` | File: no wrapping |

### G34 — Usage Message (output)
**Input:** `go run . --output test00.txt banana standard`
**Expected (stderr):**
```
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --output=<fileName.txt> something standard
```

### G36 — `go run . --output=test01.txt "hello" standard` → `cat -e test01.txt`
```
 _              _   _          $
| |            | | | |         $
| |__     ___  | | | |   ___   $
|  _ \   / _ \ | | | |  / _ \  $
| | | | |  __/ | | | | | (_) | $
|_| |_|  \___| |_| |_|  \___/  $
                               $
                               $
```

### Output Validation Command Bundle
```bash
go run . --output=test00.txt "First\nTest" shadow  && cat -e test00.txt
go run . --output=test01.txt "hello" standard      && cat -e test01.txt
go run . --output=test02.txt "123 -> #$%" standard && cat -e test02.txt
go run . --output=test07.txt "Testing long output!" standard && cat -e test07.txt
```

---

## Category 5: Alignment / Justify (`--align` flag)

### Summary

| ID | Command | Expected |
|----|---------|----------|
| G43 | `--align right "something" standard` | Usage (align); exit 1 |
| G44 | `--align=right "left" standard` | Art flush right |
| G45 | `--align=left "right" standard` | Art flush left (no padding) |
| G46 | `--align=center "hello" shadow` | Art centered |
| G47 | `--align=justify "1 Two 4" shadow` | Words span full width |
| G48 | `--align=right "23/32" standard` | Symbols flush right |
| G49 | `--align=center "#$%&\"" thinkertoy` | Symbols centered |
| G50 | `--align=right "ABCabc123" thinkertoy` | Mixed chars flush right |
| G51 | `--align=left "23Hello World!" standard` | Mixed chars flush left |
| G52 | `--align=justify 'HELLO there HOW are YOU?!' thinkertoy` | All words span full width |
| G53 | `--align=right "a -> A b -> B c -> C" shadow` | Complex string flush right |

### G43 — Usage Message (align)
**Input:** `go run . --align right "something" standard`
**Expected (stderr):**
```
Usage: go run . [OPTION] [STRING] [BANNER]

Example: go run . --align=right something standard
```

### G47 — Justify Alignment
**Input:** `go run . --align=justify "1 Two 4" shadow`
**Expected:** "1" starts at column 0; "4" ends at `terminalWidth`; "Two" positioned so all gaps are equal (or differ by at most 1 space).

### Alignment + Combined Flags
```bash
# Color + center
go run . --color=blue --align=center "testing" thinkertoy

# Output + right (uses width 80)
go run . --output=aligned.txt --align=right "world" shadow

# All three flags
go run . --output=full.txt --color=yellow --align=center "full test" shadow

# Multi-color + justify + output
go run . --output=multi.txt --color=red "H" --color=blue "W" --align=justify "Hello World" standard
```

### Warning Cases (Combined)

**Duplicate --align:**
```bash
go run . --align=center --align=right "hello" standard
# stderr: warning: previous align flag "--align=center" overridden by "--align=right"
# stdout: "hello" right-aligned
```

**Duplicate --output:**
```bash
go run . --output=first.txt --output=second.txt "hello" standard
# stderr: warning: output redirected to "second.txt"; previous flag "--output=first.txt" ignored
# result: second.txt created; first.txt not created
```

**Flag after positional args:**
```bash
go run . "hello" standard --align=right
# stdout: "hello" right-aligned (flags are position-independent)
```

---

## Category 6: Reverse (`--reverse` flag)

### Summary

| ID | Command | Expected Output |
|----|---------|-----------------|
| G54 | `--reverse example00.txt` | Usage (reverse); exit 1 |
| G55 | `--reverse=example00.txt` | `Hello World` |
| G56 | `--reverse=example01.txt` | `123` |
| G57 | `--reverse=example02.txt` | `#=\[` |
| G58 | `--reverse=example03.txt` | `something&234` |
| G59 | `--reverse=example04.txt` | `abcdefghijklmnopqrstuvwxyz` |
| G60 | `--reverse=example05.txt` | `\!" #$%&'()*+,-./` |
| G61 | `--reverse=example06.txt` | `:;{=}?@` |
| G62 | `--reverse=example07.txt` | `ABCDEFGHIJKLMNOPQRSTUVWXYZ` |
| G63 | Random mixed-case file | Exact original string |
| G64 | Random lowercase+numbers file | Exact original string |
| G65 | Random special chars file | Exact original string |

### G54 — Usage Message (reverse)
**Input:** `go run . --reverse example00.txt`
**Expected (stderr):**
```
Usage: go run . [OPTION]

EX: go run . --reverse=<fileName>
```

### G55–G62 — Known Round-Trips
```bash
go run . --reverse=example00.txt  # → Hello World
go run . --reverse=example01.txt  # → 123
go run . --reverse=example02.txt  # → #=\[
go run . --reverse=example03.txt  # → something&234
go run . --reverse=example04.txt  # → abcdefghijklmnopqrstuvwxyz
go run . --reverse=example05.txt  # → \!" #$%&'()*+,-./
go run . --reverse=example06.txt  # → :;{=}?@
go run . --reverse=example07.txt  # → ABCDEFGHIJKLMNOPQRSTUVWXYZ
```

### Round-Trip Integrity Test
```bash
# Generate art, then reverse it — must recover exact input
go run . --output=rt.txt "Hello World" standard
go run . --reverse=rt.txt            # must print: Hello World
```

---

## Combined Output + Color Golden Cases

| ID | Command | Expected |
|----|---------|----------|
| G66 | `--output=out.txt --color=red "hello" standard` | File with red ANSI sequences |
| G67 | `--output=t.txt --color=red "Hello" --color=blue "World" "Hello World" standard` | File: "Hello"=red, "World"=blue |
| G68 | `--output=s.txt --color=green "a" --color=yellow "king" "a king kitten" shadow` | "a"=green, "king"=yellow in file |
| G69 | `--output=h.txt --color=#FF5733 "123" --color=rgb(0,255,0) "456" "123-456" standard` | TrueColor ANSI in file |

### Validation
```bash
go run . --output=out.txt --color=red "hello" standard
cat -v out.txt   # verify ^[[31m before glyphs and ^[[0m after
cat -e out.txt   # verify $ line endings
```

---

## Peeking Gopher

**Trigger condition:** `text != ""` AND `containsNonASCII(text)` AND every parsed segment is blank.

```bash
go run . "こんにちは"    # triggers gopher
go run . "𝕳𝖊𝖑𝖑𝖔"       # triggers gopher
go run . "helloこんにちは" # does NOT trigger (ASCII portion "hello" is renderable)
go run . ""              # does NOT trigger (empty text)
```

**Expected:** Hardcoded ANSI gopher art printed to stdout; exit 0.

---

## Final Integration Checklist

```bash
# All test files created and valid
cat -e test00.txt test01.txt test02.txt test03.txt test04.txt test05.txt test06.txt test07.txt

# Standard packages only
go list -deps ./... | grep -v "^ascii-art-reverse"

# Build clean
go build -o ascii-art-reverse .

# Full test suite
go test ./... -v

# Coverage
go test ./... -cover

# No formatting issues
gofmt -l .

# No vet warnings
go vet ./...

# No data races
go run -race . "Hello World" standard
```
