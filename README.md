# Docker checker
The Docker Checker is small go binary, that checks if any of configured container images have a more recent versions available. Currently, only images on Dockerhub are supported. It is possible, though, to configure a custom registry host, which will be stripped away.

It notifies the configured email address when new versions are available.

## Installation
You can install the Docker Checker with its Docker container or by running the `app.go` directly. The container image is hosted on hub.docker.com under the name `iulbricht/docker-checker`.

## Configuration
The file `docker-checker.yaml` is a sample configuration for the Docker Checker.

The config file must contain a section `email` with the following config fields:

Variable   | Description
---------- | ------
`from`     | The mail address the updates should be sent to
`to`       | The mail address sending the updates
`username` | The username for the mail server
`password` | The password for the mail server
`host`     | The host of the mail server
`port`     | The port of the mail server

You can configure the images to check in the section `images`. An image can have the following configuration fields:

Variable      | Description
------------- | ------
`name`        | The name of the image in Dockerhub
`usedVersion` | The currently used version, update this field after you updated your containers
`constraint`  | A version constraint to check by. An example constraint can be found in the `docker-checker.yaml`.

## Found a bug?
If you found a bug feel free to create an issue on Github or on my personal Taiga instance: https://taiga.imanuel.dev/project/docker-version-check/

## License
Like all other projects I create, the Docker checker is distributed under the MIT License.
