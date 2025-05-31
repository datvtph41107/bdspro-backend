.PHONY: buf-generate
generate:
	buf mod update
	buf generate --template ./buf.gen.yaml --path ./organization -o ../types
