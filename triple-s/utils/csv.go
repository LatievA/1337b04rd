package utils

import (
	"encoding/csv"
	"os"
)

func ReadCSV(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) > 0 {
		return records[1:], nil
	}
	return [][]string{}, nil
}

func WriteCSV(filepath string, header []string, records [][]string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if header != nil {
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	if err := writer.WriteAll(records); err != nil {
		return err
	}
	return nil
}
