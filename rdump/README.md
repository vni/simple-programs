# rdump
Small simple hexdump toy project to play a little with rust.

### Command line arguments:
**--show-filename** *(off by default)*  show filenames before file content

**--no-filename**  do not show filename before file content

**--show-addr** *(on by default)*  show bytes addresses at the left column

**--no-addr**  do not show addr

**--show-ascii** *(on by default)*  show ascii content of the file as the rightmost column

**--no-ascii**  do not show ascii at the rightmost column

### TODO
- [ ] Add --help    | -h cli option
- [ ] Add --version | -V cli option
- [ ] Read files in 4MB chunks (currently the whole file is read into memory)
- [ ] Read files in variable size chunks
