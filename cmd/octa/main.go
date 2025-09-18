package main

import (
	"flag"
	"fmt"
	"os"

	"octa/irgen"
	"octa/lexer"
	"octa/parser"
)

func main() {
	compile := flag.Bool("c", false, "generate object file (.oo) only")
	link := flag.Bool("l", false, "link object files to executable")
	output := flag.String("o", "", "output file name")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Usage: ./octa [options] file1.octa ...")
		flag.PrintDefaults()
		return
	}

	doCompile := *compile || (!*compile && !*link)
	doLink := *link || (!*compile && !*link)

	objFiles := []string{}

	if doCompile {
		for _, fname := range files {
			src, err := os.ReadFile(fname)
			if err != nil {
				fmt.Println("Failed to read file:", fname)
				return
			}

			tokens := lexer.Lex(string(src))
			funcStmt := parser.Parse(tokens)

			ir := irgen.NewIRGen()
			ir.GenerateFunctions(funcStmt)

			outFile := *output
			if outFile == "" {
				outFile = fname[:len(fname)-5] + ".bc"
			}
			ir.WriteObject(outFile)
			fmt.Printf("Compiled %s to %s\n", fname, outFile)
			objFiles = append(objFiles, outFile)
		}
	} else {
		objFiles = files
	}

	if doLink {
		exeName := *output
		if exeName == "" && len(objFiles) > 0 {
			exeName = objFiles[0][:len(objFiles[0])-3]
		}
		err := irgen.Link(objFiles, exeName)
		if err != nil {
			fmt.Println("Linking failed:", err)
			return
		}
		fmt.Printf("Linked executable: %s\n", exeName)
	}
}
