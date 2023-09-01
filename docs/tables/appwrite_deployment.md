# Table: appwrite_deployment

Get meta information related to deployments for a function in your Appwrite project.

## Examples

### Basic query for deployments of a function

```sql
select
  id,
  status,
  build_id,
  build_stdout,
  build_stderr
from
  appwrite_deployment
where
  function_id = 'YOUR_FUNCTION_ID'
```

