# AWS CloudFormation Schema Generator

This repository contains a focused tool for generating JSON Schema for AWS CloudFormation templates.

## Purpose

This project provides a minimal, standalone CloudFormation schema generator that:

- Downloads the latest AWS CloudFormation Resource Specification
- Generates a complete JSON Schema for CloudFormation templates
- Outputs the schema for use in IDEs, validators, and other tools

## Usage

Navigate to the `cloudformation-schema-generator` directory and run:

```bash
cd cloudformation-schema-generator
go run .
```

Or build the binary:

```bash
cd cloudformation-schema-generator
go build -o cloudformation-schema-generator .
./cloudformation-schema-generator
```

This will generate `schema/cloudformation.schema.json` containing the complete CloudFormation JSON Schema.

## Testing

To run the tests:

```bash
cd cloudformation-schema-generator
go test -v
```

The test verifies that the schema generation produces consistent output by comparing checksums against a reference schema.

## Generated Schema

The generated schema includes:

- All AWS CloudFormation resources and their properties
- Proper validation rules and constraints
- Support for CloudFormation intrinsic functions
- Complete property type definitions

## Implementation

This package contains only the minimal code needed for schema generation:

- `main.go` - Entry point
- `generator.go` - Schema generation logic
- `spec.go` - CloudFormation specification data structures
- `templates/` - Template files for schema generation

The logic is adapted from the original GoFormation generate package but simplified to focus only on schema generation.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
