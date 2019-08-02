#!/bin/sh

set -ex

ran -p=1111 -r=test/valid &
pids=$!
ran -p=2222 -r=test/dead_link &
pids="$pids $!"

./muffet http://localhost:1111
! ./muffet http://localhost:2222

./muffet -c 1 http://localhost:1111
./muffet --concurrency 1 http://localhost:1111

./muffet --help

./muffet -v http://localhost:1111 | grep 200
[ $(./muffet -v http://localhost:1111 | wc -l) -eq 8 ]
./muffet --verbose http://localhost:1111 | grep 200
! ./muffet http://localhost:1111 | grep 200

./muffet -v http://localhost:1111 | sort >/tmp/muffet_1.txt
./muffet -v http://localhost:1111 | sort >/tmp/muffet_2.txt
diff /tmp/muffet_1.txt /tmp/muffet_2.txt

[ $(./muffet -rv http://localhost:1111 | wc -l) -eq 6 ]
[ $(./muffet -sv http://localhost:1111 | wc -l) -eq 6 ]

! ./muffet http://localhost:1111 | grep .

kill $pids
