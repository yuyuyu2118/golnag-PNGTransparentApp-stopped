package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/bmp"
)

func convertBMPtoPNG(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer inputFile.Close()

	img, err := bmp.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("error decoding BMP image: %v", err)
	}

	bounds := img.Bounds()
	transparentImg := image.NewRGBA(bounds)

	colorsToMakeTransparent := []color.RGBA{
		{0x81, 0x79, 0x7D, 0xFF},
		{0x69, 0x71, 0x89, 0xFF},
		{0x69, 0x89, 0x91, 0xFF},
		{0x6B, 0x8C, 0x94, 0xFF},
		{0x95, 0xA9, 0xD1, 0xFF},
	}
	transparent := color.RGBA{0x00, 0x00, 0x00, 0x00}

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c := img.At(x, y)

			makeTransparent := false
			for _, bgColor := range colorsToMakeTransparent {
				if c == bgColor {
					makeTransparent = true
					break
				}
			}

			if makeTransparent {
				transparentImg.Set(x, y, transparent)
			} else {
				transparentImg.Set(x, y, c)
			}
		}
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, transparentImg)
	if err != nil {
		return fmt.Errorf("error encoding PNG image: %v", err)
	}

	return nil
}

func convertTransparentPNG(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer inputFile.Close()

	img, err := png.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("error decoding PNG image: %v", err)
	}

	bounds := img.Bounds()
	transparentImg := image.NewRGBA(bounds)

	colorToMakeTransparent := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	transparent := color.RGBA{0x00, 0x00, 0x00, 0x00}

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c := img.At(x, y)

			if c == colorToMakeTransparent {
				transparentImg.Set(x, y, transparent)
			} else {
				transparentImg.Set(x, y, c)
			}
		}
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, transparentImg)
	if err != nil {
		return fmt.Errorf("error encoding PNG image: %v", err)
	}

	return nil
}

func processImagesInFolder(inputFolder, outputFolder, format string) error {
	err := filepath.Walk(inputFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through input folder: %v", err)
		}

		if !info.IsDir() {
			inputFilePath := path
			relPath, _ := filepath.Rel(inputFolder, path)
			outputFilePath := filepath.Join(outputFolder, relPath)

			if format == "BMP" {
				err = convertBMPtoPNG(inputFilePath, outputFilePath)
			} else {
				err = convertTransparentPNG(inputFilePath, outputFilePath)
			}

			if err != nil {
				return fmt.Errorf("error converting image: %v", err)
			}

			fmt.Printf("Successfully converted %s to %s using format %s\n", inputFilePath, outputFilePath, format)
		}

		return nil
	})

	return err
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: img-convert inputFolderPath outputFolderPath format(BMP|PNG)")
		os.Exit(1)
	}

	inputFolderPath := os.Args[1]
	outputFolderPath := os.Args[2]
	format := os.Args[3]

	if format != "BMP" && format != "PNG" {
		fmt.Println("Invalid format. Please specify either 'BMP' or 'PNG'.")
		os.Exit(1)
	}

	err := processImagesInFolder(inputFolderPath, outputFolderPath, format)

	if err != nil {
		fmt.Printf("Error converting images: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted all images in folder %s to %s using format %s\n", inputFolderPath, outputFolderPath, format)
}
