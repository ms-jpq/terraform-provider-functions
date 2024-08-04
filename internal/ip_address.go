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
	compressed      string
	exploded        string
	ipv4_mapped     string
	is_global       bool
	is_link_local   bool
	is_loopback     bool
	is_multicast    bool
	is_private      bool
	is_unspecified  bool
	reverse_pointer string
}

func (*IpAddress) New() function.Function {
	return &IpAddress{}
}

func (f *IpAddress) Metadata(_ context.Context, _ function.MetadataRequest, rsp *function.MetadataResponse) {
	rsp.Name = "ip_address"
}

func (f *IpAddress) Definition(_ context.Context, _ function.DefinitionRequest, rsp *function.DefinitionResponse) {
	rsp.Definition = function.Definition{
		Parameters: []function.Parameter{},
		Return: function.ObjectReturn{
			AttributeTypes: map[string]attr.Type{
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
			},
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
	parts := strings.Split(exploded, ":")
	slices.Reverse(parts)
	reverse_domain := ".in-addr.arpa"
	if addr.Is6() {
		reverse_domain = ".ip6.arpa"
	}

	output := IpAddress{
		compressed:      addr.String(),
		exploded:        exploded,
		ipv4_mapped:     addr.Unmap().StringExpanded(),
		is_global:       addr.IsGlobalUnicast(),
		is_link_local:   addr.IsLinkLocalUnicast(),
		is_loopback:     addr.IsLoopback(),
		is_multicast:    addr.IsMulticast(),
		is_private:      addr.IsPrivate(),
		is_unspecified:  addr.IsUnspecified(),
		reverse_pointer: strings.Join(parts, ".") + reverse_domain,
	}
	rsp.Error = function.ConcatFuncErrors(rsp.Error, rsp.Result.Set(ctx, output))
}
