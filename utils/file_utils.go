package utils

import (
	"encoding/json"
	"fmt"
	. "github.com/mcdonaldseanp/charlie/airer"
	"io"
	"os"
)

func ReadJSONFile(location string) (map[string]interface{}, *Airer) {
	json_blob, airr := ReadFileInChunks(location)
	if airr != nil {
		return nil, airr
	}

	var data map[string]interface{}
	json.Unmarshal(json_blob, &data)
	return data, nil
}

func OverwriteJSONFile(location string, data interface{}) *Airer {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed to open file:\n%s", err),
			err,
		}
	}
	defer f.Close()
	json_string, err := json.Marshal(data)
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed to marshal json:\n%s", err),
			err,
		}
	}
	_, err = f.Write([]byte(json_string))
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed to write to file:\n%s", err),
			err,
		}
	}
	return nil
}


func ReadFileInChunks(location string) ([]byte, *Airer) {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, &Airer{
			ExecError,
			fmt.Sprintf("Failed to open file:\n%s", err),
			err,
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
				return nil, &Airer{
					ExecError,
					fmt.Sprintf("Failed to read file:\n%s", err),
					err,
				}
			} else {
				break
			}
		}
	}
	return file_contents, nil
}
