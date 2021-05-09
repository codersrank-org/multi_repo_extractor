# DEPRECATED
This repo is not supported anymore. It will be merged with the repo_info_extracotor.
## What is it?
Wrapper for [repo_info_extractor](https://github.com/codersrank-org/repo_info_extractor) to extract/process multiple repositories at once.
If you have a GitHub account with multiple repos (public or private) you can
use this script to analyze them at once by providing a token.      

## How it works?
- First it is going to initialize [repo_info_extractor](https://github.com/codersrank-org/repo_info_extractor). If it is previously cloned, it will be updated.

- Secondly all of your repos (with given **provider** and **visibility**) is going to be cloned or updated (if it cloned previously). All repos are going to be processed and resulting json file will be uploaded to CodersRank. Resulting file only has metadata and don't have any code from the processed repository.

- Lastly, this program will open your browser with codersrank website to link your repositories with your account.
## Installation
### Requirements
- We are using docker version of [repo_info_extractor](https://github.com/codersrank-org/repo_info_extractor) so you need to have docker installed.
- Git
- Optional (build from source): Golang
### From source
```
$ git clone https://github.com/codersrank-org/multi_repo_extractor
$ cd multi_repo_extractor
$ go build .
```
### Binary
Executables are available here: https://github.com/codersrank-org/multi_repo_extractor/releases

## Usage
### TL;DR
This will extract all the private repos from GitHub that are available by
the given token. 
```
go run . -token="{your_actual_token}" -emails="email1@example.com,email2@example.com"
```
or using the binary
```
./multi_repo_extractor_linux -token="{your_actual_token}" -emails="email1@example.com,email2@example.com"
```
### Other options
If you want to change the default configurations you can do it like this:
```
./multi_repo_extractor_linux -token="{your_actual_token}" -emails="email1@example.com,email2@example.com" -repo_visibility="all" -provider="github.com"
```
#### Available flags 
-  `-emails` string:
        Your emails which are used when making the commits. Provide a comma separeted list for multiple emails (e.g. "one@mail.com,two@email.com")
-  `-provider` string:
        Provider for repos. Only `github.com`, `bitbucket.org` are supported now. (default "github.com")
-  `-repo_visibility` string
        Which repos do you want to get processed? Options: all, public and private. (default "private")
-  `-token` string
        Token for accessing repositories. You can also set this with TOKEN enviroment variable.


There are also two enviroment variables you can use:

- `REPO_EXTRACTOR`
    - If you want to use already downloaded [repo_info_extractor](https://github.com/codersrank-org/repo_info_extractor), provide the local path of the repo with this enviroment variable.

- `TOKEN`
    - If you don't want your token to be printed on the command line (for example if you running this program with a cron job on a remote server), you can set your token as an enviroment variable instead of providing it with a flag.
    - If this is set, program will ignore the token provided with flag.
### GitHub.com
First you have to obtain a GitHub Personal Access Token (PAT).
Navigate to [this url](https://github.com/settings/tokens) and create your token. After clicking on `Generate new token` button, select the required scope (repo) and click on `Generate token` at the bottom of the page.

![repo_scope](https://github.com/peti2001/multi_repo_extractor/blob/master/docs/github-scopes.png?raw=true)
### BitBucket.org
Right now BitBucket Cloud is supported. For authentication your have to use your username
and password. Password must be set via the `-token` flag. Example usage:
```
./multi_repo_extractor_linux -token="password1" -username="username1" -emails="email1@example.com,email2@example.com" -repo_visibility="all" -provider="bitbucket.org"
```
When you create the a new `app password` make sure you select all the necessary scopes.
![repo_scope](https://raw.githubusercontent.com/peti2001/multi_repo_extractor/master/docs/bitbucket-scope.png)
The safest way if you create an `app password` and use it instead of your user's password.
You can create it here: https://bitbucket.org/account/settings/app-passwords/
