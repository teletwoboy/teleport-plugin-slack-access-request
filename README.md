<p align="center">
  <img src="https://github.com/user-attachments/assets/fab8d00e-8d8c-41ef-9ea0-b927a6db29d9" width="300"/>
</p>

<div align=center>

# teleport-plugin-slack-access-request

[![License](https://img.shields.io/badge/License-Apache%202.0-%234b5563.svg?style=flat-square)](https://www.apache.org/licenses/LICENSE-2.0)
[![MADE BY](https://img.shields.io/badge/made%20by-teletwoboy-informational?style=flat-square)](https://github.com/teletwoboy)
![Go Version](https://img.shields.io/badge/Go-v1.24-00ADD8?logo=go&logoColor=white&labelColor=2c2c2c)

</div>

<br>
<br>

## :memo: Table of Contents

- [What is this?](#what-is-this)
- [Why Use This Plugin?](#why-use-this-plugin)
- [Installation with ArgoCD](#installation-with-argocd)
- [How to use](#how-to-use)
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

<br>

## Why Use This Plugin?

  - We have two powerful features.

  1. `/access-request`

     - Neither the requester nor the reviewer needs to log in to the Teleport Web UI. <br>
       They only need to be connected to Slack.

     - All request and review records are stored in the plugin server's database, Slack messages, and Teleport Server, making `it easy to track and audit the history`.

     - Only users in the Slack channel mapped to the requested role can review the request.
       This enables clear management of reviewers for each role.

  2. `/access-policy` : ABAC-based Access Request Auto-Review Policy
    
     - Reviewers can create auto-review policies for requests that do not require manual review every time.
       
     - Policies can be configured based on specific `channels`, `roles`, `users`, `time windows`, and whether to automatically `approve` or `deny` the request.

<br>

## Installation with ArgoCD

**Prerequisites Before Installation**
   - Kubernetes
   - ArgoCD
   - `Teleport OpenSource Cluster` and `Teleport Operator` in k8s
   - Permission to create Teleport roles, users, tbot
   - Permission to create Slack App, Channel etc 
   - Two Slack users (If testing alone, you will act as both the Requester and the Reviewer)
   - A domain name for the slack-plugin server

<details>
<summary> 1️⃣ Setting up Teleport Tbot </summary>

  <hr>

  - Define the role that the Teleport Bot will impersonate
    ```
    apiVersion: resources.teleport.dev/v1
    kind: TeleportRoleV7
    metadata:
      name: access-plugin
      namespace: teleport
    spec:
      allow:
        rules:
          - resources:
              - access_request
            verbs:
              - list
              - read
              - update
              - submit
              - create
          - resources:
              - access_plugin_data
            verbs:
              - update
          - resources:
              - user
            verbs:
              - read
              - list
          - resources:
              - role
            verbs:
              - read
              - list
          - resources:
              - event
            verbs:
              - list
              - read
          - resources:
              - user_login_state
            verbs:
              - list
              - read
              - delete
    ```
    - These rules define the minimum required permissions for the Tbot

  <hr>
  
  - Create Tbot via CLI
    ```
    kubectl exec -it -n teleport deploy/teleport-auth -- \
      tctl bots add access-plugin --roles=access-plugin
    ```

  <hr>
  
  - Copy the `bot token` value
    
    <img width="336" height="20" alt="스크린샷 2025-08-08 오전 1 38 27" src="https://github.com/user-attachments/assets/aad7880b-f2fd-4d6f-a754-5db526010a28" />

  <hr>
  
  - Deploy Tbot Server
    ```
    apiVersion: argoproj.io/v1alpha1
    kind: Application
    metadata:
      name: teleport-tbot
      namespace: argocd
    spec:
      project: <Project Name>
      destination:
        server: https://kubernetes.default.svc
        namespace: <Namespace to deploy>
      source:
        repoURL: https://charts.releases.teleport.dev
        targetRevision: <Same version as the Teleport Cluster server>
        chart: tbot
        helm:
          values: |-
            clusterName: <clusterName used during Teleport Cluster deployment>
            teleportAuthAddress: <Auth Server Service Name>.<Namespace>.svc.cluster.local:3025
    
            defaultOutput:
              enabled: false
    
            persistence: "secret"
            joinMethod: "token"
            token: <The bot token>
    
            outputs:
              - type: identity
                destination:
                  type: kubernetes_secret
                  name: <Name of the Kubernetes Secret to be created>
                roles:
                  - access-plugin
    ```

  <hr>
  
  - Check if the Secret Has Been Created
    ```
    kubectl get secrets -n teleport
    kubectl describe secret <생성된 Secret> -n teleport
    ```

    <img width="271" height="17" alt="image" src="https://github.com/user-attachments/assets/29dfd5df-e177-4ec5-962b-eb7e3bd46ed5" />

    - The `identity Data` field must be populated for it to be considered valid
    - If it doesn't work as expected, please check the `tbot pod logs`
  
</details>

<details>
<summary> 2️⃣ Creating Teleport User </summary>

  <hr>

  - Define the Role to be obtained through an Access Request
    ```
    apiVersion: resources.teleport.dev/v1
    kind: TeleportRoleV7
    metadata:
      name: dev-role
      namespace: teleport
    spec:
      allow:
        kubernetes_labels:
          "*": "*"
        kubernetes_groups:
          - "<Role the user will assume in Kubernetes>"
          # This is related to Kubernetes RBAC.
          # Please configure it manually.
      options:
        max_session_ttl: 30m
    ```
   
  <hr>
  
  - Define the Role that is allowed to request an Access Request for the target Role
    ```
    apiVersion: resources.teleport.dev/v1
    kind: TeleportRoleV7
    metadata:
      name: dev-role-requester
      namespace: teleport
    spec:
      allow:
        request:
          roles:
            - dev-role
          reason:
            mode: "optional" or "required"
            # This determines whether providing an Access Request reason is optional or required
      options:
        max_session_ttl: 1h
    ```

  <hr>
  
  - Define the User who is allowed to request the specified Role
    ```
    apiVersion: resources.teleport.dev/v2
    kind: TeleportUser
    metadata:
      name: <requester slack user email username>
      namespace: teleport
    spec:
      roles:
        - dev-role-requester
      traits:
        logins:
          - root
          - ubuntu
    ```

</details>

<details>
<summary> 3️⃣ Setting up Slack App </summary>

  <hr>

  - Add a Slack App
    
    - Go to [api.slack.com/apps](https://api.slack.com/apps) 

    - Click the `Create New App` button

      <img width="154" height="52" alt="Image" src="https://github.com/user-attachments/assets/dfd7a150-123c-48c8-b3f0-d7cd6d54422b" /> 

    - Select From scratch

      <img width="492" height="141" alt="Image" src="https://github.com/user-attachments/assets/132be0ea-fda8-49eb-9dca-1c8fa59d569c" />

    - Enter an `App Name`, select the `Workspace`, and click `Create App`

      <img width="496" height="477" alt="image" src="https://github.com/user-attachments/assets/4cd63a79-2d3b-4eb2-a667-9ca82c549b53" />

    - After creation, you’ll be redirected to the App Information page:

      - Under Basic Information, copy and securely store the Signing Secret

        <img width="638" height="118" alt="image" src="https://github.com/user-attachments/assets/727596a3-dcc4-49d4-8f1a-bc4a905ac4b6" />

      - In the left sidebar, click `OAuth & Permissions`, scroll down to the `Bot Token Scopes` section, and add the following permissions:

        <img width="645" height="300" alt="image" src="https://github.com/user-attachments/assets/9663146c-d184-4b11-a627-e584dc1fce89" />

        - `channels:read`
        - `chat:write`
        - `commands`
        - `groups:read`
        - `pins:write`
        - `users.profile:read`
        - `users:read`
        - `users:read.email`

      - Scroll up and click the `Install to <Workspace>` button to add the app to your workspace
        
        <img width="639" height="172" alt="스크린샷 2025-08-08 오전 12 14 18" src="https://github.com/user-attachments/assets/dfa4302f-1be5-4e50-8a75-2581fb4df47f" />
      - After installation, copy and securely store the `Bot User OAuth Token`

        <img width="446" height="230" alt="스크린샷 2025-08-08 오전 12 57 59" src="https://github.com/user-attachments/assets/ccdc7eab-c832-4b89-8d06-2f65d9944e25" />

  <hr>
  
  - Add Slash Commands

    - In the left sidebar, click `Slash Commands`, then click the `Create New Command` button

      <img width="711" height="289" alt="image" src="https://github.com/user-attachments/assets/2bd38cb1-0ea9-4267-80de-7a8a609b2601" />

    - Fill in the required fields and click the `Save` button:

      - Command : The command users will type in Slack chat

        - `/access-request`
        - `/access-policy`
    
      - Request URL : The endpoint that Slack will send requests to when the command is used

        - `https://<your-plugin-domain>/api/v1/access-request`
        - `https://<your-plugin-domain>/api/v1/access-policy`

      - Short Description : A brief description of the command

  <hr>
  
  - Enable Slack Interactivity

    - In the left sidebar, click `Interactivity & Shortcuts`

    - Toggle the switch from Off to On

      <img width="453" height="134" alt="image" src="https://github.com/user-attachments/assets/9b870bcd-3d82-42f9-9958-691375a1266b" />

    - In the Request URL field, enter the following and click `Save Changes` button:

      - `https://<your-plugin-domain>/api/v1/interaction`
      
</details>

<details>
<summary> 4️⃣ Setting up Slack Channel </summary>

  <hr>

  - Create Slack Channels
    
    ```
    1. Plugin server notification channel (used for general alerts)
    1. dev-role-requester # Optional, but you must add the app to at least one channel
    2. dev-role-reviewers # ‼️ Must follow the format: <Role Name + '-reviewers'> ‼️
    ```

  <hr>
  
  - Add your Slack app to each of the above channels
    ```
    /invite @<Your App Name>
    ```

  <hr>
  
  - Add reviewers to the `dev-role-reviewers` channel
    ```
    /invite @<Reviewer>
    ```

  <hr>
  
  - Copy and store the Channel ID of the `plugin server's notification channel`

    - Select the channel in Slack
    - Click the `# Channel Name` at the top of the message view
    - Scroll down and copy the `Channel ID`
  
</details>

<details>
<summary> 5️⃣ Installing the Plugin Server </summary>

  <hr>

  - **Write the `application.yaml` for ArgoCD**
    ```
    apiVersion: argoproj.io/v1alpha1
    kind: Application
    metadata:
      name: teleport-plugin-slack-access-request
      namespace: argocd
    spec:
      project: <Project Name>
      destination:
        server: https://kubernetes.default.svc
        namespace: <Namespace to deploy>
      source:
        repoURL: https://github.com/teletwoboy/teleport-plugin-slack-access-request-helm.git
        targetRevision: 0.1.0
        path: .
        helm:
          values: |-
            server:
              port: <Server Container Port>
              secret:
                slackToken: <Copied Bot User OAuth Token>
                slackSigningSecret: <Copied Signing Secret>
                slackDefaultNotifChannelID: <Copied notification channel ID>
                teleportAddress: <Teleport Cluster Server Address>
              
              teleport:
                identity:
                  secretName: <Tbot Secret generated by the Tbot Server>
    
            postgresql:
              auth:
                postgresPassword: <Admin user password>
                username: <Username to be created>
                password: <Password for the user>
                database: <Database name>
              
              primary:
                service:
                  type: ClusterIP
    ```

    If you want to pass sensitive information via Kubernetes Secrets, <br>
    configure ingress settings directly in values.yaml, <br>
    or explore other configuration options, <br>
    please refer to the [Chart's Values.yaml](https://github.com/teletwoboy/charts/blob/main/slack-access-request/values.yaml) file

  - `Note`: Ingress setup is not covered in this guide.

</details>

<br>

## How to use

<details>
<summary> Access Request </summary>

  <hr>

  - Requester

    - type `/access-request`
  
      <img width="326" height="133" alt="image" src="https://github.com/user-attachments/assets/923f5fa0-7716-4a25-a5db-5bee15c00640" />
  
    - Select `Requeted Role` and `Reviewers Channel`
  
      <img width="516" height="346" alt="image" src="https://github.com/user-attachments/assets/a7d9fe46-1ac2-4606-8966-65daaea86f23" />
  
    - Review the `summary`, fill in the `request reason`, and click the `Submit` button
      
      <img width="507" height="398" alt="image" src="https://github.com/user-attachments/assets/3239f6b9-cd57-4d07-9830-e45562cebfa4" />

    - Check the `request creation message` in the channel

      <img width="422" height="174" alt="image" src="https://github.com/user-attachments/assets/28dcf30b-e6ef-4333-88c5-e6d4b1cca48c" />

  <hr>
  
  - Reviewer (in the Reviewers Channel)

    - See the Access Request notification and click `Review Request` button

      <img width="419" height="301" alt="image" src="https://github.com/user-attachments/assets/23381625-b2db-4cbb-919f-6607d36480ad" />

    - Fill in your review and select Allow or Deny

      <img width="508" height="524" alt="image" src="https://github.com/user-attachments/assets/f6dec9e0-4ff5-4101-b929-260dc8f537c1" />

    - Check the `review result message` in the channel

      <img width="481" height="243" alt="image" src="https://github.com/user-attachments/assets/68eeb576-ae44-4be0-a1e2-ecf127f1f32b" />

  <hr>
  
  - Requester

    - Check if the request has been reviewed

      <img width="403" height="222" alt="image" src="https://github.com/user-attachments/assets/cd2125a7-2e35-43c1-93cb-7c098c50abad" />

</details>

<details>
<summary> Access Policy </summary>

  <hr>

  - Reviewer in Reviewers Channel

    - type `/access-policy`
   
      <img width="323" height="135" alt="image" src="https://github.com/user-attachments/assets/b95d7fb9-197d-42db-9213-0920339ad372" />

    - Select the target `channel`, `role`, and `user`

      <img width="511" height="445" alt="image" src="https://github.com/user-attachments/assets/01f58704-d5e6-4597-9a24-2c31f726ca5f" />

    - Choose the `start date/time`, `end date/time`, and `effect` (Allow or Deny)

      <img width="511" height="522" alt="image" src="https://github.com/user-attachments/assets/0185df86-373e-4f74-8521-e74f25e41afe" />

    - Review the `summary`, write a `title` and `reason`, then click `Submit` button

      <img width="511" height="582" alt="image" src="https://github.com/user-attachments/assets/4d4d40fc-9a2b-4a4b-a7d8-82676dfd223d" />

    - Check the `created Access Policy Message` in the channel

      <img width="397" height="244" alt="image" src="https://github.com/user-attachments/assets/b9bde956-5f22-45f3-a88c-de256b4639a7" />

      - The policy is automatically pinned for easier management.

  <hr>
  
  - Requester

    - Perform the `/access-request` process as usual

    - Confirm that the request was automatically reviewed

      <img width="327" height="185" alt="image" src="https://github.com/user-attachments/assets/d092a458-bc25-4fac-88b1-0cceaf2faa78" />

  <hr>

  - Reviewer in Reviewers Channel

    - Check the information for the `request that was auto-reviewed`

      <img width="461" height="258" alt="image" src="https://github.com/user-attachments/assets/718d7bf3-be5d-42a9-9795-64728cd2139e" />


</details>

<br>

## How to Contribute

This project was created by `university student developers`. <br>
If you find any `mistakes`, <br>
have `ideas` for new features, <br>
or `suggestions` for improvements, <br>
we welcome your `contributions`!

##### :pray: [HOW TO CONTRIBUTE](CONTRIBUTING.md)

<br>

## Directory Structure

```
.
├── db                # Database schema and migration files
├── cmd               # Application entry points (main executables)
├── internal          # Internal application packages (not for public use)
│   ├── api           # API route handlers and related logic
│   ├── app           # Application startup and lifecycle management
│   ├── config        # Configuration loading and management
│   ├── database      # Database connection and query logic
│   ├── events        # Teleport Event handling and dispatching
│   ├── logging       # Logging utilities and setup
│   ├── policy        # Access policy logic and enforcement
│   ├── seedinit      # Initial data seeding scripts
│   ├── slack         # Slack integration logic
│   ├── teleport      # Teleport integration logic
│   ├── user          # User management and related logic
│   └── util          # Utility/helper functions
```

<br>

## License

This project is licensed under the [Apache License 2.0](LICENSE).

<br>

## Reference

[Teleport Github Repository](https://github.com/gravitational/teleport?tab=readme-ov-file#support-and-contributing) <br>
[Teleport Official Website](https://goteleport.com/) <br>
[Slack API methods](https://api.slack.com/methods)

<br>

[🔝 Back to Top](#teleport-plugin-slack-access-request)
