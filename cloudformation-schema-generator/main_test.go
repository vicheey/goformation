package main

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestSchemaGeneration(t *testing.T) {
	// Get expected checksum from reference file
	expectedFile, err := os.Open("test_expected_cloudformation.schema.json")
	if err != nil {
		t.Fatalf("Failed to open expected schema file: %v", err)
	}
	defer expectedFile.Close()
	
	expectedHash := sha256.New()
	if _, err := io.Copy(expectedHash, expectedFile); err != nil {
		t.Fatalf("Failed to calculate expected checksum: %v", err)
	}
	expectedChecksum := fmt.Sprintf("%x", expectedHash.Sum(nil))
	
	// Read and decompress local test file
	file, err := os.Open("test_cloudformation_spec.json")
	if err != nil {
		t.Fatalf("Failed to open test spec file: %v", err)
	}
	defer file.Close()
	
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()
	
	var spec CloudFormationResourceSpecification
	if err := json.NewDecoder(gzReader).Decode(&spec); err != nil {
		t.Fatalf("Failed to parse test spec: %v", err)
	}
	
	sg := &SchemaGenerator{}
	if err := sg.generateJSONSchema("cloudformation", &spec); err != nil {
		t.Fatalf("Failed to generate schema: %v", err)
	}
	
	// Check checksum
	schemaFile, err := os.Open("schema/cloudformation.schema.json")
	if err != nil {
		t.Fatalf("Failed to open generated schema: %v", err)
	}
	defer schemaFile.Close()
	
	hash := sha256.New()
	if _, err := io.Copy(hash, schemaFile); err != nil {
		t.Fatalf("Failed to calculate checksum: %v", err)
	}
	
	actualChecksum := fmt.Sprintf("%x", hash.Sum(nil))
	
	if actualChecksum != expectedChecksum {
		t.Errorf("Schema checksum mismatch. Expected: %s, Got: %s", expectedChecksum, actualChecksum)
	}
}
