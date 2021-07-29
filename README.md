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

You can check these manifests yourself:

```bash
$ crane manifest docker.io/imjasonh/data | jq
{
  "schemaVersion": 2,
  "config": {
    "mediaType": "application/vnd.docker.container.image.v1+json",
    "size": 233,
    "digest": "sha256:9ad623f2020d10bb7c6f432bea337a91225de1e2bd2b282d2fb7df7ed8ffcfc6"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 10,
      "digest": "sha256:537f3fb69ba01fc388a3a5c920c485b2873d5f327305c3dd2004d6a04451659b",
      "data": "Li4uLi4uLi4uLg=="
    }
  ]
}
```

If you test this with another repository please send a PR!
