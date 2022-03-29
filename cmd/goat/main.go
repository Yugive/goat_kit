package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"goat/internal/base"
	"goat/internal/project"
	"goat/internal/proto"
	"goat/internal/run"
	"goat/internal/upgrade"
)

var rootCmd = &cobra.Command{
	Use:     "goat",
	Short:   "Goat: a toolkit for Go microservices",
	Long:    "Goat: a toolkit for Go microservices",
	Version: release,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(run.CmdRun)
}

func main() {
	/*
		if err := rootCmd.Execute(); err != nil {
			log.Fatal(err)
		}

	*/
	fmt.Printf("%s", base.GoatMod())
}
