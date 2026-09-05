# Gradle

Uses the [versions API](https://services.gradle.org/versions/all) to list distributions
and downloads the binary-only (`-bin.zip`) distribution
from [services.gradle.org/distributions](https://services.gradle.org/distributions/).

Only final releases are listed and installable: snapshots, nightlies,
release candidates and milestones are skipped.

Checksum verification (`--verify-checksum`) compares the downloaded archive with the SHA-256
published by the versions API, falling back to the `.sha256` file published next to the distribution
when the API was not used.
