package config

var (
	// AppVersion is set at build time via -ldflags
	// Example: go build -ldflags="-X github.com/varavelio/tribar/internal/config.AppVersion=1.0.0"
	AppVersion = "0.0.0-dev"

	AppName   = "Tribar Voice"
	RepoURL   = "https://github.com/varavelio/tribar"
	UserAgent = "TribarVoice/" + AppVersion
)
