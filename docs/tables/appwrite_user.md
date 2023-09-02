# Table: appwrite_user

Get information for users for your Appwrite project.

## Examples

### Basic users query for a project

```sql
select
  id,
  email,
  name
from
  appwrite_user
where
  status = true 
  and
  email_verification = true;
```

