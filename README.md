# sara_updates

### part1_photos_proto
* create directory `backend/grpc/proto`
* create file `photos.proto` and scaffold with `syntax` & `option`
* Go to google photos API explorer and show different reponses between `list` and `get` APIs
* build out appropriate protos
* create script file to build protos
    * make it executable by typing `chmod +x filename.sh`
* demonstrate compling for various languages

### part2_oauth_setup
* Explain OAuth and show [Google OAuth2 API](https://developers.google.com/identity/protocols/oauth2)
* Show the [Golang Oauth2 API](https://pkg.go.dev/golang.org/x/oauth2/clientcredentials) and discuss how using the `Config` struct will be utilized to make an authorized http client to make the requests to Google servers. 
 - The members of the `Config` struct set the basis for our protos

* Discuss the token info param of the config which sets the basis for the `Token` proto
  - Can click to the token source code from the `Config` docs or go [here](https://cs.opensource.google/go/x/oauth2/+/refs/tags/v0.4.0:token.go;drc=e07593a4c41a489556d019d1ad4d82e9ee66b4a7;l=31)

* create file `oauth.proto` and make the appropriate protos
* add `OauthConfigInfo` to `photos.proto`

### part3_grpc_photos_server
* Create file `backend/common/api.go` and establish a `SaraServer` interface
* Create file `backend/grpc/main.go` and begin coding
* Make note of missing dependencies - talk about go mod. Go to the root of the directory and type `go mod init github.com/Mooseburger1/sara` (or whatever you want to call it during the tutorial).
* Run go mod tidy to download dependencies so far.
* Create directory `backend/grpc/photos/google/server.go` and start creating server while implementing the `SaraServer` interface.