package server

import (
	"context"
	"errors"
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
