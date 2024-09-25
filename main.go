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
	flag.BoolVar(&ignoreHidden, "ignore-hidden", false, "Ignore hidden files and directories")
}

func main() {
	flag.Parse()
	root := "."

	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	if err := walkDir(root, "", true); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func walkDir(path, prefix string, isRoot bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Filter out hidden files and directories if the flag is set
	if ignoreHidden {
		var filtered []os.DirEntry
		for _, entry := range entries {
			if entry.Name()[0] != '.' {
				filtered = append(filtered, entry)
			}
		}
		entries = filtered
	}

	// Sort entries alphabetically
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for i, entry := range entries {
		connector := "├── "
		newPrefix := prefix + "│   "
		if i == len(entries)-1 {
			connector = "└── "
			newPrefix = prefix + "    "
		}

		fmt.Println(prefix + connector + entry.Name())

		if entry.IsDir() {
			if err := walkDir(filepath.Join(path, entry.Name()), newPrefix, false); err != nil {
				return err
			}
		}
	}
	return nil
}
