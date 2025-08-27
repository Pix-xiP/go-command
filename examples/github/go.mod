module github.com/pix-xip/go-command/examples/github

go 1.25.0

require (
	github.com/google/go-github/v56 v56.0.0
	github.com/pix-xip/go-command v0.0.0
)

require github.com/google/go-querystring v1.1.0 // indirect

replace github.com/pix-xip/go-command => ../..
