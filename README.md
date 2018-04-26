# registry
A Registry that stores snippets of any data format

## Why
Re-use, developers copy and paste a lot. We (developers) provide tooling and example files which are used as the base but never complete, these files require adapting. These files are stored somewhere and we have to read and extract what we need from multiple sources, how annoying...

## Motivation
When working on multiple projects within the same company I found myself copy and pasting a lot of stages from Jenkinsfiles, either keeping as is or adapting only slightly. These stages were used throughout the company, teams were using the same approach. A Jenkins shared library was created to promote reuse but didn't really solve the problem. We wrapped common stages and provide a single function for this, downside coupling..
I wanted to provide a single system which could store these useful stages.

## Terminology

### Document
A **Document** is a unique entity (name and version) which contains some re-usable/useful information (**Snippet**). Documents can also have dependencies which are also Documents.

### Snippet
A **Snippet** is just a string. Normally Snippets are useful and unique but not guarenteed.

### Usage
#### Prerequistities
* postgres database
```bash
DB_USER=<DB_USER> DB_PASS=<DB_PASS> MIGRATION_FILE_LOCATION=${PWD}/db/migration ./registry # Will start the application with defaults
# For help use the -h flag
```

### Development
#### Prerequistities
* Go 1.9.x
* [glide](https://github.com/Masterminds/glide)
* Make

```bash
go get -d github.com/haggis-io/registry
cd $GOPATH/src/github.com/haggis-io/registry
make build
```

