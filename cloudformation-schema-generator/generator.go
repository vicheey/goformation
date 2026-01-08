package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
)

// SchemaGenerator handles the generation of CloudFormation JSON schema
type SchemaGenerator struct {
	CloudFormationSpecURL string
}

// NewSchemaGenerator creates a new schema generator for CloudFormation only
func NewSchemaGenerator(cloudformationSpecURL string) (*SchemaGenerator, error) {
	return &SchemaGenerator{CloudFormationSpecURL: cloudformationSpecURL}, nil
}

// Generate downloads the CloudFormation specification and generates the JSON schema
func (sg *SchemaGenerator) Generate() error {
	fmt.Printf("Downloading CloudFormation Resource Specification...\n")
	
	spec, err := sg.downloadAndParseSpec(sg.CloudFormationSpecURL)
	if err != nil {
		return fmt.Errorf("failed to download CloudFormation spec: %w", err)
	}

	fmt.Printf("Generating CloudFormation JSON schema...\n")
	
	if err := sg.generateJSONSchema("cloudformation", spec); err != nil {
		return fmt.Errorf("failed to generate schema: %w", err)
	}

	return nil
}

// downloadAndParseSpec downloads and parses a CloudFormation specification
func (sg *SchemaGenerator) downloadAndParseSpec(url string) (*CloudFormationResourceSpecification, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body
	
	// Handle gzipped content
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gzReader.Close()
		reader = gzReader
	}

	var spec CloudFormationResourceSpecification
	if err := json.NewDecoder(reader).Decode(&spec); err != nil {
		return nil, err
	}

	return &spec, nil
}

// generateJSONSchema generates a JSON schema from a CloudFormation specification
func (sg *SchemaGenerator) generateJSONSchema(specname string, spec *CloudFormationResourceSpecification) error {
	tmpl, err := template.New("schema.template").Funcs(template.FuncMap{
		"counter": counter,
	}).ParseFiles("templates/schema.template")
	if err != nil {
		return fmt.Errorf("failed to parse schema template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, spec); err != nil {
		return fmt.Errorf("failed to generate JSON Schema: %s", err)
	}

	// Format the JSON
	var j interface{}
	if err := json.Unmarshal(buf.Bytes(), &j); err != nil {
		return fmt.Errorf("failed to unmarshal JSON Schema: %s", err)
	}

	formatted, err := json.MarshalIndent(j, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON Schema: %s", err)
	}

	// Write to file
	if err := os.MkdirAll("schema", 0755); err != nil {
		return fmt.Errorf("failed to create schema directory: %w", err)
	}

	filename := fmt.Sprintf("schema/%s.schema.json", specname)
	if err := os.WriteFile(filename, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write JSON Schema: %s", err)
	}

	return nil
}

// counter is a template function that helps with comma placement in JSON
func counter(length int) func() bool {
	i := 0
	return func() bool {
		i++
		return i < length
	}
}