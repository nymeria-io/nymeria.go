# Nymeria

[![Go Reference](https://pkg.go.dev/badge/git.nymeria.io/nymeria.go.svg)](https://pkg.go.dev/git.nymeria.io/nymeria.go)

The official golang package and command line tool to interact with the Nymeria service
and API.

Nymeria makes it easy to enrich data with contact information such as email
addresses, phone numbers and social links. The golang package wraps Nymeria's [public
API](https://www.nymeria.io/developers) so you don't have to.

![Nymeria makes finding contact details a breeze.](https://www.nymeria.io/assets/images/marquee.png)

## Go API

#### Set and Check an API Key.

```go
nymeria.SetAuth("ny_your-api-key")

if err := nymeria.CheckAuthentication(); err == nil {
  log.Println("OK!")
}
```

All API endpoints assume an api key has been set. You should set the api key
early in your program. The key will automatically be added to all future
requests.

#### Verify an Email Address

```go
nymeria.SetAuth("ny_your-api-key")

if v, err := nymeria.Verify("someone@somewhere.com"); err == nil {
  log.Println(v.Data.Result)
}
```

#### Enrich Profiles

You can enrich one or more profiles using the Enrich method. The Enrich
method takes one or more enrich params and returns one or more enrichment
records.

```go
nymeria.SetAuth("ny_your-api-key")

params := []nymeria.EnrichParams{
  {
    URL: "github.com/nymeriaio",
  },
  {
    Email: "steve@woz.org",
  },
}

if es, err := nymeria.Enrich(params...); err == nil {
  for _, enrichment := range es {
    if enrichment.Status == "success" {
      log.Println(enrichment.Meta)               /* input params, etc */

      log.Println(enrichment.Data.Bio)
      log.Println(enrichment.Data.Emails)
      log.Println(enrichment.Data.PhoneNumbers)
      log.Println(enrichment.Data.Social)
    }
  }
}
```

## Command Line Tool

The command line tool enables you to quickly test the Nymeria API.

#### Installation

You can install the command line tool with `go install`.

```bash
$ go install git.nymeria.io/nymeria.go/cmd/nymeria@v1.0.6
```

#### Set an API Key

```bash
$ nymeria --auth ny_abc-123-456
```

The API key will be cached for future commands.

#### Purge all cached data.

```bash
$ nymeria --purge
```

#### Verify an Email Address

```bash
$ nymeria --verify someone@somewhere.com
```

#### Enrich Profiles

```bash
$ nymeria --enrich '[{ "url": "github.com/nymeriaio" }, { "email": "steve@woz.org" }]'
```

## License

MIT License

Copyright (c) 2022, Nymeria LLC.

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
