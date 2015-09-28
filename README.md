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
    }

And then you can run the program:

    go run devstudio.go
    
## Disclaimer of Warranty and Limitation of Liability

This software and any compiled programs created using this software are furnished “as is” without warranty of any kind, including but not limited to the implied warranties of merchantability and fitness for a particular purpose. No oral or written information or advice given by Sabre, its agents or employees shall create a warranty or in any way increase the scope of this warranty, and you may not rely on any such information or advice.
Sabre does not warrant, guarantee, or make any representations regarding the use, or the results of the use, of this software, compiled programs created using this software, or written materials in terms of correctness, accuracy, reliability, currentness, or otherwise. The entire risk as to the results and performance of this software and any compiled applications created using this software is assumed by you. Neither Sabre nor anyone else who has been involved in the creation, production or delivery of this software shall be liable for any direct, indirect, consequential, or incidental damages (including damages for loss of business profits, business interruption, loss of business information, and the like) arising out of the use of or inability to use such product even if Sabre has been advised of the possibility of such damages.

