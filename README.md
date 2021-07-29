# testing image-spec `data` field

https://github.com/opencontainers/image-spec/pull/826

Usage:

Push an image containing inline data to a remote registry.

```
$ crane manifest $(go run ./) | jq
```
