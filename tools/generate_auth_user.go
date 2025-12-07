package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/phuslu/log"
)

// AuthUserParams defines parameters for a single user
type AuthUserParams struct {
	Username    string            `json:"username"`
	Password    string            `json:"password"`
	SpeedLimit  int64             `json:"speed_limit,omitempty"`
	AllowTunnel bool              `json:"allow_tunnel,omitempty"`
	AllowClient bool              `json:"allow_client,omitempty"`
	AllowSSH    bool              `json:"allow_ssh,omitempty"`
	AllowWebDAV bool              `json:"allow_webdav,omitempty"`
	Attrs       map[string]string `json:"attrs,omitempty"` // For any extra attributes
}

// GenerateAuthUserConfigParams parameters for the tool
type GenerateAuthUserConfigParams struct {
	Users []AuthUserParams `json:"users"`
}

// GenerateAuthUserConfig generates auth_user.csv content
func GenerateAuthUserConfig(arguments json.RawMessage) (string, error) {
	var params GenerateAuthUserConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid parameters for auth user configuration",
		)
	}

	log.Info().Int("users_count", len(params.Users)).Msg("generating auth_user.csv")

	var sb strings.Builder
	// Header: username,password,speed_limit,allow_tunnel,allow_client,allow_ssh,allow_webdav
	sb.WriteString("username,password,speed_limit,allow_tunnel,allow_client,allow_ssh,allow_webdav\n")

	for _, u := range params.Users {
		// Bool to int (0/1)
		tunnel := 0
		if u.AllowTunnel {
			tunnel = 1
		}
		client := 0
		if u.AllowClient {
			client = 1
		}
		ssh := 0
		if u.AllowSSH {
			ssh = 1
		}
		webdav := 0
		if u.AllowWebDAV {
			webdav = 1
		}

		// Handle default speed limit if not set?
		// liner example uses -1 or 0. -1 usually means no limit? 0 means default?
		// Let's assume user provides exact value, or we assume 0.
		// If user wants -1, they pass -1.

		line := fmt.Sprintf("%s,%s,%d,%d,%d,%d,%d\n",
			u.Username,
			u.Password,
			u.SpeedLimit,
			tunnel,
			client,
			ssh,
			webdav,
		)
		sb.WriteString(line)
	}

	return responses.SuccessResponse(sb.String(), "Generated auth_user.csv configuration")
}
