package cmd

import (
	"bytes"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/treeder/bump/internal/bump"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type (
	Options struct {
		Filename   string
		Input      string
		Part       string
		Extract    bool
		Format     string
		Replace    string
		Index      int
		PreRelease string
		Metadata   string
	}
)

var (
	options = &Options{}
	rootCmd = &cobra.Command{
		Use:   "bump",
		Short: "bump it dawg! See https://github.com/treeder/bump for more info.",
		Long:  `bump it, dawg! This tool will bump a version in a file. It will find the first semver string in the file and bump it.`,
		Run:   bumper,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&options.Filename, "filename", "f", "VERSION", "filename to look for version in")
	rootCmd.Flags().StringVarP(&options.Input, "input", "i", "", "use this if you want to pass in a string to pass, rather than read it from a file. Cannot be used with --filename.")
	rootCmd.Flags().StringVarP(&options.Part, "part", "p", "", "part to bump, either major, minor or patch. Default is patch.")
	rootCmd.Flags().BoolVarP(&options.Extract, "extract", "e", false, "this will just find the version and return it, does not modify anything. Safe operation.")
	rootCmd.Flags().StringVarP(&options.Format, "format", "F", "", "either M for major, M-m for major-minor or M-m-p")
	rootCmd.Flags().StringVarP(&options.Replace, "replace", "r", "", "overwrites the version with what you pass in here")
	rootCmd.Flags().IntVarP(&options.Index, "index", "x", 0, "if zero (default), uses first match. If greater than zero, uses nth match. If less than zero, starts at last match and goes backwards, ie: last match is -1.")
	rootCmd.Flags().StringVarP(&options.PreRelease, "prerelease", "P", "", "adds a prerelease tag to the version per semver spec")
	rootCmd.Flags().StringVarP(&options.Metadata, "metadata", "m", "", "adds metadata to the version per semver spec")
}

func bumper(c *cobra.Command, args []string) {
	arg := "patch"
	// fmt.Println("ARGS:", c.Args())
	if len(args) < 1 {
		// log.Fatal("Invalid arg")
	} else {
		arg = args[0]
		arg = strings.ToLower(arg)
	}

	// check for `[bump X]` in input, user can pass in git commit messages to auto bump different versions
	if strings.Contains(arg, "[bump minor]") {
		arg = "minor"
	} else if strings.Contains(arg, "[bump major]") {
		arg = "major"
	}

	var err error
	var vbytes []byte
	filename := options.Filename

	if options.Input != "" {
		vbytes = []byte(options.Input)
	} else {
		vbytes, err = os.ReadFile(filename)
		if err != nil {
			if os.IsNotExist(err) {
				err = fmt.Errorf("%v not found. Use either --filename or --input to change where to look for version", filename)
			}

			c.Println(err)
		}
	}

	bumpOptions := &bump.Options{
		Index:      options.Index,
		PreRelease: options.PreRelease,
		Metadata:   options.Metadata,
	}

	if options.Part != "" {
		bumpOptions.Part = options.Part
		c.Println("part is set", bumpOptions.Part)
	} else {
		bumpOptions.Part = arg
	}

	var oldContent, newContent string
	var newcontent []byte

	if options.Replace != "" {
		bumpOptions.Replace = options.Replace
		bumpOptions.Part = ""
		oldContent, newContent, _, newcontent, err = bump.ReplaceInContent2(vbytes, bumpOptions)
	} else {
		oldContent, newContent, _, newcontent, err = bump.BumpInContent2(vbytes, bumpOptions)
	}

	if err != nil {
		c.Println(err)
	}

	if options.Extract {
		printCommand(c, oldContent)
		return
	}

	c.Println("Old version:", oldContent)
	c.Println("New version:", newContent)

	if options.Input == "" {
		err = os.WriteFile(filename, newcontent, 0644)
		if err != nil {
			c.Println(err)
		}
	}

	printCommand(c, newContent) // write it to stdout so scripts can use it
}

func printCommand(c *cobra.Command, version string) {
	if options.Format == "" {
		fmt.Print(version)
		return
	}

	v := semver.New(version)
	// else, we format it
	var b bytes.Buffer
	for _, char := range options.Format {
		if char == 'M' {
			b.WriteString(strconv.FormatInt(v.Major, 10))
		} else if char == 'm' {
			b.WriteString(strconv.FormatInt(v.Minor, 10))
		} else if char == 'p' {
			b.WriteString(strconv.FormatInt(v.Patch, 10))
		} else {
			b.WriteRune(char)
		}
	}

	fmt.Print(b.String())
}
