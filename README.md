# MySQL R2 Backups

A Go-based MySQL backup utility that runs on Railway, takes snapshots of the database, and uploads them to Cloudflare R2.

## Features

-   Periodic MySQL database backups using a cron schedule
-   Securely uploads backups to Cloudflare R2
-   Configurable through environment variables
-   Supports running on Railway

## Requirements

Before using this tool, ensure you have the following:

-   A MySQL database (hosted on Railway or another provider).

-   A Cloudflare R2 bucket to store the backups.

## Environment Variables

The application is configured using the following environment variables:

| Variable                        | Default       | Description                                       |
| ------------------------------- | ------------- | ------------------------------------------------- |
| `CLOUD_FLARE_ACCESS_KEY_ID`     | ""            | Cloudflare R2 Access Key ID                       |
| `CLOUD_FLARE_SECRET_ACCESS_KEY` | ""            | Cloudflare R2 Secret Access Key                   |
| `CLOUD_FLARE_R2_BUCKET`         | ""            | Cloudflare R2 Bucket Name                         |
| `CLOUD_FLARE_R2_REGION`         | `auto`        | Cloudflare R2 Region                              |
| `CLOUD_FLARE_R2_ENDPOINT`       | ""            | Cloudflare R2 Endpoint                            |
| `BACKUP_DATABASE_URL`           | ""            | MySQL Database URL                                |
| `BACKUP_DATABASE_PASSWORD`      | ""            | MySQL Database Password                           |
| `BACKUP_DATABASE_USER`          | ""            | MySQL Database Username                           |
| `BACKUP_DATABASE_HOST`          | ""            | MySQL Database Host                               |
| `BACKUP_DATABASE_PORT`          | `3306`        | MySQL Database Port                               |
| `BACKUP_DATABASE_NAME`          | ""            | MySQL Database Name                               |
| `BACKUP_CRON_SCHEDULE`          | `0 0 3 * * *` | Cron schedule for backups(scheduled to 3AM daily) |
| `RUN_ON_STARTUP`                | `true`        | Run backup immediately on startup                 |
| `SINGLE_SHOT_MODE`              | `false`       | Run once and exit                                 |
| `BACKUP_FILE_PREFIX`            | `snapshot`    | Prefix for backup files                           |
| `BUCKET_SUBFOLDER`              | `snapshots`   | Subfolder in R2 bucket                            |

## Blog

You can find a blog on how to set up at [Automatic MySQL Backups](https://emiliocliff.com)
