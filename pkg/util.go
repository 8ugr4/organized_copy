package pkg

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

type Flags struct {
	SrcPath string
	DstPath string
	LogPath string
	DryRun  bool
}

func GetFlags() Flags {
	srcPath := flag.String("src", "", "Source directory path")
	dstPath := flag.String("dst", "", "Destination directory path")
	log := flag.String("log", "", "Log path")
	dryRun := flag.Bool("dry-run", false, "Dry-run option")

	// TODO: implement me: validate := flag.Bool("validate", false, "Enable SHA256 validation after copy operation")
	// TODO: implement me: verbose := flag.Bool("verbose", false, "Enable verbose output")

	flag.Parse()

	if *srcPath == "" {
		fmt.Println("source path must be provided")
		flag.Usage()
		os.Exit(1)
	}

	if *dstPath == "" {
		*dstPath = strings.Join([]string{strings.TrimSuffix(*srcPath, "/"), "_cp"}, "")
		slog.Warn("destination path is not set by user", "auto-set destination path as", *dstPath)
	}

	return Flags{
		SrcPath: *srcPath,
		DstPath: *dstPath,
		LogPath: *log,
		DryRun:  *dryRun,
	}
}
