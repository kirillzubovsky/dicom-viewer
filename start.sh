#!/bin/bash

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Function to print colored status messages
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to get absolute path
get_abs_path() {
    local path="$1"
    if [ -d "$path" ]; then
        (cd "$path" && pwd)
    elif [ -f "$path" ]; then
        if [[ $path == /* ]]; then
            echo "$path"
        else
            echo "$PWD/${path#./}"
        fi
    fi
}

# Function to check and install Go
check_go() {
    if ! command_exists go; then
        print_warning "Go is not installed. Installing Go..."
        
        case "$(uname -s)" in
            Darwin*)
                if command_exists brew; then
                    brew install go
                else
                    print_error "Please install Homebrew first: https://brew.sh/"
                    exit 1
                fi
                ;;
            Linux*)
                if command_exists apt-get; then
                    sudo apt-get update
                    sudo apt-get install -y golang-go
                elif command_exists dnf; then
                    sudo dnf install -y golang
                elif command_exists pacman; then
                    sudo pacman -S --noconfirm go
                else
                    print_error "Unsupported Linux distribution. Please install Go manually: https://golang.org/dl/"
                    exit 1
                fi
                ;;
            *)
                print_error "Unsupported operating system. Please install Go manually: https://golang.org/dl/"
                exit 1
                ;;
        esac
    fi
    print_success "Go is installed: $(go version)"
}

# Function to check and install system dependencies
check_system_deps() {
    case "$(uname -s)" in
        Darwin*)
            if ! command_exists xcode-select; then
                print_warning "Installing Xcode Command Line Tools..."
                xcode-select --install
            fi
            ;;
        Linux*)
            if command_exists apt-get; then
                print_status "Checking required packages..."
                REQUIRED_PACKAGES="gcc libgl1-mesa-dev xorg-dev"
                MISSING_PACKAGES=""
                
                for package in $REQUIRED_PACKAGES; do
                    if ! dpkg -l | grep -q "^ii  $package "; then
                        MISSING_PACKAGES="$MISSING_PACKAGES $package"
                    fi
                done
                
                if [ ! -z "$MISSING_PACKAGES" ]; then
                    print_warning "Installing missing packages:$MISSING_PACKAGES"
                    sudo apt-get update
                    sudo apt-get install -y $MISSING_PACKAGES
                fi
            elif command_exists dnf; then
                print_warning "Installing required packages for Fedora..."
                sudo dnf install -y gcc libX11-devel libXrandr-devel libXinerama-devel libXcursor-devel libXi-devel mesa-libGL-devel
            elif command_exists pacman; then
                print_warning "Installing required packages for Arch Linux..."
                sudo pacman -S --noconfirm gcc libgl xorg-server
            fi
            ;;
    esac
    print_success "System dependencies are installed"
}

# Function to set up the project
setup_project() {
    # Store the original path
    ORIGINAL_PATH="$PWD"
    
    # Create project directory if it doesn't exist
    if [ ! -d "dicom_viewer" ]; then
        print_status "Creating project directory..."
        mkdir -p dicom_viewer
    fi
    
    # Copy dicom_viewer.go to the project directory if it exists in current directory
    if [ -f "dicom_viewer.go" ]; then
        cp dicom_viewer.go dicom_viewer/
    elif [ ! -f "dicom_viewer/dicom_viewer.go" ]; then
        print_error "dicom_viewer.go not found in current directory or project directory"
        exit 1
    fi
    
    cd dicom_viewer
    
    # Initialize Go module if it doesn't exist
    if [ ! -f "go.mod" ]; then
        print_status "Initializing Go module..."
        go mod init dicom_viewer
    fi
    
    # Install dependencies
    print_status "Installing dependencies..."
    go get github.com/suyashkumar/dicom
    go get fyne.io/fyne/v2@v2.5.5
    
    # Install all required Fyne sub-modules
    print_status "Installing Fyne sub-modules..."
    modules=(
        "fyne.io/fyne/v2/internal/svg@v2.5.5"
        "fyne.io/fyne/v2/storage/repository@v2.5.5"
        "fyne.io/fyne/v2/lang@v2.5.5"
        "fyne.io/fyne/v2/internal/painter@v2.5.5"
        "fyne.io/fyne/v2/widget@v2.5.5"
        "fyne.io/fyne/v2/internal/painter/gl@v2.5.5"
        "fyne.io/fyne/v2/internal/driver/glfw@v2.5.5"
        "fyne.io/fyne/v2/internal/metadata@v2.5.5"
        "fyne.io/fyne/v2/app@v2.5.5"
    )
    
    for module in "${modules[@]}"; do
        print_status "Installing $module..."
        go get "$module"
    done
    
    # Ensure all dependencies are in go.sum
    print_status "Verifying dependencies..."
    go mod tidy
    
    cd "$ORIGINAL_PATH"
    print_success "Project setup complete"
}

# Function to check if DICOM files exist in the specified directory
check_dicom_files() {
    local dicom_path=$(get_abs_path "$1")
    
    if [ -z "$dicom_path" ]; then
        print_error "No DICOM directory specified"
        print_status "Usage: ./start.sh /path/to/dicom_files"
        exit 1
    fi
    
    if [ ! -d "$dicom_path" ]; then
        print_error "Directory does not exist: $dicom_path"
        exit 1
    fi
    
    # Check for DICOM files
    if [ -z "$(find "$dicom_path" -name "*.dcm" -o -name "*.DCM" 2>/dev/null)" ]; then
        print_warning "No .dcm files found in $dicom_path"
        print_status "Please ensure your DICOM files are in the correct format and location"
        return 1
    fi
    
    print_success "DICOM files found in $dicom_path"
    echo "$dicom_path"
}

# Function to build and run the viewer
run_viewer() {
    local dicom_path="$1"
    cd dicom_viewer || exit 1
    
    print_status "Building DICOM viewer..."
    if ! go build -o viewer dicom_viewer.go; then
        print_error "Build failed. Please check the error messages above"
        exit 1
    fi
    
    print_success "Build successful"
    print_status "Launching DICOM viewer with path: $dicom_path"
    
    # Run the viewer with the specified DICOM directory
    ./viewer "$dicom_path"
}

# Main execution
echo -e "${BOLD}DICOM Viewer Setup and Launch Script${NC}\n"

# Check for DICOM directory argument
if [ $# -eq 0 ]; then
    print_error "No DICOM directory specified"
    print_status "Usage: ./start.sh /path/to/dicom_files"
    exit 1
fi

# Get absolute path of the DICOM directory
DICOM_PATH=$(get_abs_path "$1")
if [ -z "$DICOM_PATH" ]; then
    print_error "Invalid path: $1"
    exit 1
fi

# Run all checks and setup
print_status "Checking system requirements..."
check_go
check_system_deps
setup_project

# Check DICOM files
if ! check_dicom_files "$DICOM_PATH"; then
    print_status "Proceeding anyway as the viewer has its own file validation..."
fi

run_viewer "$DICOM_PATH" 