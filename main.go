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

	s := pkg.NewStorage()
	if err := o.CreateSubdirs(o.Flags.DstPath); err != nil {
		panic(err)
	}

	var err error
	o.CsvHandler, err = pkg.NewCSVLogger(o.Flags.LogPath)
	if err != nil {
		panic(err)
	}

	extensions, subDirCount, err := o.ProcessDir(o.Flags.SrcPath, false)
	if err != nil {
		panic(err)
	}

	slog.Info("", "unique extension count", extensions)
	slog.Info("", "sub-dir count", subDirCount)
	slog.Info("", "skipped file count", len(s.Unprocessed))
	if len(s.Unprocessed) > 0 {
		for _, unprocessedFileName := range s.Unprocessed {
			slog.Warn("", "skipped", unprocessedFileName)
		}
	}
	slog.Info("", "total runtime", time.Since(startTime))

	if o.CsvHandler != nil {
		if err := o.CsvHandler.Log(time.Since(startTime).String(), "skipped file count", "total runtime", strconv.Itoa(len(s.Unprocessed))); err != nil {
			slog.Error("Failed to log:", err)
		}
	}
}
