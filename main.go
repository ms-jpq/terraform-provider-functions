package main

import (
	"context"
	"flag"
	"log"
	"main/internal"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/ms-jpq/func",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), internal.NewProvider(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
