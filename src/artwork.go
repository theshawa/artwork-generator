package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"strconv"
)

func GetOpacity(opacity float32) float32 {
	if opacity > 1 {
		return opacity / 100
	}
	return opacity
}

func GetFileName(name string, prefix string) string {
	if prefix != "" {
		return prefix + "-" + name
	}
	return name
}

func GetCompressionLevel(val uint8) png.CompressionLevel {
	cl := png.DefaultCompression
	switch val {
	case 0:
		cl = png.BestSpeed
	case 1:
		cl = png.NoCompression
	case 2:
		cl = png.DefaultCompression
	case 3:
		cl = png.BestCompression
	}
	return png.CompressionLevel(cl)

}

func CreateArtwork(dna DNA, config Config, index int) {
	container := image.Rect(0, 0, int(config.Artwork.Width), int(config.Artwork.Height))
	img := image.NewRGBA(container)
	number := strconv.Itoa(index + 1)

	if config.Artwork.Background != "" {
		bg, err := ParseHexColor(config.Artwork.Background)
		if err != nil {
			Console.Error("unable to get rgba color-value of %v due to an error: %v", config.Artwork.Background, err)
		}
		draw.Draw(img, container, &image.Uniform{bg}, image.Point{}, draw.Src)
	}

	for _, base := range dna {
		file, err := os.Open(base.Path)
		if err != nil {
			Console.Error("unable to open %v due to an error: %v", base.Path, err)
		}
		defer file.Close()

		sub_img, _, err := image.Decode(file)
		if err != nil {
			Console.Error("unable to read %v due to an error: %v", base.Path, err)
		}
		mask := image.NewUniform(color.Alpha{uint8(255 * GetOpacity(base.Opacity))})
		draw.DrawMask(img, container, sub_img, image.Point{}, mask, image.Point{}, draw.Over)
	}
	encoder := png.Encoder{CompressionLevel: GetCompressionLevel(config.Artwork.CompressionLevel)}
	file_name := GetFileName(number+".png", config.Artwork.Prefix)
	file_path := config.OutputDirectory + "images/" + file_name

	writer, err := os.Create(file_path)
	if err != nil {
		Console.Error("unable to create %v due to an error: %v", file_path, err)
	}
	defer writer.Close()

	err = encoder.Encode(writer, img)
	if err != nil {
		Console.Error("unable to encode %v due to an error: %v", img, err)
	}
	if config.Metadata.Export {
		GenerateArtworkMetaData(&config, dna, file_name, "image/png", number)
	}

	Console.Success("(%v) new artwork generated with dna: %v", index+1, dna.GetId())
}

func CreateGifArtwork(dna DNA, config Config, index int) {
	container := image.Rect(0, 0, int(config.Artwork.Width), int(config.Artwork.Height))
	delays := []int{}
	gif_images := []*image.Paletted{}
	number := strconv.Itoa(index + 1)
	disposals := []byte{}

	for _, base := range dna {
		file, err := os.Open(base.Path)
		if err != nil {
			Console.Error("unable to open %v due to an error: %v", base.Path, err)
		}
		defer file.Close()

		gif_img, err := png.Decode(file)
		if err != nil {
			Console.Error("unable to read %v due to an error: %v", base.Path, err)
		}

		palettedImg := image.NewPaletted(container, palette.Plan9)

		// Add Opacity
		mask := image.NewUniform(color.Alpha{uint8(255 * GetOpacity(base.Opacity))})
		draw.DrawMask(palettedImg, container, gif_img, image.Point{}, mask, image.Point{}, draw.Over)

		gif_images = append(gif_images, palettedImg)

		// Add Delay
		disposals = append(disposals, gif.DisposalPrevious)
		delays = append(delays, int(config.Gif.Delay))

	}
	anim := gif.GIF{
		Delay:     delays,
		Image:     gif_images,
		LoopCount: config.Gif.Repeat,
		Disposal:  disposals,
		Config: image.Config{
			Width:  int(config.Artwork.Width),
			Height: int(config.Artwork.Height),
		},
	}
	gif_file_name := GetFileName(number+".gif", config.Artwork.Prefix)
	gif_file_path := config.OutputDirectory + "gifs/" + gif_file_name
	writer, err := os.Create(gif_file_path)
	if err != nil {
		Console.Error("unable to create %v due to an error: %v", gif_file_path, err)
	}
	defer writer.Close()

	err = gif.EncodeAll(writer, &anim)
	if err != nil {
		Console.Error("unable to encode %v due to an error: %v", gif_file_path, err)
	}

	if config.Metadata.Export {
		GenerateArtworkMetaData(&config, dna, gif_file_name, "image/gif", number)
	}

	Console.Success("(%v) new gif generated with dna: %v", index+1, dna.GetId())
}
