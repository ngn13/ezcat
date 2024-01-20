# ezcat 🐱 easy shell handler
ezcat allows you to manage multiple reverse shells 
with a simple web interface

### showcase
![](/assets/showcase.gif)

### install 
an easy way install ezcat is to use docker:
```
docker run --rm --network host \
    -e PASSWORD=securepassword \
    ghcr.io/ngn13/ezcat
```
or you can download the latest binary from the [releases tab](https://github.com/ngn13/ezcat/releases),
then copy it to somewhere in your `PATH`

### build
install a recent version of go, then clone the repository and run the go build command:
```bash
go build .
```
