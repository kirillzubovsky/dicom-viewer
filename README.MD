Below is a `README.md` file that encapsulates all the instructions, explanations, and usage details for the DICOM (Digital Imaging and Communications in Medicine) viewer script we've developed. This file is written in Markdown format, making it easy to read on platforms like GitHub or in a text editor. It includes everything someone would need to set up, run, and use the viewer for their own DICOM files, along with troubleshooting tips and potential next steps.

---

# DICOM Viewer in Go

A simple, local DICOM viewer written in Go that allows you to visualize DICOM files, including multi-frame sequences, with features like zoom, frame navigation, and series selection. The viewer is designed to work with standard DICOM directory structures and supports basic metadata display.

## Understanding DICOM Imaging

### What is DICOM?
DICOM (Digital Imaging and Communications in Medicine) is the standard format for medical imaging data, used across various modalities like X-rays, CT scans, MRIs, and ultrasounds. Each DICOM file contains both image data and metadata about the image.

### How DICOM Images Are Acquired
1. **Imaging Modality**
   - CT Scanner: Uses X-rays to create detailed cross-sectional images
   - MRI: Uses magnetic fields and radio waves to image soft tissues
   - X-ray: Produces 2D images using X-rays
   - Ultrasound: Uses sound waves to image soft tissues

2. **Image Generation**
   - Raw data is captured by the imaging device
   - Data is processed to reconstruct images
   - Multiple "series" of images are often produced in a single session

3. **DICOM File Structure**
   - Hierarchical organization:
     - Patient (PAxxxxxx)
     - Study (STxxxxxx)
     - Series (SExxxxxx)
     - Image (IMxxxxxx)
   - Each file contains:
     - Patient information
     - Study details
     - Series information
     - Image data
     - Technical parameters

### Understanding Image Orientation
DICOM images can be acquired in different anatomical planes:
- **Axial**: Top-down slices (like cutting a loaf of bread horizontally)
- **Coronal**: Front-to-back slices (like slicing the bread vertically from the front)
- **Sagittal**: Side-to-side slices (like slicing the bread vertically from the side)

## Features
- **DICOM File Parsing**: Reads DICOM files and extracts metadata and pixel data
- **Multi-Frame Support**: Handles multi-frame DICOM files with "Play" functionality
- **Series Grouping**: Groups images by series using DICOMDIR or directory structure
- **Zoom Functionality**: Zoom in, zoom out, and reset zoom levels
- **Contrast Adjustment**: Adjust image contrast for better visualization
- **Simple GUI**: Built with the `fyne` library for cross-platform compatibility
- **Metadata Display**: Shows key metadata for each series
- **Orientation Information**: Displays the anatomical plane of the images

## Prerequisites
1. **Go**: Version 1.16 or later
2. **DICOM Files**: A set of DICOM files organized in a standard directory structure
3. **Terminal Access**: A terminal or command prompt to run commands

## Installation

### Step 1: Install Go
- Download Go from [golang.org/dl/](https://golang.org/dl/)
- Verify installation:
  ```bash
  go version
  ```

### Step 2: Set Up the Project
1. Create a directory:
   ```bash
   mkdir dicom_viewer_project
   cd dicom_viewer_project
   ```
2. Save the `dicom_viewer.go` script in this directory

### Step 3: Install Dependencies
1. Initialize Go module:
   ```bash
   go mod init dicom_viewer
   ```
2. Install required libraries:
   ```bash
   go get github.com/suyashkumar/dicom
   go get fyne.io/fyne/v2
   ```

### Step 4: Prepare Your DICOM Files
Ensure your DICOM files are in a standard directory structure:
```
/path/to/dicom_data/
├── DICOM/
│   └── PAxxxxxx/
│       └── STxxxxxx/
│           ├── SExxxxxx/
│           │   ├── IMxxxxxx
│           │   ├── IMxxxxxx
│           │   └── ...
│           └── ...
├── DICOMDIR
└── ...
```

### Step 5: Run the Viewer
```bash
go run dicom_viewer.go /path/to/dicom_data
```

## Usage

### Basic Navigation
- **Series Selection**: Use the dropdown to switch between different series
- **Frame Navigation**: 
  - Use "Prev Frame" / "Next Frame" to move through slices
  - Use "Play" / "Pause" to auto-cycle through frames
- **Zoom Controls**:
  - "Zoom In" / "Zoom Out" to adjust image size
  - "Reset Zoom" to return to default size

### Image Adjustment
- **Contrast Slider**: Adjust image contrast for better visualization
- **Metadata Display**: Shows information about the current series

### Understanding the Interface
1. **Series Selection**
   - Lists all available series
   - Each series represents a different view or imaging technique
   - Series are sorted numerically for easy navigation

2. **Frame Navigation**
   - Each frame represents a slice through the body
   - Use frame navigation to "travel" through the body
   - Play mode helps visualize the 3D structure

3. **Metadata Display**
   - Shows patient information
   - Displays study details
   - Shows image dimensions
   - Indicates the anatomical orientation

## Technical Details

### DICOM File Processing
1. **DICOMDIR Parsing**
   - Attempts to parse DICOMDIR for series organization
   - Falls back to directory traversal if needed

2. **Series Grouping**
   - Groups images by series directory
   - Maintains proper frame order

3. **Image Display**
   - Supports multi-frame DICOM files
   - Handles different image orientations
   - Provides contrast adjustment

### Metadata Extraction
The viewer extracts and displays:
- Patient information
- Study details
- Series information
- Image dimensions
- Anatomical orientation

## Troubleshooting
- **No Series Found**: Check directory structure and file permissions
- **DICOMDIR Issues**: Viewer falls back to directory traversal
- **Performance**: Large series may load slowly
- **Zoom Issues**: Adjust minimum size if needed

## License
This project is provided as-is, with no warranty. You are free to use, modify, and distribute it under the MIT License.

## Acknowledgments
- Built with `github.com/suyashkumar/dicom` for DICOM parsing
- Uses `fyne.io/fyne/v2` for the cross-platform GUI

---

### Including the Script in the README
If you're hosting this in a repository, you can include the `dicom_viewer.go` script directly in the repository and reference it in the README. For this README, I've assumed the script is available separately. If you'd like to include the script inline in the README, you can add it under a section like:

```markdown
## Source Code
Below is the full source code for `dicom_viewer.go`:

```go
[Insert the script here]
```
```

However, since the script is quite long, it's often better to keep it in a separate file in the repository.

### How to Use the README
1. Save the above content as `README.md` in your project directory (`dicom_viewer_project`).
2. If you're using a version control system like Git, add both `dicom_viewer.go` and `README.md` to your repository:
   ```bash
   git init
   git add dicom_viewer.go README.md
   git commit -m "Initial commit with DICOM viewer and README"
   ```
3. Share the repository or the files with others, and they can follow the instructions in the README to set up and use the viewer.

Let me know if you'd like to adjust the README further or add more details!