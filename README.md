# Expected

### Installation

```
docker-compose up
```

#### Configuration du docker-compose

| Nom | Description | Valeur par défaut |
| --- | --- | --- |
| POSTGRES_ADDR | Permet de changer l'adresse de postgres | postgres://expected:expected@postgres/expected?sslmode=disable |
| STAN_ADDR | Permet de changer l'adresse de nats | nats://localhost:4222 |
| AUTH_ADDR | Permet de changer l'adresse du service d'authentification | http://localhost:3002 |
| GITHUB_CLIENT_ID | Défini le client id pour l'oauth avec github |  |
| GITHUB_CLIENT_SECRET | Défini le client secret pour l'oauth avec github |  |
| DASHBOARD_URL | L'url du dashboard (utilisé pour definir le cookie d'authentification et rediriger l'utilisateur) | http://localhost:8080 |
| REGISTRY_AUTH_TOKEN_REALM | L'adresse du serveur d'authentification de la registry qui sera donné au client | http://localhost:3001/registry/auth |
| REGISTRY_AUTH_SERVER | L'adresse du serveur qui recoit les events de la registry | http://registryhook:3001/registry/hook |

### Lancer les migrations

Il faut installer le tool [disponible ici](https://github.com/golang-migrate/migrate/tree/master/cli).
Puis executer la commande :

```
migrate -database=$POSTGRES_ADDR -path=migrations up
```

### Utiliser les runners pour le controller

```
export DOCKER_HOST=tcp://51.15.236.158:2376 DOCKER_TLS_VERIFY=1 DOCKER_CERT_PATH="$(pwd)/certs/docker"
mkdir -p $(pwd)/certs/docker
scp "root@51.15.236.158:/root/.docker/*" $(pwd)/certs/docker
```
