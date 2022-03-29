package proto

import (
	"github.com/spf13/cobra"
	"goat/internal/proto/add"
	"goat/internal/proto/client"
	"goat/internal/proto/server"
)

var CmdProto = &cobra.Command{
	Use:   "proto",
	Short: "Generate the proto file",
	Long:  "Generate the proto file",
	Run:   run,
}

func init() {
	CmdProto.AddCommand(add.CmdAdd)
	CmdProto.AddCommand(client.CmdClient)
	CmdProto.AddCommand(server.CmdServer)
}

func run(cmd *cobra.Command, args []string) {

}
