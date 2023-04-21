# nightscout-shell

A command-line program to show the latest [Nightscout](https://github.com/nightscout/cgm-remote-monitor) values in your shell prompt.

## Configuration

The program expects a configuration file in `~/.config/nightscout-shell/config.json`. The minumum contents are below.

```
{
	"nightscout_url": "https://nightscout.domain.com",
	"nightscout_token": "token-a1b2c3d4e5f6g7h8"
}
```

## Service

Run the service in the background `./nightscout-shell service` it will read the configuration and save the assembled prompt in a cache file.

## Prompt

Set your prompt like this, with the full path to the `nightscout-shell` binary, this will use the cache file and output it correctly.

```
export PS1="\$(nightscout-shell prompt) \s-\v\$ "
```