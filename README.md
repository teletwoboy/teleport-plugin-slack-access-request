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
- [Installation with ArgoCD](#installation-with-argocd)
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

<br>

## Installation with ArgoCD

설치 이전에 필요한 것들
   - Kubernetes
   - ArgoCD
   - `Teleport OpenSource Cluster` and `Teleport Operator` in k8s
   - Teleport User 생성 권한
   - Slack Admin 권한
   - Two Slack User (만약 혼자 테스트한다면, 당신은 Requester이자 Reviewer입니다)
   - 플러그인 서버용 도메인

<details>
<summary> Teleport Tbot 설정하기 </summary>

  - Teleport Bot이 impersonate 할 Role 정의
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
    - 이 Rule들은 tbot에게 필요한 최소한의 권한들입니다.

  - CLI를 통해 Tbot 생성
    ```
    kubectl exec -it -n teleport deploy/teleport-auth -- \
      tctl bots add access-plugin --roles=access-plugin
    ```

  - `The bot token` 값 복사 후 보관하기
    
    <img width="336" height="20" alt="스크린샷 2025-08-08 오전 1 38 27" src="https://github.com/user-attachments/assets/aad7880b-f2fd-4d6f-a754-5db526010a28" />

  - Tbot Server 배포
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
        namespace: <배포될 Namespace>
      source:
        repoURL: https://charts.releases.teleport.dev
        targetRevision: <Teleport Cluster 서버와 동일한 버전>
        chart: tbot
        helm:
          values: |-
            clusterName: <Teleport Cluster 배포시 사용된 clusterName>
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
                  name: <생성할 Kubernetes Secret Name>
                roles:
                  - access-plugin
    ```

  - Secret 생성 확인하기
    ```
    kubectl get secrets -n teleport
    kubectl describe secret <생성된 Secret> -n teleport
    ```

    <img width="271" height="17" alt="image" src="https://github.com/user-attachments/assets/29dfd5df-e177-4ec5-962b-eb7e3bd46ed5" />

    - `identity Data` 값이 채워져야 정상입니다
    - 만약 제대로 수행되지 않는다면, `tbot pod log`를 확인해주세요.
  
</details>

<details>
<summary> Teleport 유저 생성 및 특정 Role에 대한 요청 권한 주기 </summary>

  - Access Request를 통해 얻고 싶은 Role 정의
    ```
    apiVersion: resources.teleport.dev/v1
    kind: TeleportRoleV7
    metadata:
      name: kubernetes-read-only-role
      namespace: teleport
    spec:
      allow:
        kubernetes_labels:
          "*": "*"
        kubernetes_groups:
          - "<쿠버네티스에서 사용자가 가지는 역할>"
          # 이것은 Kubernetes RBAC와 관련되어 있습니다
          # 직접 설정해주세요
      options:
        max_session_ttl: 30m
    ```
   
  - 해당 Role을 가진 Access Request 요청 가능한 Role 정의
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
      options:
        max_session_ttl: 1h
    ```

  - 요청 가능한 Role을 가진 User 정의
    ```
    apiVersion: resources.teleport.dev/v2
    kind: TeleportUser
    metadata:
      name: <slack user email username>
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
<summary> 슬랙 앱 설정하기 </summary>

  - 슬랙 앱 추가하기

    - [api.slack.com/apps](https://api.slack.com/apps) 접속

    - Create New App 버튼 클릭

      <img width="154" height="52" alt="Image" src="https://github.com/user-attachments/assets/dfd7a150-123c-48c8-b3f0-d7cd6d54422b" /> 

    - From scratch 클릭

      <img width="492" height="141" alt="Image" src="https://github.com/user-attachments/assets/132be0ea-fda8-49eb-9dca-1c8fa59d569c" />

    - App Name 작성 및 App이 위치할 워크스페이스 선택 및 Create App 버튼 클릭

      <img width="496" height="477" alt="image" src="https://github.com/user-attachments/assets/4cd63a79-2d3b-4eb2-a667-9ca82c549b53" />

    - 만들면 바로 해당 App에 대한 Information 페이지로 들어가지는데,

      - Basic Information 에서 `Signing Secret` 복사하여 보관하기

        <img width="638" height="118" alt="image" src="https://github.com/user-attachments/assets/727596a3-dcc4-49d4-8f1a-bc4a905ac4b6" />

      - 좌측 사이드 바에서 `OAuth & Permissions` 클릭 후 아래로 스크롤 하여 `Bot Token Scopes` 에서 아래 권한들 추가하기

        <img width="645" height="300" alt="image" src="https://github.com/user-attachments/assets/9663146c-d184-4b11-a627-e584dc1fce89" />

        - `channels:read`
        - `chat:write`
        - `commands`
        - `groups:read`
        - `pins:write`
        - `users.profile:read`
        - `users:read`
        - `users:read.email`

      - 위로 스크롤 하여 Install to <Workspace> 버튼 클릭하여 워크스페이스에 추가하기

        <img width="639" height="172" alt="스크린샷 2025-08-08 오전 12 14 18" src="https://github.com/user-attachments/assets/dfa4302f-1be5-4e50-8a75-2581fb4df47f" />
      - 추가 후 `Bot User OAuth Token` 복사하여 보관하기

        <img width="446" height="230" alt="스크린샷 2025-08-08 오전 12 57 59" src="https://github.com/user-attachments/assets/ccdc7eab-c832-4b89-8d06-2f65d9944e25" />

  - 슬랙 커맨드 추가하기

    - 좌측 사이드 바에서 `Slash Commands` 클릭 후 Create New Command 버튼 클릭

      <img width="711" height="289" alt="image" src="https://github.com/user-attachments/assets/2bd38cb1-0ea9-4267-80de-7a8a609b2601" />

    - 빈칸 채우고 `Save` 버튼 클릭

      - Command : 사용자가 슬랙 채팅으로 입력할 커맨드

        - `/access-request`
        - `/access-policy`
    
      - Request URL : 커맨드 입력 시 요청을 전송할 URL

        - `https://<플러그인 서버용 도메인>/api/v1/access-request`
        - `https://<플러그인 서버용 도메인>/api/v1/access-policy`

      - Short Description : 커맨드에 대한 간략한 설명

  - 슬랙 Interaction 추가하기

    - 좌측 사이드 바에서 `Interactivity & Shortcuts` 클릭

    - off 버튼을 on으로 변경

      <img width="453" height="134" alt="image" src="https://github.com/user-attachments/assets/9b870bcd-3d82-42f9-9958-691375a1266b" />

    - Request URL에서 아래 주소 넣고 `Save Changes` 버튼 클릭

      - `https://<플러그인 서버용 도메인>/api/v1/interaction`
      
</details>

<details>
<summary> 슬랙 채널 설정하기 </summary>

  - 슬랙 채널 생성하기
    ```
    1. 플러그인 서버 기본 알림 전용 채널
    1. dev-role-requester
    2. dev-role-reviewers
    ```

  - 위의 채널 모두에 생성한 앱 추가하기
    ```
    /invite @<생성한 App Name>
    ```

  - `dev-role-reviewers` 채널에 Reviewer 추가하기
    ```
    /invite @<Reviewer>
    ```

  - `플러그인 서버 기본 알림 전용 채널`의 `ID 값` 복사 후 보관하기

    - 슬랙에서 해당 채널 선택
    - 메시지 상단의 `# 채널이름` 클릭
    - 하단의 `채널 ID` 복사
  
</details>

이제 복사하여 보관한 내용으로 Helm Values.yaml을 작성 후, <br>
플러그인 서버를 설치합니다. <br>

<details>
<summary> 플러그인 서버 설치하기 </summary>

  - Argocd application.yaml 작성하기
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
        namespace: <배포할 Namespace>
      source:
        repoURL: https://github.com/teletwoboy/charts.git
        targetRevision: 0.1.0
        path: slack-access-request
        helm:
          values: |-
            server:
              port: <Server Container Port>
              secret:
                slackToken: <복사한 Bot User OAuth Token>
                slackSigningSecret: <복사한 Signing Secret>
                slackDefaultNotifChannelID: <복사한 플러그인 서버 기본 알림 전용 채널의 ID>
                teleportAddress: <텔레포트 클러스터 서버 주소>
              
              teleport:
                identity:
                  secretName: <Tbot 서버로부터 생성된 Tbot Secret 이름>
    
            postgresql:
              auth:
                postgresPassword: <어드민 유저 비밀번호>
                username: <생성할 유저 이름>
                password: <생성할 유저 패스워드>
                database: <데이터베이스 이름>
              
              primary:
                service:
                  type: ClusterIP
    ```

    만약 Secret을 통해 중요 정보를 넘기거나, <br>
    인그리스 설정을 values.yaml에서 설정하고 싶거나, <br>
    다른 설정들이 궁금하다면 [Chart's Values.yaml](https://github.com/teletwoboy/charts/blob/main/slack-access-request/values.yaml) 파일을 확인해주세요.

  - 인그리스 설정은 따로 설명드리지 않습니다

</details>
