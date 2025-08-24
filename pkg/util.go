package pkg

import (
	"flag"
	"fmt"
	"log/slog"
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

func Flags() (string, string, string) {
	srcPath := flag.String("src", "", "Source directory path")
	dstPath := flag.String("dst", "", "Destination directory path")
	log := flag.String("log", "", "Log path")

	// TODO: implement me: validate := flag.Bool("validate", false, "Enable SHA256 validation after copy operation")
	// TODO: implement me: verbose := flag.Bool("verbose", false, "Enable verbose output")

	flag.Parse()

	if *srcPath == "" {
		fmt.Println("source path must be provided")
		flag.Usage()
		return "", "", ""
	}

	if *dstPath == "" {
		*dstPath = strings.Join([]string{strings.TrimSuffix(*srcPath, "/"), "_cp"}, "")
		slog.Warn("destination path is not set by user", "auto-set destination path as", *dstPath)
	}

	return *srcPath, *dstPath, *log
	//if *verbose {
	//	fmt.Println("Source:", *srcPath)
	//	fmt.Println("Destination:", *dstPath)
	//	fmt.Println("Validation:", *validate)
	//}

}
