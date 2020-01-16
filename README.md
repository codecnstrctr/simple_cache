# Simple cache

## A simple cache implementation

### Setting up

1. Download the project
2. Run `go build -o cache main.go`
3. Run `./cache -p port` where port is a port number to listen to
4. Some extra run flags are available: `-t` for dry run and `-r int` for setting up expired keys remove interval

### API

Cache API is available through JSON-RCP 1.0 over TCP

* Set:

Request

``` json
{
    "method": "Command.Set",
    "params":
        [
            {
                "Key" : "string",
                "Data" : any_kind_of_data,
                "TTL" : time_to_live_in_nsec
            }
    ],
    "id": 123
}
```

Response

``` json
{
    "result": "OK",
    "error": null,
    "id": 123
}
```

* Get:

Request

``` json
{
    "method": "Command.Get",
    "params": [
        "key"
    ],
    "id": 123
}
```

Response

``` json
{
    "result":
        {
            "Key": "key",
            "Found":true,
            "Data": data_stored_at_key
        },
    "error": null,
    "id": 123
}
```

* Remove

Request

``` json
{
    "method": "Command.Remove",
    "params": [
        "key"
    ],
    "id": 123
}
```

Response

``` json
{
    "result": true,
    "error": null,
    "id": 123
}
```

* Keys

Request

``` json
{
    "method": "Command.Keys",
    "params": [],
    "id": 123
}
```

Response

``` json
{
    "result": [
        list_of_keys_stored_in_cache
    ],
    "error": null,
    "id": 123
}
```
