# PRD - ascii-art-reverse (Unified Suite)

## 1. Problem Statement

The objective is to build a unified ASCII art utility that supports **file persistence**, **color rendering**, **custom banner files**, **text alignment**, and **reverse engineering**. The tool must allow users to:
- Redirect the stylized graphic representation of a string into a physical file using the `--output` flag
- Apply ANSI color codes to full strings or targeted substrings using the `--color` flag
- Use custom banner styles beyond the original three (standard, shadow, thinkertoy) by loading additional banner files from the `banner/` directory
- Align the rendered ASCII art to the terminal width using `--align=<type>` where `type` is `left`, `right`, `center`, or `justify`
- **Reconstruct original text from an ASCII art file using the `--reverse` flag.**

All features must strictly adhere to the established 8-line-high character mapping and banner rules, with support for 855-line banner files (95 characters × 9 lines: 8 art lines + 1 separator).

**Architecture Note:** This project implements a refactored modular architecture. Logic is strictly separated into `internal/cli` (orchestration), `internal/banner` (loading), `internal/render` (pipeline), `internal/output` (I/O), and `internal/reverse` (reconstruction) to ensure high maintainability and testability.

## 2. User / Use Case

- Primary user: Developers or students using a terminal environment.
- Use case:
  - **Documentation**: Generating ASCII banners to be included in .txt or .md documentation files.
  - **Automation**: Using scripts to generate stylized logs or "MOTD" (Message of the Day) files for servers.
  - **Persistence**: Saving complex colorized art (if combined with color/output workflows) for later viewing without re-running the generator.
  - **Customization**: Using custom banner styles (doom, greek, dancing, block) for different visual aesthetics.
  - **Branding**: Applying color schemes to ASCII art for terminal-based branding or emphasis.
  - **Multi-color Art**: Creating visually rich ASCII art with multiple colors in a single output.
  - **Alignment**: Aligning ASCII art to the terminal width for clean console output or screenshots.
- **Recovery**: Decoding existing ASCII art files back into plain text for editing or verification.

---

## 3. CLI Contract
The program acts as a sophisticated CLI dispatcher. It must distinguish between standard art generation, color injection, alignment, file redirection, and reverse reconstruction based on the presence and format of specific flags.

### 3.1 Argument Patterns & Priority
The program must support the following combinations. If multiple options are present, the order of priority for parsing is:   
**Flag -> Substring (in case of color injection) -> String -> Banner.**

| Mode | Command Syntax |
|------------------|--------------------------------------------------------------|
| Standard | `go run . [STRING] [BANNER]` |
| Output Redirection | `go run . --output=<fileName.txt> [STRING] [BANNER]` |
| Full Color | `go run . --color=<color> [STRING] [BANNER]` |
| Targeted Color | `go run . --color=<color> [SUBSTRING] [STRING] [BANNER]` |
| Multi-Color Targeted | `go run . --color=<color1> [SUBSTRING1] --color=<color2> [SUBSTRING2] [STRING] [BANNER]` |
| Alignment | `go run . --align=<type> [STRING] [BANNER]` |
| 2-Flags Combined | `go run . --output=<file.txt> --color=<color> [SUBSTRING] [STRING] [BANNER]` |
| Full Combined |  `go run . --output=<file.txt> --color=<color1> [SUBSTRING1] --color=<color2> [SUBSTRING2] [STRING] [BANNER]` |
| Reverse | `go run . --reverse=<fileName>` |
| Reverse with Banner | `go run . --reverse=<fileName> [BANNER]` |

- **Note on Compatibility:** The `--color=<color>` flag remains fully functional. When used alongside `--output`, the ANSI escape codes for color must be written directly into the text file (allowing users to view color when using `cat` on the resulting file).

---

### 3.2 Strict Flag Validation
To pass the audit, the program must enforce an exact string match for the flag prefix. Spaces are not permitted between the flag and its value.
- Valid: `--output=test.txt`, `--color=red`
- Invalid: `--output test.txt`, `--color red`
- Valid: `--align=right`
- Invalid: `--align right`
- Valid: `--reverse=file.txt`
- Invalid: `--reverse file.txt`

---

### 3.3 Usage Messages (Audit Mandatory)
If the arguments are malformed, the program must print the corresponding usage message and exit with status `1`.
- **Output Flag Error Message:**

  ```plaintext
  Usage: go run . [OPTION] [STRING] [BANNER]

  EX: go run . --output=<fileName.txt> something standard
  ```
- **Color Flag Error Message:**
  ```plaintext
  Usage: go run . [OPTION] [STRING]

  EX: go run . --color=<color> <substring to be colored> "something"
  ```
  - **Banner Error Message:**
  ```plaintext
  Usage: go run . [STRING] [BANNER]
  
  EX: go run . something standard
  ```

- **Alignment Flag Error Message:**
  ```plaintext
  Usage: go run . [OPTION] [STRING] [BANNER]

  Example: go run . --align=right something standard
  ```

---

## 3.4 Logic Flow: The Single-Pass Classifier

The entire command-line interface is managed by the `internal/cli` package, which implements a high-performance, single-pass classification system. This replaces scattered parsing with a deterministic bucket-sorting approach.

**The `ParsedArgs` Struct:** This object acts as the single source of truth for the orchestrator (`main.go`).
```go
type ParsedArgs struct {
    OutputValue  string
    AlignValue   string
    ReverseValue string
    ColorRules   []RawColorRule
    StdinMode    bool
    Positional   []string      // [text] or [text, banner]
    UnknownFlags []FlagError
    Malformed    []FlagError
}
```

**Architectural Steps:**

1. **Single-Pass Classification (`ClassifyArgs`)**:
   - A single loop iterates over `os.Args[1:]`.
   - Every token is categorized exactly once into the `ParsedArgs` struct.
   - Known flags use the strict `--flag=value` syntax.

2. **Color Pairing**:
   - When a `--color` flag is encountered, the classifier peeks at the next token.
   - If it's a non-flag, it's consumed as a substring target.

3. **Positional Resolution**:
   - Remaining non-flag tokens are captured in order.
   - The orchestrator later interprets them as `[text]` or `[text, banner]`.

4. **Usage Selection**:
   - The `SelectUsage` function evaluates the `ParsedArgs` for malformed flags or unknown inputs.
   - It selects the highest-priority usage string (`Reverse > Output > Align > Color`) for the Fatal exit.

5. **Orchestration**:
   - `main.go` remains under 60 lines by simply calling `ClassifyArgs`
     and passing the results to specialized internal packages (`banner`, `render`, `output`, `reverse`).`

This architecture ensures that the program behaves predictably regardless of flag order and provides clear, prioritized feedback to the user.

---

**Audit Example**
```bash
go run . --output=text.txt --color=blue l "Hello World!" standard
```

**Internal Logic Breakdown:**
- `OutputTarget`: `text.txt`
- `ActiveColor`: blue (ANSI `\033[34m`)
- `SubString`: `l`
- `MainString`: `Hello World!`
- `Banner`: `standard.txt`

**Validation:**
Running `cat -e text.txt` should show the ASCII art for `Hello World!` with `l` wrapped in escape codes.

---

### 3.5 Extended CLI Contract: Multi-Color & Output Integration

The program must support a **"Cumulative Buffer"** approach. Each color flag is registered in a `ColorMap` before the final rendering and file-writing phase begins.

**Command Pattern**
```bash
go run . --output=result.txt --color=red "Hello" --color=blue "World" "Hello World" standard
```

**Expected Behavior:**
- The substring **"Hello"** is tagged for Red.
- The substring **"World"** is tagged for Blue.
- The resulting 8-line ASCII art is generated with mixed ANSI codes and written to `result.txt`.
 - **Invalid color values** inside multi-color mode emit a warning, are ignored, and do not block valid color rules.
 - **Malformed color flags** in multi-color mode emit a warning, are ignored, and do not block valid color rules.
 - If multiple rules overlap, the **last matching rule wins** for those rune positions.

---

### 3.6 Updated CLI Contract: Alignment Rules
Alignment introduces a strict `--align=<type>` flag where `type` is one of:
`left`, `right`, `center`, `justify`.

**Rules**
- If `--align` is missing, default to `left` alignment.
- Invalid formats (e.g., `--align right`) must print the alignment usage message and exit `1`.
- Invalid `type` values must print the alignment usage message and exit `1`.
- Alignment must adapt to terminal width using only Go standard library methods (no external terminal packages).
- Only text that fits the terminal size will be tested.
- Alignment must preserve 8-line character height integrity.
- Alignment may be combined with other correctly formatted flags (`--output`, `--color`) and remains order-independent.

**Justify Behavior**
- Distribute spaces between words to span the terminal width.
- If the input contains only one word, treat `justify` as `left`.

### 3.7 Updated CLI Contract: FS Argument Rules & Custom Banners
The FS project introduces dynamic banner loading and stricter argument validation.

**Supported Custom Banners:**
In addition to the original three banners (standard, shadow, thinkertoy), the following custom banners are supported:
- `doom` - Bold, blocky style with trailing separator format
- `greek` - Greek-inspired style with 12-character width
- `dancing` - Animated-style characters
- `block` - Block-style characters

**Banner File Requirements:**
- All banner files must be exactly 855 lines (95 characters × 9 lines)
- Each character occupies 9 lines: 8 art lines + 1 separator line
- Separator can be leading (first line) or trailing (last line)
- The program automatically detects separator format

| Mode | Command Syntax | Requirement |
| --- | --- | --- |
| FS Standard | `go run . [STRING] [BANNER]` | *If two arguments are provided (post-flags), the second MUST be a valid banner name or filename in the banner/ folder.* |
| FS Single | `go run . [STRING]` | *Must still function by defaulting to banner/standard.txt.* |
| FS Excessive Args | `go run . [STRING] [BANNER] [EXTRA]` | *Three or more arguments (e.g., "banana" standard abc) must trigger the FS usage message and exit 1.* |
| Custom Banner | `go run . [STRING] doom` | *Loads banner/doom.txt and renders using detected separator format.* |

---

### 3.8 Logic Summary Table for Audit (Advanced Integration)

This table serves as the definitive guide for how the program resolves complex argument strings.

| Priority | Feature | Parsing Rule | Handling for Combined Mode |
|----------|---------|--------------|----------------------------|
| 1. | **Flag Detection** | Scan all os.Args for `--output=`, `--color=`, or `--align=`. | Extract filename from `--output`. Populate a `ColorRegistry` slice with all `--color` pairs. Capture alignment type from `--align`. |
| 2. | **Positional Extraction** | Remove flags and their direct values (color name/file name). | The remaining arguments are assigned as `[SUBSTRING]` (if any left), then `[STRING]`, then `[BANNER]`. |
| 3. | **Conflict Resolution** | If `--output` is defined twice, the last one provided takes precedence. | Multi-color flags are additive, not exclusive. |
| 4. | **Color Mapping** | Use a bitmask or interval map to identify which runes in `[STRING]` get which ANSI code. | Ensure \"Reset\" codes (\033[0m) are inserted between different color transitions to prevent bleeding. |
| 5. | **Output Target** | Check if `OutputTarget` variable is set. | If set, use `os.O_TRUNC` to prepare the file. If not, default to `os.Stdout`. |
| 6. | **Alignment** | Apply left/right/center/justify padding using terminal width. | Preserve 8-line height and align each rendered line before writing. |

**Audit Example**
**Command**: 
```bash
go run . --output=text.txt --color=red H --color=blue W "Hello World" standard
```
- **Parser:** Extracts `text.txt` as target. Registers **Red** -> `H` and **Blue** -> `W`.
- **Colorizer:** Scans "Hello World". Marks `index 0` (H) as **Red** and `index 6` (W) as **Blue**.
- **Renderer:** Generates 8 lines of art. For each line it:
  1. injects **Red ANSI** 
  2. renders the `"H"` glyph 
  3. injects *Reset* 
  4. renders `"ello "`
  5. injects **Blue ANSI**
  6. renders the `"W"` glyph
  7. injects *Reset*
  8. finishes the line.
- **Writer:** Writes all 8 lines (including hidden ANSI codes) into log.txt.

---

## 4. Functional Requirements (Rules)
This section defines the system-level interactions for file I/O and the specific escape sequences used to satisfy the **output**, **color**, and **alignment** requirements.

### 4.1 Core Fundamentals & Color Notation
The system relies on a Modular Pipeline architecture. The components include:
- **Loader:** Reads the 855-line banner files (mapping ASCII characters 32-126).
- **Parser:** Normalizes line endings by converting `\r` and `\r\n` to `\n`, and handles tabs (`\t`).
- **Mapper:** Retrieves the corresponding 8-line glyphs.
- **Alignment:** Pads rendered lines based on terminal width for left/right/center/justify.

To meet the requirements of the color optional project, the program must support various color notations, which are translated into ANSI escape sequences before being passed to the **Renderer**:
- **Standard ANSI:** Basic colors like red, blue, green, yellow, magenta, cyan, white, and orange.
- **Hex Notation:** Supports `#RRGGBB` format (e.g., `--color=#FF5733`).
- **RGB Notation:** Supports `rgb(r, g, b)` format (e.g., `--color=rgb(255, 255, 255)`).
- **HSL Notation:** Supports `hsl(h, s, l)` format (e.g., `--color=hsl(0, 100%, 50%)`).

### 4.2 The File-Writer Module
When the `--output` flag is present:
- The program redirects the final string buffer to a file instead of `os.Stdout`.
- To ensure audit compliance (creating new files or clearing old ones), specific constant flags from the `os` package must be used.

**Required os.OpenFile Logic:**
```go
// Open the file with specific permissions:
// O_WRONLY: Open for writing only
// O_CREATE: Create the file if it doesn't exist
// O_TRUNC: Clear the file if it already exists (overwrite mode)
file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
```
- **Permissions:** `0644` (***Owner***: read/write, ***Group/Others***: read) is the standard for generated text files.

### 4.3 ANSI Map & Sequence Injection
The Colorizer sub-module uses a map to look up correct escape codes:
| Color   | Escape Code | Reset Code |
|---------|--------------|------------|
| Red     | \033[31m    | \033[0m |
| Green   | \033[32m    | \033[0m |
| Blue    | \033[34m    | \033[0m |
| Yellow  | \033[33m    | \033[0m |

If a substring is provided,
the mapper must wrap only targeted character blocks with the sequence and a reset code.

**Multi-Flag Logic (Output + Color):** When combined,
the program writes ANSI escape codes directly into the text file.
**Verification:** When running `cat <file>`, terminal interprets escape codes and displays color.
**Formatting:** The reset code (`\033[0m`) must be placed at end of every colored line to prevent "color bleeding" into terminal prompt.
- The Renderer must remain "stream-agnostic".
- If `--output` is active, target an `os.File*.*` 
- If absent, target `os.Stdout`.  
The **Colorizer** logic applies ANSI codes to strings before they are sent by the **Renderer**.

### 4.4 File Buffer Management & Writing
To prevent color "leaks" or overwrites when writing to a file with multiple flags, the Renderer must follow these rules:
- **Tagging System:** Before rendering, scan the MainString for all substring matches. Assign a "Color ID" to each rune index.
- **ANSI Boundary Injection:**
When writing the 8 horizontal lines to the file, the program must inject the ANSI start code (`\033[...m`) at the start of a matched sequence and the Reset code (`\033[0m`) immediately after the sequence ends.
- **File Stream Integrity:**
Since multiple escape codes exist in one line, the `FileWriter` must use a `strings.Builder` or `bytes.Buffer` to assemble the full 8-line block in memory before executing a single `file.Write()` call. This prevents partial writes if the program crashes.

### 4.5 Data Integrity Rules
- **Overwrite Consistency:** Every execution of the `--output` flag must result in a fresh file. Appending is forbidden.
- **Empty Strings:** If input string is empty,
the program should create an empty file (0 bytes) or one containing only newline,
depending on test case.
- **Permissions:** Files must be created with mode `0644` (owner read/write, group/others read).

### 4.6 FS & Banner Discovery (Custom Banner Support)
The program must transition from hardcoded banner loading to a dynamic file system lookup. This allows for scalability and the addition of custom templates without recompiling the source code.

**Supported Banner Files:**
- `standard.txt` - Original default banner (leading separator)
- `shadow.txt` - Shadow-style characters (leading separator)
- `thinkertoy.txt` - Compact, toy-like style (leading separator)
- `doom.txt` - Bold, blocky style (trailing separator)
- `greek.txt` - Greek-inspired style, 12-char width (leading separator)
- `dancing.txt` - Animated-style characters (leading separator)
- `block.txt` - Block-style characters (leading separator)

**Dynamic Loading Logic**
- **Path Construction:** The program must append `.txt` to the `[BANNER]` argument and look inside the `banner/` folder.
- **Rune Calculation:** Since characters in the provided templates are separated by a newline, the offset logic is:
  - Each character block = 9 lines (`8` art lines + `1` separator).
  - **Separator Detection:** The program must detect whether the separator line is at the beginning (leading) or end (trailing) of each 9-line block using the `isLeadingSeparator` heuristic.
  - **Leading Separator (standard format):** Extract lines 1-8 from each block (skip line 0).
  - **Trailing Separator (alternate format):** Extract lines 0-7 from each block (skip line 8).
  - **Target Line Start:** `(ascii_value - 32) * 9`.
- **File Validation:** All banner files must be exactly 855 lines. If the requested file does not exist in the `banner/` directory or has incorrect line count, the program must fallback safely to standard.txt with a warning to stderr.
- **Path Sanitization:** Banner names must be sanitized using `filepath.Base()` to prevent path traversal attacks (e.g., `../etc/passwd` should be rejected and fallback to standard).
- **Case Normalization:** Banner names are case-insensitive and normalized to lowercase.
- **Extension Handling:** The `.txt` extension is optional when specifying banners.

---

### 4.7 Alignment & Terminal Width
- Alignment is controlled by `--align=<type>` where `type` is `left`, `right`, `center`, or `justify`.
- Alignment must use terminal width detection from the Go standard library only, with a safe fallback when the width cannot be detected.
- Left alignment is the default when no alignment flag is provided.
- Right and center alignment must pad spaces to match terminal width.
- Justify alignment must distribute extra spaces between words to span the terminal width. If there is only one word, treat justify as left.
- Alignment must preserve 8-line character height and not mutate banner assets.

---

## 5. Non-Goals (Out of Scope)

To maintain project focus and strictly follow the audit guidelines, the following features are explicitly excluded from this implementation:

- **External Dependencies:** Use of any third-party libraries (e.g., Cobra, Viper, or color packages like fatih/color) is prohibited. Only the **Go Standard Library** is allowed.
- **Interactive Mode:** No GUI or interactive prompts for color selection or file naming.
- **Advanced File Appending:** The program will not support appending to existing files; it strictly adheres to the `O_TRUNC` (overwrite) behavior for `--output`.
- **Automatic Word Wrapping:** The program will not automatically wrap text that exceeds the terminal width. It will print the full ASCII art width as generated.
- **Multiple Banner Layering:** The program does not support combining different banners in a single execution (e.g., one word in shadow and another in standard).
- **Animation:** Static ASCII art only, no blinking text or scrolling animations.

---

## 6. Acceptance Criteria

### 6.1 Basic Functional Cases
- [ ] `go run . ""` exits `0` with no crash.
- [ ] `go run . "\n"` prints one empty line behavior correctly.
- [ ] `go run . "Hello\n"` renders `Hello` then one trailing empty line.
- [ ] `go run . "hello"` maps lowercase correctly.
- [ ] `go run . "HeLLo"` preserves case-sensitive glyph mapping.
- [ ] `go run . "Hello World"` keeps horizontal alignment with spaces.
- [ ] `go run . "123 Hello"` renders numbers and letters correctly.
- [ ] `go run . "{}"` renders symbols correctly.
- [ ] `go run . "Hello\nThere"` renders two 8-line blocks.
- [ ] `go run . "Hi\n\nBye"` keeps empty middle segment behavior.

### 6.2 Additional Golden Tests
- [ ] `go run .` prints exact section 3.4 usage message to `stderr` and exits `1`.
- [ ] Too many invalid args print exact section 3.4 usage message to `stderr` and exit `1`.
- [ ] Dependency audit confirms only Go standard library packages are imported (no external modules).
- [ ] Rendering path uses precomputed/provider-based O(1) glyph lookup (verified by code audit/perf sanity check).
- [ ] Non-ASCII-only input (100% non-ASCII characters) follows defined contract: no crash, exits `0`, and prints a multi-colored ANSI Peeking Gopher with emojis.
- [ ] Missing/non-existent banner falls back safely per implementation policy.
- [ ] Corrupted banner structure follows safe fallback policy (no usage error), no crash.
- [ ] Path traversal-like banner input is blocked/sanitized.
- [ ] Mixed supported + unsupported characters are handled deterministically without stripping/rejecting the full input.
- [ ] `\t` is expanded to exactly 3 spaces.
- [ ] Backslash collapse and newline dominance follow section 4.5.

### 6.3 Color Feature Audit Cases
- [ ] **Audit Case 1:** `go run . --color red "banana"` rejects malformed flag format and prints exact usage.
- [ ] **Audit Case 2:** `go run . --color=red` rejects missing string and prints exact usage.
- [ ] **Audit Case 3:** `go run . --color=red "hello world"` renders full output in red.
- [ ] **Audit Case 4:** `go run . --color=green "1 + 1 = 2"` renders full output in green.
- [ ] **Audit Case 5:** `go run . --color=yellow "(%&) ??"` renders full output in yellow.
- [ ] **Audit Case 6:** `go run . --color=red "ello" "hello"` colors substring from second to last letter.
- [ ] **Audit Case 7:** `go run . --color=red "e" "hello"` colors only the second letter target.
- [ ] **Audit Case 8:** `go run . --color=red "ll" "hello"` colors only the two-letter substring.
- [ ] **Audit Case 9:** `go run . --color=orange GuYs "HeY GuYs"` applies case-sensitive substring targeting.
- [ ] **Audit Case 10:** `go run . --color=blue B "RGB()"` colors only the `B` character.
- [ ] **Audit Case 11:** Random mixed-case string with random color renders correctly without crashing.
- [ ] **Audit Case 12:** Random lowercase/numbers/spaces string with random color renders correctly without crashing.
- [ ] **Audit Case 13:** Random special-character string with one-letter target colors only intended matches.
- [ ] **Audit Case 14:** Random mixed input with substring target colors only intended substring matches.
- [ ] **Audit Case 15:** Final auditor decision enforces fail reasons (Empty Work, Incomplete Work, Invalid compilation, Cheating, Crashing, Leaks).
- [ ] **Audit Case 16:** Color-targeting syntax is intuitive in practical usage.
- [ ] **Audit Case 17:** Multi-color capability in the same string is demonstrated or documented as unsupported.
- [ ] **Audit Case 18:** Runtime remains efficient with no unnecessary processing.
- [ ] **Audit Case 19:** Colored rendering output keeps correct 8-line alignment.
- [ ] **Audit Case 20:** Multiple color notations are validated (`red`, `#ff0000`, `rgb(...)`, `hsl(...)`).
- [ ] **Audit Case 21:** Multi-color, two valid targets apply different colors in one output.
- [ ] **Audit Case 22:** Multi-color, repeated target uses last rule precedence.
- [ ] **Audit Case 23:** Multi-color, overlapping targets resolve deterministically (last rule wins).
- [ ] **Audit Case 24:** Multi-color, valid + invalid colors proceeds with valid colors and warns.
- [ ] **Audit Case 25:** Multi-color, invalid notation + valid named color proceeds with valid colors and warns.
- [ ] **Audit Case 26:** Multi-color, multiple invalid colors are ignored with warnings; output still renders.
- [ ] **Audit Case 27:** Multi-color, missing target string triggers the color usage message.


### 6.4 Output Feature Audit Cases

- [ ] **Audit Case 1:** `go run . --output test00.txt banana` rejects space in flag and prints usage.
- [ ] **Audit Case 2:** `go run . --output=test00.txt "First\nTest" shadow` saves multi-line shadow art to `test00.txt`.
- [ ] **Audit Case 3:** `go run . --output=test01.txt "hello" standard` saves "hello" to `test01.txt`.
- [ ] **Audit Case 4:** `go run . --output=test02.txt "123 -> #$%" standard` handles symbols and numbers.
- [ ] **Audit Case 5:** `go run . --output=test03.txt "432 -> #$%&@" shadow` saves shadow-style symbols.
- [ ] **Audit Case 6:** `go run . --output=test05.txt "123 -> \"#$%@" thinkertoy` handles literal quotes in file output.
- [ ] **Audit Case 7:** `go run . --output=test07.txt "Testing long output!" standard` preserves horizontal integrity in file.
- [ ] **Audit Case 8:** `go run . --output=anyName.txt "MixedCase123" standard` creates custom filenames and preserves casing.
- [ ] **Audit Case 9:** `go run . --output=multi.txt --color=red "A" --color=blue "B"` persists multi-color ANSI codes into the file.

### 6.5 FS Acceptance Criteria (Custom Banner Support)
- [ ] **Dynamic Pathing:** Program successfully loads all supported banners from `banner/` directory.
- [ ] **Case Sensitivity:** Banner names are treated case-insensitively and normalized to lowercase.
- [ ] **Directory Isolation:** The program does not look for banners outside the `banner/` directory (security constraint).
- [ ] **Binary Integrity:** The character mapping logic (9-line offset with separator detection) works correctly for all provided templates.
- [ ] **Separator Detection:** Program correctly detects and handles both leading separator (standard, shadow, thinkertoy, greek, dancing, block) and trailing separator (doom) formats.
- [ ] **855-Line Validation:** Program validates that all banner files have exactly 855 lines and rejects invalid files.
- [ ] **Argument Overflow:** `go run . "banana" standard abc` triggers the FS usage message and exits 1.
- [ ] **FS Case 52:** `go run . "hello" standard` renders 8-line ASCII with correct line endings.
- [ ] **FS Case 53:** `go run . "hello world" shadow` properly renders shadow-style with space handling.
- [ ] **FS Case 54:** `go run . "nice 2 meet you" thinkertoy` correctly maps numeric characters in thinkertoy style.
- [ ] **FS Case 55:** `go run . "you & me" standard` accurately renders the ampersand symbol.
- [ ] **FS Case 56:** `go run . "123" shadow` renders clean numeric output in shadow style.
- [ ] **FS Case 57:** `go run . "/(\")\" thinkertoy` correctly maps complex punctuation (slashes, parentheses, quotes).
- [ ] **FS Case 58:** `go run . "ABCDEFGHIJKLMNOPQRSTUVWXYZ" shadow` renders complete uppercase alphabet.
- [ ] **FS Case 59:** `go run . "\"#$%&/()*+,-./\" thinkertoy` handles comprehensive special character mapping.
- [ ] **FS Case 60:** `go run . "It's Working" thinkertoy` properly handles mixed-case with apostrophe.
- [ ] **FS Case 61:** `go run . "123 abc ABC @#$" standard` renders multi-type characters without errors.
- [ ] **Custom Banner Case 62:** `go run . "ABCDE" doom` renders correctly with trailing separator format (all 8 lines visible).
- [ ] **Custom Banner Case 63:** `go run . "Hello World" greek` renders with proper alignment (12-char width normalization).
- [ ] **Custom Banner Case 64:** `go run . "Test" dancing` renders dancing-style characters correctly.
- [ ] **Custom Banner Case 65:** `go run . "Block" block` renders block-style characters correctly.
- [ ] **Fallback Case 66:** `go run . "hello" nonexistent` prints stderr warning and falls back to standard.txt.
- [ ] **Path Traversal Case 67:** `go run . "test" ../etc/passwd` is blocked and falls back to standard.txt.

### 6.6 Implementation Strategy for these Criteria
- **Flag Registry:** Implement a non-positional argument scanner that populates a `Config` struct (storing an array of colors and a single output string).
- **Writer Abstraction:** Use the `io.Writer` interface. If `--output` is present, the writer is an `*os.File`; otherwise, it defaults to `os.Stdout`.
- **ANSI Interleaving:** The Renderer must calculate color boundaries for each of the 8 lines of ASCII art to ensure ANSI codes are closed (`\033[0m`) at the end of each line within the file.

### 6.7 Testing & Validation
- **Visual Audit:** Use `cat -e <file>` to verify line endings and `cat <file>` to verify color rendering.
- **Binary Comparison:** Use `diff` to ensure that output saved to a file is identical to the output piped to a file (`go run . "str" > file`).
- **Automated Suite:** Run `go test ./unit-tests/...` to validate that the new multi-flag logic hasn't regressed the base ASCII-art character mapping.

---

### 6.8 Alignment Acceptance Criteria (Justify Project)
- [ ] **Align Case 1:** `go run . "something" standard` defaults to left alignment.
- [ ] **Align Case 2:** `go run . --align=right "left" standard` aligns to terminal right edge.
- [ ] **Align Case 3:** `go run . --align=left "right" standard` aligns to terminal left edge.
- [ ] **Align Case 4:** `go run . --align=center "hello" shadow` centers output with equal left/right padding.
- [ ] **Align Case 5:** `go run . --align=justify "1 Two 4" shadow` spreads words across terminal width.
- [ ] **Align Case 6:** `go run . --align=right "23/32" standard` right-aligns special characters.
- [ ] **Align Case 7:** `go run . --align=center "#$%&\"" thinkertoy` centers symbols.
- [ ] **Align Case 8:** `go run . --align=right "ABCabc123" thinkertoy` right-aligns mixed content.
- [ ] **Align Case 9:** `go run . --align=left "23Hello World!" standard` left-aligns mixed content.
- [ ] **Align Case 10:** Terminal resize recomputes width for right alignment.
- [ ] **Align Case 11:** Terminal resize recomputes width for center alignment.
- [ ] **Align Case 12:** `go run . --align=justify 'HELLO there HOW are YOU?!' thinkertoy` justifies long sentence.
- [ ] **Align Case 13:** `go run . --align=right "a -> A b -> B c -> C" shadow` right-aligns arrows.
- [ ] **Align Case 14:** `go run . --align right "something" standard` prints alignment usage message and exits non-zero.
- [ ] **Align Case 15:** Random mixed-case string with alignment renders correctly.
- [ ] **Align Case 16:** Random numbers/spaces string with alignment renders correctly.
- [ ] **Align Case 17:** Random special characters string with alignment renders correctly.
- [ ] **Align Case 18:** Random mixed all string with alignment renders correctly.

### 6.9 Combined Cases
- [ ] **Combined Case 1:** `go run . --output=out.txt --color=red "hello" standard` creates a file with red ASCII art.
- [ ] **Combined Case 2:** `go run . --output=aligned.txt --align=right "world" shadow` creates a file with right-aligned art.
- [ ] **Combined Case 3:** `go run . --color=blue --align=center "testing" thinkertoy` displays centered, blue art.
- [ ] **Combined Case 4:** `go run . --output=sub.txt --color=green "test" "testing" standard` creates a file with a green substring.
- [ ] **Combined Case 5:** `go run . --output=full.txt --color=yellow --align=center "full test" shadow` combines all three flags.
- [ ] **Combined Case 6:** `go run . --output=multi.txt --color=red "H" --color=blue "W" --align=justify "Hello World" standard` combines multiple colors, alignment, and output.
- [ ] **Combined Case 7:** `go run . --align=left --color=magenta --output=ordered.txt "order" thinkertoy` confirms flag order independence.

### 6.10 Usage Message & Exit Code Golden Cases
- [ ] **Usage Case 1:** `go run . --output test00.txt banana standard` prints output usage message and exits 1.
- [ ] **Usage Case 2:** `go run . --color red "banana"` prints color usage message and exits 1.
- [ ] **Usage Case 3:** `go run . --color=red` prints color usage message and exits 1.
- [ ] **Usage Case 4:** `go run . --align right "something" standard` prints alignment usage message and exits 1.

### 6.11 Alignment & Combined Warning-Tolerant Cases
- [ ] **Warning Case 1:** `go run . --align=center --align=right "hello" standard` warns about overridden align and right-aligns.
- [ ] **Warning Case 2:** `go run . --output=test.txt --align=middle "hello" standard` warns about invalid align and creates a file with default alignment.
- [ ] **Warning Case 3:** `go run . --color blue --align=center "hello" standard` warns about invalid color and center-aligns.
- [ ] **Warning Case 4:** `go run . --output=first.txt --output=second.txt "hello" standard` warns about duplicate output and writes to the second file.
- [ ] **Warning Case 5:** `go run . "hello" standard --align=right` correctly aligns with flag after arguments.
---

## 7. Implementation Approach (High Level)
This section outlines the architectural integrity, visual logic, and version control standards required to sustain the ascii-art-output project across its various optional iterations (**Color**, **Output**, and **Multi-flag** support).

### 7.1 Architecture Decision Summary
- **Decision**: Implement a **Modular Pipe-and-Filter** architecture using an `io.Writer` abstraction to unify *Terminal* and *File* output.
- **Stages**: CLI parse -> banner load -> input parse -> map -> color mask -> render -> align -> output
- **Why**:
  - **Decoupling:** The Renderer does not need to know if it is writing to a file or a screen; it simply fulfills the `io.Writer` interface.
  - **Persistence:** Using `os.O_TRUNC` ensures that every `--output` call provides a clean, predictable file state for auditors.
  - **ANSI Integrity:** By placing the Colorizer before the Renderer, we ensure that color metadata is attached to characters before they are converted into multi-line blocks, preventing "sliced" escape codes.
  - **Extensibility:** This flow accommodates alignment by adding a transformation stage before final output.

### 7.2 Flowchart

```text
                        [USER INPUT]
                              |
                              v
                  +----------------------------------+
                  |        CLI MODULE                |
                  | - `ParseCLI(osArgs)`             |
                  | - Scans all flags & args         |---- [Invalid Args] ----> [USAGE & EXIT 1]
                  | - Validates input                |
                  | - Builds unified `Config` object |
                  +----------------------------------+
                              |
                              v
                    [Valid `Comfig` Object]
                              |
            +---------------------------------------+
            |                                       |
            v                                       v
+--------------------------+          +--------------------------+
|      INPUT PARSER        |          |      BANNER LOADER       |
|  (uses `Config.Text`)    |          |  (uses `Config.Banner`)  |
+--------------------------+          +--------------------------+
            |                                       |
            +-------------------+-------------------+
                                |
                                v
                  +----------------------------+
                  |           MAPPER           |
                  |     (uses parsed text)     |
                  +----------------------------+
                                |
                                v
                  +----------------------------+                        |
                  |          COLORIZER         |
                  |   (uses `Config.Rules`)    |
                  +----------------------------+
                                |
                                v
                  +----------------------------+
                  |          RENDERER          |
                  |   (8-Line Builder + ANSI)  |
                  +----------------------------+
                                |
                                v
                  +----------------------------+
                  |         ALIGNMENT          |
                  | (Left/Right/Center/Justify)|
                  +----------------------------+
                                |
                                v
                  +----------------------------+
                  |      OUTPUT MODULE         |
                  |        (output.go)         |
                  +----------------------------+
                   /                          \
          [--output set]                [--output NOT set]
                 |                              |
                 v                              v
        +------------------+           +------------------+
        |  FILE SYSTEM     |           |     STDOUT       |
        | (os.O_TRUNC)     |           |    (Terminal)    |
        +------------------+           +------------------+
                 |                              |
                 +--------------+---------------+
                                |
                                v
                              EXIT 0
```

### 7.3 Development & Branching Strategy
- Keep `main` stable.
- Develop features in scoped branches, then merge after tests pass.
- Commits should map to milestones and acceptance criteria.

---

## 8. Milestones
This section outlines the phased development approach to integrate the Output redirection and Multi-Color functionality into the existing architecture.

### Milestone 1: Environment & CLI Foundation
- Parse `os.Args` to identify all instances of `--output=` and `--color=`.
- Implement non-positional flag extraction (allowing flags to appear in any order).
- Enforce exact usage messages and exit code 1 for malformed flags (e.g., spaces instead of `=`).
- Validate that the program still runs with a single `[STRING]` argument.

### Milestone 2: Input Parsing & Special Characters
- Implement `\n`, `\t`, and backslash collapse rules within the `parser.go` module.
- Implement deterministic handling for unsupported characters to ensure no crashes during the mapping phase.
- Maintain the "Peeking Gopher" logic for 100% non-ASCII input.

### Milestone 3: Banner Loader & Safety Fallbacks
- Enhance `loader.go` to verify banner file integrity (855 lines).
- Implement the "Safe Fallback" policy: default to `standard.txt` if a specified banner is missing or corrupt.
- Ensure the loader can handle standard, shadow, and thinkertoy dynamically.

### Milestone 4: Multi-Color Mapping & Masking
- Develop a Colorizer module capable of handling multiple `--color` flags.
- Implement substring matching logic to create a "Color Mask" for the input string.
- Support various color notations (ANSI, Hex, RGB, HSL) as defined in the technical specs.

### Milestone 5: File-Writer (`output.go`) & Redirection
- Implement the `output.go` module using `os.OpenFile` with `O_WRONLY|O_CREATE|O_TRUNC`.
- Abstract the output stream using the `io.Writer` interface to toggle between `os.Stdout` and a file pointer.
- Ensure the file is closed properly after the rendering process is complete.

### Milestone 6: High-Fidelity Rendering
- Update `renderer.go` to interleave ANSI escape codes with the 8-line ASCII glyphs.
- Ensure that the "Reset" sequence (`\033[0m`) is applied at the end of every colored segment to prevent color bleeding.
- Verify that horizontal alignment is preserved when multiple ANSI codes are injected into a single line.

### Milestone 7: Validation & Audit Readiness
- Execute all cases from `output_cases.md` (`test00.txt` through `test_multicolor.txt`) and golden_tests.md as well.
- Verify file persistence using `cat -e` to confirm identical terminal-to-file representation.
- Conduct final performance checks on long strings with multiple color targets.

### Milestone 8: FS Integration & Custom Banner Management
- Ensure all **.txt** files are in the `banner/` subdirectory.
- Refactor `loader.go` to accept a string variable for the file path instead of a constant.
- Implement file existence checks using `os.Stat` before attempting to read.
- Implement separator line detection logic (`isLeadingSeparator`) to handle both leading and trailing separator formats.
- Add support for custom banners: doom, greek, dancing, block.
- Implement 855-line validation for all banner files.
- Normalize greek.txt character widths to 12 characters for proper alignment.
- Finalize the argument counter to handle the `[STRING] [BANNER]` pattern strictly, rejecting 3+ arguments with FS usage message.
- Add path sanitization using `filepath.Base()` to prevent directory traversal attacks.
- Implement case-insensitive banner name handling with lowercase normalization.
- Support optional `.txt` extension in banner names.

---

### Milestone 9: Alignment Integration
- Add alignment parsing for `--align=<type>` with strict format validation.
- Implement terminal width detection using Go standard library only with a safe fallback for non-terminal environments.
- Implement left, right, center, and justify alignment while preserving 8-line height.
- Ensure `justify` distributes space between words and gracefully handles single-word input.
- Add alignment-specific tests and validate `docs/justify_cases.md`.

## 9. Risks / Open Questions
The following risks and open questions reflect the complexity of integrating **FS, Output** redirection, and **Multi-Color** features into a single Go CLI tool.

### 9.1 Risks
- **Flag grammar ambiguity**: Multi-flag parsing (`--output` + one or more `--color`) can become ambiguous if positional extraction is not deterministic.
- **Usage message drift**: Any mismatch with required usage texts (output and color variants) will fail audits even if runtime logic works.
- **File write safety**: Incorrect open flags may append or preserve stale content instead of strict overwrite (`O_TRUNC`) behavior.
- **File permission inconsistency**: Non-standard file modes can fail audit checks or reduce portability of generated files.
- **ANSI-in-file portability**: ANSI escape codes persisted to files may render differently across terminals/editors.
- **ANSI reset leaks**: Missing `\033[0m` boundaries can bleed colors across glyph segments, lines, or shell prompts.
- **ANSI Overhead in Files:** Large strings with multiple color flags significantly increase the file size of the `.txt` output due to the high density of escape codes; this could impact performance on extremely restricted environments.
- **Multi-color overlap conflicts**: Overlapping substring matches may produce nondeterministic results without a clear precedence policy.
- **Performance on large inputs**: Repeated substring scans with multi-color targeting can degrade runtime on long strings.
- **Path traversal and unsafe paths**: Output target and banner path handling can introduce security issues if not normalized/sanitized.
- **Banner fallback behavior drift**: Loader behavior for missing/corrupt banners may diverge from documented safe fallback policy.
- **Banner File Corruption:** If a banner file in the `banner/` directory is modified (e.g., a line is deleted), the 9-line offset logic will break, resulting in garbled ASCII art.
- **Inconsistent Character Widths:** Custom banners with inconsistent character widths (like the original greek.txt) cause misalignment when mixing different characters.
- **Separator Format Confusion:** Mixing leading and trailing separator formats without proper detection causes rendering errors (missing first or last line).
- **Terminal Width Mismatch:** While `output.go` saves the full width of the art, viewing that file with `cat` on a smaller terminal window will cause visual wrapping, potentially making the audit results look "broken."
- **Peeking Gopher Regression:** Introducing complex flag parsing must not interfere with the specific requirement to trigger the "Peeking Gopher" when 100% non-ASCII input is detected.
- **Cross-platform newline differences**: CRLF/LF handling can affect file diff checks and `cat -e`-based validation. Handling `\n` across Windows (CRLF) and Unix (LF) when writing to files might cause diff mismatches during automated testing.

### 9.2 Open Questions
- **Multi-color precedence rule**: Resolved as **last matching rule wins** for overlapping matches.
- **Duplicate `--output` handling**: Keep "last output flag wins" as mandatory behavior, or reject duplicates as invalid usage?
- **Flag order contract**: Must the parser remain fully order-independent for flags, or should one canonical flag order be enforced?
- **Color persistence policy**: Should ANSI codes always be written to output files, or should there be an optional plain-text output mode?
- **Error-stream policy**: Should all usage/validation errors be strictly `stderr`, including combined-flag failures?
- **Fallback transparency**: Should banner fallback remain silent, or emit a non-fatal warning to `stderr`?
- **Hex notation status**: Is `#RRGGBB` still bonus, or now required in baseline acceptance?
- **Output on empty input**: For `--output` with empty string, should file be 0 bytes or contain a single newline for deterministic audits?
- **Escaping in output filename**: Should filenames with spaces/special characters be supported exactly as quoted input or normalized/restricted?
- **Empty Output File:** If the user runs `--output=empty.txt ""` (empty string), should the file be created as 0 bytes or contain a single newline for deterministic `cat -e` behavior?
- **Character 127 Handling:** How should the program handle the DEL character or non-printable ASCII if they are technically within the filesystem range but lack a glyph?
