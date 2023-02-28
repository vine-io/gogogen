package main

import (
	"os"

	"github.com/spf13/pflag"
	goproto "github.com/vine-io/gogogen/gogorm-gen"
)

func main() {
	g := goproto.New()
	fs := pflag.NewFlagSet("gogorm", pflag.ExitOnError)
	g.BindFlags(fs)
	fs.Parse(os.Args)
	goproto.Run(g)
}
