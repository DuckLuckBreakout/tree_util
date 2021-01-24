package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func formFormatString(prefix string, isLast bool, isEmpty bool, isDir bool) string {
	var startRune string
	if isLast {
		startRune = "├"
	} else {
		startRune = "└"
	}

	if isDir {
		return prefix + startRune + "───%s\n"
	} else if isEmpty {
		return prefix + startRune + "───%s (empty)\n"
	} else {
		return prefix + startRune + "───%s (%db)\n"
	}

}

func getDirFiles(out io.Writer, prefix, pwd string, printFiles bool) error {
	files, err := ioutil.ReadDir(pwd)
	if err != nil {
		return err
	}

	if !printFiles {
		var printOnlyDir []os.FileInfo
		for _, file := range files {
			if file.IsDir() {
				printOnlyDir = append(printOnlyDir, file)
			}
		}
		files = printOnlyDir
	}

	length := len(files)
	for i, file := range files {
		if file.IsDir() {
			var prefixNew string
			if length > i + 1 {
				if _, err = fmt.Fprintf(out,
										formFormatString(prefix, length > i + 1,false, true),
										file.Name()); err != nil {
					return err
				}
				prefixNew = prefix + "│\t"
			} else {
				if _, err = fmt.Fprintf(out,
										formFormatString(prefix, length > i + 1,false, true),
										file.Name()); err != nil {
					return err
				}
				prefixNew = prefix + "\t"
			}
			if err = getDirFiles(out, prefixNew, filepath.Join(pwd, file.Name()), printFiles); err != nil {
				return err
			}
		} else if printFiles {
			size := file.Size()
			if size > 0 {
				if length > i + 1 {
					if _, err = fmt.Fprintf(out,
											formFormatString(prefix, length > i + 1, size == 0, false),
											file.Name(),
											file.Size()); err != nil {
						return err
					}
				} else {
					if _, err = fmt.Fprintf(out,
											formFormatString(prefix, length > i + 1, size == 0, false),
											file.Name(),
											file.Size()); err != nil {
											return err
										}
				}
			} else {
				if length > i + 1 {
					if _, err = fmt.Fprintf(out,
											formFormatString(prefix, length > i + 1, size == 0, false),
											file.Name()); err != nil {
											return err
										}
				} else {
					if _, err = fmt.Fprintf(out,
											formFormatString(prefix, length > i + 1, size == 0, false),
											file.Name()); err != nil {
											return err
										}
				}
			}
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return getDirFiles(out, "", path, printFiles)
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}