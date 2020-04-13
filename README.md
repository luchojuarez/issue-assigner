# issue-assigner
Easy use issue assigner, used to assign PRs in github

## Github token management
### Intro
Github has very low rate limit.
To increase this limit, we can use token for authentication, with this we increase up to 5000 request per hour [rate limit doc](https://developer.github.com/v3/#rate-limiting).
### use
* run app whit parameter `-t somevalidtoken`
* or previously set environment variable `GITHUB_USER_TOKEN` with token value

### Extend? Yes of course
see [tokens](./environment/tokens.go), implements your own token manager.

## Backlog
- [ ] Improve doc 🤔.
- [ ] POST to github API to assign a PR.
- [ ] EventTracer, log all events, generate report. Send email?.
- [ ] Implements special rules to skip a PR. Like "has label" "require millestone".

for bugs🐛👾 or hugs🤗 <lucho.juarez79@gmail.com>
