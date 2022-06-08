package localfile

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mcdonaldseanp/charlie/airer"
)

func ReadJSONFile(location string, data interface{}) *airer.Airer {
	json_blob, airr := ReadFileInChunks(location)
	if airr != nil {
		return airr
	}

	json.Unmarshal(json_blob, &data)
	return nil
}

func OverwriteJSONFile(location string, data interface{}) *airer.Airer {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to open file:\n%s", err),
			Origin:  err,
		}
	}
	defer f.Close()
	json_string, err := json.Marshal(data)
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to marshal json:\n%s", err),
			Origin:  err,
		}
	}
	_, err = f.Write([]byte(json_string))
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to write to file:\n%s", err),
			Origin:  err,
		}
	}
	return nil
}
