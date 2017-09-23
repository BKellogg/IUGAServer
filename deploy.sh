deploy () {

    GOOS=linux go build

    docker build -t brendankellogg/iugaserver .

    if [ "$(docker ps -aq --filter name=iuga-server)" ]; then
        docker rm -f iuga-server
    fi

    docker run -d \
    --name iuga-server \
    --network iuga-net \
    -e ADDR=:443 \
    -e TLSKEY=/tls/privkey.pem \
    -e TLSCERT=/tls/fullchain.pem \
    -e IUGASITEADDR=$IUGASITEADDR \
    -e IUGAEVENTSADDR=$IUGAEVENTSADDR \
    -v /Users/Brendan/Documents/go/src/github.com/BKellogg/iugaserver/tls/:/tls/:ro \
    -p 80:80 \
    -p 443:443 \
    brendankellogg/iugaserver
}

deploy