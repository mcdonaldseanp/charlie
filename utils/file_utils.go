package utils

import (
	"fmt"
	"os"
	"io"
	"encoding/json"
	. "github.com/McdonaldSeanp/charlie/airer"
)

func ReadJSONFile(location string) (map[string]interface{}, *Airer) {
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
	json_blob := make([]byte, 0)
	for {
		bytes_read, err := f.Read(byte_buffer)
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
		json_blob = append(json_blob, byte_buffer[:bytes_read]...)
	}
	var data map[string]interface{}
	json.Unmarshal(json_blob, &data)
	return data, nil
}

func OverwriteJSONFile(location string, data interface{}) (*Airer) {
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