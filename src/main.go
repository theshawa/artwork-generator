package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	ct "github.com/daviddengcn/go-colortext"
)

type Element struct {
	Name   string
	Path   string
	Weight int
}
type Layer struct {
	Name      string
	Opacity   float32
	ByPassDNA bool
	Elements  []Element
}

var Scanner = bufio.NewScanner(os.Stdin)

func main() {
	// seeding random
	rand.Seed(time.Now().UnixNano())

	config, layers := Init()

	// calculating maximum possible dnas count
	maximum_count := 1
	for _, layer := range layers {
		maximum_count *= len(layer.Elements)
	}

	fmt.Printf("enter artworks count (maximum is %v): ", maximum_count)
	Scanner.Scan()

	count_text := Scanner.Text()
	count, err := strconv.Atoi(count_text)

	if err != nil {
		Console.Error("unable to read count due to an err: %v", err)
	}

	if count > maximum_count {
		Console.Warning("maximum artworks count is %v", maximum_count)
		Console.Loading("defaulting value to %v", maximum_count)
		count = maximum_count
	}

	fmt.Println()
	ct.Foreground(ct.Blue, true)

	if config.Gif.Export {
		log.Print("⚙️  gif export enabled.")
	} else {
		log.Print("⚙️  gif export disabled.")
	}
	if config.Metadata.Export {
		log.Print("⚙️  json export enabled.")
	} else {
		log.Print("⚙️  json export disabled.")
	}
	if config.Debugging {
		log.Print("⚙️  debugging enabled.")
	} else {
		log.Print("⚙️  debugging disabled.")
	}

	ct.ResetColor()

	fmt.Println()

	Console.Loading("starting engine")

	fmt.Println()

	// generating dnas
	dna_queue := make(chan DNA)

	go GenerateDNAs(&layers, count, &dna_queue, &config)

	// generating_artworks
	generated_artworks_count := 0
	for dna := range dna_queue {

		CreateArtwork(&dna, &config, generated_artworks_count)

		if config.Gif.Export {
			CreateGifArtwork(&dna, &config, generated_artworks_count)
		}

		generated_artworks_count++
	}

	fmt.Println()

	Console.Success("%v unique artworks generated successfully at %v", count, config.OutputDirectory)

	ct.Foreground(ct.Cyan, true)

	fmt.Println("\nThanks for using this tool.")
	fmt.Println("Developed by Theshawa Dasun as a personal project. Visit https://theshaw.cf to see more.")
	fmt.Println("No copyrights reserved.")
	fmt.Println("\nHave a great day :)")

	ct.ResetColor()

	fmt.Print("\npress any key to continue...")
	fmt.Scanln()

	Console.Loading("closing application")
}
