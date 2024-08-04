package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const ProviderName = "functions"

var _ provider.Provider = &FnProvider{}
var _ provider.ProviderWithFunctions = &FnProvider{}

type FnProvider struct {
	version string
}

func (*FnProvider) New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FnProvider{
			version: version,
		}
	}
}

func (p *FnProvider) Metadata(_ context.Context, _ provider.MetadataRequest, rsp *provider.MetadataResponse) {
	rsp.TypeName = ProviderName
	rsp.Version = p.version
}

func (p *FnProvider) Schema(_ context.Context, _ provider.SchemaRequest, rsp *provider.SchemaResponse) {
	rsp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *FnProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *FnProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		(*IpAddress)(nil).New,
	}
}

func (p *FnProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *FnProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
