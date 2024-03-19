<p align="center">
  <img alt="logo" height="114px" src="https://github.com/krystal/go-katapult/raw/main/img/logo.svg" />
</p>

<h1 align="center">
  go-katapult
</h1>

<p align="center">
  <strong>
    Go client library for <a href="https://katapult.io">Katapult</a>.
  </strong>
</h4>

<p align="center">
  <a href="https://github.com/krystal/go-katapult/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/krystal/go-katapult/ci.yml?logo=github" alt="Actions Status">
  </a>
  <a href="https://codeclimate.com/github/krystal/go-katapult">
    <img src="https://img.shields.io/codeclimate/coverage/krystal/go-katapult.svg?logo=code%20climate" alt="Coverage">
  </a>
  <a href="https://github.com/krystal/go-katapult/commits/main">
    <img src="https://img.shields.io/github/last-commit/krystal/go-katapult.svg?style=flat&logo=github&logoColor=white"
alt="GitHub last commit">
  </a>
  <a href="https://github.com/krystal/go-katapult/issues">
    <img src="https://img.shields.io/github/issues-raw/krystal/go-katapult.svg?style=flat&logo=github&logoColor=white"
alt="GitHub issues">
  </a>
  <a href="https://github.com/krystal/go-katapult/pulls">
    <img src="https://img.shields.io/github/issues-pr-raw/krystal/go-katapult.svg?style=flat&logo=github&logoColor=white" alt="GitHub pull requests">
  </a>
  <a href="https://github.com/krystal/go-katapult/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/krystal/go-katapult.svg?style=flat" alt="License Status">
  </a>
</p>

---

**WARNING:** Work in progress; features are missing, there will be breaking
changes.

---

Documentation:

- [API Documentation](https://developers.katapult.io/api/docs/latest/)
- [Katapult Documentation](https://docs.katapult.io/)



# Next Client 
A more feature complete client is being generated in the `next` package.
The aim for this client is to be generated from an openapi spec and should 
offer access to everything that is documented / exposed in our API documentation.

## Usage Guidance

Each endpoint has multiple functions for calling it. 
Typically `FunctionName` and `FunctionNameWithResponse` are provided.

It is recommended to use the `FunctionNameWithResponse` functions as they
return a response object that contains the response data and the HTTP response
object.

The `FunctionName` functions are provided for convenience and return only the
response data.

## Example

```go
res, err := client.GetDataCenterDefaultNetworkWithResponse(ctx,
	&katapult.GetDataCenterDefaultNetworkParams{
		DataCenterPermalink: "perma-link",
	},
)
```


