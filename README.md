# Ezcat | Easy Netcat Reverse Shell Handler
Ezcat allows you to interact with mutliple netcat 
reverse shells and it lets you manage them easily 
![](assets/showcase.png) 

# Features
- HTTP/HTTPS beacon
- Payload generation
- Simple CLI interface
- File upload & download 
- Easily upgradable netcat shells

# Install
### Automated installation 
For an automated installation, you can use the install script:
```bash
curl https://raw.githubusercontent.com/ngn13/ezcat/main/scripts/install.sh | sudo bash
```

### Manual installation
Download the latest binary from the [releases tab](https://github.com/ngn13/ezcat/releases),
then copy it to somewhere in your `PATH`

# Build
Install a recent version of go (I recommend `1.20`), then clone the repository and run the go build 
command:
```
go build .
```
