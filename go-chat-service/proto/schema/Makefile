.PHONY: buf-generate
gen-buf:
	buf generate --template ./chat-buf.gen.yaml --path ./chat -o ../types
	buf generate --template ./external-buf.gen.yaml --path ./external -o ../types
	buf generate --template ./proto-buf.gen.yaml --path ./error -o ../types
	buf generate --template ./proto-buf.gen.yaml --path ./message -o ../types
