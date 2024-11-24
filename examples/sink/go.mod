module github.com/neputevshina/contraption/examples/sink

go 1.22.5

replace github.com/neputevshina/contraption => ../../
replace github.com/neputevshina/contraption/nanovgo => ../../nanovgo

require (
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240506104042-037f3cc74f2a
	github.com/neputevshina/contraption v0.0.0-20240831004727-2430eed6a7fd
	github.com/neputevshina/geom v0.0.0-20231211132238-0cb016e95618
	github.com/neputevshina/contraption/nanovgo v0.0.0-20241104171014-8329a8bbb5cf
	golang.org/x/image v0.19.0
)

require (
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/goxjs/gl v0.0.0-20230705020350-37525f4d9d35 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	golang.org/x/exp v0.0.0-20240823005443-9b4947da3948 // indirect
	golang.org/x/exp/shiny v0.0.0-20240823005443-9b4947da3948 // indirect
	golang.org/x/text v0.17.0 // indirect
	honnef.co/go/js/console v0.0.0-20150119023344-105276c43558 // indirect
)
