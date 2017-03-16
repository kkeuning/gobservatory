**Problem:**

I have too many starred repositories, need a way to organize them with tags and comments.  

**Solution:**

1. Sync the starred repository information into Ponzu using the GitHub API.
2. Add tags for categorization and additional comments.  
3. Extract to an [awesome list](https://awesome.re) formatted markdown file.
4. Or consume the enriched content from Ponzu into a web client.

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


## cmd/gobservatory
Currently just a command line app to load any github starred repositories into Ponzu for an individual github account.  By default based based on the logged in github account (cli will prompt for login), but optionally you can specify the "stargazer" of interest.  It should be possible to aggregate starred repository lists for multiple users by running the cli against each stargazer's name and loading to a common Ponzu.  The cli supports both initial creation and update of content in Ponzu.

```
cd $GOPATH/src/github.com/kkeuning/gobservatory/cmd/gobservatory
go build
./gobservatory --ponzuSecret="[redacted]" --ponzuUser="yourname@example.com"

```
In the future, gobservatory will likely just prompt you to log into Ponzu.  For now, you need to provide a secret or token.

Replace "[redacted]" in the above with the "Client Secret" from Ponzu (e.g. http://localhost:8080/admin/configure from a logged in session).  This method works if you are running the gobservatory cli and gobservatory-cms on the same server.  If not, you can alternately pass the entire token like this:

```
./gobservatory --ponzuToken="[redacted]"
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

## Building your "awesome" markdown file
```
./gobservatory markdown
```
The gobservatory cli will to extract starred repositories from the cms and organize them into a markdown file based language.  To also organize by tag, use the "useTags" flag:
```
./gobservatory markdown --ponzuHost="localhost" --ponzuPort="8080" --useTags
```

# Roadmap
- Comments.  The gobservatory cli does not yet include your personal comments in the markdown.  This option is likely in the future.
- Language overrides.  Currently the loader will trust github for identifying the primary language of a project.  Github is occasionally incorrect.  Today you can correct that information in Ponzu, but the next sync with `gobservatory load` will overwrite your changes.  A likely future enhancement will address this.   
