<p align="center">
  <img src="https://github.com/user-attachments/assets/fab8d00e-8d8c-41ef-9ea0-b927a6db29d9" width="300"/>
</p>

<div align=center>

# teleport-plugin-slack-access-request

[![License](https://img.shields.io/badge/License-Apache%202.0-%234b5563.svg?style=flat-square)](https://www.apache.org/licenses/LICENSE-2.0)
[![MADE BY](https://img.shields.io/badge/made%20by-teletwoboy-informational?style=flat-square)](https://github.com/teletwoboy)
![Version](https://img.shields.io/badge/Version-0.1.0-success?style=flat-square)
![Go Version](https://img.shields.io/badge/Go-v1.24-00ADD8?logo=go&logoColor=white&labelColor=2c2c2c)

</div>

<br>
<br>

## :memo: Table of Contents

- [What is this?](#what-is-this)
- [Why Use This Plugin?](#why-use-this-plugin)
- [Installation](#installation)
- [How to Contribute](#how-to-contribute)
- [Directory Structure](#directory-structure)
- [License](#license)
- [Reference](#reference)

<br>

## What is This?

`teleport-plugin-slack-access-request` is <br> 
`Go-based Teleport plugin server` developed by [teletwoboy](https://github.com/teletwoboy). <br>

It provides convenient `Slack-based features` <br>
integrated with the [open-source version of Teleport](https://github.com/gravitational/teleport).

Find us also at [Dockerhub](https://hub.docker.com/r/springboothate/teleport-plugin-slack-access-request)
> Latest version : `0.1.0`

<br>

## Why Use This Plugin?

   1. Creating and reviewing Access Requests entirely through Slack

      - Neither the requester nor the reviewer needs to log in to the Teleport Web UI.
      - All request and approval records are stored in the plugin server's database, Slack messages, and Teleport Server, making `it easy to track and audit the history`
      
   2. Only users in the designated Reviewers Slack channel for a given Role can review requests

      - You can `manage reviewers more clearly` by separating reviewable roles based on each channel
      
   3. Defining automatic approval rules for specific users, roles, or channels

      - Requests that don't require manual review every time `can be automatically approved`


