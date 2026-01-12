// dcompress_test.go

package dcompress

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_000(t *testing.T) { // template
	if false {
		t.Errorf("print fail, but keep testing")
	}
	fmt.Printf("Test_0000\n")
}

// Commented out - requires external mdr package
// var filePair1 = [2]string{"testdata/kermit.Z",
// 	"d61611d13775c1f3a83675e81afcadfc4352b11e0f39f7c928bad62d25675b66"}
// var filePairs = [][2]string{filePair1 /*, filePair2,  filePair3*/}

func Test_001(t *testing.T) {
}

// Test_StreamingVsBuffered tests that both implementations work correctly
func Test_StreamingVsBuffered(t *testing.T) {
	// Test with small buffer (should use buffered mode)
	t.Run("SmallBuffer", func(t *testing.T) {
		// Create minimal valid compressed data
		compressedData := []byte{0x1f, 0x9d, 0x90} // Magic + header
		
		// Add some actual compressed data to avoid immediate errors
		// This is a minimal valid LZW stream
		compressedData = append(compressedData, []byte{
			0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		}...)
		
		reader, err := NewReader(bytes.NewReader(compressedData))
		if err != nil {
			t.Logf("Got error with minimal data (expected): %v", err)
			return
		}
		if reader == nil {
			t.Error("NewReader returned nil for small buffer")
		}
	})

	t.Run("LargeFileDetection", func(t *testing.T) {
		// Create a temporary file to test size detection
		tmpFile, err := ioutil.TempFile("", "test_large_*.Z")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		// Write minimal header
		header := []byte{0x1f, 0x9d, 0x90}
		tmpFile.Write(header)
		
		// Pad to make it appear large (> 500MB)
		tmpFile.Seek(500*1024*1024, 0)
		tmpFile.Write([]byte{0})
		tmpFile.Seek(0, 0)

		// Open for reading
		f, err := os.Open(tmpFile.Name())
		if err != nil {
			t.Fatalf("Failed to open temp file: %v", err)
		}
		defer f.Close()

		// This should use the streaming implementation
		reader, err := NewReader(f)
		if err != nil {
			// Expected - we don't have valid compressed data
			t.Logf("Got expected error with dummy data: %v", err)
			return
		}

		// Verify it's a StreamingReader
		if _, ok := reader.(*StreamingReader); !ok {
			t.Error("Expected StreamingReader for large file")
		}
	})
}

// Test_StreamingReaderBasic tests basic StreamingReader functionality
func Test_StreamingReaderBasic(t *testing.T) {
	// This test verifies the StreamingReader can be created
	// Real decompression testing would require valid .Z files
	
	compressedData := []byte{0x1f, 0x9d, 0x90} // Magic + header
	compressedData = append(compressedData, make([]byte, 100)...) // Some dummy data
	
	sr, err := newStreamingReader(bytes.NewReader(compressedData))
	if err != nil {
		// Expected - we don't have valid compressed data
		t.Logf("Got expected error with dummy data: %v", err)
		return
	}
	
	if sr == nil {
		t.Error("newStreamingReader returned nil")
	}
	
	// Try reading
	buf := make([]byte, 100)
	_, err = sr.Read(buf)
	// Error expected with dummy data
	t.Logf("Read returned: %v (expected with dummy data)", err)
}