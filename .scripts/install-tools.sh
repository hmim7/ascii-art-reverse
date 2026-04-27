#!/usr/bin/env bash

echo "🛠️ Εγκατάσταση όλων των απαιτούμενων εργαλείων (Linters, Formatters, Security)..."

set -e # Σταματάει την εκτέλεση αν κάποια συγκεκριμένη εγκατάσταση αποτύχει

go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/mgechev/revive@latest
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install github.com/client9/misspell/cmd/misspell@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
go install github.com/kisielk/errcheck@latest
go install mvdan.cc/gofumpt@latest
go install go.uber.org/nilaway/cmd/nilaway@latest

echo "✅ Όλα τα εργαλεία εγκαταστάθηκαν επιτυχώς!"
echo "⚠️  Σημείωση: Για τον έλεγχο 'go test -race' απαιτείται να έχεις εγκατεστημένο το GCC (π.χ. MinGW-w64 στα Windows)."