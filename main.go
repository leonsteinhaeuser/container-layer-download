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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: download-layer <image> <output-file>")
		os.Exit(1)
	}

	imageRef := os.Args[1]
	if imageRef == "" {
		fmt.Println("Error: image reference is required.")
		os.Exit(1)
	}

	outFileName := ""
	if len(os.Args) > 2 {
		outFileName = os.Args[2]
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

	layers := map[string]layerConfig{}
	for _, layer := range manifest.Layers {
		hmb := humanize.Bytes(uint64(layer.Size))
		layers[fmt.Sprintf("Size: %s, Digest: %s", hmb, layer.Digest)] = layerConfig{
			Digest: layer.Digest,
			Size:   hmb,
		}
	}

	prompt := promptui.Select{
		Label: "Layer",
		Items: mapKeys(layers),
		Size:  10,
	}
	_, ly, err := prompt.Run()
	if err != nil {
		fmt.Println("Error selecting layer:", err)
		os.Exit(1)
	}

	selectedLayer, err := img.LayerByDigest(layers[ly].Digest)
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

	if outFileName != "" {
		outputFile, err := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY, 0664)
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
		fmt.Printf("Image %s Layer %s with size %s downloaded successfully to %s\n", imageRef, layers[ly].Digest, layers[ly].Size, outFileName)
		return
	}
	fmt.Printf("Image %s Layer %s with size %s downloaded successfully\n", imageRef, layers[ly].Digest, layers[ly].Size)
}

func mapKeys(m map[string]layerConfig) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
