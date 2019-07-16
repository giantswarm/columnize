Columnize
=========

Easy column-formatted output for golang

Note: this fork of https://github.com/ryanuber/columnize adds a cheap implementation to work with coloured text using some ANSII excape codes

[![Build Status](https://travis-ci.org/giantswarm/columnize.svg)](https://travis-ci.org/giantswarm/columnize)

Columnize is a really small Go package that makes building CLI's a little bit
easier. In some CLI designs, you want to output a number similar items in a
human-readable way with nicely aligned columns. However, figuring out how wide
to make each column is a boring problem to solve and eats your valuable time.

Here is an example:

```go
package main

import (
    "fmt"
    "github.com/giantswarm/columnize"
)

func main() {
    output := []string{
        "Name | Gender | Age",
        "Bob | Male | 38",
        "Sally | Female | 9",
    }
    result := columnize.SimpleFormat(output)
    fmt.Println(result)
}
```

As you can see, you just pass in a list of strings. And the result:

```
Name   Gender  Age
Bob    Male    38
Sally  Female  9
```

To right-align a column, edit a config's `ColumnSpec` and use `columnize.Format`
instead of `columnize.SimpleFormat`. Example:

```go
    config := columnize.DefaultConfig()
    config.ColumnSpec = []*columnize.ColumnSpecification{
		&columnize.ColumnSpecification{Alignment: columnize.AlignLeft},
		&columnize.ColumnSpecification{Alignment: columnize.AlignLeft},
		&columnize.ColumnSpecification{Alignment: columnize.AlignRight},
	}
    result := columnize.Format(output, config)
```

The result:

```
Name   Gender  Age
Bob    Male     38
Sally  Female    9
```

Columnize is tolerant of missing or empty fields, or even empty lines, so
passing in extra lines for spacing should show up as you would expect.

Configuration
=============

Columnize is configured using a `Config`, which can be obtained by calling the
`DefaultConfig()` method. You can then tweak the settings in the resulting
`Config`:

```
config := columnize.DefaultConfig()
config.Delim = "|"
config.Glue = "  "
config.Prefix = ""
config.Empty = ""
```

* `Delim` is the string by which columns of **input** are delimited
* `Glue` is the string by which columns of **output** are delimited
* `Prefix` is a string by which each line of **output** is prefixed
* `Empty` is a string used to replace blank values found in output

You can then pass the `Config` in using the `Format` method (signature below) to
have text formatted to your liking.

Usage
=====

```go
SimpleFormat(intput []string) string

Format(input []string, config *Config) string
```
