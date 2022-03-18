# Movetain

> **JBA Blockchain Hackathon 2022 Spring**  
> 第 2 期 Solana  
> 作品名「[Movetain](https://twitter.com/movetain)」  
> チーム「明太子のシャチホコ」 (
> [shmn7iii](https://github.com/shmn7iii),
> [motoha0827](https://github.com/motoha0827),
> [WATANA-be](https://github.com/WATANA-be)
> )

# About

Movetain ([@movetain](https://twitter.com/movetain)) is a Twitter BOT which mints a
**"Tweet NFT"** on Solana devnet.

# Set up

### 0. Set up Twitter API and Create Twitter App.

BOT requests "OAuth 2.0 User Context" authentication.  
 We use [_OAuth 2.0 Authorization Code Flow with PKCE_](https://developer.twitter.com/en/docs/authentication/oauth-2-0/authorization-code) 's RefreshToken, so you need to get it in advance.

### 1. Set up Ubuntu.

> Env: Azure Virtual Machine, Ubuntu 20.04.4

```bash
# dependency
$ sudo apt-get update
$ sudo apt-get upgrade -y
$ sudo apt-get install -y git vim curl

# docker
$ curl -fsSL https://get.docker.com/ | sh
$ sudo systemctl start docker

# docker-compose
$ sudo curl -L "https://github.com/docker/compose/releases/download/1.26.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose
$ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

# clone
$ git clone https://github.com/shmn7iii/movetain.git
```

### 2. Set secrets.

```bash
$ pwd
<PROJECT_DIRECTORY>

$ sftp -i ~/.ssh/xxx.pem user@ip
sftp> cd movetain
sftp> puts -f ./secrets
sftp> exit
```

- **secrets/keys.json**

  ```json
  {
    "clientId": "<Client ID>",
    "clientIdSecret": "<Client ID Secret>",
    "botUserId": "<BOT's User ID>",
    "feePayerBase58": "<Base58 of FeePayer's KeyPair>"
  }
  ```

- **secrets/refreshtoken**
  ```text
  <Refresh Token>
  ```

### 3. Build and Run!

```bash
$ cd movetain
$ docker-compose build
$ docker-compose up -d

# follow log
$ docker-compose logs -f movetain
```
