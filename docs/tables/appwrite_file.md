# Table: appwrite_file

Get meta information related to files in a bucket for your Appwrite project.

## Examples

### Basic query for files in a bucket

```sql
select
  id,
  name,
  mime_type,
  size_original
from
  appwrite_file
where
  bucket_id = 'YOUR_BUCKET_ID'
  and
  chunks_uploaded = chunks_total
```

