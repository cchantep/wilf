# Wilf

[Wilf](https://discworld.fandom.com/wiki/Gods#Wilf) is the god of astrology, or a CI friendly utility to check Python dependency updates according Pipfile.

## Usage

Locally: `wilf [OPTIONS] /path/to/Pipefile`

Options:

- `-c FILE` : Use `FILE` as path to the configuration file (see [Configuration](#configuration) thereafter)
- `-h` : Print the usage and exit
- `-r REPORTER` : Use REPORTER as the reporter; It can be specified multi time to set multiple reporters; Valid options are `monochrome-table`, `colorized-table` (the default one), `junit` (JUnit reporting on stdout), or `junit:/path/to/output/junit.xml`
- `-v` : Enable verbose output

Example:

```bash
wilf -h  # Print usage
wilf /path/to/Pipfile  # Minimal usage
wilf -c /path/to/config.toml -r junit /path/to/Pipfile
wilf -v -c /path/to/config.toml /path/to/Pipfile
wilf -r junit:/tmp/junit.xml -r colorized-table
```

> As soon as a `-r REPORTER` option is specified, the default reporter (`colorized-table` is overriden).

With Docker: ![Docker Latest Image](https://img.shields.io/docker/v/cchantep/wilf)

```bash
docker run --rm -it cchantep/wilf [arguments...]
```

## Configuration

Some settings can be defined in the configuration file.

```toml
check_dev_packages = true  # default: false
excluded_packages = ["pkg1", "pkg2"]  # default: []
update_level = "major"  # major|minor|patch; default: minor
```

A Gitlab Package registry can also be configured:

```toml
[gitlab]
project_api_packages_url = "https://gitlab.com/api/v4/projects/12345678/packages"
private_token = "YOUR_PRIVATE_TOKEN"  # Personal or CI token
```

## Integration

## Gitlab CI

Wilf can be configured as job in Gitlab CI.

```yaml
.wilf:
  image:
    name: cchantep/wilf:latest  # Recommended to set explicit version
    entrypoint: ['']
  script:
    - |
      echo "Check Pipfile dependencies in $BASEDIR ..."
      cd "$BASEDIR"
      wilf -c wilf.conf -r junit:wilf-junit.xml -r monochrome-table Pipfile
  artifacts:
    when: always
    reports:
      junit: "$BASEDIR/wilf-junit.xml"

My job:
  extends: .wilf
  variables:
    BASEDIR: "relative/path"
  # ...
```

## Build

The project is built using [Go](https://golang.org/) 1.20+.

Then to execute the incremental build:

    go build

Run the tests: [![CI](https://github.com/cchantep/wilf/actions/workflows/ci.yml/badge.svg)](https://github.com/cchantep/wilf/actions/workflows/ci.yml)

    go test

Build the [Docker image](https://hub.docker.com/r/cchantep/wilf):

```bash
export VERSION=...
docker buildx build --push --platform=linux/amd64,linux/arm64 --build-arg "version=$VERSION" -t "cchantep/wilf:$VERSION" .
```

> See [Docker documentation](https://www.docker.com/blog/faster-multi-platform-builds-dockerfile-cross-compilation-guide/)

### Release

Add a release tag:

    git tag -a 1.2.3 -m '...'

Build the Docker image:

    ./tooling/scripts/docker_build.sh

Publish the Docker image:

    docker.io/cchantep/wilf:1.2.3
