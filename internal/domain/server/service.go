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
