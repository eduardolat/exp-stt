package config

import (
	"flag"
	"fmt"
	"os"
)

// Flags holds the command-line configuration for the application.
type Flags struct {
	Debug        bool
	Host         string
	Port         int
	PortExplicit bool
	Toggle       bool
}

// DefaultHost is the default network interface to bind the HTTP server to.
const DefaultHost = "127.0.0.1"

// DefaultPort is the default port to bind the HTTP server to.
const DefaultPort = 8742

// PortRangeSize is the number of ports to try when auto-discovering.
const PortRangeSize = 100

// ParseFlags parses command-line arguments and returns the application flags.
// Exits with status 0 if --help is requested, or status 2 on parse error.
func ParseFlags() Flags {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Tribar Voice - Speech to text application")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		fs.PrintDefaults()
	}

	debug := fs.Bool("debug", false, "Enable debug mode with verbose logging")
	host := fs.String("host", DefaultHost, "Network interface to bind the HTTP server to")
	port := fs.Int("port", 0, "Port number to bind the HTTP server to (auto-discovers if not set)")
	toggle := fs.Bool("toggle", false, "Toggle recording on a running instance and exit")

	_ = fs.Parse(os.Args[1:])

	portExplicit := *port != 0
	finalPort := *port
	if !portExplicit {
		finalPort = DefaultPort
	}

	return Flags{
		Debug:        *debug,
		Host:         *host,
		Port:         finalPort,
		PortExplicit: portExplicit,
		Toggle:       *toggle,
	}
}

// Address returns the full address string in host:port format.
func (f Flags) Address() string {
	return fmt.Sprintf("%s:%d", f.Host, f.Port)
}
