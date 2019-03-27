# vugu-site
vugu.org website

## TODO:

* Finish writing docs

* Spell check everything

* Put in some basic server-side rendering

* Consider making the convention be main_wasm.go, since the same package could be used in different ways DONE
 - hm, it's not that simple - if you write some broken stuff in your component and then run your dev server,
 you're stuck - because the dev server tries to compile against the broken code and fails, but your dev server 
 is also what runs the code generator, so you're stuck.  probably stick with separate dev server for 
 quick prototyping and // +build ignore, then another one can be added that is the actual server applicationa

* What about a live preview system that lets you edit components and see the output in-browser!...

* DO a hit on the godoc..

* Deploy to production (HTTPS with LE would be nice).  Do a brief hit on the dev vs prod webserver stuff in vugu while we're at it, including an example of a build script for production.

* Clean up the github readme and make it all pretty, maybe with a graphic.
