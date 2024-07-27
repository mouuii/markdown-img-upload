package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// visit 函数用于处理每个访问到的文件或目录
func visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err) // 打印错误信息
		return err
	}
	if !info.IsDir() && filepath.Ext(path) == ".md" {
		fmt.Println("找到文件:", info.Name(), " 开始处理")
		newFunction(path, info.Name())
	}
	return nil
}

func main() {
	// 指定Markdown文件路径

	var dir string
	flag.StringVar(&dir, "dir", "./", "文件夹")

	err := filepath.Walk(dir, visit)
	if err != nil {
		fmt.Println("Error walking the path:", err)
		return
	}
}

func newFunction(markdownFilePath string, fileName string) bool {
	absolutePathPrefix := "https://raw.githubusercontent.com/cit965/k8s-learning/main"

	file, err := os.Open(markdownFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return true
	}
	defer file.Close()

	outputFile, err := os.OpenFile(fileName+"_generate.md", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return true
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, `<img src="`) {

			line = strings.Replace(line, `src="../..`, fmt.Sprintf(`src="%s`, absolutePathPrefix), 1)
		}

		outputFile.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return false
}
