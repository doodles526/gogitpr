[![BuildStatus](https://travis-ci.org/doodles526/gogitpr.svg?branch=master)](https://travis-ci.org/doodles526/gogitpr)

# gogitpr

gogitpr is a set of libraries for fetching Github API data and an abstraction
for a datastore for those data.

## Installation

```
go get -u github.com/doodles526/gogitpr
make setup
make
```

## Usage

Currently, all configuration is set via envvars

### GITPR_BASE_URL

Configures the base url of the github API. Default: `https://api.github.com`

### GITPR_GITHUB_TOKEN

Sets the Oauth token for pinging the github API. Default: blank

### GITPR_APPLICATION_NAME

Sets the Application Name to report to the github API via the `User-Agent`
header. Default: `gogitpr`

### GITPR_GITUHB_ORG

Sets the default github organization to use in fetching pull requests Default:
blank

Either the github org or the github user must be set

### GITPR_GITHUB_USER

Sets the default github user to use in fetching pull requests Default: blank

Either the github org or the github user must be set

### GITPR_PRINT

Should we print the end result from `main`
