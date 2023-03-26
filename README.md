**Warning: Until stable release, this framework is susceptible to breaking changes**

# Mango

<p>
  <img width="260px" src="https://raw.githubusercontent.com/mangopkg/Assets/main/mango.png">
</p>
<p>
	<img alt="gopher" src="https://img.shields.io/github/go-mod/go-version/mangopkg/mango?style=for-the-badge&logo=appveyor"/>
	<img alt="gopher" src="https://img.shields.io/github/license/mangopkg/mango?style=for-the-badge"/>
</p>


An easy-to-use REST API framework built on top of Go-Chi.

*New: Auto doc generation using openAPI specs and swagger-ui is now supported*

You can view docs on http://localhost:3000/api/dist

## Get Started
Getting started is easy with the handy cli tool. Make sure that you have go already installed.
```go
go install github.com/mangopkg/mng@latest
```
After installing the cli tool, Open a terminal in the directory where you would like to set up the project. Once the terminal is open perform
```go
mng new web-app
```
**web-app** can be whatever you want to call your new app.


Now, `cd` into the newly created project directory and open it in a terminal and issue the `go mod tidy` command
```go
go mod tidy
```
This will ensure all the dependencies have been successfully installed.

Now to finally run the app, Open the root of project in terminal and type

```go
go run .
```

And your app is live!

Open http://localhost:3000/book/find to see your api in action.

You can now modify the code.

## Adding new routes
New routes to your api can be easily added using the `mng` cli tool. Open the terminal in projects root and type
```go
mng add user
```
This will create a new route `/user` in your app

But before this route can work, you need to initialise it by calling `NewService` method in `/api/api.go`

Open the file and edit `initServices` function by adding `user.NewService(s)` to the function such that
```go
func initServices(s mango.Service) {
    book.NewService(s)
    user.NewService(s)
}
```
Now restart your app and open http://localhost:3000/user/find to see the new route in action!

Make sure to repeat this step each time you add a route!
### Defining routes
As you are probably well aware that implementing routes can become an issue when you have to define a lot of them. Mango solves this issue by what we call it **attributed routes**. You can define routes using comments and leave the rest for mangogic (mango + magic)

Defining a route is as easy as adding a function and adding a comment on top of it.

```go
/*
</route{
"pattern": "/find",
"func": "Find",
"method": "GET"
}route/>
*/
func (h *BookHandler) Find() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Response.Message = "Successful"
		h.Response.StatusCode = 200
		h.Response.Data = h.BookService.GetBook()
		h.Response.Send(w)
	}
}
```

Here we are adding a comment with the format `/*
<@route{} */`. Inside `@route{}` we have 3 properties
- pattern - this is the pattern of your route.
- func - func is the associated function to this route. In this case `Find`
- method - the http method to use for this route.

More or less, we define a json object inside `@route{}` which is based on the struct.

```go
type RInfo struct {
	Pattern string                                   `json:"pattern"`
	Func    string                                   `json:"func"`
	Method  string                                   `json:"method"`
	Auth    string                                   `json:"auth"`
	ReqBody string                                   `json:"reqBody"`
	Info    interface{}                              `json:"Info"`
	MountAt string                                   `json:"mountAt"`
	Handler func(http.ResponseWriter, *http.Request) `json:"handler"`
}
```

***Note: This is subjected to breaking updates until stable release***


## Response
Mango has an easy to use response utitlity that can come in handy.

```go
type Response struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Error      bool        `json:"error"`
}
```

These 4 fields are accessible in handler function. These struct has a method called `Send` that will send the API response. You need to pass the `http.ResponseWriter` to this method.

Here's an example code
```go
func (h *BookHandler) Find() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Response.Message = "Successful"
		h.Response.StatusCode = 200
		h.Response.Data = h.BookService.GetBook()
		h.Response.Send(w)
	}
}
```


## API Documentation

Api documentation is now supported in mango, The framework will auto doc your routes but you have to add info manually for more verbose docs. Future updates will try to auto document more aspects of your API.

**Warning: Currently, You must run your app on port 3000 on localhost for auto documentation support, This behaviour will be changed in a future release**

#### Specifying specs manually

Since auto doc is currently available in limited capacity, you can manually specify info for your routes.
Example code:
```go
/*
<@route{
"pattern": "/find",
"func": "Find",
"method": "GET",

	"info": {
				"get": {
					"description": "Returns all books from the system that the user has access to",
					"responses": {
					"200": {
						"description": "A list of books."
					}
					}
				}
			}
	}>
*/
func (h *BookHandler) Find() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Response.Message = "Successful"
		h.Response.StatusCode = 200
		h.Response.Data = "123"
		h.Response.Send(w)
	}
}
```

This strictly follows specs from https://github.com/go-openapi/spec Read their documetation for further information.

## Using Go-Chi
This framework is currently built on top of Go-Chi framework, You can read their documentation at https://github.com/go-chi/chi

**This framework uses chi**
