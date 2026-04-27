package main

import (
	"fmt"
	"os"

	"ascii-art-reverse/internal/banner"
	"ascii-art-reverse/internal/cli"
	"ascii-art-reverse/internal/output"
	"ascii-art-reverse/internal/render"
	"ascii-art-reverse/internal/reverse"
)

func main() {
	args := cli.ClassifyArgs(os.Args[1:])

	if len(args.UnknownFlags) > 0 ||
		(len(args.Malformed) > 0 && args.OutputValue == "" &&
			args.AlignValue == "" && args.ReverseValue == "" && len(args.ColorRules) == 0) {
		cli.Fatal(cli.SelectUsage(args))
	}
	cli.EmitWarnings(args)

	bannerName := "standard"
	if len(args.Positional) == 2 {
		bannerName = args.Positional[1]
	}

	bannerMap, err := banner.Load(bannerName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if args.ReverseValue != "" {
		result, err := reverse.Run(args.ReverseValue, bannerMap)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(result)
		return
	}

	writer, file, err := output.GetWriter(args.OutputValue)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if file != nil {
		defer file.Close()
	}

	input := ""
	if len(args.Positional) > 0 {
		input = args.Positional[0]
	}
	if args.StdinMode {
		input = cli.ReadStdin()
	}

	segments := render.ParseInput(input)
	if render.ShouldRenderGopher(input, segments) {
		render.RenderGopher(writer)
		return
	}

	rules := cli.BuildColorRules(args.ColorRules)
	width := render.ResolveWidth(args.OutputValue != "")
	render.RenderAlignedWithColorRules(writer, segments, bannerMap, render.ToColorRules(rules), args.AlignValue, width)
}
