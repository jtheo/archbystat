# Archive file by mod time

I needed something simple to archive a bunch of files, like screenshots or terraform plan/apply outputs.

archbystat, the name may change in the future, reads a dir, by default the one where it's called; it reads the list of files and dir, ignore the latter, and then it reads for each file the mod time and creates a tree structure in the `archive` directory like yyyy/mm/dd, and move the file in it.

```
Usage of archbystat
  -V    show version and exits
  -a string
        directory where to save (default "archive")
  -o int
        how many minutes older the screenshot need to be to be moved (default 60)
  -p string
        directory to process this is mandatory
  -post string
        postfix to filter the files to process
  -pre string
        prefix to filter the files to process
  -v    verbose output
```


## To Do

- exclude list
- using exif instead of modtime for images
- ???
