# Table: appwrite_document

Get meta information related to documents for a collection in your Appwrite project.

## Examples

### Basic query for documents of a collection

```sql
select
  id,
  name,
  fields
from
  appwrite_document
where
  database_id = 'YOUR_DATABASE_ID'
  and
  collection_id = 'YOUR_COLLECTION_ID';
```

