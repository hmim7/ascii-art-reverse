# Automation & Quality Control Scripts

This directory contains the automation suite designed to ensure the `ascii-art-reverse` project maintains high code quality, security, and performance standards. These scripts are aligned with the project's **Good Practices Checklists** and **PRD** requirements.

## 🚀 Quick Start

To set up your local development environment:

> **Note:** Run these commands from within the `.scripts/` directory.

1. **Install Required Tools:**
   ```bash
   chmod +x install-tools.sh && ./install-tools.sh
   # OR
   make install-tools
   ```
   *Note: Ensure your Go bin path (usually `~/go/bin`) is added to your system's PATH.*

2. **Enable Automated Verification (Git Hooks):**
   ```bash
   make install-hooks
   ```
   This installs a `pre-commit` hook that runs the full suite automatically before every commit and a `commit-msg` hook to enforce Conventional Commits.

## 🛠 Available Scripts

### `check.sh`
**Usage:** 
```bash
chmod +x check.sh
./check.sh
```
The "Gold Standard" verification script. It performs an exhaustive audit of the codebase:
- **Standardization:** Runs `goimports`, `gofmt`, and `gofumpt` to ensure idiomatic formatting.
- **Dependency Management:** Runs `go mod tidy` and verifies that **only** the Go Standard Library is used (per project constraints).
- **Static Analysis:** Executes `go vet`, `shadow`, `staticcheck`, `errcheck`, `nilaway`, and `revive`.
- **Complexity Audit:** Ensures no function exceeds a cyclomatic complexity of **15** using `gocyclo`.
- **Security Audit:** Scans for vulnerabilities using `gosec` and `govulncheck`.
- **Performance & Memory:** Checks for optimal struct field alignment and runs memory-aware benchmarks.
- **Testing:** Runs unit tests with the race detector (`-race`), fuzz tests (10s duration), and verifies code coverage (Threshold: **80%**).
- **Time Limits:** Validates that the program completes execution within the **5-minute limit** for samples.

### `install-tools.sh`
**Usage:** 
```bash
chmod +x install-tools.sh
./install-tools.sh
```
A helper script to install all Go-based linters, formatters, and security analyzers required by the pipeline.

### `Makefile`
**Usage:** 
```bash
make <target>
```

Available targets:
- `make check`: Runs the primary check pipeline (fmt, tidy, lint, audit, test).
- `make fmt`: Standardizes code formatting across the project.
- `make test`: Runs all tests, benchmarks, and generates a coverage report.
- `make install-hooks`: Sets up the Git integration.

## ⚠️ Requirements & Troubleshooting

- **GCC Requirement:** The `go test -race` check requires a C compiler. On Windows, ensure you have **MinGW-w64** installed.
- **Go Version:** The pipeline is tuned for **Go 1.25.6+**.
- **Samples Folder:** The time-limit check in `check.sh` looks for a `samples/` directory at the project root to run test cases.

## 🛑 Handling Failures

If the `check.sh` script (or the pre-commit hook) fails:
1. **Read the Output:** The script is designed to be helpful. It provides a **"👉 Handling Instructions"** section for every failure, explaining why it failed and how to fix it.
2. **Auto-Fixing:** Many formatting and field alignment issues can be fixed automatically using the `-fix` or `-w` flags mentioned in the failure logs.
3. **Manual Override:** While highly discouraged, in case of an installed pre-commit hook, you can skip the automated check during a commit using `git commit --no-verify`.

---
*Maintain the standard. Don't push broken code.*