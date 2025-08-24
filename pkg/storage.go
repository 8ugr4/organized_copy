package pkg

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	images       = "images"
	videos       = "videos"
	audios       = "audios"
	archives     = "archives"
	documents    = "documents"
	applications = "applications"
	unknown      = "unknown"
)

type Storage struct {
	images       []string
	videos       []string
	audios       []string
	archives     []string
	documents    []string
	applications []string
	unknown      []string
	Unprocessed  []string
}

func NewStorage() *Storage {
	return &Storage{
		images:       make([]string, 0),
		videos:       make([]string, 0),
		audios:       make([]string, 0),
		archives:     make([]string, 0),
		documents:    make([]string, 0),
		applications: make([]string, 0),
		unknown:      make([]string, 0),
		Unprocessed:  make([]string, 0),
	}
}

func (s Storage) AddType(ext, fp string) string {
	switch ext {
	// IMAGES
	case "jpg", "JPG", "jpeg", "png", "webp", "jfif", "HEIC", "svg", "PNG":
		s.images = append(s.images, fp)
		return images

	// VIDEOS
	case "mp4", "gif", "mpeg", "ogg":
		s.videos = append(s.videos, fp)
		return videos

	// AUDIOS
	case "wav", "asd", "mp3", "aac", "aif":
		s.audios = append(s.audios, fp)
		return audios

	// DOCUMENTS
	case "pdf", "PDF", "doc", "docx", "dotx",
		"txt", "epub", "csv", "pptx", "accdb",
		"xlsx", "bib", "sql", "json", "rtf",
		"tex", "ini", "odt":
		s.documents = append(s.documents, fp)
		return documents

	// ARCHIVES
	case "zip", "rar", "pcapng", "msix", "iso":
		s.archives = append(s.archives, fp)
		return archives

		// APPLICATIONS
	case "ipynb", "m", "exe", "py", "whl", "pcap", "msi":
		s.applications = append(s.applications, fp)
		return applications
	// UNKNOWN
	case "unknown", "rdf", "mdl", "sig", "hbs", "dat", "pkpass", "tmp", " ":
		s.unknown = append(s.unknown, fp)
		return unknown
	default:
		return unknown
	}
}

func check(dst string) error {
	if _, err := os.Stat(dst); err != nil {
		if err := os.Mkdir(dst, syscall.O_CREAT|syscall.O_EXCL|syscall.O_WRONLY); err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) CreateSubdirs(dstBasePath string) error {
	if err := check(dstBasePath); err != nil {
		return err
	}
	dirNames := []string{"images", "videos", "audios", "archives", "documents", "applications", "unknown"}
	for _, dirName := range dirNames {
		if err := os.Mkdir(path.Join(dstBasePath, dirName), syscall.O_CREAT|syscall.O_EXCL|syscall.O_WRONLY); err != nil {
			return err
		}
	}
	return nil
}
func uniqueDstPath(dstBasePath, dstDir, baseName string) string {
	ext := filepath.Ext(baseName)
	base := strings.TrimSuffix(baseName, ext)
	dstNewPath := path.Join(dstBasePath, dstDir, baseName)
	i := 1
	for {
		if _, err := os.Stat(dstNewPath); os.IsNotExist(err) {
			break
		}
		dstNewPath = path.Join(dstBasePath, dstDir, fmt.Sprintf("%s_%d%s", base, i, ext))
		i++
	}
	return dstNewPath
}

func (s Storage) Copy(dstPath, dstDir, fileAbsolutePath string, csvLogger *CSVLogger) error {
	srcFile, err := os.Open(fileAbsolutePath)
	if err != nil {
		slog.Warn("Skipping unreadable file", "path", fileAbsolutePath, "error", err)
		s.Unprocessed = append(s.Unprocessed, fileAbsolutePath)
		return nil
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close:%s:%w", fileAbsolutePath, err))
		}
	}()

	_, fileName := path.Split(fileAbsolutePath)
	destinationFile, err := os.Create(uniqueDstPath(dstPath, dstDir, fileName))
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
	if err := csvLogger.Log("SUCCESS", srcFile.Name(), fileName, destinationFile.Name()); err != nil {
		slog.Error("Failed to log:", err)
	}
	return nil
}

func (s Storage) ProcessDir(srcPath, dstPath string, csvHandler *CSVLogger, r bool) (int, int, error) {
	entries, err := os.ReadDir(srcPath)
	if err != nil {
		return 0, 0, err
	}
	slog.Info("", "entry count:", len(entries))

	total := len(entries)
	processed := 0
	subDirCount := 0
	extensions := make([]string, 0)
	for _, entry := range entries {
		fp := path.Join(srcPath, entry.Name())
		if entry.IsDir() {
			subDirCount++
			if _, _, err := s.ProcessDir(fp, dstPath, csvHandler, true); err != nil {
				return 0, 0, err
			}
			continue
		}

		info, err := os.Stat(fp)
		if err != nil {
			slog.Warn("Skipping blocked file", "path", fp, "error", err)
			s.Unprocessed = append(s.Unprocessed, fp)
			continue
		}
		if !info.Mode().IsRegular() || info.Size() == 0 {
			s.Unprocessed = append(s.Unprocessed, fp)
			continue
		}

		kind := path.Ext(fp)
		ext := ""
		if kind != "" {
			ext = kind[1:]
		}

		typeDir := s.AddType(ext, fp)
		if err := s.Copy(dstPath, typeDir, fp, csvHandler); err != nil {
			return 0, 0, err
		}
		processed++
		percentage := float64(processed) / float64(total) * 100
		if processed%max(1, total/20) == 0 && r == false {
			slog.Info("progress", "completed", fmt.Sprintf("%.1f%%", percentage))
		}
		extensions = append(extensions, ext)
	}
	extensions = RemoveDuplicateStr(extensions)

	return subDirCount, len(extensions), nil
}
