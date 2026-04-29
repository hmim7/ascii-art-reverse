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
	cli.EmitWarnings(args)
	cli.ValidateOrFatal(args)
	cli.CheckPositionalsOrFatal(args)

	bannerName := cli.ResolveBannerName(args)
	bannerMap, fallback, err := banner.Load(bannerName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cli.HandleBannerFallbackOrFatal(fallback, bannerName)

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

	input := cli.ResolveInput(args)
	segments := render.ParseInput(input)
	if render.ShouldRenderGopher(input, segments) {
		render.RenderGopher(writer)
		return
	}

	rules := cli.BuildColorRules(args.ColorRules)
	width := render.ResolveWidth(args.OutputValue != "")
	render.RenderAlignedWithColorRules(writer, segments, bannerMap, render.ToColorRules(rules), args.AlignValue, width)
}
