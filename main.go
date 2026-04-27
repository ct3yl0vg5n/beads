// Package main is the entry point for beads, a distributed message bead chain system.
// Forked from gastownhall/beads with additional features for conflict resolution and marketplace integration.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gastownhall/beads/internal/chain"
	"github.com/gastownhall/beads/internal/config"
	"github.com/gastownhall/beads/internal/server"
)

const (
	defaultPort    = 8080
	defaultHost    = "0.0.0.0"
	defaultLogLevel = "info"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var (
		host      = flag.String("host", defaultHost, "Host address to bind the server")
		port      = flag.Int("port", defaultPort, "Port to listen on")
		cfgPath   = flag.String("config", "", "Path to configuration file")
		logLevel  = flag.String("log-level", defaultLogLevel, "Log level: debug, info, warn, error")
		showVer   = flag.Bool("version", false, "Print version information and exit")
		initChain = flag.Bool("init", false, "Initialize a new bead chain and exit")
	)
	flag.Parse()

	if *showVer {
		fmt.Printf("beads %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	// Load configuration from file or defaults.
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Apply flag overrides to config.
	if *host != defaultHost {
		cfg.Host = *host
	}
	if *port != defaultPort {
		cfg.Port = *port
	}
	if *logLevel != defaultLogLevel {
		cfg.LogLevel = *logLevel
	}

	logger := log.New(os.Stdout, "[beads] ", log.LstdFlags|log.Lshortfile)

	// Initialize the bead chain storage.
	bc, err := chain.New(cfg.DataDir, logger)
	if err != nil {
		logger.Fatalf("failed to initialize bead chain: %v", err)
	}
	defer bc.Close()

	if *initChain {
		if err := bc.Init(); err != nil {
			logger.Fatalf("failed to initialize chain: %v", err)
		}
		logger.Println("bead chain initialized successfully")
		os.Exit(0)
	}

	// Start the HTTP server.
	srv := server.New(cfg, bc, logger)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Printf("starting beads server on %s (version %s)", addr, version)

	if err := srv.ListenAndServe(addr); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}
