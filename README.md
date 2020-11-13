# aaa [![GoDoc](https://godoc.org/github.com/axkit/aaa?status.svg)](https://godoc.org/github.com/axkit/aaa) [![Build Status](https://travis-ci.org/axkit/aaa.svg?branch=master)](https://travis-ci.org/axkit/aaa) [![Coverage Status](https://coveralls.io/repos/github/aaa/gonfig/badge.svg)](https://coveralls.io/github/axkit/aaa) [![Go Report Card](https://goreportcard.com/badge/github.com/axkit/aaa)](https://goreportcard.com/report/github.com/axkit/aaa)

AAA - Authentication, Authorization & Accounting

# Motivation

The package provides AAA features for http router [github.com/axkit/vatel](https://github.com/axkit/vatel) using JWT.
AAA plays a proxy between users/roles/permission storage and vatel.

# Concepts

- JTW is used.
- AAA is independent of users, roles and permissions storage structure.
- AAA does not know password encryption approach.
- Developer can extend token payload.
-

# Endpoints

### Sign in

POST /auth/sign-in
Input

```
    {
        "login" : "user1"
        "password" : "plain-or-encrypted-password",
    }
```

Output
Successfull: HTTP 200

```
{
    "access_token" : "abc.."
    "refresh_token": "zyx..."
    "permissions" : ["PermissionCode1","PermissionCode2", "PermissionCode3"...]
}
```

Access token holds following user specified payload inside:

```
{
  "user": 42,
  "login": "user1",
  "role": 1,
  "perms": "ZTY1ZmZmN2YyZmVlYzNlZmJjN2RmZmJmZGNmM2Y3ZmYzZjlmZmRmZmZmN2Y3NWJkMDE="
}
```

Refresh token holds user specified payload inside:

```
{
  "user": 42
}
```

Failed: HTTP 401

```
{
    "message" : "invalid cridentials"
}
```

### Access token validation

POST /auth/is-token-valid
Output
Successfull: HTTP 200

```
{
   "result" : "ok"
}
```

Failed: HTTP 401

```
{
    "message" : "invalid token"
}
```

### Refresh token

POST /auth/refresh-token
Input

```
    {
        "refresh_token" : "xyz.."
    }
```

Output
Successfull: HTTP 200

```
{
    "access_token" : "abc.."
    "refresh_token": "zyx..."
    "permissions" : ["PermissionCode1","PermissionCode2", "PermissionCode3"...]
}
```

Failed: HTTP 401

```
{
    "message" : "invalid token"
}
```

- Application functionality can be limited by using permissions.
- Permission (access right) represented by unique string code.
- Application can have many permissions.
- A user has a role.
- A role is set of allowed permission, it's subset of all permissions supported by application.
- As a result of succesfull sign in backend provides access and resresh tokens.
- Payload of access token shall have list of allowed permissions.
- A single permission code looks like "Customers.Create", "Customer.AttachDocuments", "Customer.Edit", etc.
- Store allowed permission codes could increase token size.
- Bitset comes here.
- Every permission shall be accociated with a single bit in the set.
- Bitset adds to the token as hexadecimal string.

## Usage Examples

Sign In

```
    var perms bitset.Bitset
    perms.Set(1)                    // 0000_0010
    perms.Set(2)                    // 0000_0110
    perms.Set(8, 10)                // 0000_0110 0000_0101
    tokenPerms := perms.String()    // returns "0605"
```

Check allowed permission in auth middleware

```
    ...
    tokenPerms := accessToken.Payload.Perms     // "0605
    bs := bitset.Parse(tokenPerms)              // returns 0000_0110 0000_0101
    if bs.AreSet(2,8) {
        // the permission allowed
    }
```

# Further Improvements

- [ ] Finalize integration BitSet with database/sql
- [ ] Add benchmarks
- [ ] Reduce memory allocations

Prague 2020

## curl examples

```
curl 127.0.0.1:8083/api/auth/sign-in -d '{"login" : "testadmin", "password":"test"}'

{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJlbGVwaGFudHNvZnQiLCJzdWIiOiJ0YXVkaXQiLCJhdWQiOiJodHRwczovL2VsZXBoYW50c29mdC5ydSIsImV4cCI6MTU4NjI0ODQ3MywiaWF0IjoxNTg2MjQ2NjczLCJqdGkiOiJ0ZXN0IiwidXNlcl9pZCI6MTEsImxvZ2luIjoidGVzdGFkbWluIiwicm9sZV9pZCI6MSwicGVybV9iaXRzZXQiOjcsImV4dHJhIjpudWxsfQ.1u_UbBAPHIg819JqJjzDHKsaW2wZBMVcEYjt92FRRWw","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJlbGVwaGFudHNvZnQiLCJzdWIiOiJ0YXVkaXQiLCJhdWQiOiJodHRwczovL2VsZXBoYW50c29mdC5ydSIsImV4cCI6MTU4ODgzODY3MywiaWF0IjoxNTg2MjQ2NjczLCJqdGkiOiJ0ZXN0IiwidXNlcl9pZCI6MTF9.x7383jbhlhk2VhABF8YfgjUY3SNp5_GFqA3lcctupjs","allowed_permissions":{"TestCreateEntity","TestDeleteEntity","TestUpdateEntity"}}

curl -X POST 127.0.0.1:8083/api/auth/is-token-valid -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJlbGVwaGFudHNvZnQiLCJzdWIiOiJ0YXVkaXQiLCJhdWQiOiJodHRwczovL2VsZXBoYW50c29mdC5ydSIsImV4cCI6MTU4NjI0ODQ3MywiaWF0IjoxNTg2MjQ2NjczLCJqdGkiOiJ0ZXN0IiwidXNlcl9pZCI6MTEsImxvZ2luIjoidGVzdGFkbWluIiwicm9sZV9pZCI6MSwicGVybV9iaXRzZXQiOjcsImV4dHJhIjpudWxsfQ.1u_UbBAPHIg819JqJjzDHKsaW2wZBMVcEYjt92FRRWw"

{"result" : "ok"}

curl 127.0.0.1:8083/api/auth/refresh-token -d '{"refresh_token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJlbGVwaGFudHNvZnQiLCJzdWIiOiJ0YXVkaXQiLCJhdWQiOiJodHRwczovL2VsZXBoYW50c29mdC5ydSIsImV4cCI6MTU4ODgzODY3MywiaWF0IjoxNTg2MjQ2NjczLCJqdGkiOiJ0ZXN0IiwidXNlcl9pZCI6MTF9.x7383jbhlhk2VhABF8YfgjUY3SNp5_GFqA3lcctupjs"}'

{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJlbGVwaGFudHNvZnQiLCJzdWIiOiJ0YXVkaXQiLCJhdWQiOiJodHRwczovL2VsZXBoYW50c29mdC5ydSIsImV4cCI6MTU4NjI0ODY1NiwiaWF0IjoxNTg2MjQ2ODU2LCJqdGkiOiJ0ZXN0IiwidXNlcl9pZCI6MTEsImxvZ2luIjoidGVzdGFkbWluIiwicm9sZV9pZCI6MSwicGVybV9iaXRzZXQiOjcsImV4dHJhIjpudWxsfQ.3wD4cfhOFFu_ZTV1jgPz_PcMPvt4MVoHLacUW2QCxG4","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJlbGVwaGFudHNvZnQiLCJzdWIiOiJ0YXVkaXQiLCJhdWQiOiJodHRwczovL2VsZXBoYW50c29mdC5ydSIsImV4cCI6MTU4ODgzODg1NiwiaWF0IjoxNTg2MjQ2ODU2LCJqdGkiOiJ0ZXN0IiwidXNlcl9pZCI6MTF9.wIh1VpkmDKkGFAY5c0IMO0SC3TVwXNsl1NufNzdITUI","allowed_permissions":["TestCreateEntity","TestDeleteEntity", "TestUpdateEntity"]}
```
