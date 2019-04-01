```
docker tag golang localhost:5000/$USER_ID/golang 
docker push localhost:5000/$USER_ID/golang 
```

En bdd on va recevoir 8 layers + le layer final avec le tag.
Si il n'y a pas de tag le tag sera latest (comportement par défaut de la registry)

Voici a quoi ca ressemble en bdd (les infos cruciales)

```
                                    digest                                  |  tag   |   size    |           
----------------------------------------------------------------------------+--------+-----------+----
... sha256:786bc4873ebcc5cc05c0ff26d6ee3f2b26ada535067d07fc98f3ddb0ef4cd7c5 |        |       124 | ...
... sha256:e5c3f8c317dc30af45021092a3d76f16ba7aa1ee5f18fec742c84d4960818580 |        |   4336053 | ...
... sha256:193a6306c92af328dbd41bbbd3200a2c90802624cccfe5725223324428110d7f |        |  10740016 | ...
... sha256:bc9ab73e5b14b9fbd3687a4d8c1f1360533d6ee9ffc3f5ecc6630794b40257b7 |        |  45309934 | ...
... sha256:a587a86c9dcb9df6584180042becf21e36ecd8b460a761711227b4b06889a005 |        |  50065549 | ...
... sha256:1bc310ac474b880a5e4aeec02e6423d1304d137f1a8990074cb3ac6386a0b654 |        |  57591433 | ...
... sha256:87ab348d90cc687a586f07c0fd275335aee1e6e52c1995d1a7ac93fc901333bc |        | 126524470 | ...
... sha256:d817ad5b9beb8ed09a78819c7d7627679c89a4aca36a3b2d47760695d49d09a0 |        |      5459 | ...
... sha256:d7edb7c08dd224178faf5fa0bf0877d796e7833fca4f5c015777dee8711ff56e | latest |      1796 | ...

```

Les layers n'ont pas de tag.

La notification du tag arrive *normalement* en dernier.

#### Cas bizangouin

1. Si le gars ctrl c, son push:
    - on va avoir des layers sans tag 
    -> `SELECT name, string_agg(tag, '') as tag FROM images GROUP BY name`
    Commande pour qui retourne un tag vide si il n'y a pas de tag (donc pas fini de push)

2. Si le gars push une nouvelle version sur un meme tag:
    - on va avoir des layers en bdd non utilisé
 
3. Si un gars a une permission denied:
    - et bien pour une raison inconnu on aura un layer vide en notification et donc en bdd ¯\_(ツ)_/¯

4. Le gars dépasse ses 1gb ou son cota max:
    - Pendant x temps on aura son image dans notre registry 
 
 
 __Proposition solution cas 1 et 3__
 
 Une app qui tourne toutes les heures (configurable). Cette app va récupérer depuis maintenant à il y a une semaine
 tous les layers qui ont été crée groupé par `name` et `owner_id` en aggregant `tag`. 
 
 Quelque chose de ce style ce présentera:
 ```
  name  |   tag    ...
--------+----------...
 golang | latest   ...
 redis  | latestv1 ...
 neo4j  |          ...
```

Ici si le tag est vide c'est que les layers ont été push mais qu'il en manque.
On va regarder la date de chaque layer et si le layer date d'il y a plus de 4 heures (configurable), alors il est supprimé.

Une ombre survient au tableau, car ce cas fonctionne si et seulement si il n'a jamais re-push l'image entièrement.

Pour réglé ce problème on peut faire un deuxième scheduler qui tourne moins fréquement, beaucoup moins fréquement.
Lui groupera par `name` et `owner_id` mais va aggreger les tag avec une `,` et les digests avec une `,` aussi.
Ensuite on va pour chaque tags regarder les layers qui compose l'image, en faire une liste, regarder la liste complète 
aggrégué et supprimer ceux qui n'y sont pas.

