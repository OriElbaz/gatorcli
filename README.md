# Blog Aggregator CLI
This project was built to explore the integration of Golang with PostgreSQL using Goose for migrations and SQLC for type-safe SQL generation. This project is part of [Boot.dev's](https://boot.dev) backend Golang course.

## Dependencies
* Golang
* PostgreSQL 18+

## How to Run
- Install project:<br>
`go install github.com/OriElbaz/gatorcli@latest`
Alternatively, clone the repo and run `go build -o gator`.

- Set up config file, just change "your-username":<br>
`
{
 "db_url": "postgres://YOUR-USERNAME:@localhost:5432/gator?sslmode=disable",
 "current_user_name": ""
}
`

## Commands
Because I really dont want to spend the time, I'll hand it off to Gemini to explain how to use the commands:<br>

This README section is designed to be clear and informative, categorizing the commands based on whether they require an active login session (as indicated by your `MiddlewareLoggedIn` wrapper).

Gator CLI allows you to manage users, follow RSS feeds, and aggregate posts. Usage follows the pattern:
`gator <command> [arguments]`

### User Management

These commands handle user creation and session switching.

* **`register <name>`** Creates a new user in the database and logs in as that user.
* **`login <name>`** Switches the current user session to the specified user.
* **`users`** Lists all registered users. The currently logged-in user is marked with an asterisk `(*)`.
* **`reset`** **Warning:** Clears all user data from the database. Use with caution.

---

### Feed Management

These commands require the user to be **logged in**.

* **`addfeed <name> <url>`** Adds a new RSS feed to the system and automatically follows it for the current user.
* **`feeds`** Displays a list of all feeds in the system along with the names of the users who added them.
* **`follow <url>`** Creates a follow relationship between the current user and an existing feed URL.
* **`following`** Lists all the feeds the current user is currently following.
* **`unfollow <url>`** Removes the follow relationship for the specified feed URL.

---

### Aggregation & Browsing

These commands handle the background processing and viewing of posts.

* **`agg <time_duration>`** Starts the aggregator. It will fetch the next pending feed every interval (e.g., `1m`, `1h`, or `30s`).
*Example: `gator agg 1m*`
* **`browse [limit]`** *(Requires Login)* Displays posts from the feeds the current user follows. You can optionally provide a limit (e.g., `gator browse 5`).
