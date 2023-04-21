# nightscout-shell

A command-line program to show the latest [Nightscout](https://github.com/nightscout/cgm-remote-monitor) values in your shell prompt.

## Configuration

The program expects a _valid_ JSON configuration file in `~/.config/nightscout-shell/config.json`, you can change the location of the file by using the `--config` argument. The minimum required contents of the configuration file are below.

```
{
	"nightscout_url": "https://nightscout.domain.com",
	"nightscout_token": "token-a1b2c3d4e5f6g7h8"
}
```

### Nigthscout token

Create a new subject, give this subject the following roles: `readable, status-only`. The token (e.g. `token-a1b2c3d4e5f6g7h8`) is then saved to the configuration file. The roles it uses are needed for the `api:status:read` permission to read the status and configuration and the `api:entries:read` permission to read the latest value.

## Service

Usage:
`nightscout-shell service --help`

The service sub-command allows the program to run as a service in the background to request the latest value from your Nightscout instance and update the cache file. You can start the service with `./nightscout-shell service` it will read the configuration and save the assembled prompt value to a cache file. This service keeps requesting the latest value from Nightscout instance at a configured interval (default 5 minutes).

If anything goes wrong at any point the program will stop and exit with a non-zero return code.

## Prompt

Usage:
`nightscout-shell prompt --help`

The prompt sub-command allows the program to request the latest value from the cache file and return it formatted according to the configuration. It is basically a `cat` of the cache file, it currently is already a sub-command so more functionally can be added in the future.

You can use it in your `.profile` (for bash) and set your prompt `PS1` variable, with the full path to the `nightscout-shell` binary.

```
export PS1="\$(nightscout-shell prompt) \s-\v\$ "
```

## Wishlist

- Support more shell prompts
- Make unicode characters and colors optional