
<details>
<summary> Project Structure(click to expand)</summary>
<pre><code>
greenlight-api/
├── bin/                # Compiled application binaries (production-ready)
├── cmd/api/            # Application entry point (main.go) and HTTP server setup
├── internal/           # Core app logic: DB, validation, services, etc.
├── migrations/         # SQL migration scripts
├── remote/             # Production deployment scripts/configs
├── go.mod              # Module definition and dependencies
├── Makefile            # Automation tasks: build, migrate, lint, etc.
└── README.md           # Project documentation
</code></pre>
</details>


**bin**  directory will contain our compiled application binaries, ready for deployment to a production server


**cmd/api** will include the code for running the server, reading and writing HTTP requests, and managing authentication

**internal** will contain the code for interacting with our database, doing data validation, sending emails etc. Any code which isn't application-specific and can potentially be reused will live in here. Our Go code under cmd/api will import the packages in the internal directory (but never the other way around)

**migrations** will contain the SQL migration files for our database

**remote** will contain the configuration files and setup scripts for the production server

**go.mod** will declare our project dependencies, versions and module path

**Makefile** will contain recipes for automating common administrative tasks- like auditing our Go code, building binaries, and executing database migrations

## we will now demonstrate how httprouter works by adding two endpoints for creating a new movie and showing the details of a specific movie


## Encapsulating the API routes
We will encapsulate all the routing rules in a new **cmd/api/routes.go**

## HTTP basics in Golang
## 1. http.Handler--Interface

This is the interface that all HTTP handlers must implement.It looks like:
<code><pre>
type Handler interface{
    serveHTTP(ResponseWriter, *Request)
}
</pre></code>

## 2. http.HandlerFunc --Adapter
This is a type that lets you turn a regular function into something that satisfies the http.Handler interface.
<code><pre>
type HandlerFunc func(ResponseWRiter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request){
    f(w,r)
}
</pre></code>


