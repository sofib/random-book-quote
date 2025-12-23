## Random Quote

Fetch a random quote.

This project is a sample setup for fetching book quotes from few sources and displaying a random one. At this point it has no persistence and can be in a few ways eg:
- directly as a Go app with `go run ./...`
- via docker:
  - `docker build -t random-quote:latest .`
  - `docker run --rm random-quote`
- via Kubernetes (Cron)Job by starting with `tilt up`. You can triggering runs later on through tilt UI or using `tilt trigger manual-run`.

Project is a sample wire up for development with tilt and kind.

### Dependencies

To run this application via `tilt` you need few dependencies on the system:
- docker
- kind
- ctlptl
- helm
- tilt