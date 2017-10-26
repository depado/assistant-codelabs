package utils

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/sirupsen/logrus"

	"github.com/Depado/assistant-codelabs/database"
	"github.com/Depado/assistant-codelabs/models"
	"github.com/pkg/errors"
)

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// ImportFromFile takes the path to a CSV file and parses it
func ImportFromFile(f string) error {
	var err error
	var fd *os.File

	if fd, err = os.Open(f); err != nil {
		return errors.Wrap(err, "ImportFromFile(): unable to open file")
	}

	return ImportCSV(fd)
}

// ImportCSV reads the CSV file
func ImportCSV(r *os.File) error {
	var err error
	var line []string
	var skipped, imported, count int
	var c *models.Car

	logrus.Info("Import action started")

	start := time.Now()

	if count, err = lineCounter(r); err != nil {
		return errors.Wrap(err, "ImportCSV(): unable to count lines")
	}
	r.Seek(0, 0)
	bar := pb.New(count)
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true
	bar.Start()

	reader := csv.NewReader(bufio.NewReader(r))
	reader.Comma = ';'

	tx, err := database.DB.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for {
		line, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if c, err = models.NewFromCSV(line); err != nil {
			skipped++
			logrus.WithError(err).Debug("Couldn't parse the line, skipping")
		} else {
			imported++
			if err = tx.Save(c); err != nil {
				logrus.WithError(err).Debug("Couldn't write to bolt, skipping")
			}
			bar.Increment()
		}
	}
	bar.Finish()
	logrus.Info("Committing transactions")
	if err = tx.Commit(); err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{"imported": imported, "skipped": skipped, "took": time.Since(start)}).Info("Done importing")
	return nil
}
