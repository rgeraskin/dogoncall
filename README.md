# DogOnCall

DogOnCall is a Go application that retrieves an on-call schedule from Datadog and sends notification to Slack about the current on-call engineer.

Slack will receive following variables from app:

- `engineer_email`: Email of the current on-call engineer.
- `schedule_name`: Schedule name.
- `schedule_link`: Link to datadog schedule page.

Engineer email in Datadog and Slack should be equal.

## Usage

### Slack workflow

In Slack app:

1. `More` button on the left screen side => `Automations` => `Workflows` => `New workflow` => `Build workflow`
1. For `Start the workflow…` choose event `From a webhook`
1. Set up variables:
   1. `engineer_email` as `Slack user email`
   1. `schedule_name` and `schedule_link` as `Text`.
1. Use `Web request URL` as `SLACK_ENDPOINT` later.
1. For `Then, do these things` choose `Send a message to...` and compose a message using variables.

   > See [this](https://stackoverflow.com/a/78514785) if you want to use `schedule_name` and `schedule_link` together as link.
1. Publish the workflow

### App configuration

The application requires the following environment variables:

- `DD_API_KEY`: Your Datadog API key.
- `DD_APP_KEY`: Your Datadog application key.
- `SCHEDULE_NAME`: The name of the Datadog on-call schedule.
- `SLACK_ENDPOINT`: The Slack workflow webhook endpoint to send notifications.

### App running

It's reasonable to run this app by a cron. For k8s I use CronJob in the Helm chart.

#### Using Helm

```sh
# put required values in helm/values-private.yaml
helm upgrade --install dogoncall oci://registry-1.docker.io/rgeraskin/helm-dogoncall -f helm/values-private.yaml
```

or from the cloned git repo:

```sh
# put required values in helm/values-private.yaml
helm upgrade --install dogoncall ./helm -f helm/values-private.yaml
```

#### Using Docker

```sh
# put required envs in .env
docker run --name dogoncall --env-file .env --rm rgeraskin/dogoncall:latest
```

or from the cloned git repo:

```sh
# put required envs in .env
docker buildx build . -t dogoncall
docker run --name dogoncall --env-file .env --rm dogoncall
```

#### From Source

```sh
go install github.com/rgeraskin/dogoncall/dogoncall@latest
# export required env vars
dogoncall
```

or from the cloned repo:

```sh
go build -C dogoncall -o ..
# export required env vars
./dogoncall
```

## Local development

### Prepare local env

1. Install [mise](https://mise.jdx.dev): `brew install mise`
2. [Activate mise](https://mise.jdx.dev/getting-started.html#_2a-activate-mise)

   Zsh example:

   ```bash
   echo 'eval "$(~/.local/bin/mise activate zsh)"' >> ~/.zshrc
   # restart shell or `source ~/.zshrc`
   ```

3. `mise install` in the repo root dir to install all required tools

### Development hints

- Use [Tilt](https://tilt.dev) for a realtime k8s experience: `tilt up`
- Use mise tasks for routines: `mise run build`, `mise run deploy`

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
