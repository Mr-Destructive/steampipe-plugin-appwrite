# Table: appwrite_database

Get meta information related to databases for your Appwrite project.

## Examples

### Basic query for databases in a project

```sql
select
  id,
  name
from
  appwrite_database
where
  id = 'YOUR_DATABASE_ID'
  OR
  name = 'YOUR_DATABASE_NAME'
```

