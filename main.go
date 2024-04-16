package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var Version = "0.0"

func main() {
	var archiveDir string
	var processDir string
	var prefix string
	var postfix string
	var older int64
	var verbose bool
	var showver bool

	flag.StringVar(&processDir, "p", "", "directory to process this is mandatory")
	flag.StringVar(&archiveDir, "a", "archive", "directory where to save")
	flag.StringVar(&prefix, "pre", "", "prefix to filter the files to process")
	flag.StringVar(&postfix, "post", "", "postfix to filter the files to process")
	flag.Int64Var(&older, "o", 60, "how many minutes older the screenshot need to be to be moved")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&showver, "V", false, "show version and exits")
	flag.Parse()

	if showver {
		log.Printf("Version %s\n", Version)
		os.Exit(0)
	}

	if processDir == "" {
		log.Printf("You need to provide a directory to process\n\n")
		flag.Usage()
		os.Exit(1)
	}

	listEntries, err := os.ReadDir(processDir)
	if err != nil {
		fmt.Printf("Error reading content of dir %s: %v\n", processDir, err)
		os.Exit(1)
	}

	for _, entry := range listEntries {
		if entry.IsDir() {
			continue
		}

		fn := entry.Name()

		fs, err := os.Stat(fn)
		if err != nil {
			log.Printf("Error loading stats of %s: %v\n", fn, err)
			os.Exit(1)
		}

		if time.Since(fs.ModTime()) < (time.Duration(older) * time.Minute) {
			if verbose {
				log.Printf("%s is newer than the threshold %d, continue", fn, older)
			}
			continue
		}

		if prefix != "" && !strings.HasPrefix(fn, prefix) {
			if verbose {
				log.Printf("%s doesn't match prefix: %s. Continue\n", fn, prefix)
			}
			continue

		}
		if postfix != "" && !strings.HasSuffix(fn, postfix) {
			if verbose {
				log.Printf("%s doesn't match postfix %s. Continue\n", fn, postfix)
			}
			continue
		}

		year, month, day := fs.ModTime().Date()
		ys := strconv.Itoa(year)
		m := int(month)
		ms := strconv.Itoa(m)
		if m < 10 {
			ms = "0" + ms
		}
		ds := strconv.Itoa(day)
		if day < 10 {
			ds = "0" + ds
		}
		archPath := filepath.Join(archiveDir, ys, ms, ds)
		err = os.MkdirAll(archPath, 0755)
		if err != nil {
			fmt.Printf("Error creating %s: %v\n", archPath, err)
		}

		err = os.Rename(fn, filepath.Join(archPath, fn))
		if err != nil {
			fmt.Printf("Error moving %s to %s: %v\n", archPath, filepath.Join(archPath, fn), err)
		}
		if verbose {
			fmt.Printf("File %s moved to %s\n", fn, archPath)
		}
	}
}
