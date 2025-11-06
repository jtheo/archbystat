// Silly program to archive files based on modtime in a structure like archive/yyyy/mm/dd
package main

import (
	"cmp"
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
	c := start()
	listEntries, err := os.ReadDir(c.processDir)
	if err != nil {
		fmt.Printf("Error reading content of dir %s: %v\n", c.processDir, err)
		os.Exit(1)
	}

	for _, entry := range listEntries {
		if entry.IsDir() {
			continue
		}

		eName := entry.Name()
		fn := filepath.Join(c.processDir, entry.Name())

		fs, err := os.Stat(fn)
		if err != nil {
			log.Printf("Error loading stats of %s: %v\n", fn, err)
			os.Exit(1)
		}

		if time.Since(fs.ModTime()) < (time.Duration(c.older) * time.Minute) {
			if c.verbose {
				log.Printf("%s is newer than the threshold %d, continue", fn, c.older)
			}
			continue
		}

		if c.prefix != "" && !strings.HasPrefix(eName, c.prefix) {
			if c.verbose {
				log.Printf("%s doesn't match prefix: %s. Continue\n", eName, c.prefix)
			}
			continue

		}
		if c.postfix != "" && !strings.HasSuffix(eName, c.postfix) {
			if c.verbose {
				log.Printf("%s doesn't match postfix %s. Continue\n", eName, c.postfix)
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
		archPath := filepath.Join(c.archiveDir, ys, ms, ds)
		if !c.dryRun {
			err = os.MkdirAll(archPath, 0o750)
			if err != nil {
				fmt.Printf("Error creating %s: %v\n", archPath, err)
			}
		}
		if !c.dryRun {
			err = os.Rename(fn, filepath.Join(archPath, eName))
			if err != nil {
				fmt.Printf("Error moving %s to %s: %v\n", archPath, filepath.Join(archPath, eName), err)
			}
		}

		if c.verbose {
			fmt.Printf("File %s moved to %s\n", fn, archPath)
		}
	}
}

type config struct {
	archiveDir            string
	processDir            string
	prefix                string
	postfix               string
	older                 int64
	verbose, superVerbose bool
	showver               bool
	dryRun                bool
}

func start() *config {
	basedir := cmp.Or(os.Getenv("HOME"), "/tmp")
	c := config{}
	flag.StringVar(&c.processDir, "p", filepath.Join(basedir, "L"), "directory to process this is mandatory")
	flag.StringVar(&c.archiveDir, "a", filepath.Join(basedir, "archive"), "directory where to save")
	flag.StringVar(&c.prefix, "pre", "", "prefix to filter the files to process")
	flag.StringVar(&c.postfix, "post", "", "postfix to filter the files to process")
	flag.Int64Var(&c.older, "o", 60, "how many minutes older the screenshot need to be to be moved")
	flag.BoolVar(&c.verbose, "v", false, "verbose output")
	flag.BoolVar(&c.superVerbose, "vv", false, "verbose output with parameters shown")
	flag.BoolVar(&c.showver, "V", false, "show version and exits")
	flag.BoolVar(&c.dryRun, "d", false, "dry-run")
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

	if c.showver {
		log.Printf("Version %s\n", Version)
		os.Exit(0)
	}

	if c.superVerbose {
		c.verbose = true
	}

	if c.dryRun {
		c.verbose = true
	}

	if c.processDir == "" {
		log.Printf("You need to provide a directory to process\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if c.superVerbose {
		fmt.Printf("Process Dir: %s\n", c.processDir)
		fmt.Printf("Archive Dir: %s\n", c.archiveDir)
		fmt.Printf("older in minutes: %d\n", c.older)
		fmt.Printf("\nVersion: %s\n", Version)
		if c.dryRun {
			fmt.Printf("Dry-Run active, command won't move files, but only shows what's happening\n")
		}
		fmt.Println()
	}
	return &c
}
