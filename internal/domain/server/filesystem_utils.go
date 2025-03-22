package server

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// sanitizePath ensures the path is safe to use in a shell command
func sanitizePath(path string) string {
	// Ensure we have a path
	if path == "" {
		return "/"
	}

	// Clean the path to remove any weird constructs like "../.."
	cleanPath := filepath.Clean(path)

	// Return the path with single quotes to prevent command injection
	return "'" + strings.ReplaceAll(cleanPath, "'", "'\\''") + "'"
}

// parseNonRecursiveFileListing parses the output of 'ls -l' command
func parseNonRecursiveFileListing(output string, basePath string) []FileSystemEntry {
	var entries []FileSystemEntry

	// Normalize the base path
	basePath = strings.Trim(basePath, "'")
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	// Split output by lines and skip the total line
	lines := strings.Split(output, "\n")
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

		// Parse the line
		entry := parseFileEntryLine(line, basePath)
		if entry != nil {
			entries = append(entries, *entry)
		}
	}

	return entries
}

// parseRecursiveFileListing parses the output of a recursive find command
func parseRecursiveFileListing(output string, basePath string) []FileSystemEntry {
	var entries []FileSystemEntry

	// Normalize the base path
	basePath = strings.Trim(basePath, "'")
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	// Split output by lines
	lines := strings.Split(output, "\n")

	// Process each line
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Get file info using stat
		fileInfo := getFileInfoFromPath(line)
		if fileInfo != nil {
			entries = append(entries, *fileInfo)
		}
	}

	return entries
}

// parseFileEntryLine parses a single line from 'ls -l' output
func parseFileEntryLine(line string, basePath string) *FileSystemEntry {
	// Regex to parse ls -l output
	// Format: perms links owner group size month day time name
	re := regexp.MustCompile(`^([d-lcrwxsS]+)\s+(\d+)\s+(\S+)\s+(\S+)\s+(\d+)\s+(\w{3})\s+(\d+)\s+(\d+:\d+|\d{4})\s+(.+)$`)

	matches := re.FindStringSubmatch(line)
	if matches == nil || len(matches) < 10 {
		return nil
	}

	permissions := matches[1]
	owner := matches[3]
	group := matches[4]
	sizeStr := matches[5]
	month := matches[6]
	day := matches[7]
	timeOrYear := matches[8]
	name := matches[9]

	// Determine if it's a directory, file, or symlink
	var entryType string
	if permissions[0] == 'd' {
		entryType = "directory"
	} else if permissions[0] == 'l' {
		entryType = "symlink"
	} else {
		entryType = "file"
	}

	// Parse the size
	size, _ := strconv.ParseInt(sizeStr, 10, 64)

	// Parse the modification time
	now := time.Now()
	var modTime time.Time

	if strings.Contains(timeOrYear, ":") {
		// Time format (current year): Jan 2 14:30
		modTime, _ = time.Parse("Jan 2 15:04", month+" "+day+" "+timeOrYear)
		modTime = time.Date(now.Year(), modTime.Month(), modTime.Day(), modTime.Hour(), modTime.Minute(), 0, 0, time.Local)

		// If the resulting time is in the future, it's probably from last year
		if modTime.After(now) {
			modTime = time.Date(now.Year()-1, modTime.Month(), modTime.Day(), modTime.Hour(), modTime.Minute(), 0, 0, time.Local)
		}
	} else {
		// Year format: Jan 2 2020
		modTime, _ = time.Parse("Jan 2 2006", month+" "+day+" "+timeOrYear)
	}

	// Handle symlink target
	if entryType == "symlink" && strings.Contains(name, " -> ") {
		parts := strings.SplitN(name, " -> ", 2)
		name = parts[0]
	}

	// Check if it's a hidden file or directory
	isHidden := strings.HasPrefix(name, ".")

	return &FileSystemEntry{
		Name:         name,
		Path:         basePath + name,
		Type:         entryType,
		Size:         size,
		Permissions:  permissions,
		Owner:        owner,
		Group:        group,
		LastModified: modTime,
		IsHidden:     isHidden,
	}
}

// getFileInfoFromPath gets file information for a given path
func getFileInfoFromPath(path string) *FileSystemEntry {
	// Extract the file/directory name from the path
	name := filepath.Base(path)

	// Check if it's a hidden file or directory
	isHidden := strings.HasPrefix(name, ".")

	// Create a basic entry (more details would require additional shell commands)
	return &FileSystemEntry{
		Name:     name,
		Path:     path,
		Type:     "unknown", // Would need to run 'stat' to get accurate type
		IsHidden: isHidden,
	}
}
