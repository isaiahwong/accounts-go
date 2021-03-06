# Get all clients
```
curl https://accounts.isaiahwong.dev/oauth2admin/clients
```
# Creating a client key
```
curl -X POST  https://accounts.isaiahwong.dev/oauth2admin/clients \
-H 'Content-Type: application/json' -H 'Accept: application/json' \
  --data '{
  "client_id": "3494ea93c8683e12dad8e918b8b6f24b8ab5df15",
  "client_secret": "d5d68a699d7b1821b7255edea14f9b0a96b5bc17",
  "grant_types": [
    "authorization_code",
    "refresh_token",
    "client_credentials",
    "implicit"
  ],
  "response_types": [
    "token",
    "code",
    "id_token"
  ],
  "scope": "openid offline photos.read",
  "redirect_uris": [
    "https://accounts.isaiahwong.dev/client/claims/noop"
  ]
}'
```

# local client key
```
curl -X POST -k http://127.0.0.1/oauth2/admin/clients \
-H 'Content-Type: application/json' -H 'Accept: application/json' \
  --data '{
  "client_id": "2dad83494ea93c861183e123138b8b6f24b8ab5d12312df15",
  "client_secret": "f9b0a96bd5d68a699d7b1821b7255edea145bc17",
  "grant_types": [
    "authorization_code",
    "refresh_token",
    "client_credentials",
    "implicit"
  ],
  "response_types": [
    "token",
    "code",
    "id_token"
  ],
  "scope": "openid offline",
  "redirect_uris": [
    "http://127.0.0.1/client/claims/noop"
  ]
}'
```

# Introspect
```
curl -X POST -k https://127.0.0.1/oauth2/admin/oauth2/introspect \
-H 'Content-Type: application/x-www-form-urlencoded' -H 'Accept: application/json' \
  --data 'token=FUr0B87IP1zLviGjqSTdsXGYn-fWVzUet0HTlPjqUcw.-NNoO7fZqDzbnq6BAWPqGf-AeGujBOM1bLO-x2t9cV8'
```
# Delete a client key
```
curl -X DELETE https://accounts.isaiahwong.dev/oauth2admin/clients/isaiahwongdev
```

# Run Oauth consumer
```
docker run --rm -it \
-p 9010:9010 \
oryd/hydra:v1.3.2 \
token user --skip-tls-verify \
  --port 9010 \
  --auth-url https://accounts.isaiahwong.dev/issuer/oauth2/auth \
  --token-url https://accounts.isaiahwong.dev/issuer/oauth2/token \
  --client-id 3494ea93c8683e12dad8e918b8b6f24b8ab5df15\
  --client-secret d5d68a699d7b1821b7255edea14f9b0a96b5bc17 \
  --scope "openid,offline,photos.read"
```