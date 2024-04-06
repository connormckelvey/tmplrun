# tmplrun

tmplrun is an open-source templating engine written in Go, designed to provide flexibility and simplicity in generating dynamic content for various applications. Tmplrun stands out by allowing the use of JavaScript within templates, thus enabling users to leverage the power and familiarity of JavaScript for templating purposes.

## Features

- **JavaScript Templating**: Utilize JavaScript directly within templates, eliminating the need for domain-specific template languages.
- **Extensibility**: Easily extend tmplrun to support additional languages for templating purposes.
- **Command Line Interface**: Integrated command-line tool for convenient template rendering.

## Installation

### CLI

```bash
go install github.com/connormckelvey/tmplrun/cmd/...
```

### Library

```bash
go get github.com/connormckelvey/tmplrun
```

## Usage

To render a template using tmplrun, you can use the command line interface:

```bash
tmplrun render -p path/to/props path/to/template
```

