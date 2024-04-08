package cmd

import (
	"github.com/spf13/cobra"
	"github.com/treeder/bump/internal/bump"
	"os"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Bump version using git tags and commit history",
	Run:   bumperGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

//set -e
//
//git fetch --tags # checkout action does not get these
//
//# git describe has issues with GitHub Actions: https://github.com/treeder/firetils/commit/160ef4560d8855c9c05f4cae207baeb71b7791f3/checks?check_suite_id=414542684
//# oldv=$(git describe --match "v[0-9]*" --abbrev=0 HEAD)
//# This new way seems to work better and avoids the issue above:
//# -v:refname is a version sort
//oldv=$(git tag --sort=-v:refname --list "v[0-9]*" | head -n 1)
//echo "oldv: $oldv"
//
//# if there is no version tag yet, let's start at 0.0.0
//if [ -z "$oldv" ]; then
//   echo "No existing version, starting at 0.0.0"
//   oldv="0.0.0"
//fi
//
//newv=$(docker run --rm -v "$PWD":/app treeder/bump --input "$oldv" patch)
//echo "newv: $newv"
//
//git tag -a "v$newv" -m "version $newv"
//git push --follow-tags
//echo "done"

func bumperGenerate(cmd *cobra.Command, args []string) {
	oldv := "0.0.0"
	newBump := bump.NewBump(oldv, "patch")
	cmd.Println("newv:", newBump.String())

	if err := newBump.Tag(); err != nil {
		cmd.Println("Error tagging:", err)
		os.Exit(1)
	}

	cmd.Println("done")
}
