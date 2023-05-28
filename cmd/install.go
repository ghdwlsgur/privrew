package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ghdwlsgur/privrew/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	installCommand = &cobra.Command{
		Use:   "install",
		Short: "Download the released software from the private repository.",
		Long:  "Download the released software from the private repository.",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err   error
				table map[string]internal.Asset
			)

			if err = cobra.NoArgs(cmd, args); err != nil {
				panicRed(err)
			}
			if err = cobra.MinimumNArgs(1)(cmd, args); err != nil {
				panicRed(err)
			}

			owner := strings.Split(args[0], "/")[0]
			repository := strings.Split(args[0], "/")[1]
			token := viper.GetString("install-token")

			if token == "" {
				panicRed(fmt.Errorf("Please enter the token value."))
			}

			repo := &internal.Repository{
				Name:     repository,
				Owner:    owner,
				Token:    token,
				LocalDir: installDir,
			}

			os := &internal.OS{
				Name: runtime.GOOS,
				Arch: runtime.GOARCH,
			}

			table, err = internal.GetReleaseLatest(repo, os)
			if err != nil {
				panicRed(err)
			}

			for key := range table {
				repo.Version = key
			}

			path, err := internal.DownloadRelease(table, repo)
			if err != nil {
				panicRed(err)
			}

			destPath := "/opt/homebrew/Cellar/" + repo.GetName() + "/" + repo.GetVersion()
			err = internal.ExtractTarGz(path, destPath, repo.GetName())
			if err != nil {
				panicRed(err)
			}

			src := "../Cellar/oops/" + repo.GetVersion() + "/bin/" + repo.GetName()
			dest := "/opt/homebrew/bin/" + repo.GetName()
			internal.CreateSymbolicLink(src, dest)

		},
	}
)

func init() {
	installCommand.Flags().StringP("token", "t", "", "Accessible token for the private repository.")
	viper.BindPFlag("install-token", installCommand.Flag("token"))

	rootCmd.AddCommand(installCommand)
}
