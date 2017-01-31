FROM hyperledger/fabric-peer:latest

RUN mkdir -p $GOPATH/src/github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/common
COPY ./common $GOPATH/src/github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/common/

RUN mkdir -p $GOPATH/src/github.com/ajmanlove/hyperledger-sandbox/vendor
COPY ./vendor $GOPATH/src/github.com/ajmanlove/hyperledger-sandbox/vendor/