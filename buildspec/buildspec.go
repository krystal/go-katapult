// Package buildspec implements the Katapult Virtual Machine build spec XML
// document format.
//
// It supports both building and parsing buld spec XML documents to/from Go
// structs, JSON and YAML.
package buildspec

import (
	"encoding/json"
	"encoding/xml"
	"io"

	"gopkg.in/yaml.v3"
)

const (
	address   = "address"
	name      = "name"
	permalink = "permalink"
)

// FromJSON parses a JSON build spec document into a *VirtualMachineSpec object.
func FromJSON(r io.Reader) (*VirtualMachineSpec, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	spec := &VirtualMachineSpec{}
	err := dec.Decode(spec)

	return spec, err
}

// FromXML parses a XML build spec document into a *VirtualMachineSpec object.
func FromXML(r io.Reader) (*VirtualMachineSpec, error) {
	dec := xml.NewDecoder(r)

	spec := &VirtualMachineSpec{}
	err := dec.Decode(spec)

	return spec, err
}

// FromYAML parses a YAML build spec document into a *VirtualMachineSpec object.
func FromYAML(r io.Reader) (*VirtualMachineSpec, error) {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)

	spec := &VirtualMachineSpec{}
	err := dec.Decode(spec)

	return spec, err
}
