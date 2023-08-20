![appserve](https://github.com/donuts-are-good/noteserver/assets/96031819/697ca17c-9e59-4f43-9049-b30d1fe6020a)
![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors)

# noteserver
an http api serving notes. uses sqlite for storage.

## setup
- acquire dependencies: `go get`
- compile: `go build`
## usage

run the server 

```shell
./noteserver
```

listens on :4096.



## logs

logging for noteserver is managed through the system's syslog service. each log entry associated with noteserver is tagged with "notes-app".

## signals

noteserver listens for the sigint (interrupt) signal. when this signal is detected, the server begins its shutdown sequence, ensuring any pending operations are completed and the database connections are closed safely.
## notes

when interfacing with the api, particularly the `/sync` endpoint, authentication is managed via a bearer token provided in the request headers. it's expecting a 64 character string, hopefully a sha256 hash. if there's any issue with your request, such as authentication failures or malformed data, refer to the returned http status codes for insight.
## api
### /sync

**get**
- fetch notes.
- authorization: bearer token header.
- response: array of note objects.
```shell
curl -H "Authorization: Bearer YOURTOKEN" https://note.rest/sync
```
**post**
- add notes.
- authorization: bearer token header.
- payload: note objects in json array.
```shell
curl -X POST -H "Authorization: Bearer YOURTOKEN" -H "Content-Type: application/json" -d '[{"date": "2023-08-18", "title": "Title", "body": "Note body"}]' https://note.rest/sync
```
### /health

**get**
- server status.
- response: 200 ok if operational.
```shell
curl https://note.rest/health
```
## data

sqlite (notes.db) schema:

```sql

CREATE TABLE notes (
    token TEXT,
    id INTEGER PRIMARY KEY,
    date TEXT,
    title TEXT,
    body TEXT
);
```


## license

mit license 2023 donuts-are-good, for more info see license.md
