::build linux amd64 
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
set GOPATH=/Users/a1234/workspace/lmdgm/deps:/Users/a1234/workspace/lmdgm:/Users/a1234/Workspaces/Go/project
go build -o lmdgm_linux