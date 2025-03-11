package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/manifoldco/promptui"
)

type layerConfig struct {
	Size   string
	Digest v1.Hash
}

func (lc layerConfig) String() string {
	return fmt.Sprintf("Size: %s, Digest: %s", lc.Size, lc.Digest)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: download-layer <image> <output-file>")
		os.Exit(1)
	}

	imageRef := os.Args[1]
	if imageRef == "" {
		fmt.Println("Error: image reference is required.")
		os.Exit(1)
	}

	outFileName := os.Args[2]
	if outFileName == "" {
		fmt.Println("Error: output file is required.")
		os.Exit(1)
	}

	ref, err := name.ParseReference(imageRef)
	if err != nil {
		fmt.Println("Error parsing image reference:", err)
		os.Exit(1)
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		fmt.Println("Error fetching image:", err)
		os.Exit(1)
	}

	manifest, err := img.Manifest()
	if err != nil {
		fmt.Println("Error fetching manifest:", err)
		os.Exit(1)
	}

	layers := map[int]layerConfig{}
	for i, layer := range manifest.Layers {
		layers[i] = layerConfig{
			Size:   humanize.Bytes(uint64(layer.Size)),
			Digest: layer.Digest,
		}
	}

	prompt := promptui.Select{
		Label: "Layer",
		Items: layerConfigDescription(layers),
		Size:  10,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Error selecting layer:", err)
		os.Exit(1)
	}
	sl := layers[idx]

	selectedLayer, err := img.LayerByDigest(sl.Digest)
	if err != nil {
		fmt.Println("Error fetching selected layer:", err)
		os.Exit(1)
	}

	layerReader, err := selectedLayer.Compressed()
	if err != nil {
		fmt.Println("Error getting layer data:", err)
		os.Exit(1)
	}
	defer layerReader.Close()

	outputFile, err := os.Create(outFileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, layerReader)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		os.Exit(1)
	}

	fmt.Printf("Image %s Layer %s with size %s downloaded successfully to %s\n", imageRef, sl.Digest, sl.Size, outFileName)
}

func layerConfigDescription(m map[int]layerConfig) []string {
	keys := make([]string, 0, len(m))
	for _, v := range m {
		keys = append(keys, v.String())
	}
	return keys
}
