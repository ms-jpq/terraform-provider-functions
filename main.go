package main

import (
	"context"
	"flag"
	"log"
	"main/internal"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const WhoAmI = "ms-jpq"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/" + WhoAmI + "/" + internal.ProviderName,
		Debug:   debug,
	}

	provider := (*internal.FnProvider)(nil).New(version)
	if err := providerserver.Serve(context.Background(), provider, opts); err != nil {
		log.Fatal(err.Error())
	}
}
