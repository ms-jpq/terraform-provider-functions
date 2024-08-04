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

var _ provider.Provider = &FuncProvider{}
var _ provider.ProviderWithFunctions = &FuncProvider{}

type FuncProvider struct {
	version string
}

func (*FuncProvider) New() provider.Provider {
	return &FuncProvider{}
}

func (p *FuncProvider) Metadata(_ context.Context, _ provider.MetadataRequest, rsp *provider.MetadataResponse) {
	rsp.TypeName = ProviderName
}

func (p *FuncProvider) Schema(_ context.Context, _ provider.SchemaRequest, rsp *provider.SchemaResponse) {
	rsp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *FuncProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *FuncProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		(*IpAddress)(nil).New,
	}
}

func (p *FuncProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *FuncProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
