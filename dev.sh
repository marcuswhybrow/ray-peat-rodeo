mkdir -p bin
cd bin

if [ ! -f "modd" ]
then
    echo "Downloading Modd Binary..."
    wget -qO- https://github.com/cortesi/modd/releases/download/v0.8/modd-0.8-linux64.tgz | tar -xz --strip-components 1
    echo ""
fi

if [ ! -f "devd" ]
then
    echo "Downloading Devd Binary..."
    wget -qO- https://github.com/cortesi/devd/releases/download/v0.9/devd-0.9-linux64.tgz | tar -xz --strip-components 1
    echo ""
fi

if [ ! -f "pagefind" ]
then
    echo "Downloading PageFile Binary..."
    wget -qO- https://github.com/CloudCannon/pagefind/releases/download/v0.12.0/pagefind-v0.12.0-x86_64-unknown-linux-musl.tar.gz | tar xvz
    echo ""
fi

cd ..

./bin/modd