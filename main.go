package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var elbow rune
var vertical rune
var horizontal rune
var tee rune
var space rune

var singleHeader string
var multiHeader string
var indentBuffer string

var totalFiles int64
var totalDirectories int64

func setChars(simple bool) {
	if !simple {
		elbow = '└'
		vertical = '│'
		horizontal = '─'
		tee = '├'
		space = ' '
	} else {
		elbow = '+'
		vertical = '|'
		horizontal = '-'
		tee = '+'
		space = ' '
	}

	singleHeader = fmt.Sprintf("%c%c%c%c", elbow, horizontal, horizontal, space)
	multiHeader = fmt.Sprintf("%c%c%c%c", tee, horizontal, horizontal, space)
	indentBuffer = fmt.Sprintf("%c%c%c%c", vertical, space, space, space)
}

func makeIndent(level int) string {
	indent := fmt.Sprintf("%c%c%c%c%c%c", 14, elbow, horizontal, horizontal, space, 15) //"+-- "
	if level > 1 {
		for i := 1; i < level; i++ {
			//indent = "|   " + indent
			indent = fmt.Sprintf("%c%c%c%c", vertical, space, space, space) + indent
		}
	}
	return indent
}

func listDir(dir string, level int) {
	if level > maxLevel {
		return
	}
	currDir, err := os.Open(dir)
	defer currDir.Close()
	if err != nil {
		fmt.Printf("Error seen opening directory %s - %v\n", dir, err)
		return
	}
	names, err := currDir.Readdirnames(0)
	if err != nil {
		fmt.Printf("Error seen scanning %s- %v\n", dir, err)
		return
	}
	sort.Strings(names)
	count := len(names)
	indent := ""
	for idx, f := range names {
		fi, err := os.Lstat(filepath.Join(dir, f))
		if err != nil {
			fmt.Printf("Error stat'ing %s - %v\n", f, err)
			return
		}
		if count == 1 || idx == count-1 {
			indent = singleHeader
		} else {
			indent = multiHeader
		}
		if level > 0 {
			for i := 0; i < level; i++ {
				indent = indentBuffer + indent
			}
		}
		// TODO: Add color highlighting for directories?
		if okayToShow(f) {
			fmt.Printf("%s%s\n", indent, f)
			if fi.IsDir() {
				totalDirectories++
				listDir(filepath.Join(dir, fi.Name()), level+1)
			} else {
				totalFiles++
			}
		}
	}
}

func okayToShow(f string) bool {
	return all || !strings.HasPrefix(f, ".")
}

var exclude string
var include string
var all bool
var size bool
var human bool
var kind bool
var maxLevel int

func main() {
	var root string
	var simple bool

	flag.StringVar(&root, "d", ".", "Directory to build tree from")
	flag.BoolVar(&simple, "s", false, "Use simple chars for line drawing")
	flag.IntVar(&maxLevel, "L", math.MaxInt16, "Maximum levels to display")
	flag.BoolVar(&all, "a", false, "Show all files")
	// Unimplemented Flags
	flag.StringVar(&include, "P", "", "Filename pattern to include in tree (e.g. *foo*)")
	flag.StringVar(&exclude, "I", "", "Filename pattern to exclude from tree (e.g. tmp*)")
	flag.BoolVar(&size, "z", false, "Show file sizes in bytes")
	flag.BoolVar(&human, "h", false, "Show file sizes in human readable format (Kib, MiB, Gib, Tib)")
	flag.BoolVar(&kind, "F", false, "Append a '/' for directories, a '=' for socket files, a '*' for executable files and a '|' for FIFO's, as per ls -F")

	flag.Parse()

	setChars(simple)

	fmt.Printf("%s\n", root)
	listDir(root, 0)
	fmt.Printf("%d directories, %d files\n", totalDirectories, totalFiles)

}
