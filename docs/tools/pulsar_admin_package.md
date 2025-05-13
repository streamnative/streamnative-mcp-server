#### pulsar_admin_package

Manage packages in Apache Pulsar. Packages are reusable components that can be shared across functions, sources, and sinks. The system supports package schemes including `function://`, `source://`, and `sink://` for different component types.

This tool provides operations across two resource types:

- **package** (A specific package):
  - **list**: List all versions of a specific package
    - `packageName` (string, required): Name of the package
  - **get**: Get metadata of a specific package
    - `packageName` (string, required): Name of the package
  - **update**: Update metadata of a specific package
    - `packageName` (string, required): Name of the package
    - `description` (string, required): Description of the package
    - `contact` (string, optional): Contact information for the package
    - `properties` (object, optional): Additional properties as key-value pairs
  - **delete**: Delete a specific package
    - `packageName` (string, required): Name of the package
  - **download**: Download a package to local storage
    - `packageName` (string, required): Name of the package
    - `path` (string, required): Path to download the package to
  - **upload**: Upload a package from local storage
    - `packageName` (string, required): Name of the package
    - `path` (string, required): Path to upload the package from
    - `description` (string, required): Description of the package
    - `contact` (string, optional): Contact information for the package
    - `properties` (object, optional): Additional properties as key-value pairs

- **packages** (Packages of a specific type):
  - **list**: List all packages of a specific type in a namespace
    - `type` (string, required): Package type (function, source, sink)
    - `namespace` (string, required): The namespace name 