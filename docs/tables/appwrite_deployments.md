# Table: appwrite_deployments

Get meta information related to deployments for a function in your Appwrite project.

## Examples

### Basic query for deployments of a function

```sql
select
  *
from
  appwrite_deployments
where
  function_id = 'YOUR_FUNCTION_ID'
```

