```bash
docker run -t --rm --privileged -d moby/buildkit:master

# vigorous_mcclintock is container name
docker exec -it vigorous_mcclintock sh

buildctl-daemonless.sh build --frontend dockerfile.v0 --local context=/app --local dockerfile=/app --output type=image,name=docker.io/pineapple217/cicd-test:latest,push=true
```

```json
# /root/.docker/config.json
{
  "auths": {
    "https://index.docker.io/v1/": {
      "auth": "REDACTED"
    }
  }
}
```

```bash
#auth
echo -n "$DOCKERHUB_USERNAME:$DOCKERHUB_TOKEN" | base64
```

later gaan we de docker config gewoon binden op de host configs
