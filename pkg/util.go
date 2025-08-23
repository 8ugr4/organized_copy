package pkg

import (
	"flag"
	"fmt"
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

	// Basic validation
	if *srcPath == "" || *dstPath == "" {
		fmt.Println("Error: both -src and -dst must be provided")
		flag.Usage()
		return "", "", ""
	}

	return *srcPath, *dstPath, *log
	//if *verbose {
	//	fmt.Println("Source:", *srcPath)
	//	fmt.Println("Destination:", *dstPath)
	//	fmt.Println("Validation:", *validate)
	//}

}
