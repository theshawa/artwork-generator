package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
)

type Base struct {
	Name    string
	Path    string
	Opacity float32
	Hidden  bool
	Type    string
}

type DNA []Base

func (dna *DNA) GetCleanString() (clean_string string) {
	for _, base := range *dna {
		if !base.Hidden {
			clean_string = fmt.Sprintf("%v.%v%v", clean_string, base.Type, base.Name)
		}
	}
	return clean_string
}

func (dna *DNA) GetId() string {
	hash := sha1.New()
	hash.Write([]byte(dna.GetCleanString()))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func CreateDNA(layers *[]Layer) DNA {
	dna := DNA{}
	for _, layer := range *layers {
		elements := layer.Elements

		// Generate Total Weight
		total_weight := 0
		for _, element := range elements {
			total_weight += element.Weight
		}
		random := rand.Intn(total_weight)

		// Get best element according to the element's weight
		for _, element := range elements {
			random -= element.Weight
			if random < 0 {
				dna = append(dna, Base{Path: element.Path, Name: element.Name, Opacity: layer.Opacity, Hidden: layer.ByPassDNA, Type: layer.Name})
				break
			}
		}
	}
	return dna
}

// Check Uniqueness
func IsUnique(dna_list *[]string, check_dna *DNA) bool {
	for _, dna := range *dna_list {
		if dna == check_dna.GetCleanString() {
			return false
		}
	}
	return true
}

func GenerateDNAs(layers []Layer, count int, dna_queue chan DNA, config Config) {
	dna_list := []string{}

	for i := 0; i < count; i++ {
		for {
			new_dna := CreateDNA(&layers)
			dna_string := new_dna.GetCleanString()
			if IsUnique(&dna_list, &new_dna) {
				dna_list = append(dna_list, dna_string)
				dna_queue <- new_dna
				if config.Debugging {
					Console.Log("(%v) dna created: %v", i+1, new_dna.GetCleanString())
				}
				break
			}
			if config.Debugging {
				Console.Warning("(%v) dna exists: %v", i+1, new_dna.GetCleanString())
			}
		}
	}
	close(dna_queue)
}
