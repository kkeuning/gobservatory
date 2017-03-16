package main

var usageHeader = `
Usage: 

gobservatory [command] 

Available Commands:

load
markdown

Flags:
`

var examples = `
Examples:

gobservatory load --ponzuSecret="redacted" --ponzuUser="redacted" --ponzuHost="localhost" --ponzuPort="8080"
gobservatory markdown --ponzuHost="localhost" --ponzuPort="8080"
`
