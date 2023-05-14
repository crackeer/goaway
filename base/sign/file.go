package sign

import (
	"os"
	"path/filepath"
	"strings"
)

// GetSignFileFromDir
//
//	@param dir
//	@return error
func GetSignCodeFromDir(dir string) (map[string]string, error) {
	list, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	retData := map[string]string{}
	for _, name := range list {
		parts := strings.Split(name.Name(), ".")
		bytes, err := os.ReadFile(filepath.Join(dir, name.Name()))
		if err != nil {
			return nil, err
		}
		retData[parts[0]] = string(bytes)
	}
	return retData, nil
}
