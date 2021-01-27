FILE=Main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o win-amd64/pdfmerge.exe $FILE
GOOS=windows GOARCH=386 go build -o win-386/pdfmerge.exe $FILE

# Linux
GOOS=linux GOARCH=amd64 go build -o linux-amd64/pdfmerge $FILE
GOOS=linux GOARCH=386 go build -o linux-386/pdfmerge $FILE

# Mac OS
GOOS=darwin GOARCH=amd64 go build -o darwin-amd64/pdfmerge $FILE


# Zip those files and folders
zip win-amd64.zip win-amd64/*
zip win-386.zip win-386/*

zip linux-amd64.zip linux-amd64/*
zip linux-386.zip linux-386/*

zip darwin-amd64.zip darwin-amd64/*
