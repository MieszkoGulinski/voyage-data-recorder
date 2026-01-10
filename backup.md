# Backup process

**TODO**: this is not fully implemented, as of now, we have hardcoded `db.sqlite` source and `backup.sqlite` destination

Backup is be done in a separate process, executed using cron and/or on demand by user.

It's worth to know that the backup is done in batches. In SQlite, each page is 4096 bytes. So, if a batch processes 100 pages, it processes 100 _ 4096 = ~400 kB data. A database having 400 MB will require 1000 batches to complete. Even if each batch takes 10 ms, a complete backup will take 10 ms _ 1000 = 10 seconds. It's not a problem, because:

1. The backup sees a snapshot of the database when it starts running, and does not include changes (insert, update, delete) after its start.
2. Unlike VACUUM INTO, using this backup method doesn't block the database from being written to.

## Backup retry

If creating a backup fails, the process retries, up to the number of times specified by the `--retry` option (default is 1 attempt = no retrying). If all retries fail, the process exits with error code 1, like that:

```go
log.Printf("ERROR: backup failed after %d attempt(s): %v", retries, err)
os.Exit(1)
```

## Backup rotation

When backup rotation is active, after adding a new backup file, no longer needed backups will be deleted. In this case, backup process must be executed exactly once every day, and files will have names containing the date, e.g. `2026-01-05.sqlite`.

## Backup integrity check

After creating a backup file, the code runs an integrity check using SQLite `PRAGMA integrity_check`. In case of no error, the response is exactly 1 row with text ok. Failing integrity check means the backup file is corrupted and backup should be retried.

## Saving logs when using cron

When using cron, it's worth to redirect output to a log file:

```
0 2 \* \* \* /usr/local/bin/backup >> /var/log/backup.log 2>&1
```

When the backup is successful, nothing is printed to stdout (unless `--diagnostics` option is used), which is good for avoiding unnecessary SD card wear.

When using systemd timers, stdout and stderr is logged by default, so redirection like that is not needed.
