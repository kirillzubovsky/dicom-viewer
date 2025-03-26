package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

type dicomSeries struct {
	seriesPath string
	frames     []image.Image
	metadata   map[string]string
}

type contrastImage struct {
	image.Image
	contrast float64
}

func (c *contrastImage) At(x, y int) color.Color {
	r, g, b, a := c.Image.At(x, y).RGBA()

	// Apply contrast adjustment
	contrast := c.contrast
	factor := (259 * (contrast + 255)) / (255 * (259 - contrast))

	r = uint32(factor*(float64(r)-128) + 128)
	g = uint32(factor*(float64(g)-128) + 128)
	b = uint32(factor*(float64(b)-128) + 128)

	return color.RGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run dicom_viewer.go <path_to_dicom_directory>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	seriesList, err := loadDICOMSeries(dirPath)
	if err != nil {
		log.Fatalf("Error loading DICOM series: %v", err)
	}

	if len(seriesList) == 0 {
		fmt.Println("No valid DICOM series found in the directory.")
		os.Exit(1)
	}

	// Initialize Fyne app
	myApp := app.New()
	w := myApp.NewWindow("DICOM Viewer")
	w.Resize(fyne.NewSize(800, 700))

	// State variables
	currentSeriesIndex := 0
	currentFrameIndex := 0
	zoomLevel := float32(1.0)
	contrastLevel := float64(0.0)

	// Loading indicator
	loadingIndicator := widget.NewLabel("")
	loadingIndicator.Hide()

	// Image display
	imgCanvas := canvas.NewImageFromImage(seriesList[currentSeriesIndex].frames[currentFrameIndex])
	imgCanvas.FillMode = canvas.ImageFillContain
	imgCanvas.SetMinSize(fyne.NewSize(512, 512))

	// Metadata label
	metadata := getMetadataFromMap(seriesList[currentSeriesIndex].metadata)
	metaLabel := widget.NewLabel(metadata)

	// Frame counter
	frameLabel := widget.NewLabel(fmt.Sprintf("Frame %d/%d", currentFrameIndex+1, len(seriesList[currentSeriesIndex].frames)))

	// Series selector
	seriesNames := make([]string, len(seriesList))
	for i, series := range seriesList {
		seriesNames[i] = filepath.Base(series.seriesPath)
	}

	// Sort series names numerically
	sort.Slice(seriesNames, func(i, j int) bool {
		// Extract numbers from series names (e.g., "SE000001" -> "000001")
		numI := strings.TrimPrefix(seriesNames[i], "SE")
		numJ := strings.TrimPrefix(seriesNames[j], "SE")
		return numI < numJ
	})

	// Create a map to find the original index for each series name
	seriesIndexMap := make(map[string]int)
	for i, series := range seriesList {
		seriesIndexMap[filepath.Base(series.seriesPath)] = i
	}

	seriesSelector := widget.NewSelect(seriesNames, func(selected string) {
		loadingIndicator.Show()
		loadingIndicator.SetText("Loading series...")
		w.Canvas().Refresh(w.Content())

		go func() {
			// Find the original index for the selected series
			if originalIndex, ok := seriesIndexMap[selected]; ok {
				currentSeriesIndex = originalIndex
				currentFrameIndex = 0
				updateDisplay(imgCanvas, metaLabel, frameLabel, seriesList, currentSeriesIndex, currentFrameIndex, contrastLevel)
				time.Sleep(100 * time.Millisecond) // Small delay for visual feedback
				loadingIndicator.Hide()
			}
		}()
	})
	seriesSelector.SetSelected(seriesNames[0])

	// Contrast slider
	contrastSlider := widget.NewSlider(-100, 100)
	contrastSlider.OnChanged = func(value float64) {
		contrastLevel = value
		updateDisplay(imgCanvas, metaLabel, frameLabel, seriesList, currentSeriesIndex, currentFrameIndex, contrastLevel)
	}
	contrastLabel := widget.NewLabel("Contrast: 0")

	// Auto-play for multi-frame
	playing := false
	var playBtn *widget.Button
	playBtn = widget.NewButton("Play", func() {
		playing = !playing
		if playing {
			playBtn.SetText("Pause")
			go func() {
				for playing && currentFrameIndex < len(seriesList[currentSeriesIndex].frames)-1 {
					currentFrameIndex++
					updateDisplay(imgCanvas, metaLabel, frameLabel, seriesList, currentSeriesIndex, currentFrameIndex, contrastLevel)
					time.Sleep(100 * time.Millisecond) // Frame rate
				}
			}()
		} else {
			playBtn.SetText("Play")
		}
	})

	// Navigation buttons
	prevFrameBtn := widget.NewButton("Prev Frame", func() {
		if currentFrameIndex > 0 {
			currentFrameIndex--
			updateDisplay(imgCanvas, metaLabel, frameLabel, seriesList, currentSeriesIndex, currentFrameIndex, contrastLevel)
		}
	})
	nextFrameBtn := widget.NewButton("Next Frame", func() {
		if currentFrameIndex < len(seriesList[currentSeriesIndex].frames)-1 {
			currentFrameIndex++
			updateDisplay(imgCanvas, metaLabel, frameLabel, seriesList, currentSeriesIndex, currentFrameIndex, contrastLevel)
		}
	})

	// Zoom controls
	zoomInBtn := widget.NewButton("Zoom In", func() {
		zoomLevel += 0.2
		imgCanvas.SetMinSize(fyne.NewSize(float32(512)*zoomLevel, float32(512)*zoomLevel))
		imgCanvas.Refresh()
	})
	zoomOutBtn := widget.NewButton("Zoom Out", func() {
		if zoomLevel > 0.2 {
			zoomLevel -= 0.2
			imgCanvas.SetMinSize(fyne.NewSize(float32(512)*zoomLevel, float32(512)*zoomLevel))
			imgCanvas.Refresh()
		}
	})
	resetZoomBtn := widget.NewButton("Reset Zoom", func() {
		zoomLevel = 1.0
		imgCanvas.SetMinSize(fyne.NewSize(512, 512))
		imgCanvas.Refresh()
	})

	// Export button
	exportBtn := widget.NewButton("Export Series", func() {
		// Show directory picker dialog
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}

			// Create series directory
			seriesDir := filepath.Join(uri.Path(), fmt.Sprintf("series_%s", filepath.Base(seriesList[currentSeriesIndex].seriesPath)))
			if err := os.MkdirAll(seriesDir, 0755); err != nil {
				dialog.ShowError(fmt.Errorf("failed to create directory: %v", err), w)
				return
			}

			// Export metadata
			metadata := getMetadataFromMap(seriesList[currentSeriesIndex].metadata)
			if err := os.WriteFile(filepath.Join(seriesDir, "metadata.txt"), []byte(metadata), 0644); err != nil {
				dialog.ShowError(fmt.Errorf("failed to write metadata: %v", err), w)
				return
			}

			// Export images
			for i, img := range seriesList[currentSeriesIndex].frames {
				// Create filename with frame number
				filename := fmt.Sprintf("frame_%03d.png", i+1)
				filepath := filepath.Join(seriesDir, filename)

				// Create file
				f, err := os.Create(filepath)
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to create file %s: %v", filename, err), w)
					continue
				}

				// Encode and save image
				if err := png.Encode(f, img); err != nil {
					f.Close()
					dialog.ShowError(fmt.Errorf("failed to encode image %s: %v", filename, err), w)
					continue
				}
				f.Close()
			}

			dialog.ShowInformation("Export Complete",
				fmt.Sprintf("Series exported to:\n%s\n\n%d images exported", seriesDir, len(seriesList[currentSeriesIndex].frames)),
				w)
		}, w)
	})

	// Layout
	controls := container.NewGridWithColumns(3,
		prevFrameBtn, nextFrameBtn, playBtn,
	)
	zoomControls := container.NewGridWithColumns(3,
		zoomInBtn, zoomOutBtn, resetZoomBtn,
	)
	contrastControls := container.NewHBox(
		widget.NewLabel("Contrast:"),
		contrastSlider,
		contrastLabel,
	)
	content := container.NewVBox(
		widget.NewLabel("Select Series:"),
		seriesSelector,
		loadingIndicator,
		metaLabel,
		frameLabel,
		imgCanvas,
		controls,
		zoomControls,
		contrastControls,
		exportBtn,
	)

	w.SetContent(container.New(layout.NewCenterLayout(), content))
	w.ShowAndRun()
}

func loadDICOMSeries(dirPath string) ([]dicomSeries, error) {
	var seriesList []dicomSeries

	// Look for DICOMDIR file
	dicomdirPath := filepath.Join(dirPath, "DICOMDIR")
	seriesPaths := make(map[string][]string) // seriesUID -> list of image paths

	log.Printf("Looking for DICOM files in: %s", dirPath)
	log.Printf("Checking for DICOMDIR at: %s", dicomdirPath)

	if _, err := os.Stat(dicomdirPath); err == nil {
		log.Println("Found DICOMDIR file, attempting to parse...")
		// Parse DICOMDIR to find series and their images
		dataset, err := dicom.ParseFile(dicomdirPath, nil)
		if err != nil {
			log.Printf("Error parsing DICOMDIR: %v", err)
		} else {
			// Traverse DICOMDIR directory records
			records, err := dataset.FindElementByTag(tag.DirectoryRecordSequence)
			if err == nil && records != nil {
				sequence := records.Value.GetValue().([]dicom.Dataset)
				log.Printf("Found %d directory records in DICOMDIR", len(sequence))
				for _, record := range sequence {
					recordType, _ := record.FindElementByTag(tag.DirectoryRecordType)
					if recordType != nil && recordType.Value.String() == "IMAGE" {
						// Get the referenced file path
						filePathElement, _ := record.FindElementByTag(tag.ReferencedFileID)
						if filePathElement != nil {
							filePathParts := filePathElement.Value.GetValue().([]string)
							filePath := filepath.Join(append([]string{dirPath}, filePathParts...)...)
							log.Printf("Found image file: %s", filePath)
							// Get the series UID
							seriesUIDElement, _ := record.FindElementByTag(tag.SeriesInstanceUID)
							if seriesUIDElement != nil {
								seriesUID := seriesUIDElement.Value.String()
								seriesPaths[seriesUID] = append(seriesPaths[seriesUID], filePath)
							}
						}
					}
				}
			}
		}
	} else {
		log.Printf("No DICOMDIR found: %v", err)
	}

	// If DICOMDIR parsing failed or wasn't found, fall back to directory traversal
	if len(seriesPaths) == 0 {
		log.Println("Falling back to directory traversal...")
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing path %s: %v", path, err)
				return err
			}
			if info.IsDir() {
				log.Printf("Found directory: %s", path)
				return nil
			}

			// Check for various DICOM file extensions
			lowerName := strings.ToLower(info.Name())
			if strings.HasSuffix(lowerName, ".dcm") ||
				strings.HasSuffix(lowerName, ".dicom") ||
				strings.HasSuffix(lowerName, ".ima") ||
				!strings.Contains(lowerName, ".") { // Some DICOM files have no extension
				log.Printf("Found potential DICOM file: %s (size: %d bytes)", path, info.Size())
				// Determine series by the parent directory (e.g., SE000007)
				seriesDir := filepath.Base(filepath.Dir(path))
				seriesPaths[seriesDir] = append(seriesPaths[seriesDir], path)
			} else {
				log.Printf("Skipping non-DICOM file: %s", path)
			}
			return nil
		})
		if err != nil {
			log.Printf("Error during directory traversal: %v", err)
			return nil, err
		}
	}

	log.Printf("Found %d series", len(seriesPaths))
	for seriesUID, paths := range seriesPaths {
		log.Printf("Series %s has %d files", seriesUID, len(paths))
	}

	// Process each series
	for seriesUID, imagePaths := range seriesPaths {
		frames := []image.Image{}
		var metadata map[string]string

		log.Printf("Processing series %s with %d files", seriesUID, len(imagePaths))
		for _, imagePath := range imagePaths {
			log.Printf("Extracting frames from: %s", imagePath)
			imgFrames, meta, err := extractFramesFromDICOM(imagePath)
			if err != nil {
				log.Printf("Error processing file %s: %v", imagePath, err)
				continue
			}
			log.Printf("Successfully extracted %d frames from %s", len(imgFrames), imagePath)
			frames = append(frames, imgFrames...)
			if metadata == nil {
				metadata = meta
			}
		}

		if len(frames) > 0 {
			log.Printf("Adding series %s with %d frames", seriesUID, len(frames))
			seriesList = append(seriesList, dicomSeries{
				seriesPath: seriesUID,
				frames:     frames,
				metadata:   metadata,
			})
		} else {
			log.Printf("No valid frames found for series %s", seriesUID)
		}
	}

	log.Printf("Total series loaded: %d", len(seriesList))
	return seriesList, nil
}

func extractFramesFromDICOM(filePath string) ([]image.Image, map[string]string, error) {
	log.Printf("Parsing DICOM file: %s", filePath)
	dataset, err := dicom.ParseFile(filePath, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse DICOM file: %v", err)
	}

	// Extract metadata
	metadata := make(map[string]string)
	patientName, _ := dataset.FindElementByTag(tag.PatientName)
	if patientName != nil {
		metadata["PatientName"] = patientName.Value.String()
	}
	studyDate, _ := dataset.FindElementByTag(tag.StudyDate)
	if studyDate != nil {
		metadata["StudyDate"] = studyDate.Value.String()
	}
	modality, _ := dataset.FindElementByTag(tag.Modality)
	if modality != nil {
		metadata["Modality"] = modality.Value.String()
	}
	rows, _ := dataset.FindElementByTag(tag.Rows)
	cols, _ := dataset.FindElementByTag(tag.Columns)
	if rows != nil && cols != nil {
		metadata["Dimensions"] = fmt.Sprintf("%d x %d", rows.Value, cols.Value)
	}

	// Extract orientation information
	orientation, _ := dataset.FindElementByTag(tag.ImageOrientationPatient)
	if orientation != nil {
		orientationValues, ok := orientation.Value.GetValue().([]string)
		if ok && len(orientationValues) >= 6 {
			// Convert string values to float64
			values := make([]float64, len(orientationValues))
			for i, v := range orientationValues {
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					values[i] = f
				}
			}

			// Calculate orientation based on the direction cosines
			// First row (x): values[0:3]
			// Second row (y): values[3:6]
			// Third row (z): cross product of first two rows
			zX := values[1]*values[5] - values[2]*values[4]
			zY := values[2]*values[3] - values[0]*values[5]
			zZ := values[0]*values[4] - values[1]*values[3]

			// Determine the plane based on the dominant component
			absX := abs(zX)
			absY := abs(zY)
			absZ := abs(zZ)

			if absX > absY && absX > absZ {
				metadata["Orientation"] = "Sagittal"
			} else if absY > absX && absY > absZ {
				metadata["Orientation"] = "Coronal"
			} else if absZ > absX && absZ > absY {
				metadata["Orientation"] = "Axial"
			} else {
				metadata["Orientation"] = "Unknown"
			}
		} else {
			log.Printf("Invalid orientation values format or insufficient values")
			metadata["Orientation"] = "Unknown"
		}
	} else {
		log.Printf("No orientation information found")
		metadata["Orientation"] = "Unknown"
	}

	// Extract series description
	seriesDesc, _ := dataset.FindElementByTag(tag.SeriesDescription)
	if seriesDesc != nil {
		metadata["SeriesDescription"] = seriesDesc.Value.String()
	}

	// Extract frames
	pixelDataElement, err := dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		return nil, metadata, fmt.Errorf("no pixel data found: %v", err)
	}

	pixelData := pixelDataElement.Value.GetValue().(dicom.PixelDataInfo)
	if len(pixelData.Frames) == 0 {
		return nil, metadata, fmt.Errorf("no frames found in pixel data")
	}

	log.Printf("Found %d frames in pixel data", len(pixelData.Frames))
	var frames []image.Image
	for i, frame := range pixelData.Frames {
		img, err := frame.GetImage()
		if err != nil {
			log.Printf("Error extracting frame %d: %v", i, err)
			continue
		}
		frames = append(frames, img)
	}

	if len(frames) == 0 {
		return nil, metadata, fmt.Errorf("no valid frames extracted")
	}
	log.Printf("Successfully extracted %d frames", len(frames))
	return frames, metadata, nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func getMetadataFromMap(metadata map[string]string) string {
	result := "--- DICOM Metadata ---\n"
	if val, ok := metadata["PatientName"]; ok {
		result += fmt.Sprintf("Patient Name: %s\n", val)
	}
	if val, ok := metadata["StudyDate"]; ok {
		result += fmt.Sprintf("Study Date: %s\n", val)
	}
	if val, ok := metadata["Modality"]; ok {
		result += fmt.Sprintf("Modality: %s\n", val)
	}
	if val, ok := metadata["Dimensions"]; ok {
		result += fmt.Sprintf("Image Dimensions: %s\n", val)
	}
	if val, ok := metadata["Orientation"]; ok {
		result += fmt.Sprintf("Orientation: %s\n", val)
	}
	if val, ok := metadata["SeriesDescription"]; ok {
		result += fmt.Sprintf("Series Description: %s", val)
	}
	return result
}

func updateDisplay(imgCanvas *canvas.Image, metaLabel, frameLabel *widget.Label, seriesList []dicomSeries, seriesIdx, frameIdx int, contrast float64) {
	// Apply contrast to the image
	contrastImg := &contrastImage{
		Image:    seriesList[seriesIdx].frames[frameIdx],
		contrast: contrast,
	}
	imgCanvas.Image = contrastImg
	metaLabel.SetText(getMetadataFromMap(seriesList[seriesIdx].metadata))
	frameLabel.SetText(fmt.Sprintf("Frame %d/%d", frameIdx+1, len(seriesList[seriesIdx].frames)))
	imgCanvas.Refresh()
}
