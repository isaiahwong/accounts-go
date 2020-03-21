package recaptcha

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/isaiahwong/accounts-go/internal/common"
	"github.com/joho/godotenv"
)

// Response structures reCAPTCHA's response
// https://developers.google.com/recaptcha/docs/verify#api_response
type Response struct {
	Success     bool        `json:"success"`
	ChallengeTS string      `json:"challenge_ts"`
	Hostname    string      `json:"hostname"`
	ErrorCodes  interface{} `json:"error-codes"`
}

// Verify sends a request to reCAPTCHA to verify client token
func Verify(token string, ip string) (*Response, error) {
	godotenv.Load()
	u := common.MapEnvWithDefaults("GOOGLE_RECAPTCHA_URL", "https://www.google.com/recaptcha/api/siteverify")
	s := common.MapEnvWithDefaults("GOOGLE_RECAPTCHA_SECRET", "_")

	resp, err := http.Post(fmt.Sprintf("%v?secret=%v&response=%v&remoteip=%v", u, s, token, ip), "text/plain", nil)
	if err != nil {
		return nil, err
	}
	r := &Response{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
