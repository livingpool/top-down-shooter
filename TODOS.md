1. make a single client to single server model

* server websocket
* server handles client updates (keystrokes); updates physics (not using RunGame? bc we dont need ui)
* so if we dont use RunGame, we need a separate timer for updating the server physics
* server sends updates to clients (set deltas)
* so i think 2 timers

* client websocket
* client sends updates asap?
* client updates physics (using RunGame? bc we need ui)
* so i think 1 timer only

* test if everything works
* optional: simulate slow connection; decide when to drop the connection
