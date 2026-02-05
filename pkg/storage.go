package pkg

import (
	"fmt"
	"github.com/barasher/go-exiftool"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	unknown = "unknown"
)

type Storage struct {
	Entries        []os.DirEntry
	Categories     map[string][]string // [categories][]extensions
	Extensions     map[string]string   // [extensions][categories]
	OutDirectories map[string][]string // []categories[files]
	SubDirs        map[string][]string // [subDir][]extensions
	Unprocessed    []string
	SortMap        map[string]string //image:year, videos:month, documents:month
	Exif           *exiftool.Exiftool
}

func NewStorage() *Storage {
	return &Storage{
		Categories:     make(map[string][]string),
		Extensions:     make(map[string]string),
		OutDirectories: make(map[string][]string),
		SubDirs:        make(map[string][]string),
		Unprocessed:    make([]string, 0),
	}
}

type Operator struct {
	Storage        Storage
	Flags          Flags
	CsvHandler     *CSVLogger
	SubDirCount    int
	ExtensionCount int
	sem            chan struct{}
	once           sync.Once
	mu             sync.Mutex
}

func (o *Operator) initPool(n int) {
	o.once.Do(func() {
		if n <= 0 {
			n = 8
		}
		o.sem = make(chan struct{}, n)
	})
}

func GetNewOperator() *Operator {
	o := &Operator{
		Storage:        *NewStorage(),
		Flags:          Flags{},
		CsvHandler:     nil,
		SubDirCount:    0,
		ExtensionCount: 0,
		sem:            nil,
		once:           sync.Once{},
		mu:             sync.Mutex{},
	}
	o.initPool(8)
	return o
}

func (o *Operator) BuildStorageMaps(c *Config) {
	for _, rule := range c.Rules {
		o.Storage.Categories[rule.Category] = make([]string, 0)
		for _, extension := range rule.Extensions {
			o.Storage.Categories[rule.Category] = append(o.Storage.Categories[rule.Category], extension)
			o.Storage.Extensions[extension] = rule.Category
		}
		if rule.SeparateExists() {
			o.Storage.SubDirs[rule.Category] = append(o.Storage.SubDirs[rule.Category], rule.Separate...)
			if err := o.initExifTool(); err != nil {
				panic(err)
			}
		}
	}
}

func (o *Operator) GetSeparateSubdirs(category, ext string) string {
	if subdirs, exists := o.Storage.SubDirs[category]; exists {
		for _, sub := range subdirs {
			if sub == ext {
				return sub
			}
		}
		return ""
	}
	return ""
}

func (o *Operator) GetSortSubDirs(category string) (string, bool) {
	if sortType, exists := o.Storage.SortMap[category]; exists {
		return sortType, true
	}
	return "", false
}

func (o *Operator) GetExtensionCategory(ext string) (string, bool) {
	if val, ok := o.Storage.Extensions[ext]; ok {
		return val, true
	}
	return unknown, false
}

// AddType adds and returns category of the file
func (o *Operator) AddType(ext, fp string) string {
	category, exists := o.GetExtensionCategory(ext)
	if !exists {
		slog.Warn("unknown extension, doesn't match to rules", "extension", ext)
		slog.Warn("copying to the unknown dir", "filepath", fp)
		return unknown
	}
	o.Storage.OutDirectories[category] = append(o.Storage.OutDirectories[category], fp)
	return category
}

func check(dst string) error {
	if _, err := os.Stat(dst); err != nil {
		if err := os.Mkdir(dst, syscall.O_CREAT|syscall.O_EXCL|syscall.O_WRONLY); err != nil {
			return err
		}
	}
	return nil
}

func (o *Operator) CreateSubdirs(dstBasePath string, rules []Rule) error {
	if o.Flags.DryRun {
		return nil
	}

	if err := check(dstBasePath); err != nil {
		return err
	}

	for _, rule := range rules {
		if err := os.Mkdir(path.Join(dstBasePath, rule.Category), syscall.O_CREAT|syscall.O_EXCL|syscall.O_WRONLY); err != nil {
			return err
		}
		if rule.SeparateExists() {
			for _, separateDir := range rule.Separate {
				if err := os.Mkdir(path.Join(dstBasePath, rule.Category, separateDir), syscall.O_CREAT|syscall.O_EXCL|syscall.O_WRONLY); err != nil {
					return err
				}
			}
		}
	}
	// TODO: drop me
	// dirNames := []string{"images", "videos", "audios", "archives", "documents", "applications", "unknown"}
	// for _, dirName := range dirNames {
	//	 if err := os.Mkdir(path.Join(dstBasePath, dirName), syscall.O_CREAT|syscall.O_EXCL|syscall.O_WRONLY); err != nil {
	//		return err
	//	}
	// }
	return nil
}

func uniqueDstPath(dstBasePath, dstDir, specialDir, baseName string) string {
	ext := filepath.Ext(baseName)
	base := strings.TrimSuffix(baseName, ext)
	dstNewPath := path.Join(dstBasePath, dstDir, baseName)
	if specialDir != "" {
		dstNewPath = path.Join(dstBasePath, dstDir, specialDir, baseName)
	}

	// TODO: improve this following idiotic logic
	original := dstNewPath
	i := 1
	for {
		if _, err := os.Stat(dstNewPath); err != nil {
			if os.IsNotExist(err) {
				break
			}
			slog.Error("stat call failed during trying to create a unique destination path", "PATH:", dstNewPath)
			panic(err)
		}
		if specialDir == "" {
			dstNewPath = path.Join(path.Dir(original), fmt.Sprintf("%s_%d%s", base, i, ext))
		} else {
			dstNewPath = path.Join(path.Dir(original), fmt.Sprintf("%s_%d%s", base, i, ext))
		}
		i++
	}
	return dstNewPath
}

func (o *Operator) initExifTool() error {
	exifTool, err := exiftool.NewExiftool()
	if err != nil {
		return err
	}
	o.Storage.Exif = exifTool
	return nil
}

// getFileDate tries EXIF -> ModTime -> regex from filename and returns either month or year as string
// periodType is "month" or "year"
func (o *Operator) getFileDate(fp, regexPattern, periodType string) (string, error) { //nolint:unused
	f, err := os.Open(fp)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	exifTool := o.Storage.Exif

	var timePeriod string
	fileInfos := exifTool.ExtractMetadata(fp)
	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			return "", fileInfo.Err
		}
		if date, exists := fileInfo.Fields["CreateDate"]; exists {
			timePeriod = date.(string) //FIXME
		}
	}

	//FIXME
	if timePeriod == "" {
		re := regexp.MustCompile(regexPattern)
		matches := re.FindStringSubmatch(fp)
		if len(matches) > 1 {
			dateStr := matches[1]

			layouts := []string{
				"20060102_150405", // YEAR MONTH DAY_HOUR MINUTE SECOND
				time.DateTime,     // "2006-01-02 15:04:05"
				"20060102",        // YEAR MONTH DAY
				time.DateOnly,     // "YEAR-MONTH-DAY"
			}
			for _, layout := range layouts {
				if t, err := time.Parse(layout, dateStr); err == nil {
					timePeriod = t.String()
					break
				}
			}
		}
	}
	parseTime := func(timePeriod, periodType string) string {
		//FIXME
		switch periodType {
		case "month":
			return fmt.Sprintf("something %s", timePeriod)
		case "year":
			return fmt.Sprintf("something %s", timePeriod)
		default:
			return ""
		}
	}(timePeriod, periodType)
	if parseTime == "" {
		return "", fmt.Errorf("invalid periodType %s, must be 'month' or 'year'", periodType)
	}
	return parseTime, nil
}

func (o *Operator) Copy(dstPath, dstDir, specialDir, fileAbsolutePath string) error {
	srcFile, err := os.Open(fileAbsolutePath)
	if err != nil {
		slog.Warn("Skipping unreadable file", "path", fileAbsolutePath, "error", err)
		// o.Storage.Unprocessed = append(o.Storage.Unprocessed, fileAbsolutePath)
		return nil
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close:%s:%w", fileAbsolutePath, err))
		}
	}()

	_, fileName := path.Split(fileAbsolutePath)
	destinationFile, err := os.Create(uniqueDstPath(dstPath, dstDir, specialDir, fileName))
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func(destinationFile *os.File) {
		err := destinationFile.Close()
		if err != nil {
			panic(err)
		}
	}(destinationFile)

	_, err = io.Copy(destinationFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy %s file to %s: %w", srcFile.Name(), destinationFile.Name(), err)
	}

	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file:%s:%w", destinationFile.Name(), err)
	}

	if o.CsvHandler != nil {
		if err := o.CsvHandler.Log("SUCCESS", srcFile.Name(), fileName, destinationFile.Name()); err != nil {
			slog.Error("failure-log", "error", err.Error())
		}
	}

	return nil
}

// skipcheck logs skipped files and adds them to unprocessed slice.
func (o *Operator) skipcheck(fp string) bool {
	info, err := os.Stat(fp)
	if err != nil {
		slog.Warn("Skipping blocked file", "path", fp, "error", err)
		o.Storage.Unprocessed = append(o.Storage.Unprocessed, fp)
		return true
	}
	if !info.Mode().IsRegular() {
		slog.Warn("Skipping blocked file", "path", fp, "error", "isn't a regular file")
		o.Storage.Unprocessed = append(o.Storage.Unprocessed, fp)
		return true
	}

	if info.Size() == 0 {
		slog.Warn("Skipping blocked file", "path", fp, "error", "has size 0")
		o.Storage.Unprocessed = append(o.Storage.Unprocessed, fp)
		return true
	}
	return false
}

func (o *Operator) AsyncProcessDir(dirpath string, r bool) (int, error) {
	entries, err := os.ReadDir(dirpath)
	if err != nil {
		return 0, err
	}
	slog.Debug("", "entry count:", len(entries))
	if o.Flags.DryRun {
		os.Exit(1)
	}
	total := len(entries)
	processed := int64(0)
	extensions := make([]string, 0)
	var extMutex, unprocMutex = sync.Mutex{}, sync.Mutex{}
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup

	for _, entry := range entries {
		fp := path.Join(dirpath, entry.Name())
		if entry.IsDir() {
			o.SubDirCount++
			if _, err := o.AsyncProcessDir(fp, true); err != nil {
				return 0, err
			}
			continue
		}
		if o.skipcheck(fp) {
			continue
		}

		kind := path.Ext(fp)
		ext := ""
		if kind != "" {
			ext = kind[1:]
		}

		typeDir := o.AddType(ext, fp)

		wg.Add(1)
		sem <- struct{}{} // get slot
		go func(fp, typeDir string, ext string) {
			defer wg.Done()
			defer func() { <-sem }() // release slot
			specialSubDir := o.GetSeparateSubdirs(typeDir, ext)
			if err := o.Copy(o.Flags.DstPath, typeDir, specialSubDir, fp); err != nil {
				unprocMutex.Lock()
				o.Storage.Unprocessed = append(o.Storage.Unprocessed, fp)
				unprocMutex.Unlock()
				return
			}

			atomic.AddInt64(&processed, 1)
			if !r {
				pct := float64(atomic.LoadInt64(&processed)) / float64(total) * 100
				if atomic.LoadInt64(&processed)%int64(max(1, total/20)) == 0 {
					slog.Info("progress", "completed", fmt.Sprintf("%.1f%%", pct))
				}
			}
			extMutex.Lock()
			extensions = append(extensions, ext)
			extMutex.Unlock()
		}(fp, typeDir, ext)
	}
	wg.Wait()
	extensions = RemoveDuplicateStr(extensions)

	return len(extensions), nil
}

func (o *Operator) ProcessDir(dirpath string, r bool) (int, error) {
	entries, err := os.ReadDir(dirpath)
	if err != nil {
		return 0, err
	}
	slog.Info("", "entry count:", len(entries))
	if o.Flags.DryRun {
		os.Exit(1)
	}

	total := len(entries)
	processed := 0
	subDirCount := 0
	extensions := make([]string, 0)
	for _, entry := range entries {
		fp := path.Join(dirpath, entry.Name())
		if entry.IsDir() {
			subDirCount++
			if _, err := o.ProcessDir(fp, true); err != nil {
				return 0, err
			}
			continue
		}
		if o.skipcheck(fp) {
			continue
		}

		kind := path.Ext(fp)
		ext := ""
		if kind != "" {
			ext = kind[1:]
		}

		typeDir := o.AddType(ext, fp)
		specialSubDir := o.GetSeparateSubdirs(typeDir, ext)
		//sortDir, exists := o.GetSortSubDirs(typeDir)

		if err := o.Copy(o.Flags.DstPath, typeDir, specialSubDir, fp); err != nil {
			return 0, err
		}
		processed++
		percentage := float64(processed) / float64(total) * 100
		if processed%max(1, total/20) == 0 && !r {
			slog.Info("progress", "completed", fmt.Sprintf("%.1f%%", percentage))
		}
		extensions = append(extensions, ext)
	}
	extensions = RemoveDuplicateStr(extensions)

	return len(extensions), nil
}

func (o *Operator) Operate() (int, error) {
	switch o.Flags.Async {
	case true:
		return o.AsyncProcessDir(o.Flags.SrcPath, false)
	case false:
		return o.ProcessDir(o.Flags.SrcPath, false)
	}
	return 0, nil
}
