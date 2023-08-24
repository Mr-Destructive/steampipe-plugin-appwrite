# Table: appwrite_accounts

Get information for accounts for your Appwrite project.

## Examples

### Basic accounts query for a project

```sql
select
  *
from
  appwrite_accounts
where
  status = true 
  and
  email_verification = true;
```

