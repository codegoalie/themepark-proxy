# Themeparks.wiki proxy

This is a quick proxy for themeparks.wiki wait times API to turn the response
into a map instead of an array. This makes handling it in Home Assistant a
little easier.

## Deployment to RPi

1. Cross-compile to ARM: `env GOOS=linux GOARCH=arm GOARM=7 go build -o themepark-proxy`
1. `scp` the binary and Dockerfile to RPi: `scp themepark-proxy
   user@server:~/themepark-proxy`
1. Restart docker compose on the server.
