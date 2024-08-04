package internal

import (
	"context"
	"net"
	"net/netip"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type IpAddress struct {
}

func (*IpAddress) New() function.Function {
	return &IpAddress{}
}

func (f *IpAddress) Metadata(_ context.Context, _ function.MetadataRequest, rsp *function.MetadataResponse) {
	rsp.Name = "ip_address"
}

func (f *IpAddress) def() map[string]attr.Type {
	return map[string]attr.Type{
		"compressed":      basetypes.StringType{},
		"exploded":        basetypes.StringType{},
		"ipv4_mapped":     basetypes.StringType{},
		"is_global":       basetypes.BoolType{},
		"is_link_local":   basetypes.BoolType{},
		"is_loopback":     basetypes.BoolType{},
		"is_multicast":    basetypes.BoolType{},
		"is_private":      basetypes.BoolType{},
		"is_unspecified":  basetypes.BoolType{},
		"reverse_pointer": basetypes.StringType{},
	}
}

func (f *IpAddress) Definition(_ context.Context, _ function.DefinitionRequest, rsp *function.DefinitionResponse) {
	rsp.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				Name: "ip",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: f.def(),
		},
	}
}

func (f *IpAddress) Run(ctx context.Context, req function.RunRequest, rsp *function.RunResponse) {
	var input string
	rsp.Error = function.ConcatFuncErrors(rsp.Error, req.Arguments.Get(ctx, &input))
	if rsp.Error != nil {
		return
	}

	ip := net.ParseIP(input)
	if ip == nil {
		_, err := netip.ParseAddr(input)
		rsp.Error = function.ConcatFuncErrors(rsp.Error, function.NewFuncError(err.Error()))
		return
	}
	addr, _ := netip.ParseAddr(ip.String())

	exploded := addr.StringExpanded()
	reverse_pointer := func() string {
		split := []string{}
		domain := ""
		if addr.Is4() {
			domain = ".in-addr.arpa"
			strings.Split(exploded, ".")
		} else {
			domain = ".ip6.arpa"
			split = strings.Split(strings.ReplaceAll(exploded, ":", ""), "")
		}
		slices.Reverse(split)
		return strings.Join(split, ".") + domain
	}()

	output, diags := basetypes.NewObjectValue(f.def(),
		map[string]attr.Value{
			"compressed":      basetypes.NewStringValue(addr.String()),
			"exploded":        basetypes.NewStringValue(exploded),
			"ipv4_mapped":     basetypes.NewStringValue(addr.Unmap().StringExpanded()),
			"is_global":       basetypes.NewBoolValue(addr.IsGlobalUnicast()),
			"is_link_local":   basetypes.NewBoolValue(addr.IsLinkLocalUnicast()),
			"is_loopback":     basetypes.NewBoolValue(addr.IsLoopback()),
			"is_multicast":    basetypes.NewBoolValue(addr.IsMulticast()),
			"is_private":      basetypes.NewBoolValue(addr.IsPrivate()),
			"is_unspecified":  basetypes.NewBoolValue(addr.IsUnspecified()),
			"reverse_pointer": basetypes.NewStringValue(reverse_pointer),
		},
	)
	rsp.Error = function.ConcatFuncErrors(rsp.Error, function.FuncErrorFromDiags(ctx, diags))
	rsp.Error = function.ConcatFuncErrors(rsp.Error, rsp.Result.Set(ctx, output))
}
