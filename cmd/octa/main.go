package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"octa/irgen"
	"octa/lexer"
	"octa/parser"
)

func main() {
	// -------------------------------
	// 命令行参数解析
	// -------------------------------
	compileFlag := flag.Bool("c", false, "compile to .oo but do not link")
	linkFlag := flag.Bool("l", false, "link .oo files to executable")
	outputFlag := flag.String("o", "", "output file name")
	flag.Parse()

	inputFiles := flag.Args()
	if len(inputFiles) == 0 {
		fmt.Println("Usage: ./octa [-c] [-l] [-o output] input files")
		os.Exit(1)
	}

	// -------------------------------
	// 缺省模式：直接从 .octa 编译链接
	// -------------------------------
	if !*compileFlag && !*linkFlag {
		compileFlag = new(bool)
		*compileFlag = true
		linkFlag = new(bool)
		*linkFlag = true
	}

	// -------------------------------
	// 编译阶段
	// -------------------------------
	var objectFiles []string
	if *compileFlag {
		for _, infile := range inputFiles {
			// 读取源码
			data, err := ioutil.ReadFile(infile)
			if err != nil {
				panic(err)
			}

			// Lexer
			l := lexer.New(string(data))
			var tokens []lexer.Token
			for tok := l.Next(); tok.Type != lexer.TokenEOF; tok = l.Next() {
				tokens = append(tokens, tok)
			}

			// Parser
			f := parser.Parse(tokens)

			// IRGen
			ir := irgen.NewIRGen()
			ir.GenerateFunctions(f)

			// 输出目标文件
			outFile := *outputFlag
			if outFile == "" {
				base := filepath.Base(infile)
				outFile = base[:len(base)-len(filepath.Ext(base))] + ".oo"
			}
			ir.WriteObject(outFile)
			fmt.Println("Compiled", infile, "to", outFile)
			objectFiles = append(objectFiles, outFile)
		}
	}

	// -------------------------------
	// 链接阶段
	// -------------------------------
	if *linkFlag {
		var exeName string
		if *outputFlag != "" {
			exeName = *outputFlag
		} else {
			// 默认生成第一个输入文件名去掉扩展名
			base := filepath.Base(inputFiles[0])
			exeName = base[:len(base)-len(filepath.Ext(base))]
		}

		err := irgen.Link(objectFiles, exeName)
		if err != nil {
			panic(err)
		}
		fmt.Println("Linked", objectFiles, "to executable", exeName)
	}
}
