package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// Ensure the implementation satisfies the expected interfaces.
var _ function.Function = &dirmapFunction{}

// NewDirmapFunction is a helper function to simplify the provider implementation.
func NewDirmapFunction() function.Function {
	return &dirmapFunction{}
}

// dirmapFunction is the function implementation.
type dirmapFunction struct{}

// Metadata returns the function type name.
func (f *dirmapFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "dirmap"
}

// Definition defines the function signature.
func (f *dirmapFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Builds a nested object from YAML/JSON files in a directory tree, returned as JSON string.",
		Description: "Recursively reads YAML (.yaml/.yml) and JSON (.json) files from the given directory (and subdirectories), merges them into a nested map using relative file paths as keys (without extension), and returns the entire structure serialized as a JSON string.\n\nUse `jsondecode(lumberyard::dirmap(...))` to get a Terraform object/map.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "path",
				Description: "Absolute or relative path to the base directory to traverse. Must exist and be readable.",
			},
		},
		VariadicParameter: function.StringParameter{
			Name:        "filter",
			Description: "Optional glob pattern to match filenames (e.g. \"**/*.yaml\", \"config/*.json\"). If omitted, includes all supported YAML and JSON files.",
		},
		Return: function.StringReturn{},
	}
}

// Run executes the function.
func (f *dirmapFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var path string
	var filter string

	// Read required path argument
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &path))

	// Read optional variadic filter argument
	// If not provided, Get will return an error which we can ignore
	if err := req.Arguments.Get(ctx, &filter); err != nil {
		// Filter not provided, use empty string as default
		filter = ""
	}

	// Build the map (reuses your existing logic)
	result, err := buildMap(path, filter)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error,
			function.NewFuncError(fmt.Sprintf("Unable to build map from directory %s: %v", path, err)),
		)
		return
	}

	// Marshal to JSON string (same as the data source)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error,
			function.NewFuncError(fmt.Sprintf("Failed to marshal result to JSON: %v", err)),
		)
		return
	}

	// Set the result
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, string(jsonResult)))
}
