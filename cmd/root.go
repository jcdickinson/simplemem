package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/jcdickinson/simplemem/internal/mcp"
)

var (
	dbPath string
)

var rootCmd = &cobra.Command{
	Use:   "simplemem",
	Short: "SimpleMem MCP Server - Memory management with RAG capabilities",
	Long: `SimpleMem is a Model Context Protocol (MCP) server that provides intelligent 
memory management with Retrieval-Augmented Generation (RAG) capabilities. 
It combines traditional file-based storage with semantic search powered by vector embeddings.`,
	Run: runServer,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", ".cache/simplemem.db", "Path to the database file")
}

func runServer(cmd *cobra.Command, args []string) {
	// Create the MCP server with custom database path
	server, err := mcp.NewServer(dbPath)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Set up error channel for server errors
	errCh := make(chan error)
	go func() {
		errCh <- server.Run()
	}()

	// Set up signal handling
	if err = signalWaiter(errCh); err != nil {
		log.Fatalf("Server error: %v", err)
		return
	}

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
}

func signalWaiter(errCh chan error) error {
	signalToNotify := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}
	if signal.Ignored(syscall.SIGHUP) {
		signalToNotify = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, signalToNotify...)

	select {
	case sig := <-signals:
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			log.Printf("Received signal: %s\n", sig)
			// graceful shutdown
			return nil
		}
	case err := <-errCh:
		return err
	}

	return nil
}