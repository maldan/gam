GOARCH=arm64 GOOS=linux go build -ldflags "-s -w" -o gam .
cp gam /mnt/orangepi/root/.gam