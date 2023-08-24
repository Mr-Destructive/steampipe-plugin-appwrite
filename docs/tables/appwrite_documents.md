# Table: appwrite_documents

Get meta information related to documents for a collection in your Appwrite project.

## Examples

### Basic query for documents of a collection

```sql
select
  *
from
  appwrite_documents
where
  database_id = 'YOUR_DATABASE_ID'
  and
  collection_id = 'YOUR_COLLECTION_ID';
```

