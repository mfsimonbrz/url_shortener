# url_shortener

This is a toy project for an url shortener. It stores the url entries on a Postgres database and uses Redis as cache database.

To add an url entry:

`curl -XPOST localhost:5000 -d '{"url" : "<your url>"}'`

To retrieve the url:

`curl localhost:5000/<short-url>`