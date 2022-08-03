package main

import (
	"encoding/json"
	"os"
	"time"
)

type Metadata struct {
	Export          bool   `json:"export"`
	Network         string `json:"network"`
	Description     string `json:"description"`
	BaseUrl         string `json:"base_url"`
	BackgroundColor string `json:"background_color"`
	YoutubeUrl      string `json:"youtube_url"`
	Solana          struct {
		Symbol               string `json:"symbol"`
		SellerFeeBasisPoints uint32 `json:"seller_fee_basis_points"`
		ExternalUrl          string `json:"external_url"`
		Creators             []struct {
			Address string `json:"address"`
			Share   uint32 `json:"share"`
		} `json:"creators"`
	} `json:"solana"`
}

type Attribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}
type File struct {
	Uri  string `json:"uri"`
	Type string `json:"type"`
}

type EthMetaData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Dna         string `json:"dna"`
	Edition     string `json:"edition"`
	Date        int64  `json:"date"`
	Attributes  []Attribute
	Compiler    string `json:"compiler"`
}

type SolProperties struct {
	Files    []File `json:"files"`
	Category string `json:"category"`
	Creators []struct {
		Address string `json:"address"`
		Share   uint32 `json:"share"`
	} `json:"creators"`
}

type SolMetaData struct {
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	Description          string `json:"description"`
	SellerFeeBasisPoints uint32 `json:"seller_fee_basis_points"`
	Image                string `json:"image"`
	ExternalUrl          string `json:"external_url"`
	Edition              string `json:"edition"`
	Attributes           []Attribute
	Properties           SolProperties `json:"properties"`
}

func GenerateArtworkMetaData(config *Config, dna *DNA, file_name *string, file_type string, edition *string) {
	metadata := config.Metadata

	var data any

	attributes := []Attribute{}

	// extracting attributes from dna
	for _, base := range *dna {
		attributes = append(attributes, Attribute{
			TraitType: base.Type,
			Value:     base.Name,
		})
	}

	if metadata.Network == "sol" {
		data = SolMetaData{
			Name:                 config.Artwork.Prefix + "-" + *edition,
			Symbol:               metadata.Solana.Symbol,
			Description:          metadata.Description,
			SellerFeeBasisPoints: metadata.Solana.SellerFeeBasisPoints,
			Image:                *file_name,
			ExternalUrl:          metadata.Solana.ExternalUrl,
			Edition:              *edition,
			Attributes:           attributes,
			Properties: SolProperties{
				Files: []File{
					{
						Uri:  *file_name,
						Type: file_type,
					},
				},
				Category: "image",
				Creators: metadata.Solana.Creators,
			},
		}
	} else {
		data = EthMetaData{
			Name:        config.Artwork.Prefix + "-" + *edition,
			Description: metadata.Description,
			Image:       *file_name,
			Dna:         dna.GetId(),
			Edition:     *edition,
			Date:        time.Now().Unix(),
			Attributes:  attributes,
			Compiler:    "Artwork Generator Developed By Theshawa Dasun",
		}
	}

	output_file_name := GetFileName(*edition+".json", config.Artwork.Prefix)

	folder_name := "json"
	if file_type == "image/gif" {
		folder_name = "gif-json"
	}

	output_path := config.OutputDirectory + folder_name + "/" + output_file_name

	file, err := os.Create(output_path)
	if err != nil {
		Console.Error("unable to create %v due to an error: %v", output_path, err)
	}
	defer file.Close()

	bytes, err := json.Marshal(data)
	if err != nil {
		Console.Error("unable to write %v due to an error: %v", output_path, err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		Console.Error("unable to write %v due to an error: %v", output_path, err)
	}

	if config.Debugging {
		Console.Log("(%v) metadata created at %v", edition, output_path)
	}
}
