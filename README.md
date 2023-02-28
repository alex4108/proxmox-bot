# Proxmox Discord Bot

[![Tests](https://github.com/alex4108/proxmox-bot/actions/workflows/test.yml/badge.svg)](https://github.com/alex4108/proxmox-bot/actions/workflows/test.yml)
[![Release](https://github.com/alex4108/proxmox-bot/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/alex4108/proxmox-bot/actions/workflows/release.yml)
[![GitHub forks](https://img.shields.io/github/forks/alex4108/proxmox-bot)](https://github.com/alex4108/proxmox-bot/network)
[![GitHub stars](https://img.shields.io/github/stars/alex4108/proxmox-bot)](https://github.com/alex4108/proxmox-bot/stargazers)
![GitHub contributors](https://img.shields.io/github/contributors/alex4108/proxmox-bot)
[![GitHub license](https://img.shields.io/github/license/alex4108/proxmox-bot)](https://github.com/alex4108/proxmox-bot/blob/main/LICENSE)
![GitHub All Releases](https://img.shields.io/github/downloads/alex4108/proxmox-bot/total)
![Docker Pulls](https://img.shields.io/docker/pulls/alex4108/proxmox-bot)
[![Discord](https://img.shields.io/discord/742969076623605830)](https://discord.gg/FpDjFEQ)

![Supports amd64](https://img.shields.io/badge/arch-amd64-brightgreen)

[![Discord Support](https://user-images.githubusercontent.com/7796475/89976812-2628c080-dc2f-11ea-92a1-fe87b6a9cf92.jpg)](https://discord.gg/FpDjFEQ)

## Purpose

Enables users in Discord to invoke Start or Stop actions on Proxmox VMs

Did I save you some time?  [Buy me a :coffee::smile:](https://venmo.com/alex-schittko)

## Configure the bot

### Env Vars

It's unsafe to store credentials in files, so we require the credentials to be injected via env variables.

For every `proxmox_host_url`, we assume there's two env vars like `HOSTPORT_USERNAME` and `HOSTPORT_PASSWORD`, where `HOST` and `PORT` are the portions defined in the below YAML.

### Config File

Modify `proxmox-bot-config.yml` as needed.

Each entry in the list has 3 properties

* logical_name: The name users will call the VM when they invoke `startvm` or `stopvm`
* proxmox_vm_id: The numberic VM ID of the target.
* proxmox_host_url: the HOST:PORT portion of a URL to connect to proxmox.

## Running the bot

### Environment Variables

* PROXMOX_BOT_DISCORD_TOKEN

### VSCode

Included is a VSCode devcontainer to help you get off the ground quickly.

Debugger is enabled in devcontainer, provided env variables below are set.

### Kubernetes

Take a look at kube-manifest.yml for an idea of what you need to deploy.

Ensure you have secrets & configmaps set.

Example Configmap:

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: proxmox-bot-live-config
data:
  config: |
    vms:
      - logical_name: A
        proxmox_vm_id: B
        proxmox_host_url: C
        proxmox_host_name: D
```

### Discord Bot Setup

* [This guide](https://www.writebots.com/discord-bot-token/) seems to have a good write up on how to generate a bot token.
* Note that during the creation of the bot, you will need to enable the "Server Members Intent" flag on the Bot page in the Discord developers portal.
* Once you have the token in step 5, replace "9999" in the `docker-compose.yml` file with your bot's token.
* Finally, craft your authorization URL.  You can copy the authorization URL from the Discord developers portal as mentioned in step 5.  
* Once the authorization URL is copied, replace the permissions integer with that from the URL given above to join the public bot to your server.
* You should now be able to visit your authorization URL and join your own bot to your Discord guild.

## Using the bot

Once the bot is joined, in any channel it has read access to, send a message `$help`.  This command will enumerate all the commands the bot currently supports.

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the project
1. Create your feature branch (`git checkout -b feature/AmazingFeature`)
1. Make changes, and update `CHANGELOG.md` to describe them.
1. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
1. Push to the branch (`git push origin feature/AmazingFeature`)
1. [Open a pull request](https://github.com/alex4108/proxmox-bot/compare)
