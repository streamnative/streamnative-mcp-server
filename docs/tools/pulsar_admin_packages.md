#### pulsar_admin_package

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_package_read` and `pulsar_admin_package_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

Supported package schemes include `function://`, `source://`, and `sink://`.

### `pulsar_admin_package_read`

Read package metadata, versions, package lists, and package contents.

- **package**
  - **list**: List versions of a package
    - `packageName` (string, required): Package name
  - **get**: Get package metadata
    - `packageName` (string, required): Package name
  - **download**: Download package contents to local storage
    - `packageName` (string, required): Package name
    - `path` (string, required): Local destination path

- **packages**
  - **list**: List packages of a type in a namespace
    - `type` (string, required): Package type: `function`, `source`, or `sink`
    - `namespace` (string, required): Namespace name

### `pulsar_admin_package_write`

Manage package metadata and package contents.

- **package**
  - **update**: Update package metadata
    - `packageName` (string, required): Package name
    - `description` (string, required): Package description
    - `contact` (string, optional): Contact information
    - `properties` (object, optional): Additional properties
  - **delete**: Delete a package
    - `packageName` (string, required): Package name
  - **upload**: Upload package contents from local storage
    - `packageName` (string, required): Package name
    - `path` (string, required): Local source path
    - `description` (string, required): Package description
    - `contact` (string, optional): Contact information
    - `properties` (object, optional): Additional properties
