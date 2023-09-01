# Table: appwrite_bucket

Get meta information related to buckets for your Appwrite project.

## Examples

### Basic bucket query for a project

```sql
select
  id,
  name,
  file_extension
from
  appwrite_bucket
where
  id = 'YOUR_BUCKET_ID'
  or
  name = 'YOUR_BUCKET_NAME';
```

