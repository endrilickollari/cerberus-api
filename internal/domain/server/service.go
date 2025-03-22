package server

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Common errors
var (
	ErrSessionNotFound = errors.New("session not found or expired")
	ErrCommandFailed   = errors.New("command execution failed")
)

// SessionRepository defines methods to access SSH sessions
type SessionRepository interface {
	// RunCommand executes a command on the SSH session
	RunCommand(ctx context.Context, sessionID string, command string) (string, error)
}

// Service defines the server details service
type Service interface {
	// GetBasicDetails retrieves basic server information
	GetBasicDetails(ctx context.Context, sessionID string) (*ServerDetails, error)

	// GetCPUInfo retrieves CPU information
	GetCPUInfo(ctx context.Context, sessionID string) ([]CPUInfo, error)

	// GetDiskUsage retrieves disk usage information
	GetDiskUsage(ctx context.Context, sessionID string) ([]DiskUsage, error)

	// GetRunningProcesses retrieves information about running processes
	GetRunningProcesses(ctx context.Context, sessionID string) ([]ProcessInfo, error)

	// GetInstalledLibraries retrieves information about installed libraries
	GetInstalledLibraries(ctx context.Context, sessionID string) ([]Library, error)

	// ListFileSystem retrieves a listing of files and directories
	ListFileSystem(ctx context.Context, sessionID string, path string, recursive bool, includeHidden bool) (*FileSystemListing, error)

	// GetFileDetails retrieves detailed information about a specific file or directory
	GetFileDetails(ctx context.Context, sessionID string, path string) (*FileSystemEntry, error)

	// SearchFiles searches for files matching a pattern
	SearchFiles(ctx context.Context, sessionID string, path string, pattern string, maxDepth int) ([]FileSystemEntry, error)
}

type service struct {
	sessionRepo SessionRepository
}

// NewService creates a new server details service
func NewService(sessionRepo SessionRepository) Service {
	return &service{
		sessionRepo: sessionRepo,
	}
}

// GetBasicDetails implements the Service interface
func (s *service) GetBasicDetails(ctx context.Context, sessionID string) (*ServerDetails, error) {
	// Execute commands to get basic server information
	hostname, err := s.sessionRepo.RunCommand(ctx, sessionID, "hostname")
	if err != nil {
		return nil, err
	}

	osInfo, err := s.sessionRepo.RunCommand(ctx, sessionID, "uname -a")
	if err != nil {
		return nil, err
	}

	kernelVersion, err := s.sessionRepo.RunCommand(ctx, sessionID, "uname -r")
	if err != nil {
		return nil, err
	}

	uptime, err := s.sessionRepo.RunCommand(ctx, sessionID, "uptime")
	if err != nil {
		return nil, err
	}

	// Create and return server details
	return &ServerDetails{
		Hostname:      strings.TrimSpace(hostname),
		OS:            strings.TrimSpace(osInfo),
		KernelVersion: strings.TrimSpace(kernelVersion),
		Uptime:        strings.TrimSpace(uptime),
	}, nil
}

// GetCPUInfo implements the Service interface
func (s *service) GetCPUInfo(ctx context.Context, sessionID string) ([]CPUInfo, error) {
	// Execute command to get CPU info
	cpuInfoOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, "cat /proc/cpuinfo")
	if err != nil {
		return nil, err
	}

	// Parse CPU info
	return parseCPUInfo(cpuInfoOutput), nil
}

// GetDiskUsage implements the Service interface
func (s *service) GetDiskUsage(ctx context.Context, sessionID string) ([]DiskUsage, error) {
	// Execute command to get disk usage
	diskUsageOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, "df -h")
	if err != nil {
		return nil, err
	}

	// Parse disk usage
	return parseDiskUsage(diskUsageOutput), nil
}

// GetRunningProcesses implements the Service interface
func (s *service) GetRunningProcesses(ctx context.Context, sessionID string) ([]ProcessInfo, error) {
	// Execute command to get running processes
	processesOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, "ps aux")
	if err != nil {
		return nil, err
	}

	// Parse processes info
	return parseProcessInfo(processesOutput), nil
}

// parseCPUInfo parses the output of 'cat /proc/cpuinfo'
func parseCPUInfo(cpuInfoData string) []CPUInfo {
	var cpuInfos []CPUInfo
	var cpuInfo CPUInfo

	lines := strings.Split(cpuInfoData, "\n")
	for _, line := range lines {
		if line == "" {
			// End of one processor info, add to slice and start a new one
			if cpuInfo.Processor != "" { // Only add if we've parsed some data
				cpuInfos = append(cpuInfos, cpuInfo)
				cpuInfo = CPUInfo{}
			}
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "processor":
			cpuInfo.Processor = value
		case "vendor_id":
			cpuInfo.VendorID = value
		case "cpu family":
			cpuInfo.CPUFamily = value
		case "model":
			cpuInfo.Model = value
		case "model name":
			cpuInfo.ModelName = value
		case "stepping":
			cpuInfo.Stepping = value
		case "microcode":
			cpuInfo.Microcode = value
		case "cpu MHz":
			cpuInfo.CPUMHz = value
		case "cache size":
			cpuInfo.CacheSize = value
		case "physical id":
			cpuInfo.PhysicalID = value
		case "siblings":
			cpuInfo.Siblings = value
		case "core id":
			cpuInfo.CoreID = value
		case "cpu cores":
			cpuInfo.CPUCores = value
		case "apicid":
			cpuInfo.APICID = value
		case "initial apicid":
			cpuInfo.InitialAPICID = value
		case "fpu":
			cpuInfo.FPU = value
		case "fpu_exception":
			cpuInfo.FPUException = value
		case "cpuid level":
			cpuInfo.CPUIDLevel = value
		case "wp":
			cpuInfo.WP = value
		case "flags":
			cpuInfo.Flags = value
		case "bogomips":
			cpuInfo.Bogomips = value
		case "clflush size":
			cpuInfo.ClflushSize = value
		case "cache_alignment":
			cpuInfo.CacheAlignment = value
		case "address sizes":
			cpuInfo.AddressSizes = value
		case "power management":
			cpuInfo.PowerManagement = value
		}
	}

	// Add the last CPU info if it exists
	if cpuInfo.Processor != "" {
		cpuInfos = append(cpuInfos, cpuInfo)
	}

	return cpuInfos
}

// parseDiskUsage parses the output of 'df -h'
func parseDiskUsage(dfOutput string) []DiskUsage {
	var diskUsages []DiskUsage
	lines := strings.Split(dfOutput, "\n")

	// Skip the first line (headers)
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		// Split the line by spaces
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		diskUsage := DiskUsage{
			Filesystem:    fields[0],
			Size:          fields[1],
			Used:          fields[2],
			Available:     fields[3],
			UsePercentage: fields[4],
			MountedOn:     fields[5],
		}

		diskUsages = append(diskUsages, diskUsage)
	}

	return diskUsages
}

// parseProcessInfo parses the output of 'ps aux'
func parseProcessInfo(psOutput string) []ProcessInfo {
	var processes []ProcessInfo
	lines := strings.Split(psOutput, "\n")

	// Skip the first line (headers)
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		// Split the line by spaces
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}

		process := ProcessInfo{
			User:  fields[0],
			PID:   fields[1],
			CPU:   fields[2],
			VSZ:   fields[3],
			RSS:   fields[4],
			TTY:   fields[5],
			Stat:  fields[6],
			Start: fields[7],
			Time:  fields[8],
			// Combine remaining fields for command
			CMD: strings.Join(fields[10:], " "),
		}

		processes = append(processes, process)
	}

	return processes
}

// parseLibrariesInfo parses the output of package manager queries
func parseLibrariesInfo(librariesOutput string) []Library {
	var libraries []Library
	lines := strings.Split(librariesOutput, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Split the line by spaces
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// Parse fields based on standard format
		name := fields[0]
		version := fields[1]

		// Set default values
		status := "unknown"
		arch := "unknown"

		// Check for status (may be formatted differently based on package manager)
		if len(fields) >= 4 && strings.Contains(fields[2], "install") {
			status = "installed"
			if len(fields) >= 4 {
				arch = fields[3]
			}
		} else if len(fields) >= 3 {
			if strings.Contains(fields[2], "install") {
				status = "installed"
				if len(fields) >= 4 {
					arch = fields[3]
				}
			} else {
				// Last field could be architecture
				arch = fields[2]
			}
		}

		library := Library{
			Name:    name,
			Version: version,
			Status:  status,
			Arch:    arch,
		}

		libraries = append(libraries, library)
	}

	return libraries
}

// GetInstalledLibraries implements the Service interface
func (s *service) GetInstalledLibraries(ctx context.Context, sessionID string) ([]Library, error) {
	// Try to detect the Linux distribution
	distroCmd := "cat /etc/os-release | grep -E '^ID=' | cut -d'=' -f2 | tr -d '\"'"
	distro, err := s.sessionRepo.RunCommand(ctx, sessionID, distroCmd)
	if err != nil {
		// Fallback if we can't detect the distribution
		distro = ""
	}
	distro = strings.TrimSpace(distro)

	var command string
	switch strings.ToLower(distro) {
	case "ubuntu", "debian", "pop", "mint", "elementary", "kali", "zorin":
		// Debian-based distributions
		command = "dpkg-query -W -f='${Package} ${Version} ${Status} ${Architecture}\n'"
	case "fedora", "rhel", "centos", "rocky", "alma":
		// Red Hat-based distributions
		command = "rpm -qa --queryformat '%{NAME} %{VERSION} installed %{ARCH}\n'"
	case "arch", "manjaro", "endeavouros":
		// Arch-based distributions
		command = "pacman -Q | awk '{print $1 \" \" $2 \" installed \"}'$(uname -m)"
	case "opensuse", "suse":
		// SUSE-based distributions
		command = "rpm -qa --queryformat '%{NAME} %{VERSION} installed %{ARCH}\n'"
	default:
		// Try to detect package manager if distro detection failed
		aptCheck, _ := s.sessionRepo.RunCommand(ctx, sessionID, "which apt &>/dev/null && echo found || echo not found")
		if strings.Contains(aptCheck, "found") {
			command = "dpkg-query -W -f='${Package} ${Version} ${Status} ${Architecture}\n'"
		} else {
			rpmCheck, _ := s.sessionRepo.RunCommand(ctx, sessionID, "which rpm &>/dev/null && echo found || echo not found")
			if strings.Contains(rpmCheck, "found") {
				command = "rpm -qa --queryformat '%{NAME} %{VERSION} installed %{ARCH}\n'"
			} else {
				pacmanCheck, _ := s.sessionRepo.RunCommand(ctx, sessionID, "which pacman &>/dev/null && echo found || echo not found")
				if strings.Contains(pacmanCheck, "found") {
					command = "pacman -Q | awk '{print $1 \" \" $2 \" installed \"}'$(uname -m)"
				} else {
					return nil, errors.New("could not detect a supported package manager")
				}
			}
		}
	}

	// Execute the appropriate command based on the detected distribution/package manager
	librariesOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, command)
	if err != nil {
		return nil, err
	}

	// Parse the libraries info
	return parseLibrariesInfo(librariesOutput), nil
}

// ListFileSystem implements the Service interface
func (s *service) ListFileSystem(ctx context.Context, sessionID string, path string, recursive bool, includeHidden bool) (*FileSystemListing, error) {
	// Sanitize the path to prevent command injection
	sanitizedPath := strings.Trim(sanitizePath(path), "'")

	var entries []FileSystemEntry
	var err error

	// Get directory contents based on recursive flag
	if recursive {
		entries, err = getRecursiveDirectoryContents(ctx, s.sessionRepo, sessionID, sanitizedPath, includeHidden)
	} else {
		entries, err = getNonRecursiveDirectoryContents(ctx, s.sessionRepo, sessionID, sanitizedPath, includeHidden)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	// Create the file system listing result
	result := &FileSystemListing{
		Path:      sanitizedPath,
		Entries:   entries,
		Recursive: recursive,
	}

	return result, nil
}

// GetFileDetails implements the Service interface
func (s *service) GetFileDetails(ctx context.Context, sessionID string, path string) (*FileSystemEntry, error) {
	// Sanitize the path to prevent command injection
	sanitizedPath := strings.Trim(sanitizePath(path), "'")

	// Get detailed file information
	fileInfo, err := getEnhancedFileInfo(ctx, s.sessionRepo, sessionID, sanitizedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file details: %w", err)
	}

	// For files, get additional information such as mime type and a preview (for text files)
	if fileInfo.Type == "file" {
		// Get file mime type
		mimeTypeCmd := fmt.Sprintf("file --mime-type -b %s", sanitizePath(sanitizedPath))
		mimeTypeOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, mimeTypeCmd)
		if err == nil {
			fileInfo.MimeType = strings.TrimSpace(mimeTypeOutput)

			// If it's a text file, get a preview (first 10 lines)
			if strings.HasPrefix(fileInfo.MimeType, "text/") {
				previewCmd := fmt.Sprintf("head -n 10 %s", sanitizePath(sanitizedPath))
				previewOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, previewCmd)
				if err == nil {
					fileInfo.Preview = previewOutput
				}
			}
		}
	}

	return fileInfo, nil
}

// SearchFiles implements the Service interface
func (s *service) SearchFiles(ctx context.Context, sessionID string, path string, pattern string, maxDepth int) ([]FileSystemEntry, error) {
	// Sanitize the path and pattern to prevent command injection
	sanitizedPath := strings.Trim(sanitizePath(path), "'")
	sanitizedPattern := strings.ReplaceAll(pattern, "'", "'\\''") // Escape single quotes

	// Build the find command for searching files
	// The -maxdepth parameter limits the search depth to avoid searching the entire filesystem
	// We're searching for files that match the pattern in either their name or content
	findNameCmd := fmt.Sprintf("find %s -maxdepth %d -type f -name '*%s*' -o -type d -name '*%s*'",
		sanitizePath(sanitizedPath), maxDepth, sanitizedPattern, sanitizedPattern)

	// Execute the find command
	nameOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, findNameCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to search files by name: %w", err)
	}

	// For content search (grep), we'll search only in text files for the pattern
	// This can be resource-intensive, so we'll limit it to files within the maxdepth
	grepCmd := fmt.Sprintf("find %s -maxdepth %d -type f -exec grep -l '%s' {} \\; 2>/dev/null",
		sanitizePath(sanitizedPath), maxDepth, sanitizedPattern)

	// Execute the grep command
	contentOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, grepCmd)
	// We don't check for error here as grep might return non-zero if no matches are found

	// Combine and deduplicate the results
	var allPaths []string

	// Process the name search results
	if nameOutput != "" {
		namePaths := strings.Split(nameOutput, "\n")
		for _, path := range namePaths {
			if path != "" {
				allPaths = append(allPaths, path)
			}
		}
	}

	// Process the content search results
	if contentOutput != "" {
		contentPaths := strings.Split(contentOutput, "\n")
		for _, path := range contentPaths {
			if path != "" {
				found := false
				for _, existingPath := range allPaths {
					if existingPath == path {
						found = true
						break
					}
				}
				if !found {
					allPaths = append(allPaths, path)
				}
			}
		}
	}

	// Get detailed information for each found file/directory
	var entries []FileSystemEntry
	for _, filePath := range allPaths {
		fileInfo, err := getEnhancedFileInfo(ctx, s.sessionRepo, sessionID, filePath)
		if err != nil {
			continue
		}
		entries = append(entries, *fileInfo)
	}

	return entries, nil
}
