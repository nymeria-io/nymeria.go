# Nymeria

The official Golang package and command line tool to interact with the Nymeria
service.

## API

## Set and Check an API Key.

```go
nymeria.SetAuth("ny_your-api-key")

if err := nymeria.CheckAuthentication(); err == nil {
  log.Println("OK!")
}
```

All API endpoints assume an auth key has been set. You should set the auth key
early in your program. The key will automatically be added to all future
requests.

## Verify an Email Address

```go
if v, err := nymeria.Verify("someone@somewhere.com"); err == nil && v.Data.Result == "valid" {
  log.Println("OK!")
}
```

At this time only professional email addresses are supported by the API.

## Enrich a Profile

```go
if v, err := nymeria.Enrich("github.com/someone"); err == nil && v.Status == "success" {
  log.Println(v.Data.Emails)
}
```

The enrich API works on a profile by profile basis. If you need to enrich
multiple profiles at once you can use the bulk enrichment API.

## Bulk Enrichment of Profiles

```go
if v, err := nymeria.BulkEnrich("github.com/someone", "linkedin.com/in/someoneelse"); err == nil && v.Status == "success" {
  log.Println(v.Data)
}
```

## License

MIT License

Copyright (c) 2021, Nymeria LLC.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
