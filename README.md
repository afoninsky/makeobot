# makeobot
```
TAG=debug skaffold build
kubectl -n default apply -f deployment.yml
```

## Roadmap

### release pipeline
- [x] hook for semantic commits [kinda](https://github.com/fteem/git-semantic-commits)
- [ ] CI pipeline: release on master [kinda](https://github.com/go-semantic-release/semantic-release)
- [ ] CI pipeline: build docker image on release [kinda](https://goreleaser.com/docker/)

### keel
- [x] notification about deployment
- [x] trigger new deployment manually
- [ ] notification about new approval
- [ ] list of pending approvals
- [ ] approve deployment

###  concourse
- [ ] notification about job result [kinda](https://github.com/mdb/concourse-webhook-resource)
- [ ] trigger resource check [kinda](https://concourse-ci.org/resources.html#resource-webhook-token)

### alertmanager
- [ ] information about alerts [kinda](https://prometheus.io/docs/alerting/configuration/#webhook_config)