# modules/notification

User notification inbox ownership for the v4 modular layout.

```
controller/   HTTP entry (gateway path parity)
services/     Count/list/read (client adapters during migration)
dto/          Request shapes for HTTP
entity/       Domain sketches for later persistence
repository/   Store interfaces (postgres later)
```

Phase-1 routes:

- `GET /v2/notification/count`
- `GET /v2/notification/list`
- `PUT /v2/notification/read/{id}`
- `PUT /v2/notification/read-many`
- `PUT /v2/notification/read-all`
- `DELETE /v2/notification/remove/{id}`

Push / admin-history / account-warning are out of scope.
