package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

var (
	debugFlag    = flag.Bool("debug", false, "Enable verbose log.")
	pkgFlag		 = flag.String("pkgs", "", "Provide a comma separated list of package roots relative to the GOPATH to include.")
)

const Usage = `instability: instability metrics for go files in terms of exported methods. 
A publicly exported method is counted as an inward connection, and any public methods used 
are counted as outward connections.

Usage:

  go-instability [flags] package

Flags:

`

/**
Instability reads all .go files in a directory and its sub-directories
and calculates the instability of the files. The instability is calculated
using the coupling metric: I = Ce / (Ce + Ca), with I being the instability,
Ca the afferent or inward coupling and Ce the efferent or outward
coupling. The coupling is counted using the number of exported methods
and structs and the number of imported methods and structs from other packages.
 */
func main() {
	flag.Parse()

	if *debugFlag {
		log.SetFlags(log.Lmicroseconds)
	}

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, Usage)
		flag.PrintDefaults()
		os.Exit(2)
	}

	goPkgPath := flag.Arg(0)
	goPath := os.Getenv("GOPATH")
	root := fmt.Sprintf("%s/%s", goPath, goPkgPath)
	pkgsInclude := make([]string, 0)
	pkgsRoots := strings.Split(*pkgFlag, ",")
	for _, pkgRoot := range pkgsRoots {
		path := fmt.Sprintf("%s/%s", goPath, pkgRoot)
		if isDir(path) {
			filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				base := filepath.Base(path)
				if !strings.Contains(filepath.Ext(base), ".") {
					pkgsInclude = append(pkgsInclude, base)
				}
				return nil
			})
		}
	}

	asts := findAndParseAllGoFiles(root)
	fmt.Println("PATH, INST.")
	for path, astFile := range asts {
		inward := exportedFuncs(path, astFile)
		outward := functionCalls(path, astFile, pkgsInclude)
		inst := instability(outward, inward)
		fmt.Printf("%s, %f\n", strings.Replace(path,  goPath + "/", "", 1), inst)
	}
}

func instability(Ce int, Ca int) float64 {
	i := float64(Ce) / (float64(Ce) + float64(Ca))
	if math.IsNaN(i) {
		i = 0.0
	}
	return i
}

func functionCalls(path string, node *ast.File, pkgsInclude []string) int {
	fCalls := 0
	ast.Inspect(node, func(n ast.Node) bool {
		switch fCall := n.(type) {
		case *ast.CallExpr:
			if fun, ok := fCall.Fun.(*ast.SelectorExpr); ok {
				x := ""
				if ident, ok := fun.X.(*ast.Ident); ok {
					x = ident.Name
				}
				if contains(pkgsInclude, x) {
					if fun.Sel.IsExported() {
						fCalls += 1
					}
				}
			}
		}
		return true
	})
	return fCalls
}

func exportedFuncs(path string, node *ast.File) int {
	exported := 0
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			if fn.Name.IsExported() {
				exported+= 1
			}
		}
		return true
	})
	return exported
}

func findAndParseAllGoFiles(root string) map[string]*ast.File {
	if !isDir(root) {
		println("Not a directory or doesn't exist: " + root)
		os.Exit(0)
	}
	logf("Trawling " + root)
	goFiles := make(map[string]*ast.File)
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ext := filepath.Ext(path)
			if ext == ".go" {
				logf("\tFound: " + path)
				fset := token.NewFileSet()
				ast, err := parser.ParseFile(fset, path, nil, 0)
				if err == nil {
					goFiles[path] = ast
				} else {
					fmt.Fprintf(os.Stderr, "Error while parsing " + path)
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return goFiles
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return fi.IsDir()
}

func logf(f string, a ...interface{}) {
	if *debugFlag {
		log.Printf(f, a...)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}