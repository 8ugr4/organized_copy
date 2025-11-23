clean:
	rm -r -f organizer testDir_cp log.csv source.txt source_hashes.txt destination.txt destination_hashes.txt

build:
	go build -o organizer main.go

test: build
	./organizer --src ./testDir --verbose

test-async: build
	./organizer --src ./testDir --verbose --async

cmp:
	./cmpDirs.sh ./testDir ./testDir_cp

integration-seq: clean build test cmp
	make clean

integration-async: clean build test-async cmp
	make clean

hyperfine: build
	rm -f benchmark.csv && touch benchmark_seq.csv
	hyperfine 'make test && make clean' -w 20 -m 20 -s full -u millisecond --export-csv benchmark_seq.csv

hyperfine-async: build
	rm -f benchmark.csv && touch benchmark_async.csv
	hyperfine 'make test-async && make clean' -w 20 -m 20 -s full -u millisecond --export-csv benchmark_async.csv

clean-benchmark:
	rm -f benchmark* organizer
