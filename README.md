# SaltSec

Security in E-Business Systems course project

## Table of Contents

- [SaltSec](#saltsec)
  * [Requirements](#requirements)
  * [Setup guide](#setup-guide)
    + [Linux / Mac](#linux---mac)
  * [Use](#use)
  * [Members](#members)

## Requirements

- Go
- node.js 14.0.0 or newer
- make

## Setup guide

### Linux / Mac

Setup above listed requirements depending on your OS / distro  

**Node**

For Node.js it is recommended to use the following [nvm] https://github.com/nvm-sh/nvm

Once nvm is installed, do the following:

```bash
# Install node 14.0.0 and set it as default
nvm install 14.0.0
nvm use 14.0.0
```

**Go**

For Go, follow the official [guide](https://golang.org/doc/install) as per your OS.

**make**

Make should come pre-installed on your OS / distro. If not, install required dev-tools as per your OS.

## Use

When running for the first time:

```bash
# Install required libraries
make install
```

After that, to run in development mode:

```bash
make dev
```

To run tests:

```bash
make test
```

## Members


| Member | Student No. |
| - | - |
| Pekez Marko | RA18/2017 |
| Farkaš Kristian | RA7/2017 |
| Knežević Milan | RA9/2017 |
| Ivošević Jovan | RA30/2017 |
