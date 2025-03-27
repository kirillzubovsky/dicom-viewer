Below is a `README.md` file that encapsulates all the instructions, explanations, and usage details for the DICOM (Digital Imaging and Communications in Medicine) viewer script we've developed. This file is written in Markdown format, making it easy to read on platforms like GitHub or in a text editor. It includes everything someone would need to set up, run, and use the viewer for their own DICOM files, along with troubleshooting tips and potential next steps.

---

# DICOM Viewer in Go

A powerful, local DICOM viewer written in Go that allows you to visualize DICOM files, including multi-frame sequences. The viewer supports advanced features like zoom, pan, contrast adjustment, frame navigation, and series selection. Built for medical professionals and researchers who need a lightweight, fast DICOM viewer.

## Quick Start
The easiest way to get started is using our automated setup script:

1. Download both `start.sh` and `dicom_viewer.go` to your computer
2. Make the script executable:
   ```bash
   chmod +x start.sh
   ```
3. Run the script with your DICOM directory:
   ```bash
   ./start.sh /path/to/dicom_files
   ```

The script will:
- Check and install all required dependencies
- Set up the project automatically
- Validate your DICOM files
- Build and launch the viewer

The script provides color-coded output to help you understand what's happening:
- 🔵 Blue: Information about current operations
- 🟢 Green: Successful steps
- 🟡 Yellow: Warnings that need attention
- 🔴 Red: Errors that need to be fixed

### Directory Structure
The viewer expects your DICOM files to be organized in this structure:
```
dicom_files/
└── DICOM/
    └── PAxxxxxx/
        └── STxxxxxx/
            └── SExxxxxx/
                ├── image1.dcm
                ├── image2.dcm
                └── ...
```

You can use either:
- Relative path: `./start.sh ./dicom_files`
- Absolute path: `./start.sh /full/path/to/dicom_files`

### Troubleshooting the Script
If you encounter issues:

1. **Directory not found**:
   - Make sure you're in the correct directory
   - Check if the path to DICOM files is correct
   - Try using the absolute path

2. **Permission errors**:
   - Make sure the script is executable: `chmod +x start.sh`
   - For system package installation: run with `sudo` if prompted

3. **Build errors**:
   - Check Go installation: `go version`
   - Ensure all dependencies are installed
   - Look for red error messages in the output

### Supported Operating Systems
The automated setup script supports:
- macOS (via Homebrew)
- Ubuntu/Debian
- Fedora
- Arch Linux

For other systems, follow the manual installation instructions below.

## Features
- **DICOM File Parsing**: Reads DICOM files and extracts metadata and pixel data
- **Multi-Frame Support**: Handles multi-frame DICOM files with frame-by-frame navigation
- **Series Grouping**: Groups images by series using DICOMDIR or directory structure
- **Advanced Image Controls**:
  - Zoom: Zoom in/out with maintained aspect ratio
  - Pan: Click and drag to pan around zoomed images
  - Contrast Adjustment: Fine-tune image contrast with a wide-range slider
  - Frame Navigation: Step through frames or use auto-play functionality
- **Metadata Display**: Shows comprehensive metadata including:
  - Patient information
  - Study details
  - Image dimensions
  - Anatomical orientation
  - Series description
- **Export Functionality**: Export series to PNG format with:
  - Current contrast settings preserved
  - Complete metadata in text format
  - Sequential frame numbering
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Prerequisites
1. **Go**: Version 1.16 or later
2. **Required Libraries**: 
   - For Windows: GCC compiler (via MinGW or MSYS2)
   - For Ubuntu/Debian: `sudo apt-get install gcc libgl1-mesa-dev xorg-dev`
   - For macOS: Xcode Command Line Tools (`xcode-select --install`)
3. **DICOM Files**: A set of DICOM files organized in a standard directory structure

## Manual Installation
If you prefer to install manually or the automated script doesn't work for your system, follow these steps:

### Step 1: Install Go
1. Download Go from [golang.org/dl/](https://golang.org/dl/)
2. Follow the installation instructions for your operating system
3. Verify installation:
   ```bash
   go version
   ```

### Step 2: Install System Dependencies

#### Windows:
1. Install MSYS2 from [msys2.org](https://www.msys2.org/)
2. Open MSYS2 terminal and run:
   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```
3. Add MinGW-w64 bin directory to your PATH

#### macOS:
```bash
xcode-select --install
```

#### Ubuntu/Debian:
```bash
sudo apt-get update
sudo apt-get install gcc libgl1-mesa-dev xorg-dev
```

### Step 3: Set Up the Project
1. Create a project directory:
   ```bash
   mkdir dicom_viewer
   cd dicom_viewer
   ```

2. Initialize Go module:
   ```bash
   go mod init dicom_viewer
   ```

3. Create `dicom_viewer.go` and copy the source code into it

4. Install dependencies:
   ```bash
   go get github.com/suyashkumar/dicom
   go get fyne.io/fyne/v2
   ```

## Running the Viewer

1. Basic usage:
   ```bash
   go run dicom_viewer.go /path/to/dicom_files
   ```

2. Building an executable:
   ```bash
   go build dicom_viewer.go
   ```
   Then run:
   - Windows: `dicom_viewer.exe path\to\dicom_files`
   - macOS/Linux: `./dicom_viewer /path/to/dicom_files`

## DICOM Directory Structure
The viewer expects one of these directory structures:

1. With DICOMDIR:
   ```
   dicom_files/
   └── DICOM/
       └── PAxxxxxx/
           └── STxxxxxx/
               └── SExxxxxx/
                   ├── IMxxxxxx
                   └── ...
   ```

## Using the Viewer

### Navigation
- **Series Selection**: Use dropdown at top to switch between series
- **Frame Navigation**: 
  - "Prev Frame"/"Next Frame" buttons
  - "Play" button for automatic frame cycling
  - Frame counter shows current position

### Image Controls
- **Zoom**:
  - "Zoom In": Increase image size
  - "Zoom Out": Decrease image size
  - "Reset Zoom": Return to original size
  - Click and drag to pan when zoomed in
- **Contrast**:
  - Slider adjusts image contrast
  - Range: -255 to +255
  - Current contrast value displayed

### Export
1. Click "Export Series" button
2. Select destination folder
3. Creates:
   - Directory named `series_SExxxxxx`
   - PNG files for each frame with current contrast
   - `metadata.txt` with series information

## Troubleshooting

### Common Issues
1. **Window System Error** (Linux):
   ```bash
   sudo apt-get install libgl1-mesa-dev xorg-dev
   ```

2. **Compilation Error** (Windows):
   - Ensure MinGW-w64 is properly installed
   - Check PATH includes MinGW-w64 bin directory

3. **No Series Found**:
   - Check directory structure
   - Verify file permissions
   - Ensure files are valid DICOM format

4. **Performance Issues**:
   - Large series may load slowly
   - Consider reducing frame rate in auto-play
   - Close other memory-intensive applications

### Error Messages
- "No valid DICOM series found": Check directory structure and file formats
- "Failed to parse DICOM file": File may be corrupted or in unsupported format
- "No frames found in pixel data": DICOM file may not contain image data

## License
This project is provided under the MIT License. You are free to use, modify, and distribute it.

## Acknowledgments
- Built with [github.com/suyashkumar/dicom](https://github.com/suyashkumar/dicom) for DICOM parsing
- Uses [fyne.io/fyne/v2](https://fyne.io/) for cross-platform GUI

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