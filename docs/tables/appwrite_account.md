# Table: appwrite_account

Get information for accounts for your Appwrite project.

## Examples

### Basic accounts query for a project

```sql
select
  id,
  email,
  name
from
  appwrite_account
where
  status = true 
  and
  email_verification = true;
```

