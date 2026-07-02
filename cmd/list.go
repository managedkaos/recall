package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mjenkins/recall/internal/config"
	"github.com/mjenkins/recall/internal/frontmatter"
	"github.com/mjenkins/recall/internal/storage"
	"github.com/spf13/cobra"
)

var tagFlag string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all recall files",
	Long:  `List all recall files in the recall directory, sorted alphabetically. Use --tag to filter by tag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := config.RecallDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := config.EnsureDir(dir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		names, err := storage.List(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "recall: cannot read directory %s: %v\n", dir, err)
			os.Exit(1)
		}

		if tagFlag == "" {
			for _, name := range names {
				fmt.Println(name)
			}
			return nil
		}

		// Filter by tag
		for _, name := range names {
			content, err := storage.Read(dir, name)
			if err != nil {
				continue
			}
			tags, _ := frontmatter.Parse(content)
			for _, t := range tags {
				if strings.EqualFold(t, tagFlag) {
					fmt.Println(name)
					break
				}
			}
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(&tagFlag, "tag", "", "Filter files by tag (case-insensitive)")
	rootCmd.AddCommand(listCmd)
}
