cp -r ../../../reinsurance_poc/common ./
docker build -t upshift/poc-fabric-peer .
rm -rf common
