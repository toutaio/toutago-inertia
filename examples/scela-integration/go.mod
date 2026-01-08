module example

go 1.22.9

require (
	github.com/toutaio/toutago-cosan-router v1.4.2
	github.com/toutaio/toutago-inertia v0.6.0
	github.com/toutaio/toutago-scela-bus v1.5.5
)

require github.com/gorilla/websocket v1.5.3 // indirect

replace github.com/toutaio/toutago-inertia => ../..

replace github.com/toutaio/toutago-cosan-router => /home/nestor/Proyects/touta-for-go/toutago-cosan-router
