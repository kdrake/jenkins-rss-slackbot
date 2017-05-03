# jenkins-rss-slackbot
Post rss feed into slack

# Usage
1. Clone source
```bash
git clone git@github.com:kdrake/jenkins-rss-slackbot.git
```
2. Put config.json to /data dir
3. Build docker image
```bash
docker build .
```

# Config example:
```bash
{
  "AssemblyURL": "url to jenkins json api",
  "JobPrefix": [
    "JobPrefix1",
    "JobPrefix2"
  ],
  "WebhookURL": "webhook url"
}
```

