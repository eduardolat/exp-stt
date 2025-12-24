package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const SVG_TEMPLATE = `
<svg width="512" height="512" viewBox="0 0 512 512" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="512" height="512" rx="112.64" fill="{bg_color}"/>
  <g stroke="{bar_color}" stroke-width="60" stroke-linecap="round">
    <path d="M136 {left_bar}"/>
    <path d="M256 {middle_bar}"/>
    <path d="M376 {right_bar}"/>
  </g>
</svg>
`

const (
	BAR_SIZE_TALL  = "106v300"
	BAR_SIZE_SHORT = "166v180"
)

type Color struct {
	Name string
	Hex  string
}

var (
	ColorWhite = Color{Name: "white", Hex: "#ffffff"}
	ColorBlack = Color{Name: "black", Hex: "#000000"}
	ColorBlue  = Color{Name: "blue", Hex: "#00a6ff"}
	ColorPink  = Color{Name: "pink", Hex: "#dd00c7"}
)

type BarVariant struct {
	Name       string
	LeftSize   string
	MiddleSize string
	RightSize  string
}

var BarVariants = []BarVariant{
	{Name: "logo", LeftSize: BAR_SIZE_SHORT, MiddleSize: BAR_SIZE_TALL, RightSize: BAR_SIZE_SHORT},
	{Name: "left", LeftSize: BAR_SIZE_TALL, MiddleSize: BAR_SIZE_SHORT, RightSize: BAR_SIZE_SHORT},
	{Name: "middle", LeftSize: BAR_SIZE_SHORT, MiddleSize: BAR_SIZE_TALL, RightSize: BAR_SIZE_SHORT},
	{Name: "right", LeftSize: BAR_SIZE_SHORT, MiddleSize: BAR_SIZE_SHORT, RightSize: BAR_SIZE_TALL},
}

func main() {
	outputDir := "./assets/logo"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}

	bgColors := []Color{ColorBlack, ColorWhite}

	for _, bg := range bgColors {
		var barColors []Color
		if bg.Name == "black" {
			barColors = []Color{ColorWhite, ColorPink, ColorBlue}
		} else {
			barColors = []Color{ColorBlack, ColorPink, ColorBlue}
		}

		for _, bar := range barColors {
			for _, variant := range BarVariants {
				fileName := fmt.Sprintf("%s-%s-%s.svg", bg.Name, bar.Name, variant.Name)
				filePath := filepath.Join(outputDir, fileName)

				svg := SVG_TEMPLATE
				svg = strings.ReplaceAll(svg, "{bg_color}", bg.Hex)
				svg = strings.ReplaceAll(svg, "{bar_color}", bar.Hex)
				svg = strings.ReplaceAll(svg, "{left_bar}", variant.LeftSize)
				svg = strings.ReplaceAll(svg, "{middle_bar}", variant.MiddleSize)
				svg = strings.ReplaceAll(svg, "{right_bar}", variant.RightSize)

				err := os.WriteFile(filePath, []byte(strings.TrimSpace(svg)), 0644)
				if err != nil {
					fmt.Printf("Error writing %s: %v\n", fileName, err)
				} else {
					fmt.Printf("Generated %s\n", fileName)
				}
			}
		}
	}
}
