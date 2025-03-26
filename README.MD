Below is a `README.md` file that encapsulates all the instructions, explanations, and usage details for the DICOM (Digital Imaging and Communications in Medicine) viewer script we’ve developed. This file is written in Markdown format, making it easy to read on platforms like GitHub or in a text editor. It includes everything someone would need to set up, run, and use the viewer for their own DICOM files, along with troubleshooting tips and potential next steps.

---

# DICOM Viewer in Go

This project provides a simple, local DICOM viewer written in Go. It allows you to visualize DICOM files, including multi-frame sequences, with features like zoom, frame navigation, and series selection. The viewer is designed to work with standard DICOM directory structures (e.g., those exported to a CD with a `DICOMDIR` file) and supports basic metadata display.

## Features
- **DICOM File Parsing**: Reads DICOM files and extracts metadata (e.g., patient name, study date, modality) and pixel data for visualization.
- **Multi-Frame Support**: Handles multi-frame DICOM files and treats each series as a sequence of frames, with "Play" functionality to auto-cycle through frames.
- **Series Grouping**: Groups images by series (e.g., `SE000001`, `SE000002`) using the `DICOMDIR` file or directory structure.
- **Zoom Functionality**: Zoom in, zoom out, and reset zoom levels for better image inspection.
- **Simple GUI**: Built with the `fyne` library, providing a cross-platform interface to navigate and view DICOM images.
- **Metadata Display**: Shows key metadata for each series (e.g., patient name, modality, image dimensions).

## Prerequisites
Before you begin, ensure you have the following:

1. **Go**: Version 1.16 or later. Download and install from [golang.org/dl/](https://golang.org/dl/).
2. **DICOM Files**: A set of DICOM files (`.dcm` extension) organized in a standard directory structure (e.g., `DICOM/PAxxxxxx/STxxxxxx/SExxxxxx/IMxxxxxx`).
3. **Terminal Access**: A terminal or command prompt to run commands.

## Installation

### Step 1: Install Go
- Download the Go installer for your operating system (Windows, macOS, or Linux) from [golang.org/dl/](https://golang.org/dl/).
- Follow the installation instructions for your OS.
- Verify the installation by opening a terminal and running:
  ```bash
  go version
  ```
  You should see output like `go version go1.21.0 ...`.

### Step 2: Set Up the Project
1. **Create a Directory** for your project:
   ```bash
   mkdir dicom_viewer_project
   cd dicom_viewer_project
   ```
2. **Save the Script**:
   - Copy the `dicom_viewer.go` script (provided in this repository or below) into a file named `dicom_viewer.go` in this directory.

### Step 3: Install Dependencies
1. **Initialize a Go Module**:
   ```bash
   go mod init dicom_viewer
   ```
2. **Install Required Libraries**:
   - Install the DICOM parsing library:
     ```bash
     go get github.com/suyashkumar/dicom
     ```
   - Install the Fyne GUI library:
     ```bash
     go get fyne.io/fyne/v2
     ```

### Step 4: Prepare Your DICOM Files
- Ensure your DICOM files are in a standard directory structure, such as:
  ```
  /path/to/dicom_data/
  ├── AUTORUN.INF
  ├── DICOM/
  │   └── PA000001/
  │       └── ST000001/
  │           ├── SE000001/
  │           │   ├── IM000001
  │           │   ├── IM000002
  │           │   └── ...
  │           ├── SE000002/
  │           │   ├── IM000001
  │           │   └── ...
  │           └── ...
  ├── DICOMDIR
  ├── README.TXT
  └── Viewer/
  ```
- The script will look for a `DICOMDIR` file to understand the patient/study/series hierarchy. If `DICOMDIR` is not present or fails to parse, it falls back to grouping images by their parent directory (e.g., `SE000001`).

### Step 5: Run the Viewer
1. **Run Directly**:
   Replace `/path/to/dicom_data` with the path to your top-level DICOM directory:
   ```bash
   go run dicom_viewer.go /path/to/dicom_data
   ```
2. **Or Compile and Run**:
   ```bash
   go build dicom_viewer.go
   ./dicom_viewer /path/to/dicom_data  # On Windows: dicom_viewer.exe /path/to/dicom_data
   ```

## Usage
Once the script runs, a GUI window will open with the following features:

- **Series Selection**:
  - A dropdown at the top lists all series (e.g., `SE000001`, `SE000002`).
  - Select a series to view its images.
- **Frame Navigation**:
  - **Prev Frame** / **Next Frame**: Move through frames in the selected series.
  - **Play** / **Pause**: Auto-play frames in a loop (100ms per frame; adjustable in the code).
- **Zoom**:
  - **Zoom In**: Increase image size.
  - **Zoom Out**: Decrease image size (minimum zoom level: 0.2).
  - **Reset Zoom**: Return to default size (1.0 scale).
- **Metadata Display**:
  - Metadata (e.g., patient name, study date, modality, image dimensions) is shown above the image.
- **Frame Counter**:
  - Displays the current frame and total frames (e.g., "Frame 1/26").

## Script Details
The script (`dicom_viewer.go`) does the following:

1. **DICOMDIR Parsing**:
   - Attempts to parse the `DICOMDIR` file to identify series and their associated images.
   - Falls back to directory traversal if `DICOMDIR` is unavailable or corrupt, grouping images by their parent directory (e.g., `SE000001`).

2. **Series Grouping**:
   - Groups all `.dcm` files within a series directory (e.g., `SE000001/IM000001`, `IM000002`, etc.) into a single sequence of frames.

3. **Image Display**:
   - Uses the `fyne` library to display images in a GUI window.
   - Supports multi-frame DICOM files and treats each `.dcm` file in a series as a frame.

4. **Metadata Extraction**:
   - Extracts key metadata (patient name, study date, modality, dimensions) from the first valid `.dcm` file in each series.

## Troubleshooting
- **No Series Found**:
  - Ensure your DICOM files are in the correct structure and contain valid pixel data.
  - Check the terminal for error messages (e.g., "Skipping file: ...").
- **DICOMDIR Issues**:
  - If the `DICOMDIR` file is corrupt or missing, the script falls back to directory traversal. You’ll see a log message like "Falling back to directory traversal...".
- **Performance**:
  - Loading many large series into memory might be slow. Consider modifying the script to load frames on demand if this is an issue.
- **Zoom Issues**:
  - If the image doesn’t scale well, adjust the `SetMinSize` value in the script (e.g., change `fyne.NewSize(512, 512)` to `fyne.NewSize(1024, 1024)`).

## Potential Enhancements
The script can be extended with additional features, such as:
- **Pan Functionality**: Add the ability to drag the image when zoomed in.
- **Window/Level Adjustments**: Add controls for brightness and contrast (common in medical imaging).
- **Hierarchy Navigation**: Add navigation by patient and study (not just series) for a more complete DICOM viewer experience.
- **Frame Rate Adjustment**: Allow users to adjust the playback speed for the "Play" feature.

To request these enhancements, modify the script directly or contact the developer.

## License
This project is provided as-is, with no warranty. You are free to use, modify, and distribute it under the MIT License (if applicable).

## Acknowledgments
- Built with the `github.com/suyashkumar/dicom` library for DICOM parsing.
- Uses the `fyne.io/fyne/v2` library for the cross-platform GUI.

---

### Including the Script in the README
If you’re hosting this in a repository, you can include the `dicom_viewer.go` script directly in the repository and reference it in the README. For this README, I’ve assumed the script is available separately. If you’d like to include the script inline in the README, you can add it under a section like:

```markdown
## Source Code
Below is the full source code for `dicom_viewer.go`:

```go
[Insert the script here]
```
```

However, since the script is quite long, it’s often better to keep it in a separate file in the repository.

### How to Use the README
1. Save the above content as `README.md` in your project directory (`dicom_viewer_project`).
2. If you’re using a version control system like Git, add both `dicom_viewer.go` and `README.md` to your repository:
   ```bash
   git init
   git add dicom_viewer.go README.md
   git commit -m "Initial commit with DICOM viewer and README"
   ```
3. Share the repository or the files with others, and they can follow the instructions in the README to set up and use the viewer.

Let me know if you’d like to adjust the README further or add more details!