package server

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// getEnhancedFileInfo gets detailed file information using stat
func getEnhancedFileInfo(ctx context.Context, sessionRepo SessionRepository, sessionID string, path string) (*FileSystemEntry, error) {
	// Use stat to get detailed file information
	statCmd := fmt.Sprintf("stat -c '%%n|%%F|%%s|%%U|%%G|%%A|%%Y' %s", sanitizePath(path))
	output, err := sessionRepo.RunCommand(ctx, sessionID, statCmd)
	if err != nil {
		return nil, err
	}

	// Parse the stat output
	parts := strings.Split(strings.TrimSpace(output), "|")
	if len(parts) < 7 {
		return nil, fmt.Errorf("invalid stat output format")
	}

	// Extract information
	name := filepath.Base(parts[0])
	fileType := parts[1]
	sizeStr := parts[2]
	owner := parts[3]
	group := parts[4]
	permissions := parts[5]
	modTimeStr := parts[6]

	// Parse size
	size, _ := strconv.ParseInt(sizeStr, 10, 64)

	// Parse modification time (Unix timestamp)
	modTimeUnix, _ := strconv.ParseInt(modTimeStr, 10, 64)
	modTime := time.Unix(modTimeUnix, 0)

	// Determine entry type
	entryType := "unknown"
	if strings.Contains(fileType, "directory") {
		entryType = "directory"
	} else if strings.Contains(fileType, "regular") {
		entryType = "file"
	} else if strings.Contains(fileType, "link") {
		entryType = "symlink"
	}

	// Check if it's a hidden file
	isHidden := strings.HasPrefix(name, ".")

	return &FileSystemEntry{
		Name:         name,
		Path:         path,
		Type:         entryType,
		Size:         size,
		Permissions:  permissions,
		Owner:        owner,
		Group:        group,
		LastModified: modTime,
		IsHidden:     isHidden,
	}, nil
}

// getRecursiveDirectoryContents gets detailed file information for all files in a directory recursively
func getRecursiveDirectoryContents(ctx context.Context, sessionRepo SessionRepository, sessionID string, path string, includeHidden bool) ([]FileSystemEntry, error) {
	var entries []FileSystemEntry

	// Build the find command
	findCmd := ""
	if includeHidden {
		findCmd = fmt.Sprintf("find %s -type f -o -type d -o -type l", sanitizePath(path))
	} else {
		findCmd = fmt.Sprintf("find %s -not -path \"*/\\.*\" -type f -o -type d -o -type l", sanitizePath(path))
	}

	// Execute the find command
	output, err := sessionRepo.RunCommand(ctx, sessionID, findCmd)
	if err != nil {
		return nil, err
	}

	// Process each line (each file path)
	paths := strings.Split(output, "\n")
	for _, filePath := range paths {
		if filePath == "" {
			continue
		}

		// Get detailed file info for this path
		fileInfo, err := getEnhancedFileInfo(ctx, sessionRepo, sessionID, filePath)
		if err != nil {
			// Skip files we can't stat
			continue
		}

		entries = append(entries, *fileInfo)
	}

	return entries, nil
}

// getNonRecursiveDirectoryContents gets detailed file information for all files in a directory (non-recursively)
func getNonRecursiveDirectoryContents(ctx context.Context, sessionRepo SessionRepository, sessionID string, path string, includeHidden bool) ([]FileSystemEntry, error) {
	var entries []FileSystemEntry

	// Build the ls command
	lsCmd := ""
	if includeHidden {
		lsCmd = fmt.Sprintf("ls -la %s", sanitizePath(path))
	} else {
		lsCmd = fmt.Sprintf("ls -l %s", sanitizePath(path))
	}

	// Execute the ls command
	output, err := sessionRepo.RunCommand(ctx, sessionID, lsCmd)
	if err != nil {
		return nil, err
	}

	// Process the ls output
	lines := strings.Split(output, "\n")

	// Find the start of the actual file listing (skip the "total" line)
	startIndex := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "total ") {
			startIndex = i + 1
			break
		}
	}

	// Process each line
	for _, line := range lines[startIndex:] {
		if line == "" {
			continue
		}

		// Parse the line to get file information
		fileEntry := parseFileEntryLine(line, path)
		if fileEntry != nil {
			entries = append(entries, *fileEntry)
		}
	}

	return entries, nil
}
