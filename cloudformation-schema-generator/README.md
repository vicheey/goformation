# CloudFormation Schema Generator

This package generates a JSON Schema for AWS CloudFormation templates from the official AWS CloudFormation Resource Specification.

## Purpose

This is a minimal, focused tool that:

- Downloads the latest AWS CloudFormation Resource Specification
- Generates a complete JSON Schema for CloudFormation templates
- Outputs the schema to `schema/cloudformation.schema.json`

## Usage

```bash
go run .
```

This will:

1. Download the CloudFormation Resource Specification from AWS
2. Generate the JSON schema using the original GoFormation template logic
3. Save the schema to `schema/cloudformation.schema.json`

## Generated Files

- `schema/cloudformation.schema.json` - The complete CloudFormation JSON Schema

## Implementation

This package contains only the minimal code needed for schema generation:

- `main.go` - Entry point
- `generator.go` - Schema generation logic
- `spec.go` - CloudFormation specification data structures
- `templates/` - Template files for schema generation

The logic is adapted from the original GoFormation generate package but simplified to focus only on schema generation.
