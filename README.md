# Dev Studio API Wrapper - Golang Version

## To Setup

This project relies on a `.env` file to load secret data.  It should look like the following:

    CLIENT_ID=V1:8888:STPS:EXT
    CLIENT_SECRET=aB0cDE12
    URL=https://api.test.sabre.com

And then you can install the go dependencies:

    go get github.com/SabreDevStudio/sabre-dev-studio-golang

Create your program, which could look something like this:

    package main
    
    import (
            devstudio "github.com/SabreDevStudio/sabre-dev-studio-golang"
            "github.com/joho/godotenv"
            "log"
    )
    
    func main() {
            if err := godotenv.Load(); err != nil {
                    log.Fatal("Error loading .env file")
            }
            client := devstudio.NewClient()
            client.GetTravelThemes()
            //client.GetFlightSearch()
    }

And then you can run the program:

    go run devstudio.go
