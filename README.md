# Nymeria

[![Go Reference](https://pkg.go.dev/badge/git.nymeria.io/nymeria.go.svg)](https://pkg.go.dev/git.nymeria.io/nymeria.go)

The official golang package and command line tool to interact with the Nymeria service
and API.

Nymeria makes it easy to enrich data with contact information such as email
addresses, phone numbers and social links. The golang package wraps Nymeria's [public
API](https://www.nymeria.io/developers) so you don't have to.

![Nymeria makes finding contact details a breeze.](https://www.nymeria.io/assets/images/marquee.png)

## Examples

#### Setting and Checking an API Key

```go
import (
	"git.nymeria.io/nymeria.go"
)

nymeria.SetAuth("YOUR API KEY GOES HERE")

if err := nymeria.CheckAuthentication(); err == nil {
  log.Println("OK!")
}
```

All actions that interact with the Nymeria service assume an API key has been
set and will fail if a key hasn't been set. A key only needs to be set once and
can be set at the start of your program.

If you want to check a key's validity you can use the CheckAuthentication
function to verify the validity of a key that has been set. If no error is
returned then the API key is valid.

#### Verifying an Email Address

```go
import (
	"git.nymeria.io/nymeria.go"
)

nymeria.SetAuth("YOUR API KEY GOES HERE")

if v, err := nymeria.Verify("dev@nymeria.io"); err == nil {
  log.Println(v.Data.Result)
}
```

You can verify the deliverability of an email address using Nymeria's service.
The response will contain a `Result` and `Tags`.

The `Result` will either be "valid" or "invalid". The `Tags` will give you additional
details regarding the email address. For example, the tags will tell you if the mail
server connection was successful, if the domain's DNS records are set up to send and
receive email, etc.

#### Enriching Profiles

```go
import (
	"git.nymeria.io/nymeria.go"
)

nymeria.SetAuth("YOUR API KEY GOES HERE")

params := []nymeria.EnrichParams{
  {
    URL: "github.com/nymeriaio", /* you can locate contact details using a supported URL */
  },
  {
    Email: "steve@woz.org",      /* you can perform an enrichment using an email address */
  },
}

if es, err := nymeria.Enrich(params...); err == nil {
  for _, enrichment := range es {
    if enrichment.Status == "success" {
      log.Println(enrichment.Meta)

      log.Println(enrichment.Data.Bio)
      log.Println(enrichment.Data.Emails)
      log.Println(enrichment.Data.PhoneNumbers)
      log.Println(enrichment.Data.Social)
    }
  }
}
```

You can enrich one or more profiles using the `Enrich` function. The Enrich
function takes a variable number of `EnrichParams`. The most common parameters
to use are `URL` and `Email`.

If you want to enrich an email address you can specify an `Email` and the
Nymeria service will locate the person and return all associated data for them.
Likewise, you can specify a supported url via the `URL` parameter if you prefer
to enrich via a URL.

At this time, Nymeria supports look ups for the following URLs:

1. LinkedIn
2. Facebook
3. Twitter
4. GitHub

If using LinkedIn URLs, please provide a public LinkedIn URL.

#### Searching for People

```go
import (
	"git.nymeria.io/nymeria.go"
)

nymeria.SetAuth("YOUR API KEY GOES HERE")

// Search for Ruby on Rails skills in the Palo Alto area, and only people that
// contain email addresses.

resp, err := nymeria.People(&nymeria.PeopleQuery{
  Skills:   []string{"Ruby on Rails"},
  Location: "Palo Alto",
  HasEmail: true,
})

if err != nil {
  log.Fatal(err)
}

uuids := []string{}

for _, preview := range resp.Data {
  uuids = append(uuids, preview.UUID) /* add to the slice of uuids to reveal */

  log.Println(preview.UUID)
  log.Println(preview.FirstName)
  log.Println(preview.LastName)
  log.Println(preview.AvailableData)
}

// Reveal the complete details for all of the people we located.

resp, err := nymeria.RevealPeople(uuids)

if err != nil {
  log.Fatal(err)
}

for _, person := range resp.Data {
  log.Println(person.Bio)
  log.Println(person.Emails)
  log.Println(person.PhoneNumbers)
  log.Println(person.Social)
}
```

You can perform searches using Nymeria's database of people. The search works
using two functions:

1. `People` which performs a search and returns a preview of each person.
2. `RevealPeople` which takes UUIDs of people and returns complete profiles.

Note, using `People` does not consume any credits but using `RevealPeople` will
consume credit for each profile that is revealed.

The `PeopleQuery` parameter enables you to specify your search criteria. In
particular, you can specify:

1. `Q` for general keyword matching text.
2. `Location` to match a specific city or country.
3. `Company` to match a current company.
4. `Title` to match current titles.
5. `HasEmail` if you only want to find people that have email addresses.
6. `HasPhone` if you only want to find people that has phone numbers.
7. `Skills` if you are looking to match specific skills.

By default, 10 people will be returned for each page of search results. You can
specify the `Page` as part of the `PeopleQuery` if you want to access
additional pages of people.

You can filter the search results and if you want to reveal the complete details
you can do so by sending the UUIDs via `RevealPeople`. Please note, credit will
be consumed for each person that is revealed.

## Command Line Tool

The command line tool enables you to quickly test the Nymeria API.

#### Installation

You can install the command line tool with `go install`.

```bash
$ go install git.nymeria.io/nymeria.go/cmd/nymeria@v1.0.7
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
