# Audit Cases: ascii-art-reverse

Functional audit cases verifying the core requirements of each feature category.
CLI-only project — all cases use `go run .` from the project root.

---

## Summary Table

| ID | Category | Description | Input | Expected |
|:--:|----------|-------------|-------|----------|
| A01 | CLI Core | Empty string | `""` | No output; exit 0 |
| A02 | CLI Core | Single newline escape | `"\n"` | One blank line |
| A03 | CLI Core | Trailing newline | `"Hello\n"` | 8-line block + blank line |
| A04 | CLI Core | Lowercase word | `"hello"` | 8-line standard art |
| A05 | CLI Core | Mixed case | `"HeLlO"` | Correct upper/lower glyph mapping |
| A06 | CLI Core | Multiple words | `"Hello There"` | Horizontal concatenation with space glyph |
| A07 | CLI Core | Numbers and letters | `"1Hello 2There"` | Numeric + alpha glyph mapping |
| A08 | CLI Core | Braces | `"{Hello There}"` | Symbols at ASCII 123/125 rendered |
| A09 | CLI Core | Multi-line | `"Hello\nThere"` | Two separate 8-line blocks |
| A10 | CLI Core | Consecutive newlines | `"Hello\n\nThere"` | Two blocks + one blank line between |
| A11 | Color | Invalid flag format | `--color red "banana"` | Usage message; exit 1 |
| A12 | Color | Full string color | `--color=red "hello world"` | All glyphs in red ANSI |
| A12b | Color | Green color, numbers and symbols | `--color=green "1 + 1 = 2"` | All glyphs in green ANSI |
| A12c | Color | Yellow color, special chars | `--color=yellow "(%&) ??"` | All glyphs in yellow ANSI |
| A13 | Color | Substring color | `--color=red "ello" "hello"` | Only "ello" glyphs colored |
| A14 | Color | Case-sensitive substring | `--color=orange GuYs "HeY GuYs"` | Only "GuYs" colored |
| A15 | Color | Single character target | `--color=blue B "RGB()"` | Only "B" colored |
| A16 | Color | Multi-color, two targets | `--color=red "he" --color=blue "lo" "hello"` | "he"=red, "lo"=blue |
| A17 | Color | Hex notation | `--color=#ff0000 "hello"` | Red via truecolor |
| A18 | Color | RGB notation | `--color=rgb(255,0,0) "hello"` | Red via truecolor |
| A19 | Color | HSL notation | `--color=hsl(0,100%,50%) "hello"` | Red via truecolor |
| A20 | FS | Standard banner | `"hello" standard` | Standard 8-line render |
| A21 | FS | Shadow banner | `"hello world" shadow` | Shadow-style glyphs |
| A22 | FS | Thinkertoy banner | `"nice 2 meet you" thinkertoy` | Thinkertoy-style glyphs |
| A23 | FS | Symbol (ampersand) | `"you & me" standard` | Correct & glyph |
| A24 | FS | Numbers (shadow) | `"123" shadow` | Shadow numeric glyphs |
| A25 | FS | Full alphabet | `"ABCDEFGHIJKLMNOPQRSTUVWXYZ" shadow` | All 26 uppercase glyphs |
| A26 | FS | Special characters | `"It's Working" thinkertoy` | Apostrophe + mixed case |
| A26b | FS | Parentheses and backslash | `"/(\")" thinkertoy` | Parens + backslash + quote glyphs |
| A26c | FS | Punctuation set (thinkertoy) | `"\"#$%&/()*+,-./" thinkertoy` | Full punctuation range rendered |
| A27 | Output | Invalid flag format | `--output test00.txt banana standard` | Usage message; exit 1 |
| A28 | Output | Multi-line to file | `--output=test00.txt "First\nTest" shadow` | File with two 8-line blocks |
| A29 | Output | Standard to file | `--output=test01.txt "hello" standard` | File with 8-line "hello" |
| A30 | Output | Numbers+symbols to file | `--output=test02.txt "123 -> #$%" standard` | File with correct mapping |
| A31 | Output | Long string to file | `--output=test07.txt "Testing long output!" standard` | No internal wrapping |
| A31b | Output | Numbers+symbols to shadow file | `--output=test03.txt "432 -> #$%&@" shadow` | File with correct shadow mapping |
| A31c | Output | Single word to shadow file | `--output=test04.txt "There" shadow` | File with shadow "There" art |
| A31d | Output | Numbers+symbols to thinkertoy file | `--output=test05.txt "123 -> \"#$%@" thinkertoy` | File with thinkertoy mapping |
| A31e | Output | Short string to thinkertoy file | `--output=test06.txt "2 you" thinkertoy` | File with thinkertoy art |
| A32 | Alignment | Invalid flag format | `--align right "something" standard` | Usage message; exit 1 |
| A33 | Alignment | Right alignment | `--align=right "left" standard` | Flush against right border |
| A34 | Alignment | Left alignment | `--align=left "right" standard` | Flush against left border |
| A35 | Alignment | Center alignment | `--align=center "hello" shadow` | Centered in terminal |
| A36 | Alignment | Justify alignment | `--align=justify "1 Two 4" shadow` | Words span full terminal width |
| A36b | Alignment | Justify long string | `--align=justify "HELLO there HOW are YOU?!" thinkertoy` | All words distributed across full width |
| A37 | Alignment | Resize adaptation | `--align=right "abcd" shadow` (resize terminal) | Width recalculated on each run |
| A37b | Alignment | Right, numbers and slash | `--align=right 23/32 standard` | Art flush against right border |
| A37c | Alignment | Right, mixed alphanumeric | `--align=right ABCabc123 thinkertoy` | Art flush against right border |
| A37d | Alignment | Right, arrows string | `--align=right "a -> A b -> B c -> C" shadow` | Art flush against right border |
| A37e | Alignment | Center, special chars | `--align=center "#$%&\"" thinkertoy` | Art centered in terminal |
| A37f | Alignment | Left, mixed string | `--align=left "23Hello World!" standard` | Art starts at column 0 |
| A38 | Reverse | Invalid flag format | `--reverse example00.txt` (space) | Usage message; exit 1 |
| A39 | Reverse | Hello World | `--reverse=example00.txt` | `Hello World` |
| A40 | Reverse | Numbers | `--reverse=example01.txt` | `123` |
| A41 | Reverse | Special chars | `--reverse=example02.txt` | `#=\[` |
| A42 | Reverse | Mixed string | `--reverse=example03.txt` | `something&234` |
| A43 | Reverse | Lowercase alphabet | `--reverse=example04.txt` | `abcdefghijklmnopqrstuvwxyz` |
| A44 | Reverse | Punctuation set | `--reverse=example05.txt` | `\!" #$%&'()*+,-./` |
| A45 | Reverse | Symbols | `--reverse=example06.txt` | `:;{=}?@` |
| A46 | Reverse | Uppercase alphabet | `--reverse=example07.txt` | `ABCDEFGHIJKLMNOPQRSTUVWXYZ` |

---

## Category 1: CLI Core

### A01 — Empty String
**Input:** `go run . "" | cat -e`
**Expected:** No output; exit status 0.
**Rules:** Empty string produces no art and no error.

---

### A02 — Single Newline Escape
**Input:** `go run . "\n" | cat -e`
**Expected:**
```
$
```
**Rules:** Literal `\n` in the argument is interpreted as a real newline, producing one blank line.

---

### A03 — Trailing Newline
**Input:** `go run . "Hello\n" | cat -e`
**Expected:** 8-line "Hello" block followed by one blank line (`$`).
**Rules:** The trailing `\n` adds an empty segment after the rendered word.

---

### A04 — Standard Lowercase
**Input:** `go run . "hello" | cat -e`
**Expected:** 8-line horizontal ASCII art for "hello" in standard banner.
**Rules:** All 5 glyphs concatenated horizontally; each line ends with `$`.

---

### A05 — Mixed Case
**Input:** `go run . "HeLlO" | cat -e`
**Expected:** Correct upper/lower glyph mapping; H, L, O from upper; e, l from lower.
**Rules:** Upper and lower case indices in the banner file are distinct.

---

### A06 — Multiple Words (Space Glyph)
**Input:** `go run . "Hello There" | cat -e`
**Expected:** "Hello" and "There" rendered with the space glyph between them.
**Rules:** Space glyph width must be consistent; horizontal concatenation must be exact.

---

### A07 — Numbers and Letters
**Input:** `go run . "1Hello 2There" | cat -e`
**Expected:** Numeric glyphs for "1" and "2" correctly mapped alongside alphabetic glyphs.

---

### A08 — Braces and Symbols
**Input:** `go run . "{Hello There}" | cat -e`
**Expected:** `{` (ASCII 123) and `}` (ASCII 125) rendered using their banner glyphs.

---

### A09 — Multi-line (`\n` in middle)
**Input:** `go run . "Hello\nThere" | cat -e`
**Expected:** Two separate 8-line blocks stacked vertically.
**Rules:** The literal `\n` splits the render into two independent segments.

---

### A10 — Consecutive Newlines
**Input:** `go run . "Hello\n\nThere" | cat -e`
**Expected:** "Hello" (8 lines), one blank line, "There" (8 lines).
**Rules:** An empty segment between two newlines produces exactly one blank line.

---

## Category 2: Color

### A11 — Invalid Flag Format
**Input:** `go run . --color red "banana"`
**Expected:**
```
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```
**Rules:** `--color` with a space instead of `=` is rejected; exit 1.

---

### A12 — Full String Coloring
**Input:** `go run . --color=red "hello world"`
**Expected:** All rendered lines wrapped in red ANSI start code + reset `\x1b[0m`.
**Rules:** Every glyph column wrapped individually; no color bleeding between characters.

---

### A12b — Green Color, Numbers and Symbols
**Input:** `go run . --color=green "1 + 1 = 2"`
**Expected:** All rendered glyph lines wrapped in green ANSI start code + reset `\x1b[0m`.

---

### A12c — Yellow Color, Special Characters
**Input:** `go run . --color=yellow "(%&) ??"`
**Expected:** All rendered glyph lines wrapped in yellow ANSI start code + reset `\x1b[0m`.

---

### A13 — Substring Coloring
**Input:** `go run . --color=red "ello" "hello"`
**Expected:** Only the 4 glyphs matching "ello" are colored red; "h" remains default.
**Rules:** Case-sensitive match; all occurrences colored.

---

### A14 — Case-Sensitive Substring
**Input:** `go run . --color=orange GuYs "HeY GuYs"`
**Expected:** Only "GuYs" colored orange; "HeY " and surrounding space remain default.

---

### A15 — Single Character Target
**Input:** `go run . --color=blue B "RGB()"`
**Expected:** Only the "B" glyph colored blue; R, G, (, ), remain default.

---

### A16 — Multi-Color, Two Targets
**Input:** `go run . --color=red "he" --color=blue "lo" "hello"`
**Expected:** "he"=red ANSI, "lo"=blue ANSI, "l" (middle) = default.

---

### A17–A19 — Color Notation Support
Each of the following must produce visually identical red output for `"hello"`:
```bash
go run . --color=red "hello"
go run . --color=#ff0000 "hello"
go run . --color=rgb(255,0,0) "hello"
go run . --color=hsl(0,100%,50%) "hello"
```
**Rules:** All four notations resolve to the same red output.

---

## Category 3: File System / Banner Selection

### A20 — Standard Banner
**Input:** `go run . "hello" standard | cat -e`
**Expected:** Same output as `go run . "hello"` (default).

---

### A21 — Shadow Banner
**Input:** `go run . "hello world" shadow | cat -e`
**Expected:** Shadow-style glyphs; underscore-based characters.

---

### A22 — Thinkertoy Banner
**Input:** `go run . "nice 2 meet you" thinkertoy | cat -e`
**Expected:** Thinkertoy-style round-character glyphs.

---

### A23 — Ampersand Symbol (Standard)
**Input:** `go run . "you & me" standard | cat -e`
**Expected:** `&` (ASCII 38) rendered using its standard glyph.

---

### A24 — Numeric String (Shadow)
**Input:** `go run . "123" shadow | cat -e`
**Expected:** Shadow-style numeric glyphs for 1, 2, 3.

---

### A25 — Full Uppercase Alphabet (Shadow)
**Input:** `go run . "ABCDEFGHIJKLMNOPQRSTUVWXYZ" shadow | cat -e`
**Expected:** All 26 uppercase glyphs rendered in shadow style.

---

### A26 — Mixed Case with Apostrophe (Thinkertoy)
**Input:** `go run . "It's Working" thinkertoy | cat -e`
**Expected:** Apostrophe (ASCII 39) and all letters rendered in thinkertoy style.

---

### A26b — Parentheses and Backslash (Thinkertoy)
**Input:** `go run . "/(\")" thinkertoy | cat -e`
**Expected:**
```
         o o    $
    o  / | | \  $
   /  o       o $
  o   |       | $
 /    o       o $
o      \     /  $
                $
                $
```

---

### A26c — Punctuation Set (Thinkertoy)
**Input:** `go run . "\"#$%&/()*+,-./" thinkertoy | cat -e`
**Expected:**
```
o o         | |                                                  $
| |  | |   -O-O-      O          o  / \  o | o                 o $
    -O-O- o | |   o  /    o     /  o   o  \|/   |             /  $
     | |   -O-O-    /    /|    o   |   | --O-- -o-           o   $
    -O-O-   | | o  /  o o-O-  /    o   o  /|\   |    o-o    /    $
     | |   -O-O-  O       |  o      \ /  o | o     o     O o     $
            | |                                    |             $
                                                                 $
```

---

## Category 4: Output (`--output` flag)

### A27 — Invalid Flag Format
**Input:** `go run . --output test00.txt banana standard`
**Expected:**
```
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --output=<fileName.txt> something standard
```
**Rules:** Space instead of `=` is rejected; exit 1.

---

### A28 — Multi-line Shadow to File
**Input:** `go run . --output=test00.txt "First\nTest" shadow`
**Validate:** `cat -e test00.txt`
**Expected:** Two 8-line shadow blocks with correct `$`-terminated line endings.
**Rules:** File created/overwritten with O_TRUNC; ANSI codes written unchanged.

---

### A29 — Standard Hello to File
**Input:** `go run . --output=test01.txt "hello" standard`
**Validate:** `cat -e test01.txt`
**Expected:** 8-line "hello" in standard art with `$`-terminated lines.

---

### A30 — Numbers and Symbols to File
**Input:** `go run . --output=test02.txt "123 -> #$%" standard`
**Validate:** `cat -e test02.txt`
**Expected:** Correct numeric and symbol glyph mapping persisted.

---

### A31 — Long String Integrity
**Input:** `go run . --output=test07.txt "Testing long output!" standard`
**Validate:** `cat -e test07.txt`
**Expected:** Full-width single-line art; no internal line wrapping.

---

### A31b — Numbers and Symbols to Shadow File
**Input:** `go run . --output=test03.txt "432 -> #$%&@" shadow`
**Validate:** `cat -e test03.txt`
**Expected:**
```
                                                                                                                  $
_|  _|   _|_|_|     _|_|                    _|             _|  _|     _|   _|_|    _|   _|           _|_|_|_|_|   $
_|  _|         _| _|    _|                    _|         _|_|_|_|_| _|_|_| _|_|  _|   _|  _|       _|          _| $
_|_|_|_|   _|_|       _|         _|_|_|_|_|     _|         _|  _|   _|_|       _|       _|_|  _| _|    _|_|_|  _| $
    _|         _|   _|                        _|         _|_|_|_|_|   _|_|   _|  _|_| _|    _|   _|  _|    _|  _| $
    _|   _|_|_|   _|_|_|_|                  _|             _|  _|   _|_|_| _|    _|_|   _|_|  _| _|    _|_|_|_|   $
                                                                      _|                           _|             $
                                                                                                     _|_|_|_|_|_| $
```

---

### A31c — Single Word to Shadow File
**Input:** `go run . --output=test04.txt "There" shadow`
**Validate:** `cat -e test04.txt`
**Expected:**
```
                                               $
_|_|_|_|_| _|                                  $
    _|     _|_|_|     _|_|   _|  _|_|   _|_|   $
    _|     _|    _| _|_|_|_| _|_|     _|_|_|_| $
    _|     _|    _| _|       _|       _|       $
    _|     _|    _|   _|_|_| _|         _|_|_| $
                                               $
                                               $
```

---

### A31d — Numbers and Symbols to Thinkertoy File
**Input:** `go run . --output=test05.txt "123 -> \"#$%@" thinkertoy`
**Validate:** `cat -e test05.txt`
**Expected:**
```
                                    o o         | |               $
  0    --  o-o            o         | |  | |   -O-O-      O   o   $
 /|   o  o    |            \            -O-O- o | |   o  /   / \  $
o |     /   oo              O            | |   -O-O-    /   o O-o $
  |    /      |       o-o  /            -O-O-   | | o  /  o  \    $
o-o-o o--o o-o            o              | |   -O-O-  O       o-  $
                                                | |               $
                                                                  $
```

---

### A31e — Short String to Thinkertoy File
**Input:** `go run . --output=test06.txt "2 you" thinkertoy`
**Validate:** `cat -e test06.txt`
**Expected:**
```
                         $
 --                      $
o  o                     $
  /        o  o o-o o  o $
 /         |  | | | |  | $
o--o       o--O o-o o--o $
              |          $
           o--o          $
```

---

## Category 5: Alignment / Justify (`--align` flag)

### A32 — Invalid Flag Format
**Input:** `go run . --align right "something" standard`
**Expected:**
```
Usage: go run . [OPTION] [STRING] [BANNER]

Example: go run . --align=right something standard
```
**Rules:** Space instead of `=` is rejected; exit 1.

---

### A33 — Right Alignment
**Input:** `go run . --align=right "left" standard`
**Expected:** Each of the 8 art lines padded on the left so the art flush-aligns to the right terminal border.

---

### A34 — Left Alignment
**Input:** `go run . --align=left "right" standard`
**Expected:** No padding; art starts at column 0.

---

### A35 — Center Alignment
**Input:** `go run . --align=center "hello" shadow`
**Expected:** Each line padded with `(terminalWidth − visibleWidth) / 2` spaces on the left.

---

### A36 — Justify Alignment
**Input:** `go run . --align=justify "1 Two 4" shadow`
**Expected:** "1" at left border, "4" at right border, "Two" distributed between them.
**Rules:** `baseGap = totalGap / gaps`; first `remainder` gaps get `baseGap + 1`.

---

### A36b — Justify Long String (Thinkertoy)
**Input:** `go run . --align=justify "HELLO there HOW are YOU?!" thinkertoy`
**Expected:** All six words distributed so the leftmost word starts at column 0 and the rightmost ends at the terminal right border.
**Rules:** Same gap-distribution formula as A36; more words = more gaps to distribute.

---

### A37 — Terminal Resize Adaptation
**Input:** `go run . --align=right "abcd" shadow` (run again after resizing terminal)
**Expected:** Width recalculated from `COLUMNS` env or ioctl on each invocation.

---

### A37b — Right, Numbers and Slash
**Input:** `go run . --align=right 23/32 standard`
**Expected:** Art padded on the left; rightmost column flush with terminal border.

---

### A37c — Right, Mixed Alphanumeric (Thinkertoy)
**Input:** `go run . --align=right ABCabc123 thinkertoy`
**Expected:** Art padded on the left; rightmost column flush with terminal border.

---

### A37d — Right, Arrow String (Shadow)
**Input:** `go run . --align=right "a -> A b -> B c -> C" shadow`
**Expected:** Art padded on the left; rightmost column flush with terminal border.

---

### A37e — Center, Special Characters (Thinkertoy)
**Input:** `go run . --align=center "#$%&\"" thinkertoy`
**Expected:** Each line padded with `(terminalWidth − visibleWidth) / 2` spaces on the left.

---

### A37f — Left, Mixed String (Standard)
**Input:** `go run . --align=left "23Hello World!" standard`
**Expected:** No padding; art starts at column 0.

---

## Category 6: Reverse (`--reverse` flag)

### A38 — Invalid Flag Format
**Input:** `go run . --reverse example00.txt`
**Expected:**
```
Usage: go run . [OPTION]

EX: go run . --reverse=<fileName>
```
**Rules:** Space instead of `=` is rejected; exit 1.

---

### A39–A46 — Known File Round-Trips

| ID | Input | Expected Output |
|:--:|-------|-----------------|
| A39 | `go run . --reverse=example00.txt` | `Hello World` |
| A40 | `go run . --reverse=example01.txt` | `123` |
| A41 | `go run . --reverse=example02.txt` | `#=\[` |
| A42 | `go run . --reverse=example03.txt` | `something&234` |
| A43 | `go run . --reverse=example04.txt` | `abcdefghijklmnopqrstuvwxyz` |
| A44 | `go run . --reverse=example05.txt` | `\!" #$%&'()*+,-./` |
| A45 | `go run . --reverse=example06.txt` | `:;{=}?@` |
| A46 | `go run . --reverse=example07.txt` | `ABCDEFGHIJKLMNOPQRSTUVWXYZ` |

**Rules:** Each example file contains standard-banner ASCII art; output must be exact string, no trailing newline artifacts.

---

### A47 — Random Mixed-Case Reverse
**Input:** File containing ASCII art for a random mixed-case string.
**Expected:** Exact original string printed; exit 0.

### A48 — Random Lowercase + Numbers + Spaces
**Input:** File containing ASCII art with lowercase letters, numbers, and spaces.
**Expected:** Exact original string printed.

### A49 — Random Special Characters Reverse
**Input:** File containing ASCII art with special characters.
**Expected:** Exact original string printed.

---

## Audit Quality Checklist

| Check | Command |
|-------|---------|
| Standard packages only | `go list -deps ./... \| grep -v "^ascii-art-reverse"` |
| Build clean | `go build -o ascii-art-reverse .` |
| Tests exist and pass | `go test ./... -v` |
| Coverage adequate | `go test ./... -cover` |
| No formatting issues | `gofmt -l .` |
| No vet warnings | `go vet ./...` |
| No data races | `go run -race . "hello"` |
