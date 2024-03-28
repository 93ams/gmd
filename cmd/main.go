package main

import (
	. "github.com/93ams/gmd"
	"github.com/samber/lo"
)

var Root = New("gmd", Add(Env))

func main() { lo.Must0(Root.Execute()) }
