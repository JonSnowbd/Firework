rm -rf ./_build
mkdir _build

mkdir ./_build/windows64

# Windows compilation.
env GOOS=windows GOARCH=amd64 go build -o ./_build/windows64/Firework.exe -ldflags "-s -H windowsgui"