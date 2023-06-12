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
* Create file `backend/grpc/photos/google/server_api.go` and stub out the googlephotoserver api methods
* Create file `backend/grpc/main.go` and implement, then run `go run main.go` to test out the code thus far.

### part4_grpc_photos_server_unit_tests
* Create file `backend/grpc/photos/google/server_test.go` and write the unit tests as shown

### part5_grpc_photos_implementation
* Begin implementing `ListAlbums` method
* Once to the "create client" error section, create the file `backend/common/errors.go` to implement the common error funcs
* Once to the headers part go to https://developers.google.com/photos/library/reference/rest/v1/albums/list and show how the pageToken and other query params are added to the URL

### part6_grpc_list_albums_unit_test
* Create file `backend/grpc/photos/google/server_api_test.go` and begin creating unit test
* You will need to explain RoundTripFunc and http.Client{Transport:}. The RoundTripFunc allows the injected client to respond to http request with the provided http.Response as returned by the RoundTripFunc

### part7_docker_setup
* Create file `sara_updated/Dockerfile.grpc` and create the dockerfile
* Buld the file by typing `docker build -t grpc-backend -f Dockerfile.grpc .`
* Run the image by typing `docker run --name gtest -it -d --rm -p 4000:4000 grpc-backend`
* Test the server with postman request
 - File -> new -> gRpc*
 - Type the server url `localhost:4000`
 - Go to service definition and use "reflection"
 - Go back to `message` section and below the form box - click "Use example Message"
* Make docker-compose file at `sara_updated/Docker-compose.yml` - explain that we don't want to keep typing all the docke cli cmds

### part8_photos_rest_service
* Create file `backend/rest/photos_service.go` - make note that we don't need to create subdirectories like grpc since REST side is provider agnostic {AWS, DropBox, etc.}
* in the photo service, once stubbing out the `RegisterGetRoutes` method, create the `backend/rest/main.go`
