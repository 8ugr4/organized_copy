package main

import (
	"backup_categorizer/pkg"
	"log/slog"
	"strconv"
	"time"
)

func main() {
	startTime := time.Now()
	srcPath, dstPath, logPath := pkg.Flags()

	if err := pkg.ValidateDir(srcPath); err != nil {
		panic(err)
	}

	s := pkg.NewStorage()
	if err := s.CreateSubdirs(dstPath); err != nil {
		panic(err)
	}

	csvHandler, err := pkg.NewCSVLogger(logPath)
	if err != nil {
		panic(err)
	}

	extensions, subDirCount, err := s.ProcessDir(srcPath, dstPath, csvHandler, false)
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

	if err := csvHandler.Log(time.Since(startTime).String(), "skipped file count", "total runtime", strconv.Itoa(len(s.Unprocessed))); err != nil {
		slog.Error("Failed to log:", err)
	}
}
