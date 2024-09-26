package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

var ignoreHidden bool

func init() {
	// Set up the command-line flag for ignoring hidden files
	flag.BoolVar(&ignoreHidden, "ignore-hidden", false, "Ignore hidden files and directories")
}

func main() {
	// Parse command-line arguments and flags
	flag.Parse()
	root := "."

	// If a path is passed as an argument, set it as the root
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	// Walk through the directory tree and print it
	if err := walkDir(root, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func walkDir(path, prefix string) error {
	// Read directory contents
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Filter out hidden files and directories if the ignoreHidden flag is set
	if ignoreHidden {
		var filtered []os.DirEntry
		for _, entry := range entries {
			if entry.Name()[0] != '.' { // Ignore files starting with a dot
				filtered = append(filtered, entry)
			}
		}
		entries = filtered
	}

	// Sort entries alphabetically
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	// Print directory contents with tree-like formatting
	for i, entry := range entries {
		connector := "├── "
		newPrefix := prefix + "│   "
		if i == len(entries)-1 {
			connector = "└── "
			newPrefix = prefix + "    "
		}

		fmt.Println(prefix + connector + entry.Name())

		// Recursively walk subdirectories
		if entry.IsDir() {
			if err := walkDir(filepath.Join(path, entry.Name()), newPrefix); err != nil {
				return err
			}
		}
	}
	return nil
}
