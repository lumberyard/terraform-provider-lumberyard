// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &dirmapDataSource{}
)

// NewDirmapDataSource is a helper function to simplify the provider implementation.
func NewDirmapDataSource() datasource.DataSource {
	return &dirmapDataSource{}
}

// dirmapDataSource is the data source implementation.
type dirmapDataSource struct{}

// dirmapDataSourceModel maps the data source schema data.
type dirmapDataSourceModel struct {
	Path   types.String  `tfsdk:"path"`
	Filter types.String  `tfsdk:"filter"`
	Result types.Dynamic `tfsdk:"result"`
}

// Metadata returns the data source type name.
func (d *dirmapDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dirmap"
}

// Schema defines the schema for the data source.
func (d *dirmapDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Traverses a nested directory of YAML and JSON files and constructs a nested object.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "The base directory to traverse.",
				Required:    true,
			},
			"filter": schema.StringAttribute{
				Description: "A glob pattern to filter files.",
				Optional:    true,
			},
			"result": schema.DynamicAttribute{
				Description: "The constructed object from the directory structure.",
				Computed:    true,
			},
		},
	}
}

// Read is called when the data source is read.
func (d *dirmapDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dirmapDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := state.Path.ValueString()
	var filter string
	if !state.Filter.IsNull() && !state.Filter.IsUnknown() {
		filter = state.Filter.ValueString()
	}

	result, err := buildMap(path, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to build map",
			fmt.Sprintf("Failed to build map from directory %s: %v", path, err),
		)
		return
	}

	tflog.Debug(ctx, "Parsed directory structure", map[string]interface{}{
		"path":   path,
		"result": result,
	})

	resultAttr, diags := convertToAttrValue(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resultDynamic := types.DynamicValue(resultAttr)
	state.Result = resultDynamic

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// convertToAttrValueSingle handles individual value conversion recursively.
func convertToAttrValueSingle(ctx context.Context, v interface{}) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch t := v.(type) {
	case map[string]interface{}:
		elements := make(map[string]attr.Value)
		for key, val := range t {
			elemVal, d := convertToAttrValueSingle(ctx, val)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}
			elements[key] = elemVal
		}
		return types.MapValueMust(types.DynamicType, elements), diags
	case []interface{}:
		elements := make([]attr.Value, len(t))
		for i, val := range t {
			elemVal, d := convertToAttrValueSingle(ctx, val)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}
			elements[i] = elemVal
		}
		return types.ListValueMust(types.DynamicType, elements), diags
	case string:
		return types.StringValue(t), diags
	case float64:
		return types.NumberValue(big.NewFloat(t)), diags
	case bool:
		return types.BoolValue(t), diags
	case nil:
		return types.DynamicNull(), diags
	default:
		diags.AddError(
			"Unsupported value type",
			fmt.Sprintf("Cannot convert type %T to attr.Value", v),
		)
		return nil, diags
	}
}

// convertToAttrValue converts the top-level interface{} to attr.Value.
func convertToAttrValue(ctx context.Context, v interface{}) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch t := v.(type) {
	case map[string]interface{}:
		return convertToAttrValueSingle(ctx, t)
	default:
		diags.AddError(
			"Invalid top-level type",
			fmt.Sprintf("Expected map[string]interface{} at top level, got %T", v),
		)
		return nil, diags
	}
}
