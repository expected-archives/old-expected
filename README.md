# Expected

### Installation

```
vagrant up
vagrant ssh
cd ~/expected
sh hack/start-services.sh
sh hack/start-apps.sh [apps...]
```

Liste des ports utilisés :
- apiserver (http: 3000, grpc: n.a)
- registryhook (http: 3001, grpc: n.a)
- authserver (http: 3002, grpc: 4002)
- controller (http: n.a, grpc: 4003)
- agent (http: n.a, grpc: n.a)

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
