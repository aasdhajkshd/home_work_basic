package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	haproxy "github.com/chrishoffman/haproxylog"
	"github.com/pkg/errors"
)

type logLevel int

const (
	_ logLevel = iota
	Error
	Warn
	Info
	Trace
)

type LogLevel struct {
	Value logLevel
}

func (l *LogLevel) setLevel(s string) {
	s = strings.ToLower(s)
	switch s {
	case "error":
		l.Value = Error
	case "warn":
		l.Value = Warn
	case "info":
		l.Value = Info
	default:
		l.Value = Trace
	}
}

func readLogFile(p string) ([]string, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open the log file")
	}
	defer f.Close()

	var rawLog []string

	b := bufio.NewScanner(f)
	for b.Scan() {
		rawLog = append(rawLog, b.Text())
	}

	if err = b.Err(); err != nil {
		return nil, errors.Wrap(err, "Failed to read the file content")
	}
	return rawLog, nil
}

func outputLogFile(r []string, p string) error {
	f, err := os.Create(p)
	if err != nil {
		return errors.Wrap(err, "Failed to create the output file")
	}
	defer f.Close()
	b := bufio.NewWriter(f)

	defer b.Flush()
	for _, j := range r {
		_, err = b.WriteString(j)
		if err != nil {
			return errors.Wrap(err, "Failed to write to the output file")
		}
	}
	return nil
}

func parseLog(r []string, l logLevel) ([]string, error) {
	var (
		parsedLog   = make([]string, 0, 1000)
		statusCodes = make(map[int64]int, 100)
	)

	for _, j := range r {
		line, err := haproxy.NewLog(j)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse the log line in file")
		}
		switch l {
		case Warn:
			if line.HTTPStatusCode >= 300 && line.HTTPStatusCode < 400 {
				parsedLog = append(parsedLog, fmt.Sprintln(line.AcceptDate, line.ClientIP, line.FrontendName, line.HTTPRequest.URL))
			}
		case Error:
			if line.HTTPStatusCode >= 400 {
				parsedLog = append(parsedLog, fmt.Sprintln(line.AcceptDate, line.ClientIP, line.FrontendName, line.HTTPRequest.URL))
			}
		case Info:
			if line.HTTPStatusCode > 0 {
				statusCodes[line.HTTPStatusCode]++
			}
		case Trace:
			if line.Message != "" {
				parsedLog = append(parsedLog, fmt.Sprintln(line.AcceptDate, line.ClientIP, line.FrontendName, line.Message))
			}
		default:
			if line.HTTPStatusCode > 100 {
				parsedLog = append(parsedLog, fmt.Sprintln(line.AcceptDate, line.ClientIP, line.FrontendName, line.BackendName, line.Tr, line.HTTPRequest.URL)) //nolint:lll
			}
		}
	}
	for i, j := range statusCodes {
		parsedLog = append(parsedLog, fmt.Sprintf("Status code %d: %s: %d\n", i, http.StatusText(int(i)), j))
	}
	return parsedLog, nil
}

func main() {
	var filePath, levelStr, output string

	env := map[string]string{"file": "LOG_ANALYZER_FILE", "level": "LOG_ANALYZER_LEVEL", "output": "LOG_ANALYZER_OUTPUT"}

	verbose := flag.Bool("verbose", false, "verbose output")
	flag.StringVar(&filePath, "file", os.Getenv(env["file"]), "path to log file to parse, mandotary")
	flag.StringVar(&levelStr, "level", os.Getenv(env["level"]), "log level \"error\", \"warn\", \"trace\" for details, defaults to statistics \"info\" level") //nolint:lll
	flag.StringVar(&output, "output", os.Getenv(env["output"]), "path to file with output parsed log")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: [LOG_ANALYZER_FILE=example.log|LOG_ANALYZER_LEVEL=error|LOG_ANALYZER_OUTPUT=output.log] %s [OPTIONS]\n", os.Args[0]) //nolint:lll
		fmt.Fprintf(os.Stdout, "Options:\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stdout, "  -%s: %s\n", f.Name, f.Usage)
		})
	}
	flag.Parse()

	if filePath == "" {
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()

	if levelStr == "" {
		levelStr = "info"
	}

	level := &LogLevel{}
	level.setLevel(levelStr)

	if *verbose {
		for i, j := range os.Args {
			fmt.Println("parsed arguments:", i, j)
		}
		fmt.Println("file:", filePath)
		fmt.Println("level:", levelStr)
		fmt.Println("output:", output)
	}

	outputLog, err := readLogFile(filePath)
	if err != nil {
		errors.Errorf("whoops: %s", err)
	}

	parsedLog, _ := parseLog(outputLog, level.Value)
	if len(output) > 0 {
		outputLogFile(parsedLog, "output.log")
	} else {
		fmt.Println()
		for _, j := range parsedLog {
			fmt.Printf("%s", j)
		}
	}
}
