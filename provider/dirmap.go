// Copyright (c) HashiCorp, Inc.

package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func buildMap(basePath string, filter string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore hidden files and directories
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			// Apply filter if provided
			if filter != "" {
				matched, err := filepath.Match(filter, info.Name())
				if err != nil {
					return err
				}
				if !matched {
					return nil
				}
			}

			// Read file content
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var data interface{}
			ext := filepath.Ext(path)

			switch ext {
			case ".yaml", ".yml":
				err = yaml.Unmarshal(content, &data)
			case ".json":
				err = json.Unmarshal(content, &data)
			default:
				// Unsupported file type, ignore
				return nil
			}

			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", path, err)
			}

			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}

			pathParts := strings.Split(relPath, string(filepath.Separator))
			fileName := pathParts[len(pathParts)-1]
			key := strings.TrimSuffix(fileName, ext)
			dirs := pathParts[:len(pathParts)-1]

			currMap := result
			for _, dir := range dirs {
				if _, ok := currMap[dir]; !ok {
					currMap[dir] = make(map[string]interface{})
				}
				currMap = currMap[dir].(map[string]interface{})
			}

			if _, ok := currMap[key]; ok {
				return fmt.Errorf("unique key violation: %s", relPath)
			}

			currMap[key] = data
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
