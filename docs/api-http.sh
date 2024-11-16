# List Toggles
## request
curl --location 'http://localhost:9002/toggles?offset=3&limit=3' \
--header 'Content-Type: application/json'

## case ok
HTTP-Status: 200
{
    "message": "Request processed successfully!",
    "data": {
        "results": [
            {
                "id": "generic/enable-foo-8",
                "status": false,
                "updatedAt": "2024-11-14T16:08:04.288654Z",
                "lastAccessedAt": null,
                "accessFreqWeekly": 0
            },
            {
                "id": "generic/enable-foo-9",
                "status": false,
                "updatedAt": "2024-11-14T16:08:06.659697Z",
                "lastAccessedAt": null,
                "accessFreqWeekly": 0
            }
        ],
        "total": 5
    }
}

## case empty
HTTP-Status: 200
{
    "message": "Request processed successfully!",
    "data": {
        "results": [],
        "total": 1
    }
}

