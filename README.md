# Testing image-spec `data` field

https://github.com/opencontainers/image-spec/pull/826

### Usage

Push an image containing inline data to a remote registry.

```
$ crane manifest $(go run ./ -ref=${image}) | jq
```

```
$ go run ./ -help
  -oci
    	if true, push as OCI image
  -ref string
    	ref to push to (default "gcr.io/imjasonh/data")
  -size int
    	size of bytes to send (default 10)
```

### Repos tested

- gcr.io/imjasonh/data
- quay.io/imjasonh/data
- docker.io/imjasonh/data
- ghcr.io/dlorenc/signed-container:foo
- registry.digitalocean.com/dlorenc/test:foo
- bundle.bar/u/danlorenc/test:v1

You can check these manifests yourself:

```bash
$ crane manifest docker.io/imjasonh/data | jq
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
  "config": {
    "mediaType": "application/vnd.docker.container.image.v1+json",
    "size": 233,
    "digest": "sha256:5c6252906e515b45399e41fda0307669c92e15d7745ea0d49670fe37c5c1b568"
  },
  "layers": [
    {
      "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
      "size": 5,
      "digest": "sha256:2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
      "data": "aGVsbG8="
    }
  ]
}
```

If you test this with another repository please send a PR!
