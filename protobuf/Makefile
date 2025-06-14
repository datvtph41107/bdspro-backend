gen:
	cd schema && buf mod update
	cd schema && buf generate --template ./notification/buf.gen.yaml --path notification
	cd schema && buf generate --template ./bdspro/buf.gen.yaml --path bdspro
	cd schema && buf generate --template ./task/buf.gen.yaml --path task
	cd schema && buf generate --template ./spec/buf.gen.yaml --path spec
	cd schema && buf generate --template ./buf.gen.yaml --path organization
	cd schema && buf generate --template ./buf.gen.yaml --path group
	cd schema && buf generate --template ./buf.gen.yaml --path appointment
	cd schema && buf generate --template ./buf.gen.yaml --path transaction


