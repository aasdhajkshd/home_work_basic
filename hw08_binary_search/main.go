package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func trimSpaceCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return bytes.TrimSpace(data[0 : len(data)-1])
	}
	return bytes.TrimSpace(data)
}

func scanComma(data []byte, atEOF bool) (advance int, token []byte, err error) {
	sepIndex := bytes.IndexByte(data, ',')
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if sepIndex > 0 {
		return sepIndex + 1, trimSpaceCR(data[:sepIndex]), nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, trimSpaceCR(data[0:i]), nil
	}
	if atEOF {
		return len(data), trimSpaceCR(data), nil
	}
	return 0, trimSpaceCR(data), bufio.ErrFinalToken
}

func readFile(filePath string) ([][]byte, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	d := make([]byte, base64.StdEncoding.DecodedLen(len(fileData)))
	n, err := base64.StdEncoding.Decode(d, fileData)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	wordsList := [][]byte{}
	s := bufio.NewScanner(bytes.NewReader(d[:n]))
	s.Buffer(d, bufio.MaxScanTokenSize)
	s.Split(scanComma)
	for s.Scan() {
		wordsList = append(wordsList, s.Bytes())
	}
	// fmt.Printf("%s", wordsList)
	err = s.Err()
	return wordsList, err
}

func findMatch(wordStr string, sliceBytes [][]byte) bool {
	for _, j := range sliceBytes {
		if bytes.Equal([]byte(strings.ToLower(wordStr)), bytes.ToLower(j)) {
			return true
		}
	}
	return false
}

func main() {
	byteWords, _ := readFile("list_test.enc")
	var strInput string
	fmt.Printf("Введите название города России, попробуем найти его в нашем списке... \n> ")
	s := bufio.NewScanner(os.Stdin)
	if s.Scan() {
		strInput = s.Text()
	}
	if len(strInput) > 0 {
		fmt.Printf("Ваше слово: %s", strInput)
		if findMatch(strInput, byteWords) {
			fmt.Println("... eсть совпадение")
		} else {
			fmt.Println("... не найдено")
		}
	}
}
