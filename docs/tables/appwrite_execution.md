# Table: appwrite_execution

Get meta information related to executions of a function in your Appwrite project.

## Examples

### Basic query for executions of a function

```sql
select
  id,
  status,
  status_code,
  response,
  stdout,
  stderr
from
  appwrite_execution
where
  function_id = 'YOUR_FUNCTION_ID'
```

