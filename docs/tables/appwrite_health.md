# Table: appwrite_health

Get health information of your Appwrite project.

## Examples

### Query for web server health

```sql
select
  ping,
  status
from
  appwrite_health
where
  service = 'web'
```

### Query for database health

```sql
select
  ping,
  status
from
  appwrite_health
where
  service = 'db'
```

### Query for function-queue 

```sql
select
  size
from
  appwrite_health
where
  service = 'function-queue'
```

### Query for server time health

```sql
select
  local_time,
  real_time,
  diff
from
  appwrite_health
where
  service = 'time'
```
