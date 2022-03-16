# Movetain

> **JBA Blockchain Hackathon 2022 Spring**
>  第2期 Solana  
> 参加作品 「[Movetain](https://twitter.com/movetain)」  
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
    We use [*OAuth 2.0 Authorization Code Flow with PKCE*](https://developer.twitter.com/en/docs/authentication/oauth-2-0/authorization-code)'s RefreshToken, so you need to get it in advance.

### 1. Clone this repository.
    ```bash
    $ git clone https://github.com/shmn7iii/movetain.git
    ```

### 2. Set secrets.  
  - [secrets/keys.json](/secrets)
      ```json
      {
        "clientId": "<Client ID>",
        "clientIdSecret": "<Client ID Secret>",
        "botUserId": "<BOT's User ID>",
        "feePayerBase58": "<Base58 of FeePayer's KeyPair>"
      }
      ```

  - [secrets/refreshtoken](/secrets)
      ```text
      <Refresh Token>
      ```

### 3. Build and Run!
    ```bash
    $ go build -o ./bin/main 
    $ ./bin/main
    ```
