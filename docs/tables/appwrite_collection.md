# Table: appwrite_collection

Get meta information related to collections for your Appwrite project.

## Examples

### Basic query for collections in a project

```sql
select
  id,
  name,
  attributes
from
  appwrite_collection
where
  database_id = 'YOUR_DATABASE_ID'
```

