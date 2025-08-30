clean:
	rm -r -f testDir_cp log.csv source.txt source_hashes.txt destination.txt destination_hashes.txt

build:
	go build -o organizer main.go

test:
	./organizer --src ./testDir --verbose

cmp: clean
	./cmpDirs.sh ./testDir ./testDir_cp

test-cmp: clean build test cmp

hyperfine: build
	rm -f benchmark.csv && touch benchmark.csv
	hyperfine 'make test && make clean' -w 20 -m 20 -s full -u millisecond --export-csv benchmark.csv
