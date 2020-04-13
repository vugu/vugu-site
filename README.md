# vugu-site
vugu.org website

# Build & Run

## Fancy New Way
`vgrun -watch-dir=app server -dev`

## Old Crusty Way
`cd server; go generate ../ && go build . && ./server -dev`

# Distribute

Make sure you have vugu.github.io checked out in an adjacent directory (the files in which will be overwritten) and then:

`cd server; go generate ../ && go build . && ./server -build`

