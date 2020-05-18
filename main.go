package main

import (
	"runtime"

	"github.com/dev-techmoe/filebrowser/v2/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
