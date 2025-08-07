package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"triple-s/flags"
	"triple-s/info"
)

func writeFileWithHeader(filePath string, header []string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return err
	}

	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}
	if fileStat.Size() == 0 {
		_, err := file.WriteString(strings.Join(header, ",") + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func InitDir() error {
	err := DirValidation()
	if err != nil {
		return err
	}
	err = os.MkdirAll(flags.Dir, 0o766)
	if err != nil {
		return err
	}
	return writeFileWithHeader(filepath.Join(flags.Dir, "buckets.csv"), info.BucketsHeader)
}

func DirValidation() error {
	prohibited := []string{"flags", "info", "routes", "storage", "utils", "home", "user", "triple-s"}
	for _, dir := range prohibited {
		if strings.Contains(flags.Dir, dir) {
			return fmt.Errorf("directory mustn't contain standard packages")
		}
	}

	if flags.Dir == "." || flags.Dir == "./" {
		return fmt.Errorf("directory mustn't contain standard packages")
	}

	return nil
}

func InitObjectFile(bucketName string) error {
	pathBucket := filepath.Join(flags.Dir, bucketName)
	err := os.MkdirAll(pathBucket, 0o666)
	if err != nil {
		return err
	}
	return writeFileWithHeader(filepath.Join(pathBucket, "objects.csv"), info.ObjectsHeader)
}
