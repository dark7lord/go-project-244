# Gendiff
Compares two configuration files (JSON/YAML) and shows a difference (stylish, plain, json).

### Hexlet tests and linter status:
[![Actions Status](https://github.com/dark7lord/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/dark7lord/go-project-244/actions)
[![test-ci](https://github.com/dark7lord/go-project-244/actions/workflows/test-ci.yml/badge.svg)](https://github.com/dark7lord/go-project-244/actions/workflows/test-ci.yml)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=dark7lord_go-project-244&metric=coverage)](https://sonarcloud.io/summary/new_code?id=dark7lord_go-project-244)
<!-- [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=dark7lord_go-project-244&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=dark7lord_go-project-244) -->

## Demo

[![asciicast](https://asciinema.org/a/lkXCgK9ft9KU2LPd.svg)](https://asciinema.org/a/lkXCgK9ft9KU2LPd)

## Install / Uninstall

### Install

```bash
make install
# or
go install ./cmd/gendiff/
```

After installation, make sure `$(go env GOPATH)/bin` is in your `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Build only (binary stays in ./bin)

```bash
make build
go build -o gendiff ./cmd/gendiff/
```

### Uninstall

```bash
make uninstall
# or
go clean -i ./cmd/gendiff/
```

## Usage

```
gendiff [global options] <old_file> <new_file>
```

The `--format` flag supports three output styles:

| Flag       | Description                                      |
|------------|--------------------------------------------------|
| `stylish`  | (default) tree view with `+`/`-` markers         |
| `plain`    | line-based text output                           |
| `json`     | JSON representation of the diff                  |

Quick examples:

```bash
gendiff fileA.json fileB.json
gendiff --format plain fileA.json fileB.json
gendiff --format json fileA.json fileB.json
```

## Examples

### Flat files

**fileA.json**
```json
{
  "host": "hexlet.io",
  "timeout": 50,
  "proxy": "123.234.53.22",
  "follow": false
}
```

**fileB.json**
```json
{
  "timeout": 20,
  "verbose": true,
  "host": "hexlet.io"
}
```

#### stylish (default)

```bash
gendiff fileA.json fileB.json
```

```diff
{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}
```

[![asciicast](https://asciinema.org/a/wje0NF90xqVkvUol.svg)](https://asciinema.org/a/wje0NF90xqVkvUol)

#### plain

```bash
gendiff --format plain fileA.json fileB.json
```

```
Property 'follow' was removed
Property 'proxy' was removed
Property 'timeout' was updated. From 50 to 20
Property 'verbose' was added with value: true
```

[![asciicast](https://asciinema.org/a/xYaJQ16N83TavNMU.svg)](https://asciinema.org/a/xYaJQ16N83TavNMU)

#### json

```bash
gendiff --format json fileA.json fileB.json
```

```json
[
  {
    "key": "follow",
    "type": "removed",
    "value": false
  },
  {
    "key": "host",
    "type": "unchanged",
    "value": "hexlet.io"
  },
  {
    "key": "proxy",
    "type": "removed",
    "value": "123.234.53.22"
  },
  {
    "key": "timeout",
    "type": "changed",
    "oldValue": 50,
    "newValue": 20
  },
  {
    "key": "verbose",
    "type": "added",
    "value": true
  }
]
```

[![asciicast](https://asciinema.org/a/wTgFtwRsjKA5wPbE.svg)](https://asciinema.org/a/wTgFtwRsjKA5wPbE)

#### YAML

Same output as JSON stylish, with `.yml` input files:

```bash
gendiff fileA.yml fileB.yml
```

[![asciicast](https://asciinema.org/a/KwRoX2hqVMcU3boD.svg)](https://asciinema.org/a/KwRoX2hqVMcU3boD)

---

### Nested files

**fileE.json**
```json
{
  "common": {
    "setting1": "Value 1",
    "setting2": 200,
    "setting3": true,
    "setting6": {
      "key": "value",
      "doge": {
        "wow": ""
      }
    }
  },
  "group1": {
    "baz": "bas",
    "foo": "bar",
    "nest": {
      "key": "value"
    }
  },
  "group2": {
    "abc": 12345,
    "deep": {
      "id": 45
    }
  }
}
```

**fileF.json**
```json
{
  "common": {
    "follow": false,
    "setting1": "Value 1",
    "setting3": null,
    "setting4": "blah blah",
    "setting5": {
      "key5": "value5"
    },
    "setting6": {
      "key": "value",
      "doge": {
        "wow": "so much"
      },
      "ops": "vops"
    }
  },
  "group1": {
    "baz": "bars",
    "foo": "bar",
    "nest": "str"
  },
  "group3": {
    "deep": {
      "id": {
        "number": 45
      }
    },
    "fee": 100500
  }
}
```

#### stylish

```bash
gendiff fileE.json fileF.json
```

```diff
{
    common: {
      + follow: false
        setting1: Value 1
      - setting2: 200
      - setting3: true
      + setting3: null
      + setting4: blah blah
      + setting5: {
            key5: value5
        }
        setting6: {
            doge: {
              - wow: 
              + wow: so much
            }
            key: value
          + ops: vops
        }
    }
    group1: {
      - baz: bas
      + baz: bars
        foo: bar
      - nest: { key: value }
      + nest: str
    }
  - group2: { abc: 12345, deep: { id: 45 } }
  + group3: { deep: { id: { number: 45 } }, fee: 100500 }
}
```

[![asciicast](https://asciinema.org/a/JDROAWxn7BHD8kRW.svg)](https://asciinema.org/a/JDROAWxn7BHD8kRW)

#### plain

```bash
gendiff --format plain fileE.json fileF.json
```

```
Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'common.setting4' was added with value: 'blah blah'
Property 'common.setting5' was added with value: [complex value]
Property 'common.setting6.doge.wow' was updated. From '' to 'so much'
Property 'common.setting6.ops' was added with value: 'vops'
Property 'group1.baz' was updated. From 'bas' to 'bars'
Property 'group1.nest' was updated. From [complex value] to 'str'
Property 'group2' was removed
Property 'group3' was added with value: [complex value]
```

[![asciicast](https://asciinema.org/a/sd1nGNVbRZUp7jEI.svg)](https://asciinema.org/a/sd1nGNVbRZUp7jEI)

#### json

```bash
gendiff --format json fileE.json fileF.json
```

```json
[
  {
    "key": "common.follow",
    "type": "added",
    "value": false
  },
  {
    "key": "common.setting1",
    "type": "unchanged",
    "value": "Value 1"
  },
  {
    "key": "common.setting2",
    "type": "removed",
    "value": 200
  },
  {
    "key": "common.setting3",
    "type": "changed",
    "oldValue": true
  },
  {
    "key": "common.setting4",
    "type": "added",
    "value": "blah blah"
  },
  {
    "key": "common.setting5",
    "type": "added",
    "value": {
      "key5": "value5"
    }
  },
  {
    "key": "common.setting6.doge.wow",
    "type": "changed",
    "oldValue": "",
    "newValue": "so much"
  },
  {
    "key": "common.setting6.key",
    "type": "unchanged",
    "value": "value"
  },
  {
    "key": "common.setting6.ops",
    "type": "added",
    "value": "vops"
  },
  {
    "key": "group1.baz",
    "type": "changed",
    "oldValue": "bas",
    "newValue": "bars"
  },
  {
    "key": "group1.foo",
    "type": "unchanged",
    "value": "bar"
  },
  {
    "key": "group1.nest",
    "type": "changed",
    "oldValue": {
      "key": "value"
    },
    "newValue": "str"
  },
  {
    "key": "group2",
    "type": "removed",
    "value": {
      "abc": 12345,
      "deep": {
        "id": 45
      }
    }
  },
  {
    "key": "group3",
    "type": "added",
    "value": {
      "deep": {
        "id": {
          "number": 45
        }
      },
      "fee": 100500
    }
  }
]
```

[![asciicast](https://asciinema.org/a/S2G5LZfO4csOrVh6.svg)](https://asciinema.org/a/S2G5LZfO4csOrVh6)

## Development

```bash
make test       # run all tests
make lint       # run linters
make cover      # run tests with coverage report
make build      # build binary to bin/gendiff
make clean      # remove build artifacts
```
