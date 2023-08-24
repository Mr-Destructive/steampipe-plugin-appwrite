# Table: appwrite_databases

Get meta information related to databases for your Appwrite project.

## Examples

### Basic query for databases in a project

```sql
select
  *
from
  appwrite_databases
where
  id = 'YOUR_DATABASE_ID'
  OR
  name = 'YOUR_DATABASE_NAME'
```

