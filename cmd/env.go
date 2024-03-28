package main

import (
	"fmt"
	. "github.com/93ams/gmd"
	"github.com/93ams/gmd/env"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var store = env.NewStore("genv")
var Env = New("env",
	PPreRun(func(cmd *cobra.Command, _ []string) { lo.Must0(store.Open()) }),
	PPostRun(func(cmd *cobra.Command, _ []string) { lo.Must0(store.Close()) }),
	Add(
		New("add",
			Run(func(cmd *cobra.Command, args []string) { lo.Must0(store.Add(args...)) }),
		),
		New("rem",
			Run(func(cmd *cobra.Command, args []string) { lo.Must0(store.Rem(args...)) }),
		),
		New("list",
			Run(func(cmd *cobra.Command, _ []string) { fmt.Println(strings.Join(store.Keys(), " ")) }),
		),
		New("show",
			Run(func(cmd *cobra.Command, _ []string) { fmt.Println(store.GetCurrent().AllSettings()) }),
		),
		New("select",
			Args(cobra.MaximumNArgs(1)),
			Run(func(cmd *cobra.Command, args []string) {
				var env string
				if len(args) > 0 {
					env = args[0]
				}
				lo.Must(store.Select(env))
			}),
		),
		New("current",
			Run(func(cmd *cobra.Command, _ []string) { fmt.Println(store.Current()) }),
		),
		New("set",
			Run(func(cmd *cobra.Command, args []string) {
				env := store.GetCurrent()
				for _, arg := range args {
					before, after, found := strings.Cut(arg, "=")
					if !found { // bool
						env.Set(strings.TrimPrefix(arg, "!"), !strings.HasPrefix(arg, "!"))
						continue
					}
					parts := strings.Split(before, ":")
					var value any = after
					if len(parts) == 2 {
						switch parts[1] {
						case "int":
							value = lo.Must(strconv.Atoi(after))
						}
					}
					env.Set(parts[0], value)
				}
				lo.Must0(store.Save(env.Name()))
			}),
		),
		New("unset",
			Run(func(cmd *cobra.Command, args []string) {
				env := store.GetCurrent()
				for _, arg := range args {
					env.Set(arg, nil)
				}
				lo.Must0(store.Save(env.Name()))
			}),
		),
	))
