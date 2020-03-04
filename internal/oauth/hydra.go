package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/isaiahwong/accounts-go/internal/util"
)

type Hydra struct {
	hydraUrl string
}

func (h *Hydra) get(flow string, challenge string) (*HydraResponse, error) {
	url := fmt.Sprintf("%v/%v?=%v_challenge=%v", h.hydraUrl, flow, flow, challenge)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 302 {
		return nil, errors.New("An error while making request to hydra " + string(b))
	}
	r := &HydraResponse{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (h *Hydra) put(flow string, action string, challenge string, body interface{}) (*HydraRedirect, error) {
	url := fmt.Sprintf("%v/%v/%v?=%v_challenge=%v", h.hydraUrl, flow, action, flow, challenge)
	d, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 302 {
		return nil, errors.New("An error while making request to hydra " + string(b))
	}

	r := &HydraRedirect{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Login processes hydra oauth login
func (h *Hydra) Login(challenge string) (*HydraResponse, error) {
	return h.get("login", challenge)
}

func (h *Hydra) AcceptLogin(challenge string, body *HydraLoginAccept) (*HydraRedirect, error) {
	return h.put("login", "accept", challenge, body)
}

func (h *Hydra) RejectLogin(challenge string, body *HydraReject) (*HydraRedirect, error) {
	return h.put("login", "reject", challenge, body)
}

func (h *Hydra) Consent(challenge string) (*HydraResponse, error) {
	return h.get("consent", challenge)
}

func (h *Hydra) AcceptConsent(challenge string, body *HydraConsentAccept) (*HydraRedirect, error) {
	return h.put("consent", "accept", challenge, body)
}

func (h *Hydra) RejectConsent(challenge string, body *HydraReject) (*HydraRedirect, error) {
	return h.put("consent", "reject", challenge, body)
}

func (h *Hydra) Logout(challenge string) (*HydraResponse, error) {
	return h.get("logout", challenge)
}

func (h *Hydra) AcceptLogout(challenge string) (*HydraRedirect, error) {
	return h.put("logout", "accept", challenge, nil)
}

func (h *Hydra) RejectLogout(challenge string, body *HydraReject) (*HydraRedirect, error) {
	return h.put("logout", "reject", challenge, nil)
}

// NewHydraClient
func NewHydraClient() *Hydra {
	url := util.MapEnvWithDefaults("HYDRA_ADMIN_URL", "http://localhost:9000")
	url += "/oauth2/auth/requests"
	return &Hydra{
		hydraUrl: url,
	}
}
