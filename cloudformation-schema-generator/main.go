package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("CloudFormation Schema Generator\n")

	// Fetch and process the AWS published CloudFormation Resource Specification
	cloudformationSpec := "https://d1uauaxba7bl26.cloudfront.net/latest/gzip/CloudFormationResourceSpecification.json"

	sg, err := NewSchemaGenerator(cloudformationSpec)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	if err := sg.Generate(); err != nil {
		fmt.Printf("ERROR: Failed to generate CloudFormation schema: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated CloudFormation schema: schema/cloudformation.schema.json\n")
}