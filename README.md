# SE'23 School Case Task

## Installation
### Docker (preferable)
1. build an image
```
docker build -t bitcoin-rate-app . --progress=plain
```

2. run a container (app uses port 3333 by default)
```
docker run --rm -p 3333:3333 bitcoin-rate-app
```
