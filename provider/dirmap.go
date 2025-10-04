// Copyright (c) HashiCorp, Inc.

package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// convertKeysToStrings recursively converts map[interface{}]interface{} to map[string]interface{}.
func convertKeysToStrings(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := make(map[string]interface{})
		for k, v := range x {
			m2[fmt.Sprint(k)] = convertKeysToStrings(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convertKeysToStrings(v)
		}
		return x
	}
	return i
}

// buildMap traverses the directory at path and builds a nested map from YAML/JSON files.
func buildMap(basePath string, filter string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path %s: %w", path, err)
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
					return fmt.Errorf("invalid filter pattern %s: %w", filter, err)
				}
				if !matched {
					return nil
				}
			}

			// Read file content
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			var data interface{}
			ext := filepath.Ext(path)

			switch ext {
			case ".yaml", ".yml":
				decoder := yaml.NewDecoder(strings.NewReader(string(content)))
				if err := decoder.Decode(&data); err != nil {
					return fmt.Errorf("failed to parse YAML file %s: %w", path, err)
				}
			case ".json":
				if err := json.Unmarshal(content, &data); err != nil {
					return fmt.Errorf("failed to parse JSON file %s: %w", path, err)
				}
			default:
				// Unsupported file type, ignore
				return nil
			}

			data = convertKeysToStrings(data)

			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
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
				var ok bool
				currMap, ok = currMap[dir].(map[string]interface{})
				if !ok {
					return fmt.Errorf("type mismatch in nested structure for directory %s", dir)
				}
			}

			if _, ok := currMap[key]; ok {
				return fmt.Errorf("unique key violation for file %s", relPath)
			}

			currMap[key] = data
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to traverse directory %s: %w", basePath, err)
	}

	return result, nil
}
