# rssreader

Simple RSS reader. Done with ChatGPT 3.5

## Usage

Edit the `reader.conf` file:

```json
{
    "output": [ "file", "stdout" ],
    "feeds": [
        {
            "name": "lorem-rss test feed, updates every minute",
            "url": "http://lorem-rss.herokuapp.com/feed?unit=second&interval=30",
            "notify": "email"
        },
        {
            "name": "databreaches.net",
            "url": "https://www.databreaches.net/feed/",
            "notify": "email"
        },
        {
            "name": "niebiezpiecznik.pl",
            "url": "https://feeds.feedburner.com/niebezpiecznik",
            "notify": "email"
        }
    ]
}
```

Run application
```
go run main.go
```