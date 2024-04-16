package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

type color string

const (
	ColorBlack  color  = "\u001b[30m"
	ColorRed    color  = "\u001b[31m"
	ColorGreen  color  = "\u001b[32m"
	ColorYellow color  = "\u001b[33m"
	ColorBlue   color  = "\u001b[34m"
	ColorReset  color  = "\u001b[0m"
	InfoStr     string = "info"
)

type logLevel int

const (
	_ logLevel = iota
	ERROR
	WARN
	INFO
	TRACE
	EVENT
)

type LogLevel struct {
	Value logLevel
}

func (l *LogLevel) setLevel(s string) {
	s = strings.ToLower(s)
	switch s {
	case "error", "err":
		l.Value = ERROR
	case "warning", "warn":
		l.Value = WARN
	case InfoStr:
		l.Value = INFO
	case "trace", "debug":
		l.Value = TRACE
	}
}

func (l LogLevel) String() string {
	switch l.Value {
	case ERROR:
		return "error"
	case WARN:
		return "warn"
	case INFO:
		return InfoStr
	case TRACE:
		return "trace"
	case EVENT:
		return "event"
	default:
		return "unknown"
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

func countLogLevels(r []string, l LogLevel) int {
	c := 0
	for _, j := range r {
		f := func(c rune) bool {
			return unicode.IsSpace(c) || c == ':' // нюанс лог файла из примера "WARNING:.."
		}
		w := strings.FieldsFunc(j, f)[4]

		if l.String() == strings.ToLower(w) {
			c++
		}
	}
	return c
}

func main() {
	var (
		filePath, levelStr, output string
		outputText                 []string
	)

	env := map[string]string{"file": "LOG_ANALYZER_FILE", "level": "LOG_ANALYZER_LEVEL", "output": "LOG_ANALYZER_OUTPUT"}

	verbose := flag.Bool("verbose", false, "verbose output")
	useColor := flag.Bool("color", false, "display colorized output")
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
		levelStr = InfoStr
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
		log.Fatalf("whoops, check the file name: %s", err)
	}

	result := countLogLevels(outputLog, *level)

	outputText = append(outputText, fmt.Sprintf("Found total messages with log level "+
		"\"%s\": %d of lines %d in file \"%s\"",
		level.String(), result, len(outputLog), filePath))

	if len(output) > 0 {
		outputLogFile(outputText, output)
	} else {
		if *useColor {
			print(string(ColorBlue))
		}
		for _, j := range outputText {
			fmt.Printf("%s\n", j)
		}
	}
	print(string(ColorReset))
}
