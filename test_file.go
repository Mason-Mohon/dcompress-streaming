// +build ignore

package main

import (
	"fmt"
	"io"
	"os"
	
	"github.com/Mason-Mohon/dcompress-streaming"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_file.go <file.Z>")
		os.Exit(1)
	}
	
	filename := os.Args[1]
	
	// Check file size
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("File: %s (%.2f MB)\n", filename, float64(info.Size())/(1024*1024))
	
	// Open and decompress
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	
	dcompress.VerboseFlag = true
	
	reader, err := dcompress.NewReader(f)
	if err != nil {
		fmt.Printf("NewReader error: %v\n", err)
		os.Exit(1)
	}
	
	// Try to read all
	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("ReadAll error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Success! Decompressed %d bytes\n", len(data))
}
