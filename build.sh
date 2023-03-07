echo "Building Pages..."
echo ""

go run main.go

echo ""

mkdir -p bin
cd bin

if [ ! -f "pagefind" ]
then
    echo "Downloading PageFile Excecutable..."
    wget -qO- https://github.com/CloudCannon/pagefind/releases/download/v0.12.0/pagefind-v0.12.0-x86_64-unknown-linux-musl.tar.gz | tar xvz
    echo ""
fi

cd ..

echo "Building Search Index (PageFind)..."

./bin/pagefind --source build