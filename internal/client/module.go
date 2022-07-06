package client

type Module interface {
	Load()
}

var modules map[string]Module

func moduleInit() {
	modules = make(map[string]Module)
}
