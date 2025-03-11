# container-layer-download

This CLI tool allows you to download a single layer from a container image. This is useful if you have to debug a container image and want to inspect a single layer.

## How to use

```bash
cld <image> <output-file (optional)>
```

The interactive CLI provides a list of layers in the container image. From this list, you can select the layer you want to download. The layer is saved to the output file. If the output file is not provided, the layer is **not** saved to disk.

### Example

```bash
user@pc % cld nginx:latest layer.tar.gz
Use the arrow keys to navigate: ↓ ↑ → ←
? Action:
  ▸ Size: 406 B, Digest: sha256:9dd21ad5a4a6a856d82bb6bb6147c30ad90a9768c3651c55775354e7649bc74d
    Size: 1.2 kB, Digest: sha256:943ea0f0c2e42ccacc72ac65701347eadb2b0cb22828fac30f1400bba3d37088
    Size: 1.4 kB, Digest: sha256:103f50cb3e9f200431b555078cce5e8df3db6ddc2e54d714a10b994e430e98a3
    Size: 28 MB, Digest: sha256:7cf63256a31a4cc44f6defe8e1af95363aee5fa75f30a248d95cae684f87c53c
    Size: 44 MB, Digest: sha256:bf9acace214a6c23630803d90911f1fd7d1ba06a3083f0a62fd036a6d1d8e274
    Size: 626 B, Digest: sha256:513c3649bb1480ca9a04c73f320b6b5a909e24e4ac18ae72fd56b818241d6730
    Size: 957 B, Digest: sha256:d014f92d532d416c7b9eadb244f14f73fdb3d2ead120264b749e342700824f3c
```

Navigate with the arrow keys and select the layer you want to download. Press enter to download the layer.

## Using a proxy

Since in some environments the direct download of the layer is not possible, you can use a proxy to download the layer. For this, you have to set the environment variable `HTTP_PROXY` or `HTTPS_PROXY` with the proxy URL.

```bash
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=http://proxy.example.com:8080
```
