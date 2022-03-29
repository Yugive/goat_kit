package upgrade

import (
	"fmt"
	"github.com/spf13/cobra"
	"goat/internal/base"
)

var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the goat tool kit",
	Long:  "Upgrade the goat tool kit. Example: goat upgrade",
	Run:   Run,
}

func Run(cmd *cobra.Command, args []string) {
	err := base.GoInstall(
		"github.com/ngyugive/go-web-framework",
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
		"github.com/google/gnostic/cmd/protoc-gen-openapi@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}
