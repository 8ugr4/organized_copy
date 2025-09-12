# organizer
Lightweight CLI tool written in Go that copies files from a source directory to a destination directory while automatically organizing them into folders by type (audio, images, video, documents, other).

### Example Usage
```shell
  make build
  # Dry-run: show what will be copied
  ./organizer --src ~/Downloads --dst ~/Sorted --dry-run

  # Perform copy with log file including status of copy process of every single file and dir
  ./organizer --src ~/Downloads --dst ~/Sorted --log=~/logfile.csv

  # sha256-validation (optional)
  ./cmpDirs.sh ~/Downloads ~/Sorted
```


## Features
- if multiple files exist with same name + extension, new files get `_number` after first one.
- if user doesn't set a destination path, auto destination path is source path + `_cp` in same directory.
- User can set a rule-set, defining which files will go to which destination.
- **WIP** 'name_contains' and 'priority_order' is still under development.

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
