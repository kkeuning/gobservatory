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

# Components:

## gobservatory-cms
This is the Ponzu server to store and manage the starred repository records.  Tags and comments should be edited here.

## gobservatory-loader
Currently just a command line app to load any github starred repositories into Ponzu for an individual github account.  Initially based on the logged in github account (cli will prompt for login).  I have some ideas about aggregating starred repository lists for multiple users that are not yet implemented.  Initially there is no support for bulk update, once a star is added to Ponzu and exists it would be skipped on the next load.

## gobservatory-awesome
Not yet implemented.   Will be similar to gobservatory-loader, a command line app to extract starred repositories from the cms and organize them into a markdown file based on tags and language.

## gobservatory-buffalo
Not yet implemented.  Initially will take the gobservatory-awesome output and render it using buffalo.  Later I'd like to include an example of pulling the CMS json content live into a web app and may do this with Buffalo and/or React.
