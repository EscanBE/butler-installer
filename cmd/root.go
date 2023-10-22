package cmd

import (
	"github.com/EscanBE/butler-installer/constants"
	"github.com/EscanBE/butler-installer/types"
	"github.com/EscanBE/butler-installer/utils"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.BINARY_NAME,
	Short: constants.APP_DESC,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true    // hide the 'completion' subcommand
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true}) // hide the 'help' subcommand

	operationUserInfo, err := utils.GetOperationUserInfo()
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to get operation user info:", err)
		os.Exit(1)
	}

	ctx := types.NewContext(operationUserInfo)

	err = rootCmd.ExecuteContext(types.WrapAppContext(ctx))
	if err != nil {
		os.Exit(1)
	}
}
