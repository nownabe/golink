# Design Doc

## Context

See [PRD](./PRD.md).

## Terminology

- **Golink**: the whole system of this Golink
- **golink**: a definition of a redirect, which is composed of its name and a target url to redirect

## Overview

Golink consists of six components: Chrome extension, API, Redirector, database, Identity-Aware Proxy, Console.

[![architecture](./architecture.png)](https://googlecloudcheatsheet.withgoogle.com/architecture?link=0038ca10-313a-11ee-82a6-f9c22e3525da)

Server-side components are built on Google Cloud.
Clients, Chrome Extension and Console, communicate with API via gRPC.

Chrome extension has two roles: First, the service worker interrupts requests to `http://go/*` and redirects to Redirector. Second, the popup provides quick operations for golinks, such as creating a new golink.

## Detailed Design

### Infrastructure

As the diagram above shows, all server-side components are built on Google Cloud.

Redirector, API, and Console are hosted on Google App Engine and `dispatch.yaml` dispatches requests to each service.
The database is Cloud Firestore, which is used by Redirector and API.
These App Engine services are protected by Identity-Aware Proxy, so Golink users needs to sign in.

### gRPC

Methods and the Golink message are as follow.

```protobuf
service GolinkService {
  rpc CreateGolink(CreateGolinkRequest) returns (CreateGolinkResponse) {}
  rpc GetGolink(GetGolinkRequest) returns (GetGolinkResponse) {}
  rpc ListGolinks(ListGolinksRequest) returns (ListGolinksResponse) {}
  rpc ListGolinksByUrl(ListGolinksByUrlRequest) returns (ListGolinksByUrlResponse) {}
  rpc UpdateGolink(UpdateGolinkRequest) returns (UpdateGolinkResponse) {}
  rpc DeleteGolink(DeleteGolinkRequest) returns (DeleteGolinkResponse) {}
  rpc AddOwner(AddOwnerRequest) returns (AddOwnerResponse) {}
  rpc DeleteOwner(DeleteOwnerRequest) returns (DeleteOwnerResponse) {}
}

message Golink {
  string name = 1;
  string url = 2;
  repeated string owners = 3;
}
```

#### `CreateGolink` method

```proto
message CreateGolinkRequest {
  string name = 1;
  string url = 2;
}

message CreateGolinkResponse {
  Golink golink = 1;
}
```

`CreateGoink` creates a new golink by a name and a URL. If the given name is already taken, `CreateGolink` returns an `ALREADY_EXISTS` error.
`CreateGolink` is used by both Console and Extension.

#### `GetGolink` method

```proto
message GetGolinkRequest {
  string name = 1;
}

message GetGolinkResponse {
  Golink golink = 1;
}
```

`GetGolink` gets a golink by its name. If no golink is found by the given name, `GetGolink` returns a `NOT_FOUND` error.
`GetGolink` is used by Console.

#### `ListGolinks` method

```proto
message ListGolinksRequest {
}

message ListGolinksResponse {
  repeated Golink golinks = 1;
}
```

`ListGolinks` returns golinks that the user owns.
`ListGolinks` is used by Console.

#### `ListGolinksByUrl` method

```proto
message ListGolinksByURLRequest {
  string url = 1;
}

message ListGolinksByURLResponse {
  repeated Golink golinks = 1;
}
```

`ListGolinksByUrl` returns golinks associated to the given URL.
`ListGolinksByUrl` is used by Extension.

**Note**: "Url" or "URL"? It is mentiond in [Google's API design guide](https://cloud.google.com/apis/design/naming_convention#camel_case) that except for field names and enum values, all names must be UpperCamelCase, as defined by [Google Java Style](https://google.github.io/styleguide/javaguide.html#s5.3-camel-case).

#### `UpdateGolink` method

```proto
message UpdateGolinkRequest {
  string name = 1;
  string url = 2;
}

message UpdateGolinkResponse {
  Golink golink = 1;
}
```

`UpdateGolink` updates the URL of a specified golink. If the user is not the owner of the golink, `UpdateGolink` returns a `PERMISSION_DENIED` error.
`UpdateGolink` is used by Console.

#### `DeleteGolink` method

```proto
message DeleteGolinkRequest {
  string name = 1;
}

message DeleteGolinkResponse {
}
```

`DeleteGolink` deletes a golink by its name. If the user is not the owner of the golink, `DeleteGolink` returns a `PERMISSION_DENIED` error.
`DeleteGolink` is used by Console.

#### `AddOwner` method

```proto
message AddOwnerRequest {
  string name = 1;
  string owner = 2;
}

message AddOwnerResponse {
  Golink golink = 1;
}
```

`AddOwner` adds a new owner given as an email. If the request user is not an owner of the golink, `AddOwner` returns a `PERMISSION_DENIED` error.
`AddOwner` is used by Console.

#### `DeleteOwner` methods

```proto
message DeleteOwnerRequest {
  string name = 1;
  string owner = 2;
}

message DeleteOwnerResponse {
  Golink golink = 1;
}
```

`DeleteOwner` remove a specified owner given as an email from a golink. If the request user is not an owner of the golink, `DeleteOwner` returns a `PERMISSION_DENIED` error.
If the request user is the last owner of the golink, `DeleteOwner` returns a `FAILED_PRECONDITION` error.
`DeleteOwner` is used by Console.

### Database

There are only an object type, Golink. Golink documents are stored `golinks/` collection and identified by golink name.

```json
// golinks/ collection
{
  "linkname1": {
    "url": "http://example.com/foo",
    "redirect_count": 10,
    "created_at": "...",
    "updated_at": "...",
    "owners": [
      "owner1@example.com",
      "owner2@example.com"
    ]
  },
  "linkname2": { ... }
}
```

Each gRPC method runs following queries:

```go
// CreateGolink, GetGolink, UpdateGolink, DeleteGolink, AddOwner, DeleteOwner
client.Collection("golinks").Doc(linkName);

// ListGolinks
client.Collection("golinks").Where("owners", "array-contains-any", userEmail);

// ListGolinksByUrl
client.Collection("golinks").Where("url", "==", url);
```

### Identity-Aware Proxy

### API

### Redirector

### Chrome Extension

### Console (Web Frontend)

## Alternatives Considered
