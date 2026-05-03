package main

import (
	"fmt"
	"strings"

	"ascii-art-reverse/internal/cli"
)

func main() {
	args := cli.ClassifyArgs([]string{"--ouput=out.txt", "hello"})
	fmt.Printf("UnknownFlags: %+v\n", args.UnknownFlags)
	for _, f := range args.UnknownFlags {
		raw := strings.ToLower(f.Raw)
		fmt.Printf("  raw=%q HasPrefix(--out)=%v\n", raw, strings.HasPrefix(raw, "--out"))
	}
	fmt.Printf("SelectUsage: %q\n", cli.SelectUsage(args))
}
