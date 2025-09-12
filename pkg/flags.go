package pkg

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Flags struct {
	SrcPath  string
	DstPath  string
	RulePath string
	LogPath  string
	DryRun   bool
	Async    bool
	Verbose  bool
}

func GetFlags() Flags {
	srcPath := flag.String("src", "./testDir", "Source directory path")
	dstPath := flag.String("dst", "", "Destination directory path")
	rulePath := flag.String("rules", "./rules.yaml", "output category rules")
	log := flag.String("log", "", "Log path")
	dryRun := flag.Bool("dry-run", false, "Dry-run option")
	async := flag.Bool("async", true, "Faster async option, uses goroutines")
	verbose := flag.Bool("verbose", false, "Set to debug mode")
	// TODO: implement me: validate := flag.Bool("validate", false, "Enable SHA256 validation after copy operation")

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

	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if *rulePath == "" {
		slog.Error("path for rules file is empty, going to use default settings from 'rules.yaml'")
		*rulePath = "./rules.yaml"
	}

	return Flags{
		SrcPath:  *srcPath,
		DstPath:  *dstPath,
		LogPath:  *log,
		DryRun:   *dryRun,
		Async:    *async,
		Verbose:  *verbose,
		RulePath: *rulePath,
	}
}
