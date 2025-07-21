
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

**DSN- data source name, string that contains the necessary connection parameters**


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

<code><pre>
Additional Information
API versioning
APIs which support real-world businesses and users often need to change their functionality
and endpoints over time — sometimes in a backwards-incompatible way. So, to avoid
problems and confusion for clients, it’s a good idea to always implement some form of API
versioning.
There are two common approaches to doing this:
1. By prefixing all URLs with your API version, like /v1/healthcheck or /v2/healthcheck.
2. By using custom Accept and Content-Type headers on requests and responses to
convey the API version, like Accept: application/vnd.greenlight-v1 .
From a HTTP semantics point of view, using headers to convey the API version is the ‘purer’
approach. But from a user-experience point of view, using a URL prefix is arguably better. It
makes it possible for developers to see which version of the API is being used at a glance,
and it also means that the API can still be explored using a regular web browser (which is
harder if custom headers are required).
Throughout this book we’ll version our API by prefixing all the URL paths with /v1/ — just
like we did with the /v1/healthcheck endpoint in this chapter.
</pre></code>


<code><pre>
Note: The readIDParam() method doesn’t use any dependencies from our
application struct so it could just be a regular function, rather than a method on
application. But in general, I suggest setting up all your application-specific handlers
and helpers so that they are methods on application. It helps maintain consistency in
your code structure, and also future-proofs your code for when those handlers and
helpers change later and they do need access to a dependency.
</pre><code>

#### Parentheses {} define a JSON object, which is made up of comma-separated key/value pairs.. we need to set a Content-Type: application/json header on the response, so that the client knows it's receiving JSON and can interpret it 


