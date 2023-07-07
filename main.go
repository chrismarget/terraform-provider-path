package main

import (
	"context"
	"flag"
	tfpath "github.com/chrismarget/terraform-provider-path/path"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), tfpath.NewProvider, providerserver.ServeOpts{
		Address: "registry.terraform.io/chrismarget/path",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
