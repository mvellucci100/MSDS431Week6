package main

import (
	"fmt"
	"os"
	"testing"
	"strings"
	"image"
	"image/color"
	"github.com/stretchr/testify/assert"
	"goroutines_pipeline/image_processing"
)

// Mock image processing functions for testing

// MockReadImage simulates reading an image (just returns a red 1x1 image).
func mockReadImage(path string) image.Image {
	// Create a simple 1x1 image for testing
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255}) // Red pixel
	return img
}

// MockWriteImage simulates writing an image (no-op for testing).
func mockWriteImage(outPath string, img image.Image) {
	// Print the output path for testing
	fmt.Println("Mock writing image to", outPath)
}

// Replace the original ReadImage and WriteImage functions with mocks for testing.
var ReadImage = mockReadImage
var WriteImage = mockWriteImage

// Error checking: Validate image file (ensure it's a JPEG)
func mockValidateImage(filePath string) error {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Check if the file has a valid JPEG extension
	if !strings.HasSuffix(filePath, ".jpeg") && !strings.HasSuffix(filePath, ".jpg") {
		return fmt.Errorf("invalid file type: %s. Only .jpeg and .jpg files are allowed", filePath)
	}

	return nil
}

// Output error checking: Function to check if a directory exists
func mockDirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func mockLoadImage(paths []string, outputDir string) <-chan Job {
	out := make(chan Job)
	go func() {
		for _, p := range paths {
			// Validate the image file before processing
			if err := validateImage(p); err != nil {
				fmt.Println(err)
				continue // Skip the file if it's invalid
			}
			job := Job{
				InputPath: p,
				OutPath:   strings.Replace(p, "images/", outputDir, 1),
				Image:     ReadImage(p), // Use the mocked ReadImage function
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func mockResize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			// Simulate resizing the image
			job.Image = imageprocessing.Resize(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func mockConvertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			// Simulate converting the image to grayscale
			job.Image = imageprocessing.Grayscale(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func mockSaveImage(input <-chan Job) <-chan bool {
	out := make(chan bool)
	go func() {
		for job := range input {
			// Simulate saving the image
			WriteImage(job.OutPath, job.Image) // Use the mocked WriteImage function
			out <- true
		}
		close(out)
	}()
	return out
}

// Unit Test for validateImage function
func TestValidateImage(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-image-*.jpg")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up

	// Test valid JPEG file
	err = validateImage(tempFile.Name())
	assert.NoError(t, err)

	// Test non-existing file
	err = validateImage("nonexistent.jpg")
	assert.Errorf(t, err, "file does not exist")

	// Test invalid file type
	err = validateImage("invalid.txt")
	assert.Errorf(t, err, "invalid file type")
}

// Unit Test for directoryExists function
func TestDirectoryExists(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "test-dir")
	assert.NoError(t, err)
	defer os.RemoveAll(dir) // Clean up

	// Test an existing directory
	exists := directoryExists(dir)
	assert.True(t, exists)

	// Test a non-existing directory
	exists = directoryExists("nonexistent-dir")
	assert.False(t, exists)
}

// Unit Test for loadImage function
func TestLoadImage(t *testing.T) {
	// Define test input paths
	paths := []string{
		"images/watermelon.jpg", // Valid path
		"invalid/path/to/image.jpg", // Invalid path
	}

	outputDir := "./images/output/"

	// Mocking the loadImage behavior by directly invoking it
	resultChan := loadImage(paths, outputDir)

	// Test the results from the loadImage function
	for result := range resultChan {
		if result.InputPath == "images/watermelon.jpg" {
			assert.NotNil(t, result.Image, "Expected image to be loaded")
		} else {
			assert.Nil(t, result.Image, "Expected invalid image path to return nil image")
		}
	}
}

// Unit Test for resize function
func TestResize(t *testing.T) {
	input := make(chan Job)
	go func() {
		// Create a mock job with a simple image
		job := Job{
			InputPath: "test.jpg",
			Image:     mockReadImage("test.jpg"),
			OutPath:   "./output/test_resized.jpg",
		}
		input <- job
		close(input)
	}()

	output := resize(input)

	for job := range output {
		assert.NotNil(t, job.Image, "Expected image to be resized")
	}
}

// Unit Test for convertToGrayscale function
func TestConvertToGrayscale(t *testing.T) {
	input := make(chan Job)
	go func() {
		// Create a mock job with a simple image
		job := Job{
			InputPath: "test.jpg",
			Image:     mockReadImage("test.jpg"),
			OutPath:   "./output/test_grayscale.jpg",
		}
		input <- job
		close(input)
	}()

	output := convertToGrayscale(input)

	for job := range output {
		// Assert the image has been converted (mocking the process here)
		assert.NotNil(t, job.Image, "Expected image to be converted to grayscale")
	}
}




