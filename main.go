package main

import (
	"backup_categorizer/pkg"
	"time"
)

func main() {
	startTime := time.Now()

	o := pkg.GetNewOperator()
	o.Flags = pkg.GetFlags()

	if err := pkg.ValidateDir(o.Flags.SrcPath); err != nil {
		panic(err)
	}

	if err := o.CreateSubdirs(o.Flags.DstPath); err != nil {
		panic(err)
	}

	var err error
	o.CsvHandler, err = pkg.NewCSVLogger(o.Flags.LogPath)
	if err != nil {
		panic(err)
	}

	extensions, err := o.Operate()
	if err != nil {
		panic(err)
	}

	pkg.ResultLog(extensions, o, startTime)
}
