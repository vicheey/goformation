# AWS CloudFormation Schema Generator

This repository contains a focused tool for generating JSON Schema for AWS CloudFormation templates.

## Purpose

This project provides a minimal, standalone CloudFormation schema generator that:
- Downloads the latest AWS CloudFormation Resource Specification
- Generates a complete JSON Schema for CloudFormation templates
- Outputs the schema for use in IDEs, validators, and other tools

## Usage

Navigate to the `schema-generator` directory and run:

```bash
cd schema-generator
go run .
```

Or build the binary:

```bash
cd schema-generator
./build.sh
./cfn-schema-generator
```

This will generate `schema/cloudformation.schema.json` containing the complete CloudFormation JSON Schema.

## Generated Schema

The generated schema includes:
- All AWS CloudFormation resources and their properties
- Proper validation rules and constraints
- Support for CloudFormation intrinsic functions
- Complete property type definitions

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
