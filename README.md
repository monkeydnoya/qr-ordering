# QR Ordering Service

## Description: Ordering service stores all created orders by clients and has methods to update and add new items to existing order.

Configration file loaded from ./etc/ordering.yaml

Example:
```yaml
Name: ordering
Host: 0.0.0.0
Port: 8888

DatabaseType: postgres
Postgres:
  Host: localhost
  Port: "5432"
  Username: postgres
  Password: postgres
  DbName: qr
```