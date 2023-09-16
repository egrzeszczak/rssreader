# rssreader

Simple RSS reader as a learning Go project. Done with ChatGPT 3.5

**Note:** This just reads feeds for changes. I'm going to add e-mail notifications in the future

## Usage

Edit the `reader.conf` file:

```json
{
    "interval": 15, # in seconds
    "output": [ "file", "stdout" ], 
    "feeds": [
        {
            "name": "lorem-rss test feed, updates every minute",
            "url": "http://lorem-rss.herokuapp.com/feed?unit=second&interval=30",
            "notify": ["email"]
        },
        {
            "name": "databreaches.net",
            "url": "https://www.databreaches.net/feed/",
            "notify": ["email", "telegram"]
        },
        {
            "name": "niebiezpiecznik.pl",
            "url": "https://feeds.feedburner.com/niebezpiecznik",
            "notify": ["syslog"]
        }
    ]
}
```

Run application
```
go run main.go
```