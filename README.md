<h1 align="center">
  <br>
  <a href="https://deeplooklabs.com"><img src="static/banner.png" width="100%" alt="Gitsearch"></a>
</h1>

Searching github the way it worked...
## Install

```bash
git clone https://github.com/deeplooklabs/gitsearch.git
cd gitsearch; go install

```

## Config:

Set enviroment GITHUB_TOKEN

```bash
export GITHUB_TOKEN=XXX; gitsearch SEARCH
```

## How to search:

> Search AWS keys on source python:
```bash
gitsearch "\"target\" \"AKIA\" boto language:python"
```

## Todo:

- [ ] Get 30+ Results 