#!/bin/bash

curl  --data "password=happyGiraffe" http://localhost:8080/hash &
curl  --data "password=angryMonkey"  http://localhost:8080/hash &
curl  --data "password=sillyBear"    http://localhost:8080/hash &
curl  --data "password=friskyKitten" http://localhost:8080/hash &
curl  --data "password=lazyWalrus"   http://localhost:8080/hash &
curl  --data "password=amicableFish" http://localhost:8080/hash &
curl  --data "password=sadPanda"     http://localhost:8080/hash &
curl  --data "password=feistyPuppy"  http://localhost:8080/hash &
curl  --data "password=blueChipmunk" http://localhost:8080/hash &
curl  --data "password=honestWolf"   http://localhost:8080/hash &

# Now that a bunch of requests are in progress, send the shutdown request
curl http://localhost:8080/shutdown &

wait

echo "Done"
