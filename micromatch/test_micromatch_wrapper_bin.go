package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func test_micromatch_bin_help() {
	fmt.Println("Testing micromatch-wrapper-win!!")

	out, err := exec.Command("./micromatch-wrapper-win.exe", "-h").Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	fmt.Println(string(out))
}

func test_micromatch_match(files, patterns []string) {
	// Print file and pattern before slash
	fmt.Printf("Files: [%s]\n", strings.Join(files, ","))
	fmt.Printf("Patterns: [%s]\n", strings.Join(patterns, ","))

	// slash file and pattern and print
	var files_toslash []string
	for _, file := range files {
		file_toslash := filepath.ToSlash(file)
		files_toslash = append(files_toslash, file_toslash)
		// fmt.Printf("File to slash: %s\n", file_toslash)
	}

	var patterns_toslash []string
	for _, pattern := range patterns {
		pattern_toslash := filepath.ToSlash(pattern)
		patterns_toslash = append(patterns_toslash, pattern_toslash)
		// fmt.Printf("Pattern to slash: %s\n", pattern_toslash)
	}

	var cmdArgs []string
	cmdArgs = append(cmdArgs, "match")
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, files_toslash...)
	cmdArgs = append(cmdArgs, "-p")
	cmdArgs = append(cmdArgs, patterns_toslash...)

	cmd := exec.Command("./micromatch-wrapper-win.exe", cmdArgs...)
	out, err := cmd.Output()
	//exitCode := cmd.ProcessState.ExitCode()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	output := string(out)
	fmt.Println(output)
	//fmt.Println(exitCode)
	fmt.Println()
}

func test_micromatch_ismatch(files, patterns []string) {
	// Print file and pattern before slash
	fmt.Printf("Files: [%s]\n", strings.Join(files, ","))
	fmt.Printf("Patterns: [%s]\n", strings.Join(patterns, ","))

	// slash file and pattern and print
	file_toslash := filepath.ToSlash(strings.Join(files, " "))
	// fmt.Printf("File to slash: %s\n", file_toslash)

	var patterns_toslash []string
	for _, pattern := range patterns {
		pattern_toslash := filepath.ToSlash(pattern)
		patterns_toslash = append(patterns_toslash, pattern_toslash)
		// fmt.Printf("Pattern to slash: %s\n", pattern_toslash)
	}

	var cmdArgs []string
	cmdArgs = append(cmdArgs, "ismatch")
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, file_toslash)
	cmdArgs = append(cmdArgs, "-p")
	cmdArgs = append(cmdArgs, patterns_toslash...)

	cmd := exec.Command("./micromatch-wrapper-win.exe", cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	output := string(out)
	fmt.Println(output)
	fmt.Println()
}

func test_micromatch_getPackagedFiles(source, target string, patterns []string) {
	// Print file and pattern before slash
	fmt.Printf("Source: %s\n", source)
	fmt.Printf("Target: %s\n", target)
	fmt.Printf("Patterns: [%s]\n", strings.Join(patterns, ","))

	var patterns_toslash []string
	for _, pattern := range patterns {
		pattern_toslash := filepath.ToSlash(pattern)
		patterns_toslash = append(patterns_toslash, pattern_toslash)
		// fmt.Printf("Pattern to slash: %s\n", pattern_toslash)
	}

	var cmdArgs []string
	cmdArgs = append(cmdArgs, "getPackagedFiles")
	cmdArgs = append(cmdArgs, "-s")
	cmdArgs = append(cmdArgs, source)
	cmdArgs = append(cmdArgs, "-t")
	cmdArgs = append(cmdArgs, target)
	if len(patterns_toslash) > 0 {
		cmdArgs = append(cmdArgs, "-p")
		cmdArgs = append(cmdArgs, patterns_toslash...)
	}

	cmd := exec.Command("./micromatch-wrapper-win.exe", cmdArgs...)
	out, err := cmd.Output()
	//exitCode := cmd.ProcessState.ExitCode()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	output := string(out)
	fmt.Println(output)

	fmt.Println()
}

func main() {
	// Hello
	test_micromatch_bin_help()

	// Test 1
	fmt.Printf("Test %d\n", 1)
	files := []string{"node_modules/braces/lib/parse.js"}
	patterns := []string{"node_modules/braces/**"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 2
	fmt.Printf("Test %d\n", 2)
	files = []string{"node_modules\\braces\\lib\\parse.js"}
	patterns = []string{"node_modules\\braces\\**"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 3
	fmt.Printf("Test %d\n", 3)
	files = []string{"node_modules/commander/lib/help.js"}
	patterns = []string{"node_modules/!(braces)/**"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 4
	fmt.Printf("Test %d\n", 4)
	files = []string{"/node_modules/npm_module/pkg/young.yang03.js"}
	patterns = []string{"/**/pkg/**"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 5
	fmt.Printf("Test %d\n", 5)
	files = []string{"/node_modules/npm_module/pkg/young.yang03.js"}
	patterns = []string{"/**/!(pkg)/**"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 6
	fmt.Printf("Test %d\n", 6)
	files = []string{"/node_modules/npm_module/pkg/young.yang03.js"}
	patterns = []string{"/**/!(pkg)/!(pkg)/**"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 7
	fmt.Printf("Test %d\n", 7)
	files = []string{"/node_modules/npm_module/pkg/young.yang03.js"}
	patterns = []string{"/**/!(pkg)/!(pkg)/!(pkg)/**"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 8
	fmt.Printf("Test %d\n", 8)
	files = []string{"/node_modules/npm_module/notpkg/young.yang03.js"}
	patterns = []string{"/**/!(pkg)/**"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 9
	fmt.Printf("Test %d\n", 9)
	files = []string{"licence.js"}
	patterns = []string{"*.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 10
	fmt.Printf("Test %d\n", 10)
	files = []string{"node_module/licence.js"}
	patterns = []string{"node_module/*"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 11
	fmt.Printf("Test %d\n", 11)
	files = []string{"node_module/licence.js"}
	patterns = []string{"node_module/**"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 12
	fmt.Printf("Test %d\n", 12)
	files = []string{"node-js\\gulpfile.js"}
	patterns = []string{"node-js\\*.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 13
	fmt.Printf("Test %d\n", 13)
	files = []string{"C:\\Workspace\\Test_Project\\go_glob_patter_test\\test_micromatch\\node-js\\gulpfile.js"}
	patterns = []string{"C:\\Workspace\\Test_Project\\go_glob_patter_test\\test_micromatch\\node-js\\*.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 14
	fmt.Printf("Test %d\n", 14)
	files = []string{"node_module/licence.js"}
	patterns = []string{"node_module/?icence.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 15
	fmt.Printf("Test %d\n", 15)
	files = []string{"node_module/braces/lib"}
	patterns = []string{"node_module/b[r,e]aces/*"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 16
	fmt.Printf("Test %d\n", 16)
	files = []string{"node_modules/@sap/excluded_nodejs_packages"}
	patterns = []string{"node_modules/@sap/!(included_nodejs_packages)"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 17
	fmt.Printf("Test %d\n", 17)
	files = []string{"node_modules/@sap/included_nodejs_packages"}
	patterns = []string{"node_modules/@sap/!(included_nodejs_packages)"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 18
	fmt.Printf("Test %d\n", 18)
	files = []string{"node_modules/@sap/excluded_nodejs_packages"}
	patterns = []string{"!node_modules/@sap/included_nodejs_packages"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 19
	fmt.Printf("Test %d\n", 19)
	files = []string{"node_modules/@sap/included_nodejs_packages"}
	patterns = []string{"!node_modules/@sap/included_nodejs_packages"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 20
	fmt.Printf("Test %d\n", 20)
	files = []string{"node_modules/@sap/included_nodejs_packages"}
	patterns = []string{"!node_modules/@sap/included_nodejs_packages"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 21
	fmt.Printf("Test %d\n", 21)
	files = []string{"node_modules/@sap/included_nodejs_packages"}
	patterns = []string{"node_modules/@sap/**", "!node_modules/@sap/included_nodejs_packages", "!node_modules/@sap/excluded_nodejs_packages"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 22
	fmt.Printf("Test %d\n", 22)
	files = []string{"a/b/3.js"}
	patterns = []string{"a/b/**", "!a/b/3.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 23
	fmt.Printf("Test %d\n", 23)
	files = []string{"a/b/3.js", "a/b/4.js", "a/b/5.js", "a/b/6.js"}
	patterns = []string{"a/b/**", "!a/b/3.js"}
	test_micromatch_match(files, patterns)
	files = []string{"a/b/3.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/4.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/5.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/6.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 24
	fmt.Printf("Test %d\n", 24)
	files = []string{"a/b/3.js"}
	patterns = []string{"a/b/**", "!a/b/3.js", "!a/b/4.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 25
	fmt.Printf("Test %d\n", 25)
	files = []string{"a/b/3.js", "a/b/4.js", "a/b/5.js", "a/b/6.js"}
	patterns = []string{"a/b/**", "!a/b/3.js", "!a/b/4.js"}
	test_micromatch_match(files, patterns)
	files = []string{"a/b/3.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/4.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/5.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/6.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 26
	fmt.Printf("Test %d\n", 26)
	files = []string{"a/b/3.js", "a/b/4.js", "a/b/5.js", "a/b/6.js"}
	patterns = []string{"a/b/!(3.js)", "a/b/!(4.js)"}
	test_micromatch_match(files, patterns)
	files = []string{"a/b/3.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/4.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/5.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/6.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 27
	fmt.Printf("Test %d\n", 27)
	files = []string{"a/b/3.js", "a/b/4.js", "a/b/5.js", "a/b/6.js"}
	patterns = []string{"a/b/**", "a/b/!(3.js)", "!a/b/3.js", "!a/b/4.js"}
	test_micromatch_match(files, patterns)
	files = []string{"a/b/3.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/4.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/5.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/6.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 28
	fmt.Printf("Test %d\n", 28)
	files = []string{"a/b/3.js"}
	patterns = []string{"!a/b/3.js", "a/b/**"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 29
	fmt.Printf("Test %d\n", 29)
	files = []string{"a/b/3.js", "a/b/4.js", "a/b/5.js", "a/b/6.js"}
	patterns = []string{"!a/b/4.js", "a/b/**", "!a/b/3.js"}
	test_micromatch_match(files, patterns)
	files = []string{"a/b/3.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/4.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/5.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/6.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 30
	fmt.Printf("Test %d\n", 30)
	files = []string{"a/b/3.js", "a/b/4.js", "a/b/5.js", "a/b/6.js"}
	patterns = []string{"!a/b/3.js", "!a/b/4.js", "a/b/**"}
	test_micromatch_match(files, patterns)
	files = []string{"a/b/3.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/4.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/5.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/6.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 31
	fmt.Printf("Test %d\n", 31)
	files = []string{"a/b/3.js", "a/b/4.js", "a/b/5.js", "a/b/6.js"}
	patterns = []string{"a/b/!(4.js)", "a/b/!(3.js)"}
	test_micromatch_match(files, patterns)
	files = []string{"a/b/3.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/4.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/5.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/6.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 32
	fmt.Printf("Test %d\n", 32)
	files = []string{"a/b/3.js", "a/b/4.js", "a/b/5.js", "a/b/6.js"}
	patterns = []string{"a/b/!(3.js)", "a/b/**", "!a/b/3.js", "!a/b/4.js"}
	test_micromatch_match(files, patterns)
	files = []string{"a/b/3.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/4.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/5.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/6.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 33
	fmt.Printf("Test %d\n", 33)
	files = []string{"a/b/3.js", "a/b/4.js", "a/b/5.js", "a/b/6.js"}
	patterns = []string{"a/b/**", "!a/b/3.js", "!a/b/4.js", "a/b/!(3.js)"}
	test_micromatch_match(files, patterns)
	files = []string{"a/b/3.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/4.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/5.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)
	files = []string{"a/b/6.js"}
	test_micromatch_match(files, patterns)
	test_micromatch_ismatch(files, patterns)

	// Test 34
	fmt.Printf("Test %d\n", 34)
	files = []string{"node_modules/@sap/cds", "node_modules/@sap/credential-store-client-node",
		"node_modules/@sap/external-service-mashup", "node_modules/@sap/low-code-event-handler", "node_modules/@sap/xssec"}
	patterns = []string{"node_modules/@sap/**", "!node_modules/@sap/credential-store-client-node", "!node_modules/@sap/external-service-mashup"}
	test_micromatch_match(files, patterns)

	// Test 35
	fmt.Printf("Test %d\n", 35)
	files = []string{"node_modules/@sap/cds", "node_modules/@sap/credential-store-client-node",
		"node_modules/@sap/external-service-mashup", "node_modules/@sap/low-code-event-handler", "node_modules/@sap/xssec"}
	patterns = []string{"node_modules/@sap/!(credential-store-client-node)", "node_modules/@sap/!(external-service-mashup)"}
	test_micromatch_match(files, patterns)

	// Test 36
	fmt.Printf("Test %d\n", 36)
	files = []string{"node_modules/braces/lib/utils.js"}
	patterns = []string{"node_modules/**", "!node_modules/braces/**"}
	test_micromatch_match(files, patterns)

	// Test 37
	fmt.Printf("Test %d\n", 37)
	source, _ := os.Getwd()
	target := filepath.Join(source, "tmpfile")
	patterns = []string{}
	test_micromatch_getPackagedFiles(source, target, patterns)

	// Test 38
	fmt.Printf("Test %d\n", 38)
	source, _ = os.Getwd()
	target = filepath.Join(source, "tmpfile")
	patterns = []string{"node_modules/**", "!node_modules/lodash/**"}
	test_micromatch_getPackagedFiles(source, target, patterns)
}
