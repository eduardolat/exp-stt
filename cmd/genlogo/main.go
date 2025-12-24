package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
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
	if err := generateSVGS(); err != nil {
		panic(err)
	}

	if err := generatePNGS(); err != nil {
		panic(err)
	}

	if err := generateICOS(); err != nil {
		panic(err)
	}
}

func generateSVGS() error {
	outputDir := "./assets/logo/svg"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
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
					fmt.Printf("Generated SVG: %s\n", fileName)
				}
			}
		}
	}
	return nil
}

func generatePNGS() error {
	svgDir := "./assets/logo/svg"
	pngDir := "./assets/logo/png"
	if err := os.MkdirAll(pngDir, 0755); err != nil {
		return err
	}

	files, err := os.ReadDir(svgDir)
	if err != nil {
		return err
	}

	sizes := []int{16, 32, 48, 64, 128, 256, 512}
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(5)

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".svg") {
			continue
		}

		svgFile := file.Name()
		baseName := strings.TrimSuffix(svgFile, ".svg")

		for _, size := range sizes {
			g.Go(func() error {
				outputName := fmt.Sprintf("%s_%d.png", baseName, size)
				outputPath := filepath.Join(pngDir, outputName)

				cmd := exec.Command("rsvg-convert", "-w", fmt.Sprintf("%d", size), "-h", fmt.Sprintf("%d", size), filepath.Join(svgDir, svgFile), "-o", outputPath)
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("error generating png %s: %v", outputPath, err)
				}
				fmt.Printf("Generated PNG: %s\n", outputName)
				return nil
			})
		}
	}

	return g.Wait()
}

func generateICOS() error {
	svgDir := "./assets/logo/svg"
	pngDir := "./assets/logo/png"
	icoDir := "./assets/logo/ico"
	if err := os.MkdirAll(icoDir, 0755); err != nil {
		return err
	}

	files, err := os.ReadDir(svgDir)
	if err != nil {
		return err
	}

	sizes := []int{16, 32, 48, 64, 128, 256, 512}
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(5)

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".svg") {
			continue
		}

		baseName := strings.TrimSuffix(file.Name(), ".svg")
		g.Go(func() error {
			icoPath := filepath.Join(icoDir, baseName+".ico")

			args := []string{}
			for _, size := range sizes {
				args = append(args, filepath.Join(pngDir, fmt.Sprintf("%s_%d.png", baseName, size)))
			}
			args = append(args, icoPath)
			cmd := exec.Command("magick", args...)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error generating ico %s: %v", icoPath, err)
			}
			fmt.Printf("Generated ICO: %s.ico\n", baseName)
			return nil
		})
	}

	return g.Wait()
}
