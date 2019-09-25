# drone-sonar-qualitygate-plugin

The plugin of Drone CI to integrate the response of the SonarScanner against SonarQube (previously called Sonar), which is an open source code quality management platform.

Detail tutorials: [DOCS.md](DOCS.md).

## Build process

build go binary file
`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o drone-sonar-qualitygate`

build docker image
`docker build -t dxas90/drone-sonar-qualitygate-plugin .`

## Testing the docker image

```commandline
docker run --rm \
  -e DRONE_REPO=test \
  -e PLUGIN_SOURCES=. \
  -e SONAR_TOKEN=60878847cea1a31d817f0deee3daa7868c431433 \
  dxas90/drone-sonar-qualitygate-plugin
```

## Pipeline example

```yaml
steps
- name: code-analysis
  image: dxas90/drone-sonar-qualitygate-plugin
  settings:
      sonar_token:
        from_secret: sonar_token
```
