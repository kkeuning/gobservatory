In its current state, gobservatory is useful for creating an [awesome list](https://awesome.re) formatted markdown file containing all of your starred GitHub repositories and optionally organizing them using custom tags in addition to categorization by language.

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
Currently just a command line app to load any GitHub starred repositories into Ponzu for an individual GitHub account or extract a markdown file.

# Usage

## Loading and/or updating your CMS:
By default, gobservatory will load starred repositories from the logged in github account (cli will prompt for login), but optionally you can specify the "stargazer" of interest.  It should be possible to aggregate starred repository lists for multiple users by running the cli against each stargazer's name and loading to a common Ponzu.  The cli supports both initial creation and update of content in Ponzu.

```
cd $GOPATH/src/github.com/kkeuning/gobservatory/cmd/gobservatory
go build
./gobservatory load --ponzuSecret="[redacted]" --ponzuUser="yourname@example.com"

```
In the future, gobservatory will likely just prompt you to log into Ponzu.  For now, you need to provide a secret or token.

Replace "[redacted]" in the above with the "Client Secret" from Ponzu (e.g. http://localhost:8080/admin/configure from a logged in session).  This method works if you are running the gobservatory cli and gobservatory-cms on the same server.  

### Accessing a remote server:
Instead of using the secret you can pass in a token remotely:

To obtain the token, you need to log in to Ponzu:
```
curl -v -X "POST" "http://localhost:8080/admin/login" \
     -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
     --data-urlencode "email=yourname@example.com" \
     --data-urlencode "password=redacted"
```
You will find the token in curl output as the value after "Set-Cookie: _token="

Do not include the trailing semi-colon.  


```
./gobservatory load --ponzuToken="[redacted]" --ponzuUser="yourname@example.com"
```

## Managing your content:
Use the Ponzu interface (e.g. http://localhost:8080/admin/contents?type=Star) to manage your content prior to publishing, adding tags or comments to your database of starred repositories.

GitHub will continue to be your authoritative data source, and updates or additions to your starred repositories in GitHub will be pulled into gobservatory on your next update, potentially overwriting your changes with some exceptions:
1.  GitHub doesn't support tags for starred repositories, so your tags won't be overwritten when you update from GitHub.
2.  If you "Unstar" a repository in GitHub, you will need to delete it manually from gobservatory at this point in time.  
3.  If you correct/alter the programming language associated with a project, currently that change will be overwritten by `gobservatory load` on the next update unless its been fixed in GitHub.

## Building your "awesome" markdown file:
```
./gobservatory markdown
```
Since creating the markdown is a read operation, authentication using ponzuToken or ponzuSecret is not required.

The gobservatory cli will to extract starred repositories from the cms and organize them into a markdown file categorized by language.  To also organize by tag, add the "useTags" flag:
```
./gobservatory markdown --ponzuHost="localhost" --ponzuPort="8080" --useTags
```

# Roadmap
- Comments.  The gobservatory cli does not yet include your personal comments in the markdown.  This option is likely in the future.
- Language overrides.  Currently the loader will trust github for identifying the primary language of a project.  Github is occasionally incorrect.  Today you can correct that information in Ponzu, but the next sync with `gobservatory load` will overwrite your changes.  A likely future enhancement will address this.   
