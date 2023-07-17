# Обертки для метрик

## HTTP Server

Пример

```

```

## SQL Database

```
import "github.com/happywbfriends/metrics/metrics"

db, _ := sqlx.Open("postgres", "...connection string...")

m := metrics.NewDbMetrics("MyDatabase")
go metrics.DbMetricsHelper(m, db.DB, 5*time.Second, context.TODO())

```

## SQL Database Query

```
import "github.com/happywbfriends/metrics/metrics"

type db struct {
    conn  *sqlx.DB
    
    metrics struct {
      someQuery  metrics.IDbQueryMetrics
    }
}

func newDb(c *sqlx.DB) *db {
  res := &db{c}
  
  res.metrics.someQuery = metrics.NewDbQueryMetrics("MyDatabase", "SomeQuery") 
  
  return res
}



func (d *db) someQuery(...) (e error) {
	defer func(from time.Time) {
		metrics.DbQueryMetricsHelper(d.metrics.someQuery, from, e)
	}(time.Now())
    
    ...
    return nil
}
```
