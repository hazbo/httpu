# httpu

httpu is a terminal first, general purpose HTTP client, designed to get you
interacting with all aspects of an HTTP API within just a few seconds. This
makes it a good tool for testing out various endpoints, methods, payloads - so
you are able to see what's being sent, and what response you are then getting.

![httpu](docs/demo.gif)

## Getting started

### Installing

#### macOS

```
brew tap httpu/httpu
brew install httpu
```

#### Building from source
```
cd $GOPATH
mkdir -p src/github.com/hazbo
cd src/github.com/hazbo
git clone git@github.com:hazbo/httpu.git
cd httpu
make
```

### Basic usage

I started writing this project whilst working with the [Moltin API][1], which is
a headless commerce API, so there are a few examples using that, if you want to
play around with that, just head to the website, get your API keys and you are
good to go!

Once httpu has been installed, you can get started either by first pulling down
the preconfigured packages:

```
httpu pull
```

And then loading the configuration of a package into httpu like so:

```
httpu new moltin
```

*or* by creating your own configuration. A basic project may look like the
following:

```
mkdir -p httpbin/requests
touch httpbin/project.json httpbin/requests/ip.json
```

> httpbin/project.json
```
{
  "project": {
    "url": "https://nghttp2.org/httpbin",
    "resourceFiles": [
      "httpbin/requests/ip.json"
    ]
  }
}
```

> httpbin/requests/ip.json
```
{
  "kind": "request",
  "name": "ip",
  "spec": {
    "uri": "/ip",
    "method": "GET"
  }
}
```

Once you have that setup, you're ready to run httpu!

```
httpu new httpbin
```

When the UI has loaded, typing in `ip` (the `name` of the request) into the
prompt, followed by hitting enter will run that request and you will be able to
see the response body in the right-hand window. In this case, it will simply
just return your IP address from which the request has been made.

Keybinding                              | Description
----------------------------------------|---------------------------------------
<kbd>Up</kbd>                           | Switch to command mode
<kbd>Down</kbd>                         | Switch to default mode
<kbd>Left</kbd>                         | Move cursor to request view
<kbd>Right</kbd>                        | Move cursor to response view
<kbd>Ctrl+w</kbd>                       | Move cursor from request / response view to the prompt
<kbd>Ctrl+s</kbd>                       | Switch the cursor from request view to response view
<kbd>Ctrl+c</kbd>                       | Quit

To see what commands are available, switch to command mode, then type in `list-commands`.



### Advanced usage

httpu is able to look at a JSON response, take a given value and store it in
memory, to then be used for another request. This in-memory store is called the
stash. So for example, in the previous example, if you wanted to store the IP
address that is returned, and use it else where, the request file may look like
this:

> httpbin/requests/ip.json
```
{
  "kind": "request",
  "name": "ip",
  "spec": {
    "uri": "/ip",
    "method": "GET",
    "stashValues": [
      {
        "name": "my-ip",
        "jsonPath": "origin"
      }
    ]
  }
}
```

and then in a seperate request you can then access it, once the `ip` request has
been made, like so:

> httpbin/requests/get.json
```
{
  "kind": "request",
  "name": "get",
  "spec": {
    "uri": "/get?ip=${stash[ip]}",
    "method": "GET"
  }
}
```

with `${stash[ip]}` being a variable created after running the `ip` request.

For more examples for advanced usage including the stash, sending request data,
using environment variables etc... head over to the [packages repo][2] and check
out the example I've started creating for the [Moltin API][3].

## Coming soon
  - A view to display response headers (priority)
  - Creating a request flow
  - Creating tests for a request or set of requests
  - Various UI tweaks
  - An HTTP API interface

## Known issues

This project is in its very early stages, so there will be things that need
fixing. One issue at the moment is a problem with parsing stashed variables
into request data before a request has been made. If this happens, and you are
seeing `${stash[var_name]}`, just fire the request again and it should work.

## Contributing
  - Fork httpu
  - Create a new branch (`git checkout -b my-feature`)
  - Commit your changes (`git commit`)
  - Push to your new branch (`git push origin my-feature`)
  - Create new pull request

[1]: https://moltin.com/
[2]: https://github.com/httpu/packages
[3]: https://github.com/httpu/packages/tree/master/moltin/requests
