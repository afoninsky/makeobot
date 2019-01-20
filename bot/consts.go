package main

const deploymentEventTemplate = "*keel.sh: {{ .Name }}*\n{{ .Message }}"

const helpMessage = "" +
	"[ChatOps bot](https://github.com/afoninsky/makeobot) welcomes you. Available commands are\n" +
	"`/ping` - check my liveness\n" +
	"`/release image tag` - deploy new image using keel.sh" +
	""
