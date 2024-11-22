# MSDS431Week6

# Image Processing Pipeline with Goroutines

This project demonstrates an image processing pipeline in Go, designed to process images sequentially or in parallel using goroutines. The program supports basic operations such as resizing, converting to grayscale, and saving images to a specified output directory. It also benchmarks the throughput time for each stage and the entire pipeline.

---

## Features
- **Sequential Mode**: Processes images one by one, stage by stage.
- **Parallel Mode**: Uses goroutines for concurrent processing of images through the pipeline.
- **Benchmarking**: Captures and reports the execution time for each pipeline stage and the overall process.
- **Error Handling**: Validates input files and checks output directory existence.

---

## Pipeline Stages
1. **Load Images**: Reads images from a specified input path.
2. **Resize Images**: Resizes images to a standard size.
3. **Convert to Grayscale**: Converts images to grayscale.
4. **Save Images**: Writes the processed images to the output directory.

---

## Requirements
- Go 1.18 or higher
- JPEG or JPG image files for input

---

## Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/mvellucci100/MSDS431Week6.git
   cd image-processing-pipeline

## Usage
1. Run the Program
 ```bash
   go run main.go
2. Test the Program
```go test -v

