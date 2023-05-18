package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	sourcePath, _ := os.Getwd()

	// ignorePattern := []string{"tools/node_modules/commander/**", "!tools/node_modules/commander/package-support.json"}
	// ignorePattern := []string{"tools/**", "!tools/node_modules/commander/package-support.json"}
	// ignorePattern := []string{"tools/node_modules/commander/**", "!tools/node_modules/commander/lib/**"}
	// ignorePattern := []string{"tools/node_modules/commander/*", "!tools/node_modules/commander/lib/**"}
	ignorePattern := []string{"tools/node_modules/**", "!tools/node_modules/commander/lib/**"}

	MBTIgnoreTest(sourcePath, ignorePattern)
}

func getIgnorePatterns(sourcePath string, ignorePattern []string) ([]string, error) {
	var ignorePaths []string
	for _, ign := range ignorePattern {
		ignore_path := filepath.ToSlash(ign)
		ignorePaths = append(ignorePaths, ignore_path)
		fmt.Printf("ignore path: %s \n", ignore_path)
		fmt.Println()
	}
	return ignorePaths, nil
}

func MBTIgnoreTest(sourcePath string, ignorePattern []string) error {
	ignorePatterns, err := getIgnorePatterns(sourcePath, ignorePattern)
	if err != nil {
		return err
	}

	return filepath.Walk(sourcePath, func(fullpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// get relative file path and slash
		filePath, err := filepath.Rel(sourcePath, fullpath)
		if err != nil {
			return err
		}
		filePath = filepath.ToSlash(filePath)

		// Print file and pattern
		fmt.Printf("file: %s \n", filePath)
		fmt.Printf("ignore patterns: %v \n", ignorePatterns)

		var cmdArgs []string
		cmdArgs = append(cmdArgs, "ismatch")
		cmdArgs = append(cmdArgs, "-f")
		cmdArgs = append(cmdArgs, filePath)
		cmdArgs = append(cmdArgs, "-p")
		cmdArgs = append(cmdArgs, ignorePatterns...)

		out, err := exec.Command("micromatch-wrapper-win.exe", cmdArgs...).Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return err
		}
		fmt.Println(string(out))
		fmt.Println()

		return nil
	})
}

/*
func MBTIgnoreTest(ignore []string, sourcePath string) error {
	ignoreMap, err := getIgnoredEntries(ignore, sourcePath)
	if err != nil {
		return err
	}

	err = walk(sourcePath, ignoreMap)
	return err
}

func getIgnoredEntries(ignore []string, sourcePath string) (map[string]interface{}, error) {
	regularSourcePath := sourcePath
	ignoredEntriesMap := map[string]interface{}{}
	for _, ign := range ignore {
		fmt.Printf("ignore in MTA: %s \n", ign)
		path := filepath.Join(regularSourcePath, ign)
		ignoredEntriesMap[path] = nil
	}

	fmt.Println("ignore entry map:")
	for key, value := range ignoredEntriesMap {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}

	return ignoredEntriesMap, nil
}

func walk(sourcePath string, ignore map[string]interface{}) error {
	// fmt.Printf("walking source path: %s \n", sourcePath)

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("walking path: %s \n", path)

		for ignore_pattern, _ := range ignore {
			path = filepath.ToSlash(path)
			ignore_pattern = filepath.ToSlash(ignore_pattern)
			result, err := exec.Command("node", "ignoreparser.js", "-f "+path, "-p "+ignore_pattern).Output()
			if err != nil {
				fmt.Println("Error executing command:", err)
				return err
			}

			res := strings.ReplaceAll(string(result), "\n", "")
			if res == "true" {
				fmt.Printf("path '%s' is match ignore pattern '%s', it will be skipped \n", path, ignore_pattern)
				if info.IsDir() {
					// fmt.Printf("path '%s' is match ignore pattern '%s', and it is an dir, it will be skipped \n", path, ignore_pattern)
					return filepath.SkipDir
				}
				return nil
			}
			// fmt.Printf("path '%s' match pattern '%s' result is '%s' \n", path, ignore_pattern, res)
			// fmt.Printf("path '%s' is NOT match ignore pattern '%s' \n", path, ignore_pattern)
		}

		return nil
	})
}
*/

/* func getIgnoredEntries(ignore []string, sourcePath string) (map[string]interface{}, error) {
	regularSourcePath := sourcePath
	ignoredEntriesMap := map[string]interface{}{}
	for _, ign := range ignore {
		fmt.Printf("ignore in MTA: %s \n", ign)

		path := filepath.Join(regularSourcePath, ign)
		entries, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			fmt.Printf("ignore entry: %s \n", entry)
			ignoredEntriesMap[entry] = nil
		}
	}

	fmt.Println("ignore entry map:")
	for key, value := range ignoredEntriesMap {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}

	return ignoredEntriesMap, nil
}

func walk(sourcePath string, ignore map[string]interface{}) error {
	fmt.Printf("walking source path: %s \n", sourcePath)

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("walking path: %s \n", path)

		if _, ok := ignore[path]; ok {
			fmt.Printf("path %s is in ignore map \n", path)
			if info.IsDir() {
				fmt.Printf("path %s is in ignore map, because it is an dir, it is skipped \n", path)
				return filepath.SkipDir
			}
			return nil
		}

		return nil
	})
} */
