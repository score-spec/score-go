# JSON schemas for Score files

| version | file |
| --- | --- |
| v1-beta1 | score-v1b1.json |

## Embed schemas into project

Add Score schemas into projects with `git subtree add` command:

```
git subtree add \
  --prefix schemas \
  git@github.com:score-spec/schema.git main \
  --squash
```

> **Note:** To avoid storing the entire history of the sub-project in the main repository, make sure to include `--squash` flag.

## Update schemas from upstream

Get the latest versions of the schemas `git subtree pull` command:

```
git subtree pull \
  --prefix schemas \
  git@github.com:score-spec/schema.git main \
  --squash
```

## Contribute changes to upstream

All changes to `score-spec/schema` should be done via pull requests and comply with the review and sign-off policies.

