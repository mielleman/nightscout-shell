# nightscout-shell
[![Test](https://github.com/mielleman/nightscout-shell/actions/workflows/test.yml/badge.svg)](https://github.com/mielleman/nightscout-shell/actions/workflows/test.yml)

A command-line program to show the latest [Nightscout](https://github.com/nightscout/cgm-remote-monitor) values in your shell prompt.

## Configuration

The program expects a _valid_ JSON configuration file in `~/.config/nightscout-shell/config.json`, you can change the location of the file by using the `--config` argument. The minimum required contents of the configuration file are below.

```
{
	"nightscout_url": "https://nightscout.domain.com",
	"nightscout_token": "token-a1b2c3d4e5f6g7h8"
}
```

The following values can be used

| Setting          | Required | Type    | Description                                                                  |
| ---------------- | -------- | ------- | ---------------------------------------------------------------------------- |
| nightscout_url   | Yes      | String  | The URL to your Nighscout instance                                           |
| nightscout_token | Yes      | String  | The token to be used to connect to your Nighscout instance                   |
| cache_file       | No       | String  | The location for the cache file                                              |
| service_interval | No       | Integer | What time interval (in minutes) should the service retrieve the latest value |

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

### Usage 

You can use the latest value in your `.profile` (for bash) and set your prompt `PS1` variable. Make sure you use the full path to the `nightscout-shell` binary and you escape the `\$` so it is evaluated during each time the prompt is printed.

In the following example we take the default Ubuntu `$PS1` value of `\u@\h:\w\$` and add the nightscout prompt value to the beginning with `\$(${HOME}/bin/nightscout-shell prompt)`. This adds the `${HOME}/bin/nightscout-shell` command with the `prompt` sub-command to the beginnen of `PS1`.

```
export PS1="\$(${HOME}/bin/nightscout-shell prompt) \u@\h:\w\$ "
```

## Wishlist

- Support more shell prompts
- Make unicode characters and colors optional