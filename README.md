# Nymeria

[![Go Reference](https://pkg.go.dev/badge/github.com/nymeriaio/nymeria.go.svg)](https://pkg.go.dev/github.com/nymeriaio/nymeria.go)

The official golang package and command line tool to interact with the Nymeria service
and API.

Nymeria makes it easy to enrich data with contact information such as email
addresses, phone numbers and social links. The golang package wraps Nymeria's [public
API](https://www.nymeria.io/developers) so you don't have to.

![Nymeria makes finding contact details a breeze.](https://www.nymeria.io/assets/images/marquee.png)

## Examples

#### Setting an API Key

```go
package main

import (
  "github.com/nymeriaio/nymeria.go"
)

func main() {
  nymeria.ApiKey = "YOUR API KEY GOES HERE"
}
```

All actions that interact with the Nymeria service assume an API key has been
set and will fail if a key hasn't been set. A key only needs to be set once and
can be set at the start of your program.

#### Verifying an Email Address

```go
package main

import (
  "github.com/nymeriaio/nymeria.go"
  "github.com/nymeriaio/nymeria.go/email"
)

func main() {
  nymeria.ApiKey = "YOUR API KEY GOES HERE"

  if v, err := email.Verify("dev@nymeria.io"); err == nil {
    log.Println(v.Result)
  }
}
```

You can verify the deliverability of an email address using Nymeria's service.
The response will contain a `Result` and `Flags`.

The `Result` will either be "valid" or "invalid". The `Flags` will give you additional
details regarding the email address. For example, the tags will tell you if the mail
server connection was successful, if the domain's DNS records are set up to send and
receive email, etc.

#### Enriching Profiles

```go
package main

import (
  "github.com/nymeriaio/nymeria.go"
  "github.com/nymeriaio/nymeria.go/person"
)

func main() {
  nymeria.ApiKey = "YOUR API KEY GOES HERE"

  params := person.EnrichParams{
    Profile: "github.com/nymeriaio", /* you can locate contact details using a supported URL */
    Email: "steve@woz.org",          /* you can perform an enrichment using an email address */
  }

  if person, err := person.Enrich(params); err == nil {
    log.Println(person.Emails, person.PhoneNumbers)
  }
}
```

If you want to enrich an email address you can specify an `Email` and the
Nymeria service will locate the person and return all associated data for them.
Likewise, you can specify a supported URL via the `Profile` parameter if you prefer
to enrich via a URL.

At this time, Nymeria supports look ups for the following sites:

1. LinkedIn
2. Facebook
3. Twitter
4. GitHub

Please note, if using LinkedIn URLs provide the public profile LinkedIn URL.

Two other common parameters are `Filter` and `Require`. If you wish to filter
out professional emails (only receive personal emails) you can do so by
specifying `professional-emails` as the Filter parameter.

The `Require` parameter works by requiring certain kinds of data. For example,
you can request an enrichment but only receive a result if the profile contains
a phone number (or an email, personal email, professional email, etc). The
following are all valid requirements:

1. `email`
2. `phone`
3. `professional-email`
4. `personal-email`

You can specify multiple requirements by using a comma between each
requirement. For example you can require a phone and personal email with:
`phone,personal-email` as the Require parameter.

#### Searching for People

```go
package main

import (
  "github.com/nymeriaio/nymeria.go"
  "github.com/nymeriaio/nymeria.go/person"
)

func main() {
  nymeria.ApiKey = "YOUR API KEY GOES HERE"

  people, err := person.Search(&nymeria.SearchQuery{
    Query: `location_name:"New York" has_email:true`,
  })

  if err != nil {
    log.Fatal(err)
  }
}
```

By default, 10 people will be returned for each page of search results. You can
specify the `Size` as part of the `SearchQuery` if you want to access
additional pages of people.

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
