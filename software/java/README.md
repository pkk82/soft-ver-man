# JAVA

Use API to download **zulu** packages from **Azul Systems**.

API is described in [Swagger](https://api.azul.com/metadata/v1/docs/swagger).

## Maven toolchains

Installing a JDK registers it in `~/.m2/toolchains.xml` (the file is created if missing)
as a `jdk` toolchain with the major version and vendor `zulu`.
An existing `jdk` entry with the same major version is overridden.
Uninstalling repoints the entry to the newest remaining installation
with the same major version, or removes the entry when none is left.
Other toolchain entries are preserved as they are.

