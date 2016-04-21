metha
=====

Command line OAI-PMH incremental harvester. Data is harvested in chunks.

```
$ metha-sync http://export.arxiv.org/oai2
...
```

All downloaded files are written to

```
$ METHA_DIR=/tmp/harvest metha-sync -dir http://export.arxiv.org/oai2
/tmp/harvest/I29haV9kYyNodHRwOi8vZXhwb3J0LmFyeGl2Lm9yZy9vYWky
```

The default `METHA_DIR` is `$HOME/.metha`.

Harvesting can be CTRL-C'd any time. The data is harvested up to the last full
day, so there is a small latency. The HTTP client is resilient.

Example: If the current date would be *Thu Apr 21 14:28:10 CEST 2016*, the harvester
would request all data since the repositories earliest date and *2016-04-20 23:59:59*.

You can stream records to stdout, too.

```
$ metha-cat http://export.arxiv.org/oai2
```

This will stream all harvested records to stdout. You can emit records based on datestamp as well:

```
$ metha-cat -from 2016-01-01 http://export.arxiv.org/oai2
```

This will only stream records with a datestamp equal or after 2016-01-01.

To just stream all data really fast, use find and zcat on the harvesting dir.

To display basic repository information:

```
$ metha-id http://export.arxiv.org/oai2
```

To list all harvested endpoints:

```
$ metha-ls
```

Installation
------------

Use a [release](https://github.com/miku/metha/releases) or

```
$ go get github.com/miku/metha/cmd/...
```
