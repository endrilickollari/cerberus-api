package main

import (
	"log"
	"net/http"
	"os"
	"remote-server-api/pkg/login"
	"remote-server-api/pkg/server/details"
	"remote-server-api/pkg/server/details/cpu_info"
	"remote-server-api/pkg/server/details/disk_usage"
	"remote-server-api/pkg/server/details/running_processes"
	"remote-server-api/pkg/server/docker"
)

func main() {
	// login
	http.HandleFunc("/login", login.LoginHandler)

	// details
	http.HandleFunc("/server-details", login.TokenValidationMiddleware(details.ServerDetailsHandler))
	http.HandleFunc("/server-details/cpu-info", login.TokenValidationMiddleware(cpu_info.GetCPUInfo))
	http.HandleFunc("/server-details/disk-usage", login.TokenValidationMiddleware(disk_usage.GetDiskUsageInfo))
	http.HandleFunc("/server-details/running-processes", login.TokenValidationMiddleware(running_processes.GetRunningProcessesInfo))

	// docker
	http.HandleFunc("/docker/container-details", login.TokenValidationMiddleware(docker.GetContainerInfo))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback for local development
	}
	addr := ":" + port
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
