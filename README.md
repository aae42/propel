# propel

tool for triggering pre-defined actions on a remote system

## wat?

Use propel as part of your remote management system.

It can be particularly useful for ad hoc management tasks and deployment processes.
It exposes webhooks that can be used to run a predetermined command on the local
system.

## quick start

Set up a `propel_config.yaml` that looks like this:

```yaml
---
port: 42424
endpoints:
  f4883ed8-070c-4acc-9e1e-c67c2e3d471d:
    start_in: .
    command: whoami
  881ce173-b02c-4bb5-b1a3-233bae46723a:
    start_in: .
    command: ls -la
```

Provide any endpoints you want (i have randomly generated UUIDs in the example),
along with the start directory and the commands to run.

Run `propel`, it will start an http server on the port, and you can `POST` to
those endpoints to run the commands:

```bash
curl -X POST http://localhost:42424/881ce173-b02c-4bb5-b1a3-233bae46723a
```

results in "ls -la" being run in the working directory `propel` was run from.
