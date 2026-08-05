**Languages:** [English](../README.md) | [日本語](README.ja.md) | [فارسی](README.fa.md) | [Deutsch](README.de.md) | [简体中文](README.zh-CN.md) | [Русский](README.ru.md) | [Español](README.es.md) | [Português](README.pt-BR.md) | [Français](README.fr.md) | [한국어](README.ko.md)

# DCB — Docker Compose Builder

> **Écrivez moins de YAML. Exécutez de meilleurs conteneurs.**

DCB lit un `dcb.yaml` simplifié et génère un `docker-compose.yml` prêt pour la production
avec des health checks automatiques, des réseaux nommés, des politiques de redémarrage et un ordonnancement des dépendances.

[![CI](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml/badge.svg)](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/JrNovaEX/DCB)](https://goreportcard.com/report/github.com/JrNovaEX/DCB)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)

---

## Pourquoi DCB ?

`docker-compose.yml` est puissant mais verbeux. Un stack de production typique répète le même
boilerplate pour les health checks, les politiques de redémarrage, les déclarations de réseau et le nommage des volumes.

Avec DCB, vous écrivez ceci :

```yaml
# dcb.yaml
project: webapp

services:
  api:
    build: ./api
    port: 8080
    depends_on: [db]
    healthcheck: true
    env_file: .env

  db:
    image: postgres:15-alpine
    port: 5432
    env:
      POSTGRES_DB: webapp
    volume: pgdata:/var/lib/postgresql/data
    healthcheck: true
```

Et vous obtenez un `docker-compose.yml` entièrement configuré avec :

- ✅ Health check correct `pg_isready` pour postgres (conscient de l'image)
- ✅ `depends_on` avec `condition: service_healthy`
- ✅ Réseau nommé `webapp_default`
- ✅ Déclaration top-level de `volumes:`
- ✅ `restart: unless-stopped` sur chaque service
- ✅ Détection des dépendances circulaires au moment du parsing

---

## Installation

### Homebrew (macOS / Linux)
```bash
brew install JrNovaEX/tap/dcb
```

### Go install
```bash
go install github.com/JrNovaEX/DCB/cmd/dcb@latest
```

### Télécharger le binaire

Récupérez la dernière version pour votre plateforme depuis la
[page Releases](https://github.com/JrNovaEX/DCB/releases).

```bash
# Exemple Linux amd64
curl -Lo dcb https://github.com/JrNovaEX/DCB/releases/latest/download/dcb_linux_amd64
chmod +x dcb && sudo mv dcb /usr/local/bin/
```

### Compiler depuis les sources
```bash
git clone https://github.com/JrNovaEX/DCB
cd DCB
make install
```

---

## Démarrage rapide

```bash
# 1. Créer interactivement un nouveau dcb.yaml
dcb init

# 2. Valider la configuration et vérifier la disponibilité des ports
dcb check

# 3. Générer docker-compose.yml
dcb build

# 4. Démarrer tout en arrière-plan
dcb up --detach

# 5. Suivre les logs de tous les services
dcb logs -f

# 6. Arrêter et nettoyer
dcb down
```

---

## Référence dcb.yaml

```yaml
project: myapp          # Obligatoire. Utilisé pour le namespace des réseaux et volumes.
version: "3.8"          # Optionnel. Version du schema Docker Compose (par défaut : 3.8).

services:
  <service-name>:
    # Source — choisissez UN parmi image ou build :
    image: nginx:alpine

    build: ./myservice          # Forme courte : chemin du context
    build:                      # Forme longue :
      context: ./myservice
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production

    port: 8080           # Exposé en host:container (même valeur des deux côtés).
    env_file: .env       # Chemin vers un fichier env.
    env:                 # Variables d'environnement en ligne.
      KEY: value

    depends_on:          # Noms des services à attendre (condition: service_healthy).
      - db

    healthcheck: true    # Génère automatiquement un health check basé sur le nom de l'image.
                         # Supportés : postgres, mysql, mariadb, redis, nginx,
                         #             mongo, rabbitmq. Les autres reçoivent CMD-SHELL exit 0.

    volume: pgdata:/var/lib/postgresql/data   # Named ou bind mount.
    restart: unless-stopped                   # Par défaut. Peut être surchargé.
```

---

## Commandes

| Commande | Description |
|---------|-------------|
| `dcb init` | Créer interactivement un nouveau `dcb.yaml` |
| `dcb check` | Valider `dcb.yaml` + vérifier la disponibilité des ports de l'hôte |
| `dcb build` | Générer `docker-compose.yml` à partir de `dcb.yaml` |
| `dcb up` | Construire puis démarrer les services (`docker compose up`) |
| `dcb down` | Arrêter et supprimer les services (`docker compose down`) |
| `dcb logs` | Diffuser les logs des services |
| `dcb version` | Afficher la version, le commit et la date de build |

### Flags globaux
```
--verbose, -V    Activer les logs de débogage
--json           Sortie des logs en JSON (pour les pipelines CI/CD)
```

### dcb build
```
--file,     -f   Fichier d'entrée (par défaut : dcb.yaml)
--output,   -o   Fichier de sortie (par défaut : docker-compose.yml)
--env,      -e   Environnement cible (par défaut : dev)
--validate       Exécuter docker compose config après génération (par défaut : true)
```

### dcb up
```
--file,     -f   dcb.yaml d'entrée (par défaut : dcb.yaml)
--env,      -e   Environnement cible
--detach,   -d   Exécuter en arrière-plan
--build          Forcer la reconstruction des images
```

### dcb down
```
--file,     -f   Fichier Compose (par défaut : docker-compose.yml)
--project,  -p   Remplacer le nom du projet
--volumes,  -v   Supprimer les volumes nommés
--remove-orphans Supprimer les conteneurs orphelins
```

### dcb logs
```
--file           Fichier Compose (par défaut : docker-compose.yml)
--follow,   -f   Diffuser la sortie des logs
--tail,     -n   Nombre de lignes à afficher (-1 = toutes)
[services...]    Filtrer par noms de services spécifiques
```

---

## Support multi-environnements

DCB fusionne un `dcb.yaml` de base avec un overlay spécifique à l'environnement :

```
dcb build --env prod
```

Charge `dcb.yaml`, puis fusionne `dcb.prod.yaml` par-dessus. Les clés de `dcb.prod.yaml` écrasent la base ; les clés absentes retombent sur la base.

---

## Architecture

DCB suit la **Clean Architecture** (Hexagonale) :

```
cmd/dcb/
  commands/         Couche CLI — commandes Cobra, parsing des flags, câblage DI
internal/
  domain/           Entités principales — ProjectConfig, Service, Volume, Port
  ports/            Interfaces — ComposeGenerator, ConfigParser, DockerRunner
  usecases/         Logique applicative — Build, Check, Up, Down, Logs
  adapters/
    parser/         Parseur YAML (gopkg.in/yaml.v3, mode strict)
    generator/      Générateur de templates docker-compose.yml (text/template)
    runner/         Exécuteur Docker CLI (os/exec)
    config/         Chargeur de configuration multi-environnements (Viper)
  infrastructure/
    logger/         Initialiseur slog
    validator/      Wrapper go-playground/validator
pkg/version/        Informations de version au moment du build
```

Les dépendances circulent vers l'intérieur : **cmd → usecases → domain**. Les adapters dépendent des ports, pas les uns des autres.

---

## Développement

```bash
# Lancer tous les tests
make test

# Tests avec race detector et coverage
make cover

# Lint
make lint

# Format
make fmt

# Construire le binaire local
make build
./dist/dcb version
```

### Lancer les tests

```bash
go test ./...                          # tous les packages
go test ./internal/domain/...          # domain uniquement
go test ./internal/adapters/parser/... # parser uniquement
```

---

## Contribuer

1. Forkez le dépôt
2. Créez une branche de fonctionnalité (`git checkout -b feat/my-feature`)
3. Committez en utilisant [Conventional Commits](https://www.conventionalcommits.org/)
4. Ouvrez une pull request

Assurez-vous que `make test lint` passe avant d'ouvrir une PR.

---

## Licence

MIT — voir [LICENSE](../LICENSE).
