package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/resoul/travel/internal/infrastructure/cache"
	"github.com/spf13/cobra"
)

func newCacheCmd(c *cache.FileCache, cacheDir string) *cobra.Command {
	cacheCmd := &cobra.Command{
		Use:   "cache",
		Short: "Manage local file cache",
	}

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all cached responses",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.Clear(); err != nil {
				return fmt.Errorf("failed to clear cache: %w", err)
			}
			fmt.Println("Cache cleared successfully.")
			return nil
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show cache directory and entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := os.ReadDir(cacheDir)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("Cache dir: %s (empty / not created yet)\n", cacheDir)
					return nil
				}
				return err
			}

			var count int
			var totalSize int64
			for _, e := range entries {
				if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
					count++
					if info, err := e.Info(); err == nil {
						totalSize += info.Size()
					}
				}
			}

			fmt.Printf("Cache Directory: %s\n", cacheDir)
			fmt.Printf("Cached Files:    %d\n", count)
			fmt.Printf("Total Size:      %.2f KB\n", float64(totalSize)/1024.0)
			return nil
		},
	}

	cacheCmd.AddCommand(clearCmd)
	cacheCmd.AddCommand(infoCmd)

	return cacheCmd
}
