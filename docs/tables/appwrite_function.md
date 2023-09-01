# Table: appwrite_function

Get meta information related to functions in your Appwrite project.

## Examples

### Basic query for functions

```sql
select
  id,
  name,
  runtime,
  variable,
  deployment
from
  appwrite_function
where
  name = 'YOUR_FUNCTION_NAME'
```

