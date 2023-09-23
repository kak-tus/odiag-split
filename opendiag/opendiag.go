package opendiag

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	nameFormat = "appLog-2006-01-02-15-04-05.log"
	timeFormat = "Time:	15:04:05,000"
	maxEntries = 850
)

type Entries []Entry

type Entry struct {
	Time    time.Time
	Send    string
	Receive string
}

type Log struct {
	Header        string
	Entries       Entries
	FileCreatedAt time.Time
}

func SupportedFileName(name string) bool {
	return filepath.Ext(name) == ".log"
}

// Assume that timezone at device where log was created is same as
// timezone at current code executor
func DateFromFileName(name string) (time.Time, error) {
	parsed, err := time.ParseInLocation(nameFormat, name, time.Local)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

func Decode(fileCreatedAt time.Time, data string) (Log, error) {
	splitted := strings.Split(data, "\n")

	var (
		headers                   []string
		entries                   Entries
		isHeader                  bool = true
		isTime, isSend, isReceive bool
		entry                     Entry
	)

	for _, line := range splitted {
		switch {
		case strings.HasPrefix(line, "Time:	"):
			if !isHeader && !isReceive {
				return Log{}, fmt.Errorf("can't parse, time not after header or receive: '%s'", line)
			}

			isHeader = false

			if isReceive {
				isReceive = false

				entries = append(entries, entry)

				entry = Entry{}
			}

			parsed, err := time.Parse(timeFormat, line)
			if err != nil {
				return Log{}, err
			}

			entryTime := time.Date(
				fileCreatedAt.Year(), fileCreatedAt.Month(), fileCreatedAt.Day(),
				parsed.Hour(), parsed.Minute(), parsed.Second(),
				parsed.Nanosecond(), fileCreatedAt.Location(),
			)

			entry.Time = entryTime
			isTime = true
		case strings.HasPrefix(line, "Send:	"):
			if !isTime {
				return Log{}, fmt.Errorf("can't parse, send not after time: '%s'", line)
			}

			entry.Send = line
			isTime = false
			isSend = true
		case strings.HasPrefix(line, "Receive: "):
			if !isSend {
				return Log{}, fmt.Errorf("can't parse, receive not after send: '%s'", line)
			}

			entry.Receive = line
			isSend = false
			isReceive = true
		default:
			// Additional line
			switch {
			case isHeader:
				headers = append(headers, line)
			case isTime:
				// Strange but may be possible
			case isSend:
				entry.Send += "\n" + line
			case isReceive:
				entry.Receive += "\n" + line
			}
		}
	}

	if isReceive {
		entries = append(entries, entry)
	}

	res := Log{
		Header:  strings.Join(headers, "\n"),
		Entries: entries,
	}

	return res, nil
}

func (hdl Log) Encode() (fileName string, data string) {
	data = hdl.Header + "\n"

	for _, entry := range hdl.Entries {
		data += entry.Time.Format(timeFormat) + "\n" + entry.Send + "\n" + entry.Receive + "\n"
	}

	fileDate := hdl.FileCreatedAt
	if len(hdl.Entries) != 0 {
		fileDate = hdl.Entries[0].Time
	}

	return fileDate.Format(nameFormat), data
}

func (hdl Log) NeedSplit() bool {
	return len(hdl.Entries) > maxEntries
}

func (hdl Log) Split() []Log {
	if len(hdl.Entries) <= maxEntries {
		return []Log{hdl}
	}

	// https://github.com/golang/go/wiki/SliceTricks#batching-with-minimal-allocation
	batchSize := maxEntries
	batches := make([]Log, 0, (len(hdl.Entries)+batchSize-1)/batchSize)

	for batchSize < len(hdl.Entries) {
		val := hdl
		val.Entries = hdl.Entries[0:batchSize:batchSize]

		hdl.Entries, batches = hdl.Entries[batchSize:], append(batches, val)
	}

	batches = append(batches, hdl)

	return batches
}
