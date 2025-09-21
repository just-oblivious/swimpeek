# SwimPeek for Swimlane Turbine

SwimPeek "connects the dots" between playbooks, components, applications, and other resources in Swimlane Turbine.

With SwimPeek you can browse relationships between resources and answer questions like:
- What playbooks are triggered when a record is created or modified in this application?
- What playbook-actions touch records in this application?
- Where is this component used?

*This tool works by requesting configuration data from a Swimlane Turbine tenant and turning it into a graph-like data structure, this graph can then be navigated in a [fancy terminal UI](https://charm.land/).*


## Limitations

- A cloud-hosted instance is assumed, on-prem deployments are not supported (I don't have access to one for testing);
- Only Turbine content is supported, legacy content is ignored;
- SwimPeek was developed with Turbine v25.3.1 in mind, there's no guarantee that this tool keeps working for newer releases;
- This tool was created by ~~reading the tea leaves~~ analyzing API responses, the output may not be 100% accurate.


## Installation & configuration

1.  Download and install Go (>=1.24.1) from https://go.dev/dl/

1.  Run Go install to download, build, and install SwimPeek:
    ```sh
    go install github.com/just-oblivious/swimpeek@latest
    ```

1. Create a Personal Access Token for your Swimlane account ([docs](https://docs.swimlane.com/docs/introduction/customize-your-user-profile.htm))

1.  Configure SwimPeek:
    ```sh
    swimpeek config
    ```

*Command not found?
Add the following line to your shell config to ensure that the Go bin directory is included in the system path:*
  ```sh
  export PATH=${PATH}:`go env GOPATH`/bin
  ```


## Usage

1.  Download the contents of a Swimlane tenant to a JSON file:
    ```sh
    swimpeek dump
    ```

1.  Launch the analyzer:
    ```sh
    swimpeek analyze -infile path_to_dump.json
    ```

Run `swimpeek cmd -help` to learn more about the usage of each subcommand


## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## Roadmap

Features to be added soon™:

- Filtering in list views
- Navigation between flow nodes and references
- Support for triggers (view all playbook triggers and track flow events)
- Support for assets (see what assets exists and where they are used)
- SUpport for connector actions (list available connectors and see where they are used)
- Configuration details for individual action nodes (i.e. the input parameters)
- Global search
