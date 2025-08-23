#!/usr/bin/env bash
# $1 = source dir, $2 = copied (dest) dir
find "$1" -type f -print0 | xargs -0 sha256sum | sort > source.txt
find "$2" -type f -print0 | xargs -0 sha256sum | sort > destination.txt

cut -d' ' -f1 source.txt | sort | uniq -c > source_hashes.txt
cut -d' ' -f1 destination.txt | sort | uniq -c > destination_hashes.txt

while IFS= read -r line; do
    hash="${line%% *}"
    file="${line#* }"
    found=$(grep -c "^$hash " destination.txt)
    if [ "$found" -lt 1 ]; then
        echo "MISSING: $file ($hash)"
    fi
done < source.txt
