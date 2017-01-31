cp -r ../../../reinsurance_poc/common ./
cp -r ../../../vendor ./
docker build -t upshift/poc-fabric-peer .
rm -rf common
rm -rf vendor
