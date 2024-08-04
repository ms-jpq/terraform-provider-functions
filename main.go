package main

import (
	"context"
	"flag"
	"log"
	"main/internal"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const WhoAmI = "ms-jpq"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/" + WhoAmI + "/" + internal.ProviderName,
		Debug:   debug,
	}

	provider := (*internal.FuncProvider)(nil).New
	if err := providerserver.Serve(context.Background(), provider, opts); err != nil {
		log.Fatal(err.Error())
	}
}
