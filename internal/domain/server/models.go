package server

// ServerDetails contains basic information about the server
type ServerDetails struct {
	Hostname      string `json:"hostname"`
	OS            string `json:"os"`
	KernelVersion string `json:"kernel_version"`
	Uptime        string `json:"uptime"`
}

// CPUInfo contains detailed information about a CPU core
type CPUInfo struct {
	Processor       string `json:"processor"`
	VendorID        string `json:"vendor_id"`
	CPUFamily       string `json:"cpu_family"`
	Model           string `json:"model"`
	ModelName       string `json:"model_name"`
	Stepping        string `json:"stepping"`
	Microcode       string `json:"microcode"`
	CPUMHz          string `json:"cpu_mhz"`
	CacheSize       string `json:"cache_size"`
	PhysicalID      string `json:"physical_id"`
	Siblings        string `json:"siblings"`
	CoreID          string `json:"core_id"`
	CPUCores        string `json:"cpu_cores"`
	APICID          string `json:"apicid"`
	InitialAPICID   string `json:"initial_apicid"`
	FPU             string `json:"fpu"`
	FPUException    string `json:"fpu_exception"`
	CPUIDLevel      string `json:"cpuid_level"`
	WP              string `json:"wp"`
	Flags           string `json:"flags"`
	Bogomips        string `json:"bogomips"`
	ClflushSize     string `json:"clflush_size"`
	CacheAlignment  string `json:"cache_alignment"`
	AddressSizes    string `json:"address_sizes"`
	PowerManagement string `json:"power_management"`
}

// DiskUsage contains information about disk usage
type DiskUsage struct {
	Filesystem    string `json:"filesystem"`
	Size          string `json:"size"`
	Used          string `json:"used"`
	Available     string `json:"available"`
	UsePercentage string `json:"use_percentage"`
	MountedOn     string `json:"mounted_on"`
}

// ProcessInfo contains information about a running process
type ProcessInfo struct {
	User  string `json:"user"`
	PID   string `json:"process_id"`
	CPU   string `json:"cpu_consumption"`
	VSZ   string `json:"vsz"`
	RSS   string `json:"rss"`
	TTY   string `json:"tty"`
	Stat  string `json:"stat"`
	Start string `json:"started"`
	Time  string `json:"time"`
	CMD   string `json:"command"`
}
