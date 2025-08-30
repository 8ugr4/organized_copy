package main

import (
	"backup_categorizer/pkg"
	"log/slog"
	"strconv"
	"time"
)

func main() {
	startTime := time.Now()
	o := pkg.GetNewOperator()

	o.Flags = pkg.GetFlags()

	if err := pkg.ValidateDir(o.Flags.SrcPath); err != nil {
		panic(err)
	}

	if !o.Flags.DryRun {
		if err := o.CreateSubdirs(o.Flags.DstPath); err != nil {
			panic(err)
		}
	}

	var err error
	o.CsvHandler, err = pkg.NewCSVLogger(o.Flags.LogPath)
	if err != nil {
		panic(err)
	}

	var extensions int
	extensions, err = o.Operate()

	slog.Debug("", "unique extension count", extensions)
	slog.Debug("", "sub-dir count", o.SubDirCount)
	slog.Debug("", "skipped file count", len(o.Storage.Unprocessed))
	if len(o.Storage.Unprocessed) > 0 {
		for _, unprocessedFileName := range o.Storage.Unprocessed {
			slog.Warn("", "skipped", unprocessedFileName)
		}
	}
	slog.Info("", "total runtime", time.Since(startTime))

	if o.CsvHandler != nil {
		if err := o.CsvHandler.Log(time.Since(startTime).String(), "skipped file count", "total runtime", strconv.Itoa(len(o.Storage.Unprocessed))); err != nil {
			slog.Error("Failed to log:", err)
		}
	}
}
