# vugu-site
vugu.org website

# Build & Run

`cd server; go generate ../app && go build . && ./server -dev`

# Distribute

Make sure you have vugu.github.io checked out in an adjacent directory (the files in which will be overwritten) and then:

`cd server; go generate ../app && go build . && ./server -build`

