# User statistics service

Tiny service to fetch aggregated statistics from user accounts

## Handled metrics

- users_count (count): Count of user accounts (total)
- users_active (count): Count of activated accounts
- users_weekly (count): Count of activated accounts with a weekly reminded subscription
- users_newsletter (count): Count of activated accounts with a newsletter subscription
- users_weekday (map): Map of count of activated accounts by weekday assigned for weekly reminder (key of the map=weekday number)
- profiles_count (count): Count of profiles (total)
- profiles_active (count): Count of profiles among active accounts
- profiles_histogram (map): Count of accounts (value) by number of profiles in the account (key), distribution of accounts by profile count.
## Usage

The server is loaded by default on :8080 port (:3252 on docker image)

Endpoints :

### $baseURI/fetch/:instance

Fetch stats for the provided instance name (replace :instance by instance name)

Optional parameters:

- `from` and `until`: timestamp for users activated from resp. until this timestamp (using account.accountConfirmedAt field)

Response: an array of Counter results

Each Counter has 4 possible fields: 

- name: counter name
- value: object value (depend on type see below)
- type: 'count' or 'map'
- error: if this field is provided then an error occurred

value field is 
- a number if type is 'count'
- an object (key value pair) if type is 'map'

Example

```json
[
    {
        "name": "users_count",
        "type": "count",
        "value": 9911
    },
    {
        "name": "users_weekday",
        "type": "map",
        "value": {
                "1": 292,
                "2": 861,
                "3": 934,
                "4": 938,
                "5": 928,
                "6": 620
        }
    }
]
```

## Configuration

Environments:

- AUTH_PASSWORD : a bcrypt hashed password to use for Basic Authentication (no auth if empty)
- Db Connection vars: see in [User Management Service](https://github.com/grippenet/user-management-service/blob/master/build/docker/example/user-management-env.list), accept User Db variables and general db client settings
- PORT: change the listening port
- GIN_MODE: can be configured ('release' will be less verbose)