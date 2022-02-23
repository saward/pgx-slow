package main

import (
	"context"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func main() {

	config, err := pgxpool.ParseConfig("user=slow dbname=slow password=slow host=127.0.0.1 sslmode=disable port=5434")
	if err != nil {
		panic(err)
	}
	config.MaxConns = 100
	// config.ConnConfig.PreferSimpleProtocol = true

	spew.Dump(config)

	dbx, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		panic(err)
	}

	for i := 0; i < 8; i++ {
		// Grab a random uuid:
		var id string
		err := dbx.QueryRow(context.Background(), "select request_id from request order by random() limit 1").Scan(&id)

		if err != nil {
			panic(err)
		}

		// Run pagination
		start := time.Now()
		log.Println("starting query...")
		rows, err := dbx.Query(context.Background(), qry, nil, nil, 5, id, true)
		if err != nil {
			panic(err)
		}
		rows.Close()

		log.Printf("...duration: %s", time.Now().Sub(start))
	}
	_ = dbx
}

var qryShort = "select $1, $2, $3::integer, $4, $5::bool"

var qry = `
SELECT * FROM (SELECT
  *,
  count(*) OVER () > $3 AS has_more,
  row_number() OVER ()
  FROM (
    WITH counted AS (
      SELECT count(*) AS total
      FROM   (select request_id::text, data
from request
where
  (cast($1 as text) is null or lower(data) like '%' || lower($1) || '%')
and
  (cast($2 as text) is null or lower(request_id::text) like '%' || lower($2) || '%')) base
    ), cursor_row AS (
      SELECT base.request_id -- order fields
      FROM   (select request_id::text, data
from request
where
  (cast($1 as text) is null or lower(data) like '%' || lower($1) || '%')
and
  (cast($2 as text) is null or lower(request_id::text) like '%' || lower($2) || '%')) base
      WHERE  base.request_id = $4 -- primary key field name
    )
    SELECT counted.*, base.*
      FROM   (select request_id::text, data
from request
where
  (cast($1 as text) is null or lower(data) like '%' || lower($1) || '%')
and
  (cast($2 as text) is null or lower(request_id::text) like '%' || lower($2) || '%')) base
      LEFT JOIN   cursor_row ON true
      LEFT JOIN   counted ON true
      WHERE ((
            $4 IS NULL OR cast($5 as bool) IS NULL
          ) OR (
            (base.request_id)
              > (cursor_row.request_id)
          ))
      -- Reverse the order of each item if after = false
      ORDER BY base.request_id ASC
      LIMIT $3 + 1
) xy LIMIT $3 ) z ORDER BY row_number ASC
`
