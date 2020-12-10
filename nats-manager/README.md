# nats-manager

nats-manager manages accounts for embedded Timeterm (frontend-embedded) devices, referred to as emdevs to keep topics
short.

## Lifecycle

When the nats-manager starts up, it attempts to connect with the database. If connecting fails, it terminates. When the
database is empty, it runs its migrations and automatically sets up some of the scaffolding for creating key pairs and
signing JWTs. This happens as follows:

1. Create a public key for the system account
2. Create a new operator and set the `system_account` property to the public key of the system account
   > All accounts created after the creation of this operator are issued with it.

3. Create a new account and user for the operator
4. Create a system account and system account user
5. Run JWT migrations (which involves creating several users required by the Timeterm system itself)

## Account Structure

Convention is to use uppercase (SCREAMING_SNAKE_CASE) for operator and account names, and to use lowercase (snake_case)
for usernames.

| Entity | Name            | Description                                                            |
|--------|-----------------|------------------------------------------------------------------------|
| `O`    | `TIMETERM`&ast; | Operator which issues all accounts                                     |
| `-A`   | `TIMETERM`&ast; | Issues `timeterm` user                                                 |
| `--U`  | `timeterm`&ast; | Currently none                                                         |
| `-A`   | `BACKEND`       | Issues all users that are its children                                 | 
| `--U`  | `backend`       | Used by the Timeterm backend (only pub/req on allowed service methods) |
| `--U`  | `superuser`     | Used by nats-manager to control JetStream                              |
| `--U`  | `nats-manager`  | Used by nats-manager (only pub/sub on `NATS-MANAGER.>`)                |
| `-A`   | `EMDEVS`        | Used by the Timeterm backend to send messages to devices, has many private exports |
| `--U`  | `backend`       | Used by the Timeterm backend (pub/sub on `EMDEV.>`, required JetStream topics) |
| `-A`   | `EMDEV-?`       | Used by a device, imports from `EMDEVS` with activations&ast;&ast; |
| `--U`  | `emdev`         | Used by the device itself |
| `--U`  | `superuser`     | Used by nats-manager to control JetStream (if necessary) |

&ast; = name can be specified by user in configuration  
&ast;&ast; = JetStream topics are not imported as `$JS.>`, but instead under other subjects as not to conflict with account-specific JetStream topics. Importing JetStream streams from other accounts is possible since the nightly NATS build of 11/20/2020.

Individual accounts for individual devices are created to allow for revocation of non-expiring accounts without making the `EMDEVS` account a big fat mess. 
In the case where every device would be given its own user under the `EMDEVS` account with a non-expiring token, revocations of user tokens must be noted in the account JWT. The token would increase in size as more device users are revoked.
In the first case, revoking access for a specific device is trivial: remove the account and activation, and the user and all access to streams from the `EMDEVS` account is void.
