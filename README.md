# GqlAppSync

## Description

GqlAppSync is a code-generating tool built on top of [gqlgen](https://github.com/99designs/gqlgen), where the server and resolvers have been removed.
It was created to add strongly-typed models for [AWS AppSync](https://aws.amazon.com/appsync/).

## Prerequisites

1. Ensure Git is authorized to access https://github.com/nuuday/
2. Add "github.com/nuuday" to the `GOPRIVATE` environment variable.
   - powershell: `$env:GOPRIVATE = "github.com/nuuday"`
   - linux: `export GOPRIVATE=github.com/nuuday`

## Quick start

1. Get the latest version of GqlAppSync
   `go get github.com/nuuday/gqlappsync`

2. Define a `schema.graphql` file.
   E.g.:

```graphql
type Book {
  title: String!
  author: Author!
}
type Query {
  books: [Book!]!
}

type Author {
  name: String!
}
```

3. Define a `gqlgen.yml` file.
   E.g.:

```yaml
# gqlgen.yml
# Where are all the schema file located?
schema:
  - ./schema.graphql

# Where should any generated models go?
model:
  filename: ./generated/models_gen.go
  package: generated

# This section declares type mapping between the GraphQL and go type systems
#
# The first line in each type will be used as defaults for resolver arguments and
# modelgen, the others will be allowed when binding to fields. Configure them to
# your liking
models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  AWSDateTime:
    model:
      - github.com/99designs/gqlgen/graphql.Time
```

4. Generate models
   `go run github.com/nuuday/gqlappsync`
   
   The `--config 'path'` flag can be added to specify the path to the `gqlgen.yml`-filename

5. Import and use the models in the lambdas

```go
package main

import (
  "context"
  "todo/generated" // reference to the generated models
  "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) ([]generated.Book, error) {
  books := []generated.Book{
    {
      Title: "Clean Code",
      Author: &generated.Author{
        Name: "Robert Cecil Martin",
      },
    },
    {
      Title: "The Pragmatic Programmer",
      Author: &generated.Author{
        Name: "Andy Hunt and Dave Thomas",
      },
    },
  }
  return books, nil
}

func main() {
  lambda.Start(handler)
}
```

## Working with interfaces and unions in AppSync and Go

Appsync can't differentiate between interface implementations or unions unless the __typename field is provided with the name of the type.
Therefore all graphql types that implement an interface or are part of a union will have a Typename field generated with a `json:"__typename"`-tag. Currently, Go doesn't support custom default values, so the Typename field has to be assigned through the `SetTypenameRecursively(x)`method.

```go
package main

import (
  "context"
  "todo/generated" // reference to the generated models
  "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) ([]generated.Book, error) {
  books := []generated.Book{ // Book is an interface
    generated.TextBook{// Textbook is an implementation of the interface
      Title: "Clean Code",
      Author: &generated.Author{
        Name: "Robert Cecil Martin",
      },
      SupplementaryMaterial: []generated.MediaItem{
        generated.AudioClip{
          Duration: 120,
        }
      },
    }
  }
  return generated.SetTypenameRecursively(books), nil // This is required for Appsync to know which type is being returned.
}

func main() {
  lambda.Start(handler)
}
```

### Tip: Middleware

If the handler returns an interface that is implemented by a number of types that each would require invoking the `SetTypenameRecursively(x)`method, you could instead move the invocation to a lambda middleware.

```go
func handler(ctx context.Context) (generated.Book, error) {
  if foo {
    return generated.Textbook{...}, nil
  }
  return generated.CookingBook{...}, nil
}
type handlerFunc func(context.Context) (generated.Book, error)

func SetTypename(f handlerFunc) handlerFunc {
  return func(ctx context.Context) (generated.Book, error) {
    response, err := f(ctx)
    response = generated.SetTypenameRecursively(response) // The invocation
    return response, err
  }
}

func main() {
  lambda.Start(SetTypename(handler))
}
```
