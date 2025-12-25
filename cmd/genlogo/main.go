// This program generates the logo assets for the STT application.
// It creates SVG files from a template, converts them to PNG in multiple sizes
// using rsvg-convert, and bundles them into ICO files using ImageMagick.
// Finally, it generates a Go file with embed directives to expose all assets
// as structured variables.
package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

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
	ColorGray  = Color{Name: "gray", Hex: "#808080"}  // Model unmounted
	ColorWhite = Color{Name: "white", Hex: "#ffffff"} // Model mounted
	ColorAmber = Color{Name: "amber", Hex: "#ffd700"} // Model mounting
	ColorBlack = Color{Name: "black", Hex: "#000000"} // Background
	ColorPink  = Color{Name: "pink", Hex: "#dd00c7"}  // Recording
	ColorBlue  = Color{Name: "blue", Hex: "#00a6ff"}  // Transcribing
	ColorGreen = Color{Name: "green", Hex: "#00ff6a"} // Post-processing
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

var (
	BlackBarColors = []Color{ColorWhite, ColorGray, ColorAmber, ColorPink, ColorBlue, ColorGreen}
	WhiteBarColors = []Color{ColorBlack, ColorGray, ColorAmber, ColorPink, ColorBlue, ColorGreen}
)

func main() {
	if err := os.RemoveAll("./assets/logo"); err != nil {
		panic(err)
	}

	var total int

	qty, err := generateSVGS()
	if err != nil {
		panic(err)
	}
	total += qty

	qty, err = generatePNGS()
	if err != nil {
		panic(err)
	}
	total += qty

	qty, err = generateICOS()
	if err != nil {
		panic(err)
	}
	total += qty

	qty, err = generateEmbedGo()
	if err != nil {
		panic(err)
	}
	total += qty

	fmt.Printf("Total generated files: %d\n", total)
}

func generateSVGS() (int, error) {
	outputDir := "./assets/logo/svg"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return 0, err
	}

	bgColors := []Color{ColorBlack, ColorWhite}
	var count int

	for _, bg := range bgColors {
		barColors := WhiteBarColors
		if bg.Name == "black" {
			barColors = BlackBarColors
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
					count++
				}
			}
		}
	}
	fmt.Printf("SVG files generated: %d\n", count)
	return count, nil
}

func generatePNGS() (int, error) {
	svgDir := "./assets/logo/svg"
	pngDir := "./assets/logo/png"
	if err := os.MkdirAll(pngDir, 0755); err != nil {
		return 0, err
	}

	files, err := os.ReadDir(svgDir)
	if err != nil {
		return 0, err
	}

	sizes := []int{16, 32, 48, 64, 128, 256, 512}
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(5)

	var count int
	var mu sync.Mutex

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".svg") {
			continue
		}

		svgFile := file.Name()
		baseName := strings.TrimSuffix(svgFile, ".svg")

		for _, size := range sizes {
			g.Go(func() error {
				outputName := fmt.Sprintf("%s-%d.png", baseName, size)
				outputPath := filepath.Join(pngDir, outputName)

				cmd := exec.Command("rsvg-convert", "-w", fmt.Sprintf("%d", size), "-h", fmt.Sprintf("%d", size), filepath.Join(svgDir, svgFile), "-o", outputPath)
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("error generating png %s: %v", outputPath, err)
				}
				mu.Lock()
				count++
				mu.Unlock()
				return nil
			})
		}
	}

	err = g.Wait()
	fmt.Printf("PNG files generated: %d\n", count)
	return count, err
}

func generateICOS() (int, error) {
	svgDir := "./assets/logo/svg"
	pngDir := "./assets/logo/png"
	icoDir := "./assets/logo/ico"
	if err := os.MkdirAll(icoDir, 0755); err != nil {
		return 0, err
	}

	files, err := os.ReadDir(svgDir)
	if err != nil {
		return 0, err
	}

	sizes := []int{16, 32, 48, 64, 128, 256, 512}
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(5)

	var count int
	var mu sync.Mutex

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".svg") {
			continue
		}

		baseName := strings.TrimSuffix(file.Name(), ".svg")
		g.Go(func() error {
			icoPath := filepath.Join(icoDir, baseName+".ico")

			args := []string{}
			for _, size := range sizes {
				args = append(args, filepath.Join(pngDir, fmt.Sprintf("%s-%d.png", baseName, size)))
			}
			args = append(args, icoPath)
			cmd := exec.Command("magick", args...)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error generating ico %s: %v", icoPath, err)
			}
			mu.Lock()
			count++
			mu.Unlock()
			return nil
		})
	}

	err = g.Wait()
	fmt.Printf("ICO files generated: %d\n", count)
	return count, err
}

func generateEmbedGo() (int, error) {
	var buf bytes.Buffer

	buf.WriteString(`
		// Code generated by cmd/genlogo/main.go
		// DO NOT EDIT

		package logo

		import _ "embed"
	`)

	type Resource struct {
		VarName string
		Path    string
	}

	var svgResources []Resource
	var pngResources []Resource
	var icoResources []Resource

	bgColors := []Color{ColorBlack, ColorWhite}
	sizes := []int{16, 32, 48, 64, 128, 256, 512}

	capitalize := func(s string) string {
		if len(s) == 0 {
			return ""
		}
		return strings.ToUpper(s[:1]) + s[1:]
	}

	for _, bg := range bgColors {
		barColors := WhiteBarColors
		if bg.Name == "black" {
			barColors = BlackBarColors
		}

		for _, bar := range barColors {
			for _, variant := range BarVariants {
				baseName := fmt.Sprintf("%s-%s-%s", bg.Name, bar.Name, variant.Name)
				camelBase := bg.Name + capitalize(bar.Name) + capitalize(variant.Name)

				// SVG
				svgResources = append(svgResources, Resource{
					VarName: camelBase + "SVG",
					Path:    "svg/" + baseName + ".svg",
				})

				// PNG
				for _, size := range sizes {
					pngResources = append(pngResources, Resource{
						VarName: fmt.Sprintf("%sPNG%d", camelBase, size),
						Path:    fmt.Sprintf("png/%s-%d.png", baseName, size),
					})
				}

				// ICO
				icoResources = append(icoResources, Resource{
					VarName: camelBase + "ICO",
					Path:    "ico/" + baseName + ".ico",
				})
			}
		}
	}

	buf.WriteString("var (\n")
	for _, r := range svgResources {
		fmt.Fprintf(&buf, "\t//go:embed %s\n\t%s []byte\n", r.Path, r.VarName)
	}
	for _, r := range pngResources {
		fmt.Fprintf(&buf, "\t//go:embed %s\n\t%s []byte\n", r.Path, r.VarName)
	}
	for _, r := range icoResources {
		fmt.Fprintf(&buf, "\t//go:embed %s\n\t%s []byte\n", r.Path, r.VarName)
	}
	buf.WriteString(")\n\n")

	// Types
	buf.WriteString(`
		type ResourceSet struct {
			Left   []byte
			Middle []byte
			Right  []byte
			Logo   []byte
		}

		type PNGResources struct {
			Size16  ResourceSet
			Size32  ResourceSet
			Size48  ResourceSet
			Size64  ResourceSet
			Size128 ResourceSet
			Size256 ResourceSet
			Size512 ResourceSet
		}

		type LogoResources struct {
			SVG ResourceSet
			PNG PNGResources
			ICO ResourceSet
		}
	`)

	// Exported variables
	buf.WriteString("var (\n")
	for _, bg := range bgColors {
		barColors := WhiteBarColors
		if bg.Name == "black" {
			barColors = BlackBarColors
		}

		for _, bar := range barColors {
			varName := "Logo" + capitalize(bg.Name) + capitalize(bar.Name)
			fmt.Fprintf(&buf, "\t%s = LogoResources{\n", varName)

			// SVG
			buf.WriteString("\t\tSVG: ResourceSet{\n")
			for _, v := range BarVariants {
				camelBase := bg.Name + capitalize(bar.Name) + capitalize(v.Name)
				fmt.Fprintf(&buf, "\t\t\t%s: %sSVG,\n", capitalize(v.Name), camelBase)
			}
			buf.WriteString("\t\t},\n")

			// PNG
			buf.WriteString("\t\tPNG: PNGResources{\n")
			for _, size := range sizes {
				fmt.Fprintf(&buf, "\t\t\tSize%d: ResourceSet{\n", size)
				for _, v := range BarVariants {
					camelBase := bg.Name + capitalize(bar.Name) + capitalize(v.Name)
					fmt.Fprintf(&buf, "\t\t\t\t%s: %sPNG%d,\n", capitalize(v.Name), camelBase, size)
				}
				buf.WriteString("\t\t\t},\n")
			}
			buf.WriteString("\t\t},\n")

			// ICO
			buf.WriteString("\t\tICO: ResourceSet{\n")
			for _, v := range BarVariants {
				camelBase := bg.Name + capitalize(bar.Name) + capitalize(v.Name)
				fmt.Fprintf(&buf, "\t\t\t%s: %sICO,\n", capitalize(v.Name), camelBase)
			}
			buf.WriteString("\t\t},\n")

			buf.WriteString("\t}\n")
		}
	}
	buf.WriteString(")\n")

	// Format code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return 0, fmt.Errorf("error formatting generated code: %w", err)
	}

	err = os.WriteFile("./assets/logo/embed.go", formatted, 0644)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Go embed file generated: 1\n")
	return 1, nil
}
