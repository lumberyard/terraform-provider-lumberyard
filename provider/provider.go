// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure LumberyardProvider satisfies various provider interfaces.
var _ provider.Provider = &LumberyardProvider{}

type LumberyardProvider struct {
	version string
}

func (p *LumberyardProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lumberyard"
	resp.Version = p.version
}

func (p *LumberyardProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A terraform provider for interacting with lumberyard.",
	}
}

func (p *LumberyardProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *LumberyardProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *LumberyardProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDirmapDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LumberyardProvider{
			version: version,
		}
	}
}
