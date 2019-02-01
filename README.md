# RollEs

Manage Elasticsearch temporal indices with rollover strategy

## Usage

```
Usage:
  rolles [command]

Available Commands:
  alias       Manage ES indices via their aliases
  help        Help about any command
  template    Manage ES templates

Flags:
  -h, --help   help for rolles

Use "rolles [command] --help" for more information about a command.
```

### rolles alias

```
Manage ES indices via their aliases

Usage:
  rolles alias [command]

Available Commands:
  del         Delete index by alias
  put         Put index by alias

Flags:
  -c, --config string   indices configuration file (default "./index_conf.json")
      --es string       Elasticsearch address (default "http://localhost:9200")
  -h, --help            help for alias
  -n, --name string     alias name (all if not specified)
  -p, --prefix string   alias name prefix (default "default")

Use "rolles alias [command] --help" for more information about a command.
```


### rolles template

```
Manage ES templates

Usage:
  rolles template [command]

Available Commands:
  del         Delete template
  put         Put template

Flags:
      --es string         Elasticsearch address (default "http://localhost:9200")
  -h, --help              help for template
  -n, --name string       template name (all if not specified)
  -p, --prefix string     template name prefix (default "default")
  -d, --temp-dir string   root template directory (default "./templates")

Use "rolles template [command] --help" for more information about a command.
```
