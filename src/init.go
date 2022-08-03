package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Config struct {
	LayersDirectory string `json:"layers_directory"`
	OutputDirectory string `json:"output_directory"`
	Delimitter      string `json:"delimitter"`
	Debugging       bool   `json:"debugging"`
	Gif             struct {
		Export bool `json:"export"`
		Repeat int  `json:"repeat"`
		Delay  uint `json:"delay"`
	} `json:"gif"`
	Artwork struct {
		Prefix           string  `json:"prefix"`
		Width            float64 `json:"width"`
		Height           float64 `json:"height"`
		CompressionLevel uint8   `json:"compression_level"`
		Background       string  `json:"background"`
	} `json:"artwork"`
	Layers []struct {
		Name      string  `json:"name"`
		Opacity   float32 `json:"opacity"`
		ByPassDNA bool    `json:"bypass_dna"`
	} `json:"layers"`
	Metadata Metadata
}

func ReadConfig() (config Config) {
	Console.Loading("reading %v", CONFIG_FILE_PATH)

	// reading configuration file
	data, err := ioutil.ReadFile(CONFIG_FILE_PATH)
	if err != nil {
		Console.Error("unable to read %v due to an error: %v", CONFIG_FILE_PATH, err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		Console.Error("unable to read %v due to an error: %v", CONFIG_FILE_PATH, err)
	}

	// validating delimitter
	if slices.Contains(INVALID_DELIMITTERS, config.Delimitter) {
		Console.Warning("invalid delimitter detected")
		Console.Loading("defaulting value to `%v`", DEFAULT_DELIMITTER)
		config.Delimitter = DEFAULT_DELIMITTER
	}

	// validating artwork
	if config.Artwork.CompressionLevel > 3 {
		Console.Warning("invalid artwork-compression-level detected")
		Console.Loading("defaulting value to %v", 2)
		config.Artwork.CompressionLevel = 2
	}
	if config.Artwork.Width < 0 {
		config.Artwork.Width = DEFAULT_ARTWORK_WIDTH
	}
	if config.Artwork.Height < 0 {
		config.Artwork.Height = DEFAULT_ARTWORK_HEIGHT
	}

	// validating and formatting layers-directory & output-directory
	if !ValidathPath(config.LayersDirectory) {
		Console.Warning("invalid layers-directory detected")
		Console.Loading("defaulting value to %v", DEFAULT_LAYERS_DIRECTORY)
		config.LayersDirectory = DEFAULT_LAYERS_DIRECTORY
	}
	config.LayersDirectory = FormatPath(config.LayersDirectory)

	if !ValidathPath(config.OutputDirectory) {
		Console.Warning("invalid output-directory detected")
		Console.Loading("defaulting value to %v", DEFAULT_OUTPUT_DIRECTORY)
		config.OutputDirectory = DEFAULT_OUTPUT_DIRECTORY
	}
	config.OutputDirectory = FormatPath(config.OutputDirectory)

	// validating gif
	if config.Gif.Delay < 10 {
		config.Gif.Delay = DEFAULT_GIF_DELAY
	}
	if config.Gif.Repeat < 0 {
		config.Gif.Repeat = DEFAULT_GIF_REPEAT
	}

	Console.Success("configuration done")

	return config
}

func ReadLayers(config Config) (layers []Layer) {

	Console.Loading("reading %v directory", config.LayersDirectory)

	root_layers_directory := config.LayersDirectory
	layer_directories := config.Layers

	total_count := 0

	for _, layer := range layer_directories {
		layer_directory := root_layers_directory + layer.Name

		files, err := os.ReadDir(layer_directory)
		if err != nil {
			Console.Error("unable to read %v directory due to an error: %v", config.LayersDirectory, err)
		}

		elements := []Element{}

		for _, file := range files {
			if !file.IsDir() {
				// get element path
				file_path := layer_directory + "/" + file.Name()

				// get element name & weight
				file_name_partials := strings.Split(file.Name(), ".")
				file_name_with_weight := strings.Split(file_name_partials[0], config.Delimitter)
				file_name := file_name_partials[0]

				weight := DEFAULT_ELEMENT_WEIGHT

				if len(file_name_with_weight) == 2 && file_name_with_weight[1] != "" {
					file_name = file_name_with_weight[0]
					converted_weight, err := strconv.Atoi(file_name_with_weight[1])
					if err != nil {
						log.Fatal(err)
					}

					if converted_weight > 0 {
						weight = converted_weight
					}

				}
				elements = append(
					elements,
					Element{
						Name:   strings.ToLower(file_name),
						Path:   file_path,
						Weight: weight,
					},
				)
			}
		}

		sort.Slice(elements, func(i, j int) bool { return elements[i].Weight < elements[j].Weight })

		total_count += len(elements)

		opacity := layer.Opacity

		// validating layer opacity
		if layer.Opacity < 0 || layer.Opacity > 1 {
			opacity = DEFAULT_LAYER_OPACITY
		}

		layers = append(
			layers,
			Layer{
				Name:      strings.ToLower(layer.Name),
				Elements:  elements,
				Opacity:   opacity,
				ByPassDNA: layer.ByPassDNA,
			},
		)
	}

	Console.Success("%v layers and %v elements detected.", len(layers), total_count)

	return layers
}

func Init() (config Config, layers []Layer) {

	config = ReadConfig()

	// create directories
	if _, err := os.Stat(config.OutputDirectory); err != nil {
		os.Mkdir(config.OutputDirectory, os.ModeAppend)
	}
	if _, err := os.Stat(config.OutputDirectory + "images"); err != nil {
		os.Mkdir(config.OutputDirectory+"images", os.ModeAppend)
	}
	if config.Metadata.Export {
		if _, err := os.Stat(config.OutputDirectory + "json"); err != nil {
			os.Mkdir(config.OutputDirectory+"json", os.ModeAppend)
		}
	}
	if config.Gif.Export {
		if _, err := os.Stat(config.OutputDirectory + "gifs"); err != nil {
			os.Mkdir(config.OutputDirectory+"gifs", os.ModeAppend)
		}
	}
	if config.Gif.Export && config.Metadata.Export {
		if _, err := os.Stat(config.OutputDirectory + "gif-json"); err != nil {
			os.Mkdir(config.OutputDirectory+"gif-json", os.ModeAppend)
		}
	}

	layers = ReadLayers(config)

	fmt.Println()

	return config, layers
}
