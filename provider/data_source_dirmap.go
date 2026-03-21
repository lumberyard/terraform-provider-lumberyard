package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	Path   types.String `tfsdk:"path"`
	Filter types.String `tfsdk:"filter"`
	Result types.String `tfsdk:"result"`
}

// Metadata returns the data source type name.
func (d *dirmapDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dirmap"
}

// Schema defines the schema for the data source.
func (d *dirmapDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Traverses a nested directory of YAML and JSON files and constructs a nested object, returned as a JSON string.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "The base directory to traverse.",
				Required:    true,
			},
			"filter": schema.StringAttribute{
				Description: "A glob pattern to filter files.",
				Optional:    true,
			},
			"result": schema.StringAttribute{
				Description: "The constructed object from the directory structure, as a JSON string.",
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

	// Marshal the result to a JSON string
	jsonResult, err := json.Marshal(result)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to marshal result to JSON",
			err.Error(),
		)
		return
	}

	state.Result = types.StringValue(string(jsonResult))

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
