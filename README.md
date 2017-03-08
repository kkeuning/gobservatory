**Problem:**

I have too many starred repositories, need a way to organize them with tags and comments.  

**Potential Solution:**

1. Sync the starred repository information into Ponzu using the GitHub API.
2. Add labels and more detailed comments.  
3. Extract to an [awesome list](https://awesome.re) formatted markdown file.
4. Or consume the enriched content from Ponzu into a web page.

**Other goals:**

This is deliberately an experiment with using Ponzu and Buffalo, along with a few other intentional dependencies of projects I have an interest in.  After the initial POC this list may change:
* [Ponzu](https://github.com/ponzu-cms/ponzu) CMS
* [go-github](https://github.com/google/go-github) SDK for GitHub API
* [structs](https://github.com/fatih/structs) package of utilities for working with go structs
* [buffalo](https://github.com/gobuffalo/buffalo) Web framework for fronting the CMS content with a web app

**Important: This project requires Go 1.8 or higher.**

# Components:

## gobservatory-cms
This is the Ponzu server to store and manage the starred repository records.  Tags and comments should be edited here.

```
go get github.com/ponzu-cms/ponzu/...
go get github.com/kkeuning/gobservatory/...
cd $GOPATH/src/github.com/kkeuning/gobservatory/gobservatory-cms
ponzu build
ponzu run
```


## gobservatory-loader
Currently just a command line app to load any github starred repositories into Ponzu for an individual github account.  Initially based on the logged in github account (cli will prompt for login).  I have some ideas about aggregating starred repository lists for multiple users that are not yet implemented.  Supports both initial creation and update of content in Ponzu.

```
cd $GOPATH/src/github.com/kkeuning/gobservatory/gobservatory-loader
go build
./gobservatory-loader --ponzuSecret="[redacted]" --ponzuUser yourname@example.com

```
In the future, gobservatory-loader will likely just prompt you to log into Ponzu.  For now, you need to provide a secret or token.

Replace "[redacted]" in the above with the "Client Secret" from Ponzu (e.g. http://localhost:8080/admin/configure from a logged in session).  This method works if you are running gobservatory-loader and gobservatory-cms on the same server.  If not, you can alternately pass the entire token like this:

```
./gobservatory-loader --ponzuToken="[redacted]"
```
To obtain the token, you need to log in to Ponzu:
```
curl -v -X "POST" "http://localhost:8080/admin/login" \
     -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
     --data-urlencode "email=yourname@example.com" \
     --data-urlencode "password=redacted"
```
You will find the token in curl output as the value after "Set-Cookie: _token="

Do not include the trailing semi-colon.  

## gobservatory-awesome
Not yet implemented.   Will be similar to gobservatory-loader, a command line app to extract starred repositories from the cms and organize them into a markdown file based on tags and language.

## gobservatory-buffalo
Not yet implemented.  Initially will take the gobservatory-awesome output and render it using buffalo.  Later I'd like to include an example of pulling the CMS json content live into a web app and may do this with Buffalo and/or React.
