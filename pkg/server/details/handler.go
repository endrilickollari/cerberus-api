package details

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/ssh"
	"net/http"
	"remote-server-api/pkg/login"
	"remote-server-api/pkg/utils"
	"strings"
)

// ServerDetailsHandler retrieves and returns server details using an authenticated SSH session.
//
// @Summary Get server details
// @Description Retrieves server details using an authenticated SSH session.
// @Tags server
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} ServerDetails "Server details"
// @Failure 401 {object} string "Invalid token or session expired"
// @Failure 500 {object} string "Failed to get server details"
// @Router /server-details [post]
func ServerDetailsHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &login.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return login.JwtKey, nil
	})

	if err != nil || !tkn.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Retrieve SSH session using the session ID from the token
	client, exists := login.GetSession(claims.SessionID)
	if !exists {
		http.Error(w, "Session expired or not found", http.StatusUnauthorized)
		return
	}

	serverDetails, err := GetServerDetails(client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get server details: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serverDetails)
}

func GetServerDetails(client *ssh.Client) (*ServerDetails, error) {
	hostname, _ := utils.RunCommand(client, "hostname")
	osInfo, _ := utils.RunCommand(client, "uname -a")
	kernelVersion, _ := utils.RunCommand(client, "uname -r")
	uptime, _ := utils.RunCommand(client, "uptime")
	//cpuModel, _ := utils.RunCommand(client, "cat /proc/cpuinfo | grep 'model name' | head -1")
	//cpuLoad, _ := utils.RunCommand(client, "uptime | awk -F 'load average: ' '{print $2}'")
	//totalMem, _ := utils.RunCommand(client, "free -m | grep Mem | awk '{print $2}'")
	//usedMem, _ := utils.RunCommand(client, "free -m | grep Mem | awk '{print $3}'")
	//freeMem, _ := utils.RunCommand(client, "free -m | grep Mem | awk '{print $4}'")
	//diskUsage, _ := utils.RunCommand(client, "df -h")
	//ipAddresses, _ := utils.RunCommand(client, "ip addr | grep inet")
	//networkInterfaces, _ := utils.RunCommand(client, "ip link show")
	//openPorts, _ := utils.RunCommand(client, "ss -tuln")
	//loggedInUsers, _ := utils.RunCommand(client, "w")
	//runningProcesses, _ := utils.RunCommand(client, "ps aux")
	//firewallRules, _ := utils.RunCommand(client, "sudo iptables -L")

	return &ServerDetails{
		Hostname:      hostname,
		OS:            osInfo,
		KernelVersion: kernelVersion,
		Uptime:        uptime,
	}, nil
}
