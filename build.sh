mkdir bin
mkdir -p bin/data/downloads
cp -r src/gui/frontend bin

cd src/gui/
go build server.go
cd ../../
mv src/gui/server bin