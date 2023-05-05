# Solomon RESTful API
A RESTful API for the Solomon project

## Getting Started
### Prerequisites
Before you can build and run the file, you'll need to have the following installed:
- Go
- Mongo DBMS
### Step 1: Clone the repository
Open a terminal and navigate to the directory where you want to clone the repository. Then, run the following command to clone the repository:
```bash
git clone https://github.com/Solom-Onion/API/
```
### Step 2: Edit the config file
In the root directory of the cloned repository, you'll find a file named config.solo. Open this file in a text editor and modify the following variables according to your preferences:
```env
# Mongo URI, always starts with mongodb://HOST:IP
mongo_uri=mongodb://localhost:27017
# Server port, you already know
server_port=8080
# Mongo Database
database=dev
# Mongo Collection
collection=my_collection
# Limit for database results
results_limit=15
```
### Step 3: Build the file
Navigate to the root directory of the cloned repository in the terminal. Then, run the following command to build the file:
```bash
go build
```
### Step 4: Run the file
```
./api
```
_______________

## Examples
### `GET /api/search`

Search stored records, returns a list of results.

#### Parameters

- `query` (required): the search query you want to perform
- `field` (optional): the field in which to perform the search. Valid options are `login`, `password`, `soft`, and `url`. If not provided, the search will be performed on all fields.

#### Response

The response will be a JSON object containing the results of the search. Each result will have the following fields:

- `_id`: the ID of the result in the database
- `database.log_date`: the date the login information was added to the database
- `information.device_country`: the country of the device
- `information.device_image`: base64 encoded image of the devices desktop
- `information.device_type`: the type of device (Unix, Windows, MacOS)
- `information.device_ip`: the IP address of the device
- `device._logins`: an array of login objects containing the login information
- `device._cookies`: an array of cookie objects containing the cookie information
- `device._cards`: an array of credit card objects containing the credit card information
- `cost.credit_cost`: the estimated cost of the login information in USD

#### Example

Request: `GET /api/search?query=johndoe&field=login`

Response:

```json
[
  {
    "_id": "mKeuVVzIToWXHcJDnixZwlswsWvoXHwkuRbWBTITPWKAUagDLcTHCCOllvdb",
    "database.log_date": "Thu Jan 26 2023 03:48:56 GMT+0100 (Central European Standard Time)",
    "information.device_country": "US",
    "information.device_image": "ENCODED_IMAGE",
    "information.device_type": "Windows 8 Pro",
    "information.device_ip": "192.168.1.1",
    "device._logins": [
      {
        "Login": "johndoe@example.com",
        "Password": "password123",
        "Soft": "Firefox",
        "URL": "https://example.com"
      }
    ],
    "device._cookies": [],
    "device._cards": [],
    "cost.credit_cost": 0.02
  }
]
```


### `GET /api/view`

Retrieves a single document from the database by its `_id`.

#### Parameters

- `id` (required): the `_id` of the document to retrieve.

#### Example Request
Request: `GET /api/view?id=6137ee027d501c001651a677`

Response:
```json
{
    "_id": "mKeuVVzIToWXHcJDnixZwlswsWvoXHwkuRbWBTITPWKAUagDLcTHCCOllvdb",
    "database": {
        "log_date": "2021-09-01T10:12:03Z"
    },
    "information": {
        "device_country": "US",
        "device_image": "ENCODED_IMAGE",
        "device_type": "Unix - Linux 6.1.19-1-GENTOO",
        "device_ip": "192.168.1.10"
    },
    "device": {
        "_logins": [
            {
                "Login": "user1",
                "Password": "pass123",
                "Soft": "Firefox",
                "URL": "https://example.com/login"
            },
            {
                "Login": "user2",
                "Password": "pass456",
                "Soft": "Chrome",
                "URL": "https://example.com/login"
            }
        ],
        "_cookies": [
            {
                "Name": "session_id",
                "Value": "123456",
                "Domain": "example.com",
                "Path": "/",
                "Expires": "2021-10-01T10:12:03Z"
            }
        ],
        "_cards": [
            {
                "Number": "4111111111111111",
                "Holder": "John Doe",
                "Expiration": "12/24",
                "CVC": "123"
            }
        ]
    },
    "cost": {
        "credit_cost": 6.75
    }
}
```
