## Random Quote

Fetch a random quote.

This project is a sample setup for fetching book quotes from few sources and displaying a random one. At this point it has no persistence and can be run in a few ways eg:
- directly as a Go app with `go run ./...`
- via docker:
  ```bash
  docker build -t random-quote:latest .
  docker run --rm random-quote
  ```
  
- via Kubernetes (Cron)Job by:
  ```bash
  TOOLS_PATH= make cluster-up
  TOOLS_PATH= tilt up
  ```

Project is a sample wire up for development with tilt and kind. Initially there was an attempt to wire up kind custer bootstrap within Tiltfile but this was abandoned as the approach is cumbersome and not a standard one.

## Architecture

```mermaid
%%{init: {'theme':'neutral'}}%%
graph TB
    subgraph "Main Application"
        Main[main.go]
        App[App]
        Settings[Settings]
        Collector[QuoteCollector]
        Randomizer[QuoteRandomizer]
    end

    subgraph "Infrastructure Layer"
        Emitter[Event Emitter]
        Factory[Crawler Factory]

        subgraph "Crawlers"
            GR[Goodreads Crawler]
            LQ[Libquotes Crawler]
        end
    end

    subgraph "External Sources"
        GRS[Goodreads.com]
        LQS[Libquotes.org]
    end

    Main --> App
    App --> Settings
    App --> Emitter
    App --> Collector
    App --> Randomizer
    App --> Factory

    Factory --> GR
    Factory --> LQ

    GR --> Emitter
    LQ --> Emitter

    GR -->|HTTP Fetch| GRS
    LQ -->|HTTP Fetch| LQS

    Emitter -->|QuoteFoundEvent| Collector
    Collector --> Randomizer
    Randomizer -->|Random Quote| Main

    style Main fill:#e1f5ff
    style App fill:#e1f5ff
    style Emitter fill:#fff4e1
    style Factory fill:#fff4e1
    style GR fill:#ffe1e1
    style LQ fill:#ffe1e1
    style Collector fill:#e1ffe1
    style Randomizer fill:#e1ffe1
```


## Dependencies

To run this application via `tilt` you need few dependencies on the system:
- docker
- kind
- ctlptl
- helm
- tilt