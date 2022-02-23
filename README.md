1. `docker-compose up -d`
2. `docker exec -ti pgx-slow_postgres_1 psql -U slow`
3.  Create table:

```sql
create table request (
  request_id uuid primary key default gen_random_uuid(),
  data text
);

insert into request (data) (select data from (select generate_series(1,10000) AS id, md5(random()::text) AS data) x);
```
4.  Run via `go run main.go`

You should observe that it gets stuck after 5 iterations of the loop.

To speed this up try one of:

a) qryShort instead of qry
b) enable simple protocol via `config.ConnConfig.PreferSimpleProtocol = true`

Example output:

```
mark@Marks-MacBook-Air pgx-slow % go run main.go
(*pgxpool.Config)(0xc00012c410)({
 ConnConfig: (*pgx.ConnConfig)(0xc00019a000)({
  Config: (pgconn.Config) {
   Host: (string) (len=9) "127.0.0.1",
   Port: (uint16) 5434,
   Database: (string) (len=4) "slow",
   User: (string) (len=4) "slow",
   Password: (string) (len=4) "slow",
   TLSConfig: (*tls.Config)(<nil>),
   ConnectTimeout: (time.Duration) 0s,
   DialFunc: (pgconn.DialFunc) 0x12229a0,
   LookupFunc: (pgconn.LookupFunc) 0x1222a80,
   BuildFrontend: (pgconn.BuildFrontendFunc) 0x12190e0,
   RuntimeParams: (map[string]string) {
   },
   Fallbacks: ([]*pgconn.FallbackConfig) {
   },
   ValidateConnect: (pgconn.ValidateConnectFunc) <nil>,
   AfterConnect: (pgconn.AfterConnectFunc) <nil>,
   OnNotice: (pgconn.NoticeHandler) <nil>,
   OnNotification: (pgconn.NotificationHandler) <nil>,
   createdByParseConfig: (bool) true
  },
  Logger: (pgx.Logger) <nil>,
  LogLevel: (pgx.LogLevel) info,
  connString: (string) (len=76) "user=slow dbname=slow password=slow host=127.0.0.1 sslmode=dis
able port=5434",
  BuildStatementCache: (pgx.BuildStatementCacheFunc) 0x12d0660,
  PreferSimpleProtocol: (bool) false,
  createdByParseConfig: (bool) true
 }),
 BeforeConnect: (func(context.Context, *pgx.ConnConfig) error) <nil>,
 AfterConnect: (func(context.Context, *pgx.Conn) error) <nil>,
 BeforeAcquire: (func(context.Context, *pgx.Conn) bool) <nil>,
 AfterRelease: (func(*pgx.Conn) bool) <nil>,
 MaxConnLifetime: (time.Duration) 1h0m0s,
 MaxConnIdleTime: (time.Duration) 30m0s,
 MaxConns: (int32) 100,
 MinConns: (int32) 0,
 HealthCheckPeriod: (time.Duration) 1m0s,
 LazyConnect: (bool) false,
 createdByParseConfig: (bool) true
})
2022/02/23 14:33:19 starting query...
2022/02/23 14:33:19 ...duration: 10.314167ms
2022/02/23 14:33:19 starting query...
2022/02/23 14:33:19 ...duration: 7.097084ms
2022/02/23 14:33:19 starting query...
2022/02/23 14:33:19 ...duration: 10.007625ms
2022/02/23 14:33:19 starting query...
2022/02/23 14:33:19 ...duration: 12.393417ms
2022/02/23 14:33:19 starting query...
2022/02/23 14:33:19 ...duration: 12.168709ms
2022/02/23 14:33:19 starting query...
2022/02/23 14:33:35 ...duration: 16.34915775s
2022/02/23 14:33:35 starting query...
```
