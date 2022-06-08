package remotefile

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localfile"
)

func readBody(resp http.Response) ([]byte, *airer.Airer) {
	// Create a buffer, read 32 bytes at a time
	byte_buffer := make([]byte, 32)
	file_contents := make([]byte, 0)
	for {
		bytes_read, err := resp.Body.Read(byte_buffer)
		if bytes_read > 0 {
			file_contents = append(file_contents, byte_buffer[:bytes_read]...)
		}
		if err != nil {
			if err != io.EOF {
				return nil, &airer.Airer{
					Kind:    airer.ExecError,
					Message: fmt.Sprintf("Failed to read body:\n%s", err),
					Origin:  err,
				}
			} else {
				break
			}
		}
	}
	return file_contents, nil
}

func DownloadTemp(url string) (string, *airer.Airer) {
	f, err := os.CreateTemp("", "kind_config.yaml")
	if err != nil {
		return "", &airer.Airer{
			Kind:    airer.ShellError,
			Message: "Could not create tmp file!",
			Origin:  err,
		}
	}
	filename := f.Name()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, arr := readBody(*resp)
	if arr != nil {
		return "", arr
	}
	localfile.OverwriteFile(filename, data)
	return filename, nil
}
