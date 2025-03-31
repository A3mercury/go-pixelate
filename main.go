package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
)

type pixel struct {
	R, G, B, A uint32
}

func quantize(img image.Image, numColors int) image.Image {
	bounds := img.Bounds()
	pixels := make([]pixel, 0, bounds.Dx()*bounds.Dy())

	// Collect all pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels = append(pixels, pixel{r, g, b, a})
		}
	}

	// K-means clustering
	centroids := make([]pixel, numColors)
	// Initialize centroids (simple: spread across pixel list)
	for i := 0; i < numColors && i < len(pixels); i++ {
		centroids[i] = pixels[i*len(pixels)/numColors]
	}

	for iter := 0; iter < 5; iter++ {
		clusters := make([][]pixel, numColors)
		for i := range clusters {
			clusters[i] = []pixel{}
		}

		// Assign pixels to nearest centroid
		for _, p := range pixels {
			minDist := uint64(1<<63 - 1)
			clusterIdx := 0
			for i, c := range centroids {
				dist := sqDist(p, c)
				if dist < minDist {
					minDist = dist
					clusterIdx = i
				}
			}
			clusters[clusterIdx] = append(clusters[clusterIdx], p)
		}

		// Update centroids
		for i, cluster := range clusters {
			if len(cluster) == 0 {
				continue
			}
			var r, g, b, a uint64
			for _, p := range cluster {
				r += uint64(p.R)
				g += uint64(p.G)
				b += uint64(p.B)
				a += uint64(p.A)
			}
			count := uint64(len(cluster))
			centroids[i] = pixel{
				R: uint32(r / count),
				G: uint32(g / count),
				B: uint32(b / count),
				A: uint32(a / count),
			}
		}
	}

	// Create new image with quantized colors
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			p := pixel{r, g, b, a}
			minDist := uint64(1<<63 - 1)
			var nearest pixel
			for _, c := range centroids {
				if dist := sqDist(p, c); dist < minDist {
					minDist = dist
					nearest = c
				}
			}
			result.Set(x, y, color.RGBA{
				R: uint8(nearest.R >> 8),
				G: uint8(nearest.G >> 8),
				B: uint8(nearest.B >> 8),
				A: uint8(nearest.A >> 8),
			})
		}
	}
	return result
}

func sqDist(p1, p2 pixel) uint64 {
	dr := int64(p1.R) - int64(p2.R)
	dg := int64(p1.G) - int64(p2.G)
	db := int64(p1.B) - int64(p2.B)
	da := int64(p1.A) - int64(p2.A)
	return uint64(dr*dr + dg*dg + db*db + da*da)
}

func pixelate(img image.Image, pixelSize int) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	result := image.NewRGBA(bounds)

	// Manual pixelation to avoid interpolation
	for y := 0; y < height; y += pixelSize {
		for x := 0; x < width; x += pixelSize {
			// Sample one pixel per block
			c := img.At(x, y)
			// Fill the block with that exact color
			for dy := 0; dy < pixelSize && y+dy < height; dy++ {
				for dx := 0; dx < pixelSize && x+dx < width; dx++ {
					result.Set(x+dx, y+dy, c)
				}
			}
		}
	}
	return result
}

func main() {
	if len(os.Args) != 5 {
		log.Fatal("Usage: go run main.go <input-image> <output-image> <pixel-size> <num-colors>")
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	pixelSize := 10
	numColors := 16 // Default back to 16 for flexibility

	if n, err := strconv.Atoi(os.Args[3]); err == nil {
		pixelSize = n
	}
	if n, err := strconv.Atoi(os.Args[4]); err == nil {
		numColors = n
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Error opening input file: %v", err)
	}
	defer inputFile.Close()

	img, format, err := image.Decode(inputFile)
	if err != nil {
		log.Fatalf("Error decoding image: %v", err)
	}

	// Apply quantization first, then pixelation
	quantizedImg := quantize(img, numColors)
	pixelatedImg := pixelate(quantizedImg, pixelSize)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	if numColors == 2 {
		log.Println("Note: For exactly 2 colors, PNG output is recommended to avoid compression artifacts.")
	}

	switch format {
	case "jpeg":
		err = jpeg.Encode(outputFile, pixelatedImg, &jpeg.Options{Quality: 100})
	case "png":
		err = png.Encode(outputFile, pixelatedImg)
	default:
		log.Fatalf("Unsupported image format: %s", format)
	}

	if err != nil {
		log.Fatalf("Error encoding output image: %v", err)
	}

	log.Printf("Pixelated image with %d colors saved to %s", numColors, outputPath)
}
