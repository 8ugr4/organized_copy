clean:
	rm -r organizer testDir_cp	log.csv source.txt source_hashes.txt destination.txt destination_hashes.txt
build:
	go build -o organizer main.go
test: build
	./organizer --src ./testDir
	./cmpDirs.sh ./testDir ./testDir_cp
