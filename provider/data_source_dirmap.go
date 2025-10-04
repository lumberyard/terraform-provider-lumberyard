// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	req.Config.Get(ctx, &state)

	result, _ := buildMap(state.Path.ValueString(), state.Filter.ValueString())

	resultMap, _ := types.MapValueFrom(ctx, types.MapType{ElemType: types.DynamicType}, result)

	state.Result = types.DynamicValue(resultMap)

	// Set state
	resp.State.Set(ctx, &state)
}