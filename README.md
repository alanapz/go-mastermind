# GOMASTER

Test Go application, used to test the following technologies:

- Basic Go REST API (using Gin framework)
- Front-end using Knockout
- Testing building single static images


Image is statically built: Only 1 binary and assets (CSS/HTML).

In Dive:

```
Permission     UID:GID       Size  Filetree
drwxr-xr-x         0:0      11 MB  └── app
drwxr-xr-x         0:0      11 kB      ├── assets
-rwxr-xr-x         0:0      237 B      │   ├── gomaster.css
-rwxr-xr-x         0:0      11 kB      │   └── index.html
-rwxr-xr-x         0:0      11 MB      └── gomaster
```

For a total image size of 11MB.

# To Run

Via Docker:

```
$ docker run --init -p 8080:8080 docker.io/alanmpinder/gomaster
```

Then visit http://localhost:8080 in your browser.


# To Build (Docker)

From the root folder:

```
$ cd app
$ docker build -f Dockerfile gomaster
```

# To Build (from source)

From the root folder:

```
$ cd app\gomaster
$ go build .

# Or, to build a static binary
# (Dont forget to copy assets folder if necessary)
$ CGO_ENABLED=0 go build -ldflags="-extldflags=-static"
```

