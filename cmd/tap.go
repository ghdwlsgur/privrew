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
		Short: "...",
		Long:  "...",
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
				panicRed(fmt.Errorf("..."))
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
	tapCommand.Flags().StringP("token", "t", "", "...")
	viper.BindPFlag("tap-token", tapCommand.Flag("token"))

	rootCmd.AddCommand(tapCommand)
}
