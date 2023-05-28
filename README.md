<div align="center">

This tool is designed to allow downloading of software programs released in private repositories that are not accessible through the brew package manager. It can be useful for distributing or sharing software within internal organizations in a private manner.

</div>

### How to Install

```bash
brew tap ghdwlsgur/privrew
brew install privrew
```

### How to Use

1. It issues a token exclusive to private repositories that have released software

   > Settings > Developer settings > Personal access tokens

2. ```bash
       privrew install [OWNER]/[REPO] -t [REPO_TOKEN]
       privrew install ghdwlsgur/example -t github_pat_1234
   ```
