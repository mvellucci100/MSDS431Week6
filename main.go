package main

import (
	"fmt"
	"goroutines_pipeline/image_processing"
	"image"
	"os"
	"strings"
	"time"
)

// Job structure to hold image paths and image object
type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
}

// Error checking: Validate image file (ensure it's a JPEG)
func validateImage(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	if !strings.HasSuffix(filePath, ".jpeg") && !strings.HasSuffix(filePath, ".jpg") {
		return fmt.Errorf("invalid file type: %s. Only .jpeg and .jpg files are allowed", filePath)
	}

	return nil
}

// Function to check if a directory exists
func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// Sequential implementation of pipeline stages
func runSequential(imagePaths []string, outputDir string) {
	startPipeline := time.Now() // Start timer for sequential mode

	for _, path := range imagePaths {
		if err := validateImage(path); err != nil {
			fmt.Println(err)
			continue
		}

		start := time.Now()
		job := Job{
			InputPath: path,
			OutPath:   strings.Replace(path, "images/", outputDir, 1),
		}
		job.Image = imageprocessing.ReadImage(path)
		fmt.Printf("Load stage took %v\n", time.Since(start))

		start = time.Now()
		job.Image = imageprocessing.Resize(job.Image)
		fmt.Printf("Resize stage took %v\n", time.Since(start))

		start = time.Now()
		job.Image = imageprocessing.Grayscale(job.Image)
		fmt.Printf("Grayscale stage took %v\n", time.Since(start))

		start = time.Now()
		imageprocessing.WriteImage(job.OutPath, job.Image)
		fmt.Printf("Save stage took %v\n", time.Since(start))

		fmt.Println("Job completed successfully!")
	}

	totalElapsed := time.Since(startPipeline)
	fmt.Printf("Total sequential pipeline time: %v\n", totalElapsed)
}

// Parallel implementation of pipeline stages
func runParallel(imagePaths []string, outputDir string) {
	startPipeline := time.Now() // Start timer for parallel mode

	channel1 := loadImage(imagePaths, outputDir)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	writeResults := saveImage(channel3)

	// Wait for pipeline to complete
	for success := range writeResults {
		if success {
			fmt.Println("Success!")
		} else {
			fmt.Println("Failed!")
		}
	}

	totalElapsed := time.Since(startPipeline)
	fmt.Printf("Total parallel pipeline time: %v\n", totalElapsed)
}

func loadImage(paths []string, outputDir string) <-chan Job {
	out := make(chan Job)
	go func() {
		for _, p := range paths {
			if err := validateImage(p); err != nil {
				fmt.Println(err)
				continue
			}
			job := Job{InputPath: p, OutPath: strings.Replace(p, "images/", outputDir, 1)}
			job.Image = imageprocessing.ReadImage(p)
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			start := time.Now()
			job.Image = imageprocessing.Resize(job.Image)
			fmt.Printf("Resize stage took %v\n", time.Since(start))
			out <- job
		}
		close(out)
	}()
	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			start := time.Now()
			job.Image = imageprocessing.Grayscale(job.Image)
			fmt.Printf("Grayscale stage took %v\n", time.Since(start))
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan bool {
	out := make(chan bool)
	go func() {
		for job := range input {
			start := time.Now()
			imageprocessing.WriteImage(job.OutPath, job.Image)
			fmt.Printf("Save stage took %v\n", time.Since(start))
			out <- true
		}
		close(out)
	}()
	return out
}

func main() {
	outputDir := "./images/output/"
	if !directoryExists(outputDir) {
		fmt.Printf("Output directory does not exist: %s\n", outputDir)
		return
	}

	imagePaths := []string{"images/watermelon.jpg", "images/apple.jpg", "images/blueberry.jpg", "images/lemon.jpg"}

	var mode string
	fmt.Println("Enter mode (sequential/parallel):")
	fmt.Scanln(&mode)

	if mode == "sequential" {
		fmt.Println("Running pipeline in sequential mode...")
		runSequential(imagePaths, outputDir)
	} else if mode == "parallel" {
		fmt.Println("Running pipeline in parallel mode...")
		runParallel(imagePaths, outputDir)
	} else {
		fmt.Println("Invalid mode. Please enter 'sequential' or 'parallel'.")
	}
}
