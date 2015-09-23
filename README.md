# Dev Studio API Wrapper - Golang Version

## To Setup

This project relies on a `.env` file to load secret data.  It should look like the following:

    CLIENT_ID=V1:8888:STPS:EXT
    CLIENT_SECRET=aB0cDE12
    URL=https://api.test.sabre.com

And then you can install the go dependencies:

    go get golang.org/x/oauth2
    go get github.com/joho/godotenv

And then you can run the program:

    go run devstudio.go
