# vc

[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=vinyl-linux_vc&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=vinyl-linux_vc)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=vinyl-linux_vc&metric=security_rating)](https://sonarcloud.io/dashboard?id=vinyl-linux_vc)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=vinyl-linux_vc&metric=sqale_index)](https://sonarcloud.io/dashboard?id=vinyl-linux_vc)

`vc` is a `cloud-init`-alike tool for running a script once, and only once, to bootstrap a server. Unlike `cloud-init`, or similar tools, it is distributed as a statically compiled binary and does sod-all beyond run a script and store the output.

## Running

This tool can either be invoked directly on the CLI as per:

```bash
$ vc run some-job.toml
```

Where `some-job.toml` is a file following the below spec.

It can also be run as a [vinit](https://github.com/vinyl-linux/vinit) job with a job [such as this](https://github.com/jspc/vinit-bootscripts/tree/master/vinit/99-vc).

## Job Spec

The following sample job is a good jumping off point:

```toml
name = "My simple vc script"
description = "Run a single script to say hello <3"

[script]
interpreter = "/bin/sh"
body = '''
#!/usr/bin/env bash
#

set -axe

echo "Hello from vc!"
'''
```

The most important section is the `[script]` section. It has the following keys:

1. `interpreter` - this must be set; this tool expects an interpreter to be set to encourage clear, easily grokkable scripts
1. `body` - a script which can be run via the interpretter.
1. `url` - vc can also download a script to run

Only one of body or url can be included.

Jobs can be validated with:

```bash
$ vc validate some-job.yaml
```

## Other commands

```bash
$ vc

Usage:
  vc [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Run a vc script
  status      Show the status of either the specified script, or all scripts
  validate    Validate a vc config is correct
  version     Return server and client version information

```
