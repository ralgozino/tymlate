# tymlate

Folder structure aware template engine, that mirrors the source folder structure to target, while processing go template files.

You give me a structure source folder, that looks like a desired target, and I will:
1. recreate the folder structure on target directory
1. copy the files if it is not a template
1. process template files to target

It is based on go template with addition of Sprig template functions (http://masterminds.github.io/sprig/).

## Install 

Under release (https://github.com/ralgozino/tymlate/releases) pick the latest archive of apropritate flavor 
untar and copy `tymlate` to a folder that is on your `PATH` (on Linux e.g. `/usr/local/bin`). 

Alternatively you can use go:
```
go get github.com/ralgozino/tymlate
```

## Usage

```text
Usage:
  tymlate [flags]

Flags:
  -c, --configuration string   path to configuration file
  -d, --datasource strings     Datasource in name=file format
  -h, --help                   help for tymlate
  -s, --source string          path to template source folder
  -t, --target string          path to target folder

```

In order to use `tymlate` you need:
- folder containing source tree with files to be templated
    - from `v1.1.0` even file or folder names can be templated
- target folder where the files will be generated
- configuration YAML inside of source folder with name 
`.tymlate.yml` or provided with `-c` flag
- optionally data YAML files 

## Configuration file

Here is an example configuration file:
```yaml
data: #direct data
  meta:
    package: "github.com/ralgozino/tymlate"
    name: MyEzample
  Dat1: # the name of the map
    key1: "value1"
    key2: 22
    key3: true
  Dat2:
    key11: "value eleven"
    key2: 23
    key3: false
include:
  cfgName1: subcfg1.yml
  cfgNameTwo: subconf/subcfg2.yml
templates:
  includes:
  - ".*\\.tpl"
  - ".*\\.go"
  - LICENSE
  excludes:
  - ".*\\.inv"
  suffix: ".tpl"
  processFilename: true
```
Where:

1. `data` (optional): contains the variables to be used in context,
from the example you would access {{.meta.name}} from your template
1. `include` (optional): include additional yaml data files,
the key will be the key from the yaml file (e.g. `cfgName1`)
1. `templates` (required)
    1. `includes` - files to include (regexp)
    1. `excludes` - files to ignore (regexp)
    1. `suffix` (required): the template file suffix
    1. `processFilename` (optional): true if filenames and 
    folder names should be processed as go templates (by default false)

## Examples included

See: [generator/testdata](https://github.com/lorands/tymlate/tree/master/generator/testdata)

Under the [`source`](https://github.com/lorands/tymlate/tree/master/generator/testdata/source) you will find a source we use to test,
and under [`target`](https://github.com/lorands/tymlate/tree/master/generator/testdata/target) you can see the desired output.
The configuration is provided in [`conf.yml`](https://github.com/lorands/tymlate/tree/master/generator/testdata/conf.yml)

> Note that the testdata folder has been cancelled in this fork because it was breaking the `go install` command.
> Please refer to upstream for the testdata.

## Step by step usage

1. Create a template folder structure. 
1. Create a configuration that defines how to process
 the structure.
1. Define folder structure
1. Define context data file(s)
1. Run tymlate with configuration

Let's assume you would like to generate this structure 
each time you have a new project:

```
.
|___ cmd /
|  |___ main.go
|___ model /
|  |___ model.go
|___ README.md
|___ LICENSE
|___ Dockerfile
|___ go.mod

```

### Template folder structure

Just create a structure that will be mirrored to target.

Decide on what will be the template file extension, 
we use `.tpl` in this example (see: `tymlate/generator/testdata/source`).

```
.
|___ cmd /
|  |___ main.go.tpl
|___ model /
|  |___ model.go.tpl
|___ README.md.tpl
|___ LICENSE
|___ Dockerfile.tpl
|___ go.mod.tpl
```

Each template file (.tpl in this case), is a go template file.

## It's just a go template

Each template file is "just" a go template.

## Template functions included

We included Sprig (http://masterminds.github.io/sprig/)
functions for your convenience. 

## Command line (cli) examples

The simplest case:
```shell script
tymlate -s path/to/source -t path/to/target
```
In this case the configuration should be provided
in `.tymlate.yml` file inside the source.

Simple case, with an external config: 
```shell script
tymlate -c path/to/config.yml -s path/to/source -t path/to/target
```

Additional data sources
```shell script
tymlate -c path/to/config.yml -d ds=path/to/data.yml -s path/to/source -t path/to/target
```

More data sources
```shell script
tymlate -c path/to/config.yml -d ds1=path/to/data1.yml -d ds2=path/to/data2.yml -s path/to/source -t path/to/target
```

