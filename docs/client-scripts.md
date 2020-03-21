# Creating a client key
```
curl -X POST  http://127.0.0.1/admin/clients \
-H 'Content-Type: application/json' -H 'Accept: application/json' \
  --data '{ "client_id": "test1", "client_secret": "test12345678", "grant_types": ["authorization_code","refresh_token","client_credentials","implicit"], "response_types": ["token","code","id_token"], "scope": "openid offline photos.read", "redirect_uris": ["http://127.0.0.1:9010/callback"] }'
```

# Run Oauth consumer
```
docker run --rm -it \
-p 9010:9010 \
oryd/hydra:v1.3.2 \
token user --skip-tls-verify \
  --port 9010 \
  --auth-url http://127.0.0.1/issuer/oauth2/auth \
  --token-url http://127.0.0.1/issuer/oauth2/token \
  --client-id test1\
  --client-secret test12345678 \
  --scope "openid,offline,photos.read"
```