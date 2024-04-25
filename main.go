package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
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
	var verbose, superVerbose bool
	var showver bool

	flag.StringVar(&processDir, "p", "", "directory to process this is mandatory")
	flag.StringVar(&archiveDir, "a", "archive", "directory where to save")
	flag.StringVar(&prefix, "pre", "", "prefix to filter the files to process")
	flag.StringVar(&postfix, "post", "", "postfix to filter the files to process")
	flag.Int64Var(&older, "o", 60, "how many minutes older the screenshot need to be to be moved")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&superVerbose, "vv", false, "verbose output with parameters shown")
	flag.BoolVar(&showver, "V", false, "show version and exits")
	flag.Parse()

	if Version == "0.0" {
		Version = func() string {
			if info, ok := debug.ReadBuildInfo(); ok {
				ver := []string{}
				for _, setting := range info.Settings {
					if setting.Key == "vcs.revision" {
						ver = append(ver, setting.Value[:7])
					}
					if setting.Key == "vcs.time" {
						ver = append(ver, setting.Value)
					}
				}
				return strings.Join(ver, " ")
			}
			return ""
		}()
	}

	if showver {
		log.Printf("Version %s\n", Version)
		os.Exit(0)
	}

	if superVerbose {
		verbose = true
	}

	if processDir == "" {
		log.Printf("You need to provide a directory to process\n\n")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Process Dir: %s\n", processDir)
	fmt.Printf("Archive Dir: %s\n", archiveDir)
	fmt.Printf("older in minutes: %d\n", older)
	fmt.Printf("\nVersion: %s\n\n", Version)

	listEntries, err := os.ReadDir(processDir)
	if err != nil {
		fmt.Printf("Error reading content of dir %s: %v\n", processDir, err)
		os.Exit(1)
	}

	for _, entry := range listEntries {
		if entry.IsDir() {
			continue
		}

		eName := entry.Name()
		fn := filepath.Join(processDir, entry.Name())

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

		if prefix != "" && !strings.HasPrefix(eName, prefix) {
			if verbose {
				log.Printf("%s doesn't match prefix: %s. Continue\n", eName, prefix)
			}
			continue

		}
		if postfix != "" && !strings.HasSuffix(eName, postfix) {
			if verbose {
				log.Printf("%s doesn't match postfix %s. Continue\n", eName, postfix)
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

		err = os.Rename(fn, filepath.Join(archPath, eName))
		if err != nil {
			fmt.Printf("Error moving %s to %s: %v\n", archPath, filepath.Join(archPath, eName), err)
		}
		if verbose {
			fmt.Printf("File %s moved to %s\n", fn, archPath)
		}
	}
}
