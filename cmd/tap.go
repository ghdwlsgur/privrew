package cmd

import (
	"fmt"
	"strings"

	"github.com/ghdwlsgur/privrew/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tapCommand = &cobra.Command{
		Use:   "tap",
		Short: "It provides the same functionality as the tap feature in Homebrew.",
		Long:  "allows users to register and manage additional repositories, enabling them to access and install packages from external sources beyond the default Homebrew/core repository.",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err error
			)

			if err = cobra.MinimumNArgs(1)(cmd, args); err != nil {
				panicRed(err)
			}

			owner := strings.Split(args[0], "/")[0]
			repository := strings.Split(args[0], "/")[1]
			token := viper.GetString("tap-token")

			if token == "" {
				panicRed(fmt.Errorf("Please enter the token value."))
			}

			repo := &internal.Repository{
				Name:     repository,
				Owner:    owner,
				Token:    token,
				LocalDir: tapDir,
			}

			internal.CloneReleaseRepo(repo)
		},
	}
)

func init() {
	tapCommand.Flags().StringP("token", "t", "", "Accessible token for the private repository")
	viper.BindPFlag("tap-token", tapCommand.Flag("token"))

	rootCmd.AddCommand(tapCommand)
}
