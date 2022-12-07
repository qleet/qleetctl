# Developer Docs

## Subcommand Code Generation

There is a lot of common code between subcommands.  This utility will scaffold out
the common stuff so you can begin writing logic right away.

### Prerequisites

The following examples presuppose there related additions have been made to the
threeport-rest-api and threeport-go-client projects:

* a new `Widget` object in the API
* a new `CreateWidget` function in the go client

### Build

First you need to build the subcommand tool from the source that lives in
`codegen/subcommand`:

```bash
make codegen-subcommand
```

### How To

Now you can use it to add source code for qleetctl.

For example, if you want to use qleetctl to create a new object called a "widget,"
you will run:

```bash
./subcommand create widget -l
```

This will add a new file `cmd/create_widget.go`.  It will have the scaffolding
to add the `widget` subcommand to `qleetctl create` so that when you're finished
users will be able to run `qleetctl create widget --config /path/to/widget.yaml`.

### Config Package

Next, you will need to add a `WidgetConfig` object to the `internal/config`
package.  This will define the attributes in the `widget.yaml` config file.
For example, create a new file:

```go
// internal/config/widget.go
package config

import (
	"encoding/json"
	"io/ioutil"

	tpclient "github.com/threeport/threeport-go-client"
	tpapi "github.com/threeport/threeport-rest-api/pkg/api/v0"
)

type WidgetConfig struct {
	Name      string `yaml:"Name"`
    Sprockets int    `yaml:"Sprockets"`
}

func (wc *WidgetConfig) Create() (*tpapi.Widget, error) {
	// construct widget object
	widget := &tpapi.Widget{        // assumes a new Widget object has been created in API
		Name:      &wc.Name,
		Sprockets: &wc.Sprockets,
	}

	// create widget widget in API
	wcJSON, err := json.Marshal(&widget)
	if err != nil {
		return nil, err
	}
    // assumes the CreateWidget function has been added to the go client
    // http://localhost:1323 API endpoint is hardcoded for local dev in this example
	wc, err := tpclient.CreateWidget(wcJSON, "http://localhost:1323", "")
	if err != nil {
		return nil, err
	}

	return wc, nil
}
```
