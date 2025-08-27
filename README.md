# organized_copy
A go tool that copies a dir and organizes it in sub-dirs

## Usage
```shell 
  make build
  ./oc --src="/path/to/your/source/dir" --dst="/path/to/your/destination/dir" --log="/path/to/your/logfile.csv"
  ./cmpDirs.sh # sha256-validation (optional)
```


## Features
- use ``make clean-all`` to get rid of test directory, log file, sha256 comparison log files.
- if multiple files exist with same name + extension, new files get `_number` after first one.
- if user doesn't set a destination path, auto destination path is source path + `_cp` in same directory.
## Example (after run)
```shell
├── testDir
│   ├── dir1
│   │   ├── dfgdsfg.pdf
│   │   ├── dummy.jpg
│   │   ├── dummy.png
│   │   ├── dummy.txt
│   │   └── feedback.pdf
│   ├── dir2
│   │   ├── asdfgherh4413.pdf
│   │   ├── dfgdsfg.pdf
│   │   ├── dummy.jpg
│   │   ├── dummy.png
│   │   ├── dummy.txt
│   │   ├── feedback.pdf
│   │   └── fsdfgkn.pdf
│   ├── dummasdasdy.png
│   ├── dummy2.jpg
│   ├── feedback.pdf
│   ├── feedback2.pdf
│   ├── fil1
│   ├── sofjngongf24214443.pdf
│   ├── xyz1.csv
│   └── xyz2.csv
└── testDir_cp
    ├── applications
    ├── archives
    ├── audios
    ├── documents
    │   ├── asdfgherh4413.pdf
    │   ├── dfgdsfg.pdf
    │   ├── dfgdsfg_1.pdf
    │   ├── dummy.txt
    │   ├── dummy_1.txt
    │   ├── feedback.pdf
    │   ├── feedback2.pdf
    │   ├── feedback_1.pdf
    │   ├── feedback_2.pdf
    │   ├── fsdfgkn.pdf
    │   ├── sofjngongf24214443.pdf
    │   ├── xyz1.csv
    │   └── xyz2.csv
    ├── images
    │   ├── dummasdasdy.png
    │   ├── dummy.jpg
    │   ├── dummy.png
    │   ├── dummy2.jpg
    │   ├── dummy_1.jpg
    │   └── dummy_1.png
    ├── unknown
    │   └── fil1
    └── videos

```