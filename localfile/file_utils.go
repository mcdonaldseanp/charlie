package localfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mcdonaldseanp/charlie/airer"
)

const STDIN_IDENTIFIER string = "__STDIN__"

func readFromStdin() string {
	var builder strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		builder.WriteString(scanner.Text() + "\n")
	}
	return builder.String()
}

func ReadFileOrStdin(maybe_file string) ([]byte, error) {
	var raw_data []byte
	var arr error
	if maybe_file == STDIN_IDENTIFIER {
		raw_data = []byte(readFromStdin())
	} else {
		raw_data, arr = ReadFileInChunks(maybe_file)
		if arr != nil {
			return nil, arr
		}
	}
	return raw_data, nil
}

func ReadFileInChunks(location string) ([]byte, error) {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to open file:\n%s", err),
			Origin:  err,
		}
	}
	defer f.Close()

	// Create a buffer, read 32 bytes at a time
	byte_buffer := make([]byte, 32)
	file_contents := make([]byte, 0)
	for {
		bytes_read, err := f.Read(byte_buffer)
		if bytes_read > 0 {
			file_contents = append(file_contents, byte_buffer[:bytes_read]...)
		}
		if err != nil {
			if err != io.EOF {
				return nil, &airer.Airer{
					Kind:    airer.ExecError,
					Message: fmt.Sprintf("Failed to read file:\n%s", err),
					Origin:  err,
				}
			} else {
				break
			}
		}
	}
	return file_contents, nil
}

func OverwriteFile(location string, data []byte) error {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to open file:\n%s", err),
			Origin:  err,
		}
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to write to file:\n%s", err),
			Origin:  err,
		}
	}
	return nil
}

func TempFile(tmpname string, data []byte) (string, error) {
	f, err := os.CreateTemp("", tmpname)
	if err != nil {
		return "", &airer.Airer{
			Kind:    airer.ShellError,
			Message: "Could not create tmp file!",
			Origin:  err,
		}
	}
	filename := f.Name()
	OverwriteFile(filename, data)
	return filename, nil
}
