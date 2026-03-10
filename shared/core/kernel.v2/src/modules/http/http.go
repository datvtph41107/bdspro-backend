package http

import contract_runtime "core/kernel.v2/src/contracts/runtime"

type HTTPServer interface {
	contract_runtime.Component
	Use(middleware any) HTTPServer
	Handle(method, path string, handler any) HTTPServer
	GET(path string, handler any) HTTPServer
	POST(path string, handler any) HTTPServer
}
