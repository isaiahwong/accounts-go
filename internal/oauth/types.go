package oauth

import "time"

// HydraResponse defines hydra's data response
// https://www.ory.sh/docs/hydra/sdk/api
type HydraResponse struct {
	Challenge string `json:"challenge"`
	Client    struct {
		AllowedCorsOrigins                []string  `json:"allowed_cors_origins"`
		Audience                          []string  `json:"audience"`
		BackchannelLogoutSessionRequired  bool      `json:"backchannel_logout_session_required"`
		BackchannelLogoutURI              string    `json:"backchannel_logout_uri"`
		ClientID                          string    `json:"client_id"`
		ClientName                        string    `json:"client_name"`
		ClientSecret                      string    `json:"client_secret"`
		ClientSecretExpiresAt             int       `json:"client_secret_expires_at"`
		ClientURI                         string    `json:"client_uri"`
		Contacts                          []string  `json:"contacts"`
		CreatedAt                         time.Time `json:"created_at"`
		FrontchannelLogoutSessionRequired bool      `json:"frontchannel_logout_session_required"`
		FrontchannelLogoutURI             string    `json:"frontchannel_logout_uri"`
		GrantTypes                        []string  `json:"grant_types"`
		Jwks                              struct {
			Keys []struct {
				Alg string   `json:"alg"`
				Crv string   `json:"crv"`
				D   string   `json:"d"`
				Dp  string   `json:"dp"`
				Dq  string   `json:"dq"`
				E   string   `json:"e"`
				K   string   `json:"k"`
				Kid string   `json:"kid"`
				Kty string   `json:"kty"`
				N   string   `json:"n"`
				P   string   `json:"p"`
				Q   string   `json:"q"`
				Qi  string   `json:"qi"`
				Use string   `json:"use"`
				X   string   `json:"x"`
				X5C []string `json:"x5c"`
				Y   string   `json:"y"`
			} `json:"keys"`
		} `json:"jwks"`
		JwksURI  string `json:"jwks_uri"`
		LogoURI  string `json:"logo_uri"`
		Metadata struct {
		} `json:"metadata"`
		Owner                        string    `json:"owner"`
		PolicyURI                    string    `json:"policy_uri"`
		PostLogoutRedirectUris       []string  `json:"post_logout_redirect_uris"`
		RedirectUris                 []string  `json:"redirect_uris"`
		RequestObjectSigningAlg      string    `json:"request_object_signing_alg"`
		RequestUris                  []string  `json:"request_uris"`
		ResponseTypes                []string  `json:"response_types"`
		Scope                        string    `json:"scope"`
		SectorIdentifierURI          string    `json:"sector_identifier_uri"`
		SubjectType                  string    `json:"subject_type"`
		TokenEndpointAuthMethod      string    `json:"token_endpoint_auth_method"`
		TosURI                       string    `json:"tos_uri"`
		UpdatedAt                    time.Time `json:"updated_at"`
		AccountinfoSignedResponseAlg string    `json:"userinfo_signed_response_alg"`
	} `json:"client"`
	Context struct {
		Property1 []byte `json:"property1"`
		Property2 []byte `json:"property2"`
	} `json:"context"`
	LoginChallenge string `json:"login_challenge"`
	LoginSessionID string `json:"login_session_id"`
	OidcContext    struct {
		AcrValues         []string `json:"acr_values"`
		Display           string   `json:"display"`
		IDTokenHintClaims struct {
			Property1 struct {
			} `json:"property1"`
			Property2 struct {
			} `json:"property2"`
		} `json:"id_token_hint_claims"`
		LoginHint string   `json:"login_hint"`
		UILocales []string `json:"ui_locales"`
	} `json:"oidc_context"`
	RequestURL                   string   `json:"request_url"`
	RequestedAccessTokenAudience []string `json:"requested_access_token_audience"`
	RequestedScope               []string `json:"requested_scope"`
	SessionID                    string   `json:"session_id"`
	RPInitiated                  bool     `json:"rp_initiated"`
	SID                          string   `json:"sid"`
	Skip                         bool     `json:"skip"`
	Subject                      string   `json:"subject"`
}

type HydraLoginAccept struct {
	Acr     string `json:"acr"`
	Context struct {
		Property1 []byte `json:"property1"`
		Property2 []byte `json:"property2"`
	} `json:"context"`
	ForceSubjectIdentifier string `json:"force_subject_identifier"`
	Remember               bool   `json:"remember"`
	RememberFor            int    `json:"remember_for"`
	Subject                string `json:"subject"`
}

type HydraError struct {
	ErrorName        string `json:"error"`
	ErrorDebug       string `json:"error_debug"`
	ErrorDescription string `json:"error_description"`
	ErrorHint        string `json:"error_hint"`
	StatusCode       int    `json:"status_code"`
}

func (h *HydraError) Error() string {
	return h.ErrorName
}

type AccessToken struct {
	Property1 struct {
	} `json:"property1"`
	Property2 struct {
	} `json:"property2"`
}

type Session struct {
	AccessToken AccessToken       `json:"access_token"`
	IDToken     map[string]string `json:"id_token"`
}

type HydraConsentAccept struct {
	GrantAccessTokenAudience []string  `json:"grant_access_token_audience"`
	GrantScope               []string  `json:"grant_scope"`
	HandledAt                time.Time `json:"handled_at"`
	Remember                 bool      `json:"remember"`
	RememberFor              int       `json:"remember_for"`
	Session                  Session   `json:"session"`
}

type HydraRedirect struct {
	RedirectTo string `json:"redirect_to"`
}

type Ext struct {
	Property1 struct {
	} `json:"property1"`
	Property2 struct {
	} `json:"property2"`
}

type IntrospectRequest struct {
	Token string `json:"token"`
	Scope string `json:"scope"`
}

type InstrospectResponse struct {
	Active            bool     `json:"active"`
	Aud               []string `json:"aud"`
	ClientID          string   `json:"client_id"`
	Exp               int64    `json:"exp"`
	Ext               Ext      `json:"ext"`
	Iat               int64    `json:"iat"`
	Iss               string   `json:"iss"`
	Nbf               int64    `json:"nbf"`
	ObfuscatedSubject string   `json:"obfuscated_subject"`
	Scope             string   `json:"scope"`
	Sub               string   `json:"sub"`
	TokenType         string   `json:"token_type"`
	Username          string   `json:"username"`
}
