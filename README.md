#### API for tweets and statistics from twitter.

### Dependencies

1. [Golang](https://golang.org/dl/)
1. [Docker](https://www.docker.com/)
1. [docker-compose](https://docs.docker.com/compose/install/)
1. [direnv](https://direnv.net/)  

### Configure host aliases
Add `127.0.0.1	tweethub.io` to `/etc/hosts`.

### Start the system
```bash
make start
```

### Development
Setup development utilities.
```bash
make setup
```