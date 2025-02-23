# crosswordgame-go

### Local dev

Using `make templ` will execute the templ code-generation.

For CSS, you need to get the Tailwind CLI with `make tailwind-get-cli`; for now
this downloads the MacOS ARM binary specifically. Then, you can use
`make tailwind` to build the CSS, or `make tailwind-watch` to build with hot
reload.

`make lint` will build and lint the code.

There are ... some ... tests with `make test`.

`make run-api` will start the server locally.

### Release

`make docker-build docker-push` will push to the `latest` tag (for now).

Then see https://github.com/mcoot/crosswordgame-go-config for deployment.