package main

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
	LeftSize   string
	MiddleSize string
	RightSize  string
}

var (
	BarVariantLogo   = BarVariant{LeftSize: BAR_SIZE_SHORT, MiddleSize: BAR_SIZE_TALL, RightSize: BAR_SIZE_SHORT}
	BarVariantLeft   = BarVariant{LeftSize: BAR_SIZE_TALL, MiddleSize: BAR_SIZE_SHORT, RightSize: BAR_SIZE_SHORT}
	BarVariantMiddle = BarVariant{LeftSize: BAR_SIZE_SHORT, MiddleSize: BAR_SIZE_TALL, RightSize: BAR_SIZE_SHORT}
	BarVariantRight  = BarVariant{LeftSize: BAR_SIZE_SHORT, MiddleSize: BAR_SIZE_SHORT, RightSize: BAR_SIZE_TALL}
)

type LogoConfig struct {
	BgColor   string
	BarColor  string
	LeftBar   string
	MiddleBar string
	RightBar  string
}

func main() {

}
