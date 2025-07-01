package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/TheEntropyCollective/randomfs-core"
	"github.com/spf13/cobra"
)

var (
	ipfsAPI   string
	dataDir   string
	cacheSize int64
	verbose   bool
)

var rootCmd = &cobra.Command{
	Use:   "randomfs-cli",
	Short: "RandomFS CLI - Owner Free File System command line interface",
	Long: `RandomFS CLI provides command line access to the Owner Free File System.
Store and retrieve files using randomized blocks on IPFS with rd:// URLs.`,
}

var storeCmd = &cobra.Command{
	Use:   "store [file]",
	Short: "Store a file in RandomFS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		// Read file
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		// Detect content type
		contentType := detectContentType(filename)

		// Create RandomFS instance
		rfs, err := randomfs.NewRandomFS(ipfsAPI, dataDir, cacheSize)
		if err != nil {
			log.Fatalf("Failed to initialize RandomFS: %v", err)
		}

		// Store file
		randomURL, err := rfs.StoreFile(filename, data, contentType)
		if err != nil {
			log.Fatalf("Failed to store file: %v", err)
		}

		fmt.Printf("File stored successfully!\n")
		fmt.Printf("rd:// URL: %s\n", randomURL.String())
		fmt.Printf("Rep Hash:  %s\n", randomURL.RepHash)
		fmt.Printf("File Size: %d bytes\n", randomURL.FileSize)

		if verbose {
			stats := rfs.GetStats()
			fmt.Printf("\nSystem Stats:\n")
			fmt.Printf("  Files Stored: %d\n", stats.FilesStored)
			fmt.Printf("  Blocks Generated: %d\n", stats.BlocksGenerated)
			fmt.Printf("  Total Size: %d bytes\n", stats.TotalSize)
		}
	},
}

var retrieveCmd = &cobra.Command{
	Use:   "retrieve [hash] [output-file]",
	Short: "Retrieve a file from RandomFS by representation hash",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		repHash := args[0]
		outputFile := args[1]

		// Create RandomFS instance
		rfs, err := randomfs.NewRandomFS(ipfsAPI, dataDir, cacheSize)
		if err != nil {
			log.Fatalf("Failed to initialize RandomFS: %v", err)
		}

		// Retrieve file
		data, rep, err := rfs.RetrieveFile(repHash)
		if err != nil {
			log.Fatalf("Failed to retrieve file: %v", err)
		}

		// Write to output file
		err = os.WriteFile(outputFile, data, 0644)
		if err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}

		fmt.Printf("File retrieved successfully!\n")
		fmt.Printf("Original Name: %s\n", rep.FileName)
		fmt.Printf("Content Type:  %s\n", rep.ContentType)
		fmt.Printf("File Size:     %d bytes\n", rep.FileSize)
		fmt.Printf("Block Count:   %d\n", len(rep.BlockHashes))
		fmt.Printf("Output File:   %s\n", outputFile)

		if verbose {
			stats := rfs.GetStats()
			fmt.Printf("\nSystem Stats:\n")
			fmt.Printf("  Cache Hits: %d\n", stats.CacheHits)
			fmt.Printf("  Cache Misses: %d\n", stats.CacheMisses)
		}
	},
}

var parseCmd = &cobra.Command{
	Use:   "parse [rd-url]",
	Short: "Parse a rd:// URL and show its components",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rdURL := args[0]

		randomURL, err := randomfs.ParseRandomURL(rdURL)
		if err != nil {
			log.Fatalf("Failed to parse rd:// URL: %v", err)
		}

		fmt.Printf("Parsed rd:// URL:\n")
		fmt.Printf("  Scheme:    %s\n", randomURL.Scheme)
		fmt.Printf("  Host:      %s\n", randomURL.Host)
		fmt.Printf("  Version:   %s\n", randomURL.Version)
		fmt.Printf("  File Name: %s\n", randomURL.FileName)
		fmt.Printf("  File Size: %d bytes\n", randomURL.FileSize)
		fmt.Printf("  Rep Hash:  %s\n", randomURL.RepHash)
		fmt.Printf("  Timestamp: %d\n", randomURL.Timestamp)
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download [rd-url] [output-file]",
	Short: "Download a file using its rd:// URL",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		rdURL := args[0]
		outputFile := args[1]

		// Parse rd:// URL
		randomURL, err := randomfs.ParseRandomURL(rdURL)
		if err != nil {
			log.Fatalf("Failed to parse rd:// URL: %v", err)
		}

		// Create RandomFS instance
		rfs, err := randomfs.NewRandomFS(ipfsAPI, dataDir, cacheSize)
		if err != nil {
			log.Fatalf("Failed to initialize RandomFS: %v", err)
		}

		// Retrieve file using rep hash
		data, rep, err := rfs.RetrieveFile(randomURL.RepHash)
		if err != nil {
			log.Fatalf("Failed to download file: %v", err)
		}

		// Write to output file
		err = os.WriteFile(outputFile, data, 0644)
		if err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}

		fmt.Printf("File downloaded successfully!\n")
		fmt.Printf("Original Name: %s\n", rep.FileName)
		fmt.Printf("File Size:     %d bytes\n", rep.FileSize)
		fmt.Printf("Output File:   %s\n", outputFile)
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show RandomFS system statistics",
	Run: func(cmd *cobra.Command, args []string) {
		// Create RandomFS instance
		rfs, err := randomfs.NewRandomFS(ipfsAPI, dataDir, cacheSize)
		if err != nil {
			log.Fatalf("Failed to initialize RandomFS: %v", err)
		}

		stats := rfs.GetStats()

		if verbose {
			jsonData, _ := json.MarshalIndent(stats, "", "  ")
			fmt.Printf("RandomFS Statistics:\n%s\n", jsonData)
		} else {
			fmt.Printf("RandomFS Statistics:\n")
			fmt.Printf("  Files Stored:     %d\n", stats.FilesStored)
			fmt.Printf("  Blocks Generated: %d\n", stats.BlocksGenerated)
			fmt.Printf("  Total Size:       %d bytes\n", stats.TotalSize)
			fmt.Printf("  Cache Hits:       %d\n", stats.CacheHits)
			fmt.Printf("  Cache Misses:     %d\n", stats.CacheMisses)
		}
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&ipfsAPI, "ipfs", "http://localhost:5001", "IPFS API endpoint")
	rootCmd.PersistentFlags().StringVar(&dataDir, "data", "./data", "Data directory")
	rootCmd.PersistentFlags().Int64Var(&cacheSize, "cache", 500*1024*1024, "Cache size in bytes")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Add subcommands
	rootCmd.AddCommand(storeCmd)
	rootCmd.AddCommand(retrieveCmd)
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(statsCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func detectContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	case ".tar.gz", ".tgz":
		return "application/gzip"
	default:
		// Try to detect from file content
		file, err := os.Open(filename)
		if err != nil {
			return "application/octet-stream"
		}
		defer file.Close()

		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return "application/octet-stream"
		}

		return http.DetectContentType(buffer[:n])
	}
}
