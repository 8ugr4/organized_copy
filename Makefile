clean-cmpDirs:
	rm source.txt source_hashes.txt destination.txt destination_hashes.txt
clean-all:
	rm -r oc testDir_cp	log.csv source.txt source_hashes.txt destination.txt destination_hashes.txt
build:
	go build -o oc main.go