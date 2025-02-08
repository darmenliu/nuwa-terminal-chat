#!/bin/bash

# Define variables
REPO="darmenliu/nuwa-terminal-chat"
INSTALL_DIR="$HOME/.nuwa-terminal"
BIN_DIR="/usr/local/bin"
TEMP_DIR=$(mktemp -d)

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run this script with sudo"
  exit 1
fi

# Get the latest version
echo "Fetching the latest version..."
LATEST_RELEASE=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# Get system information
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Adjust download filename based on architecture
case $ARCH in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64|arm64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Build download URL
ASSET_NAME="nuwa-terminal-$OS-$ARCH"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$ASSET_NAME"

# Download the binary
echo "Downloading $ASSET_NAME..."
curl -L -o "$TEMP_DIR/nuwa-terminal" "$DOWNLOAD_URL"
if [ $? -ne 0 ]; then
  echo "Download failed"
  exit 1
fi

# Create installation directory
echo "Creating installation directory..."
mkdir -p "$INSTALL_DIR"
if [ $? -ne 0 ]; then
  echo "Failed to create installation directory"
  exit 1
fi

# Install the binary
echo "Installing the binary..."
sudo install -m 755 "$TEMP_DIR/nuwa-terminal" "$BIN_DIR/nuwa-terminal"
if [ $? -ne 0 ]; then
  echo "Installation failed"
  exit 1
fi

# Create envs.sh file
echo "Creating envs.sh configuration file..."
cat > "$INSTALL_DIR/envs.sh" << 'EOF'
#!/bin/bash

# Configure LLM backend
export LLM_BACKEND=deepseek
export LLM_MODEL_NAME=deepseek-coder
export LLM_API_KEY="your_api_key_here"
export LLM_TEMPERATURE=0.8
export LLM_BASE_URL=https://api.deepseek.com
export OLLAMA_SERVER_URL="no_needed"

# Other configurations
export NUWA_TERMINAL_DIR="$HOME/.nuwa-terminal"
EOF

# Set permissions
chmod 644 "$INSTALL_DIR/envs.sh"

# Clean up temporary files
rm -rf "$TEMP_DIR"

echo "Installation complete!"
echo "Please edit $INSTALL_DIR/envs.sh to configure your API key and other settings."
echo "Then run 'source $INSTALL_DIR/envs.sh' to apply the configuration."
echo "You can now use the 'nuwa-terminal' command to start the program." 
