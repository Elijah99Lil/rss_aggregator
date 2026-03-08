# Welcome to the Really Simple Syndication Aggregator Project! 

## This is a semi-guided project by Boot.dev to give users the ability to follow news directly in a CLI. 

## Prerequisites:

#### In order to use this project, you will need:
- PostgreSQL v15 or later with a running instance 
- Go 1.23.0 or later
- Unix-based terminal (e.g., Linux or macOS) at the very least installed to run it. 

### We also use Goose for the migrations so you'll need to install goose with the following command:

- go install github.com/pressly/goose/v3/cmd/goose@latest

## Database Setup

Enter the psql shell:

- **Mac:** `psql postgres`
- **Linux:** `sudo -u postgres psql`

Then run the following commands:

```sql
CREATE DATABASE gator;
\c gator
```

#### Linux users will need to also set up a password with this command:

ALTER USER postgres PASSWORD 'your-password';

## Configuration

#### You will also need a config file called "~/.gatorconfig.json" with the following content:

```json
{
  "db_url": "postgres://example"
}
```

## Commands

#### There are a variety of commands you can use in this project, all prefixed with the command:

```go
go run . <command>
```

#### Here are the commands available:

- login <user>
- register <user>
- reset
- users
- agg (Use this in a different terminal to auto-scrape feeds in the background)
- addfeed <url>
- feeds
- follow <url>
- following
- unfollow <url>
- browse