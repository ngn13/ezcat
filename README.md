# ezcat 🐱 easy shell handler written in go
![](/assets/showcase.gif)

---

### 🚀 Install 
An easy way install ezcat is to use docker:
```
docker run --rm --network host \
    -e PASSWORD=securepassword \
    ghcr.io/ngn13/ezcat
```
or you can download the latest binary from the [releases tab](https://github.com/ngn13/ezcat/releases),
extract it with `tar` and then copy it to somewhere in your `PATH`.

### ⚒️ Build
To build ezcat, install a recent version go. Then clone the repository and run the go build command:
```bash
go build .
```
For development, login to the web interface with the default password, `ezcat`.
