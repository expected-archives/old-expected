# expected

### lancer les migrations

Il faut installer le tool [disponible ici](https://github.com/golang-migrate/migrate/tree/master/cli).
Puis executer la commande :

```
migrate -database=$POSTGRES_ADDR -path=migrations up
```
