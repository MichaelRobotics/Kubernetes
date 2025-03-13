#!/bin/bash

# Define source files and target directories
SOURCE_CONTROLLER="podtracker_controller.go"
TARGET_CONTROLLER_DIR="projects/guestbook/controllers/"
TARGET_CONTROLLER_PATH="${TARGET_CONTROLLER_DIR}${SOURCE_CONTROLLER}"

SOURCE_TYPES="podtracker_types.go"
TARGET_TYPES_DIR="projects/guestbook/api/v1/"
TARGET_TYPES_PATH="${TARGET_TYPES_DIR}${SOURCE_TYPES}"

# Function to copy a file
copy_file() {
    local source_file=$1
    local target_dir=$2
    local target_path=$3

    # Ensure the target directory exists
    mkdir -p "$target_dir"

    # Copy the file
    if [[ -f "$source_file" ]]; then
        cp "$source_file" "$target_path"
        echo "File '$source_file' successfully copied to '$target_path'."
    else
        echo "Error: Source file '$source_file' does not exist."
        exit 1
    fi
}

# Copy podtracker_controller.go
copy_file "$SOURCE_CONTROLLER" "$TARGET_CONTROLLER_DIR" "$TARGET_CONTROLLER_PATH"

# Copy podtracker_types.go
copy_file "$SOURCE_TYPES" "$TARGET_TYPES_DIR" "$TARGET_TYPES_PATH"
