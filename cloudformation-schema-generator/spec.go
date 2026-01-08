package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"
)

// CloudFormationResourceSpecification represents a resource specification document
type CloudFormationResourceSpecification struct {
	ResourceSpecificationVersion   string
	ResourceSpecificationTransform string
	Resources                      map[string]Resource `json:"ResourceTypes"`
	Properties                     map[string]Resource `json:"PropertyTypes"`
	Globals                        Globals             `json:"Globals"`
}

// Resource represents an AWS CloudFormation resource
type Resource struct {
	Properties map[string]Property
}

// Required returns a comma separated list of the required properties for this resource
func (r Resource) Required() string {
	required := []string{}
	for name, property := range r.Properties {
		if property.Required {
			required = append(required, `"`+name+`"`)
		}
	}
	sort.Strings(required)
	return strings.Join(required, ", ")
}

// Schema returns a JSON Schema for the resource (as a string)
func (r Resource) Schema(name string, isCustomProperty bool) string {
	tmpl, err := template.New("schema-resource.template").Funcs(template.FuncMap{
		"counter": counter,
	}).ParseFiles("templates/schema-resource.template")
	if err != nil {
		fmt.Printf("Error: Failed to load resource schema template: %s\n", err)
		return ""
	}

	templateData := struct {
		Name             string
		Resource         Resource
		IsCustomProperty bool
	}{name, r, isCustomProperty}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		fmt.Printf("Error: Failed to generate resource schema for %s: %s\n", name, err)
		return ""
	}
	return buf.String()
}

// Property represents an AWS CloudFormation resource property
type Property struct {
	ItemType                    string   `json:"ItemType"`
	PrimitiveItemType           string   `json:"PrimitiveItemType"`
	PrimitiveType               string   `json:"PrimitiveType"`
	Required                    bool     `json:"Required"`
	Type                        string   `json:"Type"`
	PrimitiveTypes              []string `json:"PrimitiveTypes"`
	PrimitiveItemTypes          []string `json:"PrimitiveItemTypes"`
	ItemTypes                   []string `json:"ItemTypes"`
	Types                       []string `json:"Types"`
	InclusivePrimitiveItemTypes []string `json:"InclusivePrimitiveItemTypes"`
	InclusiveItemTypes          []string `json:"InclusiveItemTypes"`
	InclusiveItemPattern        bool     `json:"InclusiveItemPattern"`
}

var typeToJSON = map[string]string{
	"String": "string", "Long": "number", "Integer": "number", "Double": "number",
	"Boolean": "boolean", "Timestamp": "string", "Json": "object", "Map": "object",
}

// IsPolymorphic checks whether a property can be multiple different types
func (p Property) IsPolymorphic() bool {
	return len(p.PrimitiveTypes) > 0 || len(p.PrimitiveItemTypes) > 0 ||
		len(p.ItemTypes) > 0 || len(p.Types) > 0 ||
		len(p.InclusivePrimitiveItemTypes) > 0 || len(p.InclusiveItemTypes) > 0
}

// IsPrimitive checks whether a property is a primitive type
func (p Property) IsPrimitive() bool {
	return p.PrimitiveType != ""
}

// IsMap checks whether a property should be a map
func (p Property) IsMap() bool {
	return p.Type == "Map"
}

// IsList checks whether a property should be a list
func (p Property) IsList() bool {
	return p.Type == "List"
}

// IsCustomType checks whether a property is a custom type
func (p Property) IsCustomType() bool {
	return p.PrimitiveType == "" && p.ItemType == "" && p.PrimitiveItemType == "" && 
		p.Type != "" && p.Type != "List" && p.Type != "Map"
}

// GetJSONPrimitiveType returns the correct primitive property type for a JSON Schema
func (p Property) GetJSONPrimitiveType() string {
	return p.convertTypeToJSON()
}

// HasJSONPrimitiveType checks if GetJSONPrimitiveType is not ""
func (p Property) HasJSONPrimitiveType() bool {
	return p.convertTypeToJSON() != ""
}

func (p Property) convertTypeToJSON() string {
	if p.PrimitiveType != "" {
		return convertTypeToJSON(p.PrimitiveType)
	} else if p.PrimitiveItemType != "" {
		return convertTypeToJSON(p.PrimitiveItemType)
	} else {
		return convertTypeToJSON(p.ItemType)
	}
}

func convertTypeToJSON(name string) string {
	if t, ok := typeToJSON[name]; ok {
		return t
	}
	return ""
}

// Schema returns a JSON Schema for the property (as a string)
func (p Property) Schema(name, parent string) string {
	tmpl, err := template.New("schema-property.template").Funcs(template.FuncMap{
		"counter":           counter,
		"convertToJSONType": convertTypeToJSON,
	}).ParseFiles("templates/schema-property.template")
	if err != nil {
		fmt.Printf("Error: Failed to load property schema template: %s\n", err)
		return ""
	}

	parentpaths := strings.Split(parent, ".")
	templateData := struct {
		Name     string
		Parent   string
		Property Property
	}{name, parentpaths[0], p}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		fmt.Printf("Error: Failed to generate property schema for %s: %s\n", name, err)
		return ""
	}
	return buf.String()
}

// Globals represents the Globals section of a SAM template
type Globals struct {
	Children map[string]Global `json:",inline"`
}

// Global represents a global configuration in SAM
type Global struct {
	Reference string   `json:"Reference"`
	Exclude   []string `json:"Exclude"`
}

// Schema returns a JSON Schema for the global (as a string)
func (g Global) Schema(name string, resources map[string]Resource) string {
	tmpl, err := template.New("schema-globals.template").Funcs(template.FuncMap{
		"counter": counter,
	}).ParseFiles("templates/schema-globals.template")
	if err != nil {
		fmt.Printf("Error: Failed to load globals schema template: %s\n", err)
		return ""
	}

	templateData := struct {
		Name      string
		Global    Global
		Resources map[string]Resource
	}{name, g, resources}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		fmt.Printf("Error: Failed to generate global schema for %s: %s\n", name, err)
		return ""
	}
	return buf.String()
}

func (g Global) isExcluded(propertyName string) bool {
	for _, excluded := range g.Exclude {
		if excluded == propertyName {
			return true
		}
	}
	return false
}
