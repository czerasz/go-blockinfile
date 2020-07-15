![test](https://github.com/czerasz/go-blockinfile/workflows/test/badge.svg)
![lint](https://github.com/czerasz/go-blockinfile/workflows/lint/badge.svg)

# Blockinfile the Golang version

Makes sure provided content block is inside the specified file.

> **NOTE**
> <br />The logic uses a very primitive approach. Do **NOT** use with **large** files.

## Usage

### Basic Usage

- view initial `config.yml` content:

  ```bash
  $ cat config.yml
  image: nginx:latest
  ```

- run `blockinfile`:

  ```bash
  blockinfile -path config.yml -content 'name: app'
  ```

- examine content:

  ```bash
  $ cat config.yml
  image: nginx:latest
  # BEGIN MANAGED BLOCK
  name: app
  # END MANAGED BLOCK
  ```

- run `blockinfile` again:

  ```bash
  blockinfile -path config.yml -content 'name: server'
  ```

- examine updated content:

  ```bash
  $ cat config.yml
  image: nginx:latest
  # BEGIN MANAGED BLOCK
  name: server
  # END MANAGED BLOCK
  ```

### Use with Pipe

```bash
terraform-docs markdown table . | blockinfile -path README.md -marker '# {{ .Mark }} TERRAFORM DOCS' -content -
```