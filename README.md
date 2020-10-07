# What is it?
Wrapper for [repo_info_extractor](https://github.com/codersrank-org/repo_info_extractor) to extract/process multiple repositories at once.

# How to use?
You can use pre-built binaries or you can run `go run .` in the folder after cloning this repository.

We are using docker version of [repo_info_extractor](https://github.com/codersrank-org/repo_info_extractor) so you need to have docker installed.

While using the program, you need to provide some variables as flags:

- provider (Where should we get the repositories?)
    - Example: *github.com* 
    - Currently only github.com is supported and it is the default value.

- repoVisibility (Which repositories should be processed?)
    - Values: *all*, *public* and *private*
    - Default value is private.

- token (How we can access to the repositories?)
    - You can create your personal access token from [here](https://github.com/settings/tokens).
    - **Repo** scope is needed.

- emails (Which commits are yours?)
    - You can provide multiple emails as a comma separated list (e.g. "one@mail.com, two@mail.com").
    - At least one email address is needed.

There are also two enviroment variables you can use:

- REPO_EXTRACTOR
    - If you want to use already downloaded [repo_info_extractor](https://github.com/codersrank-org/repo_info_extractor), provide the local path of the repo with this enviroment variable.

- TOKEN
    - If you don't want your token to be printed on the command line (for example if you running this program with a cron job on a remote server), you can set your token as an enviroment variable instead of providing it with a flag.
    - If this is set, program will ignore the token provided with flag.

# How it works?
- First it is going to initialize [repo_info_extractor](https://github.com/codersrank-org/repo_info_extractor). If it is previously cloned, it will be updated.

- Secondly all of your repos (with given **provider** and **visibility**) is going to be cloned or updated (if it cloned previously). All repos are going to be processed and resulting json file will be uploaded to CodersRank. Resulting file only has metadata and don't have any code from the processed repository.

- Lastly, this program will open your browser with codersrank website to link your repositories with your account.