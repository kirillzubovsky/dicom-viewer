I fell while skiing and damanged my knee. Infuriated by the experience of reading my MRI results, I quickly built a DicomViewer with the help of Cursor and Grok. You can read about the experienec on [Novice.Media](https://novice.media). The app is not perfect, just like the rest of us. I hope that the next tech person that needs to see their MRI would be to use it, maybe fix bugs or contribute their features. Reach out on [X](https://x.com/kirillzubovsky) if you would like to talk about it.

---

Below is a `README.md` file that encapsulates all the instructions, explanations, and usage details for the DICOM (Digital Imaging and Communications in Medicine) viewer script . This file includes everything someone would need to set up, run, and use the viewer for their own DICOM files, along with troubleshooting tips and potential next steps.

---

# DICOM Viewer in Go

A powerful, local DICOM viewer written in Go that allows you to visualize DICOM files, including multi-frame sequences. The viewer supports zoom, pan, contrast adjustment, frame navigation, and series selection. Built for anyone with a computer who wants to load some DICOM files and view them, without the need to dowload any professional software, or to register for an account. The tool is lightweight and fast. No frills, yet.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.16-blue)
![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey)

## Table of Contents
- [Quick Start](#quick-start)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Code of Conduct](#code-of-conduct)
- [Security](#security)
- [Acknowledgments](#acknowledgments)

## Quick Start
The easiest way to get started is using our automated setup script. The script handles all dependencies, building, and launching the viewer.

### Prerequisites
You only need:
- A Unix-like system (macOS, Linux)
- Basic terminal access
- DICOM files to view
- (Theoretically it will work on Windows too, but I don't have one, haven't tried)

### Setup Steps
1. Download both files to your computer:
   - `start.sh`: The setup script
   - `dicom_viewer.go`: The viewer source code

2. Make the script executable:
   ```bash
   chmod +x start.sh
   ```

3. Run the script with your DICOM directory:
   ```bash
   ./start.sh /path/to/dicom_files
   ```

### What the Script Does
The script automatically:
1. Checks your system requirements
2. Installs Go if not present (using your system's package manager)
3. Installs all necessary dependencies
4. Sets up the project structure
5. Validates your DICOM files
6. Builds and launches the viewer

### Understanding the Output
The script uses color-coded messages to keep you informed:
- Blue [INFO]: Normal progress information
- Green [SUCCESS]: Completed steps
- Yellow [WARNING]: Non-critical issues
- Red [ERROR]: Critical issues that need attention

### Supported Systems
The automated setup works on:
- macOS (using Homebrew)
- Ubuntu/Debian (using apt)
- Fedora (using dnf)
- Arch Linux (using pacman)

For other systems, follow the manual installation below.

### Troubleshooting the Script
If you encounter issues:

1. **Path Issues**:
   ```bash
   # Try using absolute path
   ./start.sh $(pwd)/dicom_files
   ```

2. **Permission Issues**:
   ```bash
   # Make sure script is executable
   chmod +x start.sh
   # For system packages
   sudo ./start.sh /path/to/dicom_files
   ```

3. **Build Issues**:
   ```bash
   # Clean and retry
   rm -rf dicom_viewer/
   ./start.sh /path/to/dicom_files
   ```

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

## Contributing
We welcome contributions! Here's how you can help:

### Getting Started
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### Development Setup
```bash
# Clone your fork
git clone https://github.com/yourusername/dicom-viewer.git
cd dicom-viewer

# Install dependencies
go mod download

# Run tests
go test ./...
```

### Guidelines
- Follow Go coding standards and conventions
- Add tests for new features
- Update documentation as needed
- Keep pull requests focused and atomic
- Write meaningful commit messages

## Security
### Reporting Security Issues
Please report security issues to [your-email@example.com] privately. Do NOT create public issues for security vulnerabilities.

### Data Privacy
This viewer:
- Processes DICOM files locally only
- Does not transmit any data
- Does not store or cache patient information
- Allows anonymization of DICOM data (planned feature)

## Code of Conduct
We follow the [Contributor Covenant](https://www.contributor-covenant.org/version/2/0/code_of_conduct/) Code of Conduct. By participating, you are expected to uphold this code.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Third-party Licenses
- [github.com/suyashkumar/dicom](https://github.com/suyashkumar/dicom) - MIT License
- [fyne.io/fyne/v2](https://github.com/fyne-io/fyne) - BSD 3-Clause License

## Versioning
We use [SemVer](http://semver.org/) for versioning. For available versions, see the [tags on this repository](https://github.com/yourusername/dicom-viewer/tags).

## Roadmap
Planned features and improvements:
- [ ] DICOM data anonymization
- [ ] Support for more DICOM formats
- [ ] Advanced visualization tools
- [ ] Measurement tools
- [ ] DICOM network protocol support
- [ ] Plugin system for extensions

## Support
- Create an issue for bug reports or feature requests
- Join our [Discord/Slack] community (coming soon)
- Check out the [Wiki](https://github.com/yourusername/dicom-viewer/wiki) for detailed documentation

## Citation
If you use this software in your research, please cite it as:
```bibtex
@software{dicom_viewer_go,
  author = {Your Name},
  title = {DICOM Viewer in Go},
  year = {2024},
  url = {https://github.com/yourusername/dicom-viewer}
}
```