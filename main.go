package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"

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

	//ここに色の追加をする
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

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: bmp2png inputDirectory outputDirectory")
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputDir := os.Args[2]

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".bmp" {
			relPath, _ := filepath.Rel(inputDir, path)
			outputPath := filepath.Join(outputDir, strings.TrimSuffix(relPath, filepath.Ext(relPath))+".png")

			err = os.MkdirAll(filepath.Dir(outputPath), 0755)
			if err != nil {
				return fmt.Errorf("error creating output directory: %v", err)
			}

			err = convertBMPtoPNG(path, outputPath)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully converted %s to %s\n", path, outputPath)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error processing files:", err)
		os.Exit(1)
	}
}
