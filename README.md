# goexpert-rate-limiter <img src="https://www.svgrepo.com/show/353830/gopher.svg" width="40" height="40">

------
## Challenge description
### Objective
Develop a rate limiter in Go that can be configured to limit the maximum number of requests per second based on a specific IP address or an access token.

### Description
The goal of this challenge is to create a rate limiter in Go that can be used to control the traffic of requests to a web service. The rate limiter should be able to limit the number of requests based on two criteria:

- **IP Address:** The rate limiter should restrict the number of requests received from a single IP address within a defined time interval.
- **Access Token:** The rate limiter should also be able to limit requests based on a unique access token, allowing different expiration time limits for different tokens. The Token should be provided in the header in the following format:
  `API_KEY: <TOKEN>`
  The access token limit settings should override those of the IP. For example, if the limit per IP is 10 req/s and that of a specific token is 100 req/s, the rate limiter should use the token information.

### Requirements
- The rate limiter should be able to work as a middleware that is injected into the web server.
- The rate limiter should allow the configuration of the maximum number of requests allowed per second.
- The rate limiter should have the option to choose the blocking time for the IP or Token if the number of requests has been exceeded.
- Limit settings should be done via environment variables or in a ".env" file in the root folder.
- It should be possible to configure the rate limiter for both IP and access token limitation.
- The system should respond appropriately when the limit is exceeded:
    - HTTP Code: 429
    - Message: you have reached the maximum number of requests or actions allowed within a certain time frame
- All limiter information should be stored and queried from a Redis database. You can use docker-compose to spin up Redis.
- Create a "strategy" that allows easily swapping Redis for another persistence mechanism.
- The limiter logic should be separated from the middleware.

### Examples
- **IP Limitation:** Suppose the rate limiter is configured to allow a maximum of 5 requests per second per IP. If IP 192.168.1.1 sends 6 requests in one second, the sixth request should be blocked.
- **Token Limitation:** If a token abc123 has a configured limit of 10 requests per second and sends 11 requests within that interval, the eleventh should be blocked.
  In both cases above, subsequent requests can only be made after the total expiration time has passed. For example, if the expiration time is 5 minutes, a specific IP can only make new requests after the 5 minutes have passed.

### Tips
- Test your rate limiter under different load conditions to ensure it functions as expected under high traffic situations.

### Delivery
- The complete source code of the implementation.
- Documentation explaining how the rate limiter works and how it can be configured.
- Automated tests demonstrating the effectiveness and robustness of the rate limiter.
- Use docker/docker-compose so we can test your application.
- The web server should respond on port 8080.

------ 
## How to run locally

### Using Go
- Clone the repository;
- Open terminal on the project folder;
- Make sure you have redis running locally in port `6379`;
- Run `go run cmd/main.go`;

### Using docker-compose
- Clone the repository;
- Run `docker-compose up --build`;

### How to run tests
- Make sure redis is running locally in port `6379`;
- Open folder `tests` in terminal and run `go test -v`;

### How to test the rate limiter
- Configure parameters on the file `configs/config.env`;
- You can specify the list of tokens and their limits;
- You can specify the limit to IPs;
- You can specify the block time to ips and tokens, if not specified, the default is 1 minute;
- After this you can run the app using Go or docker-compose
- The app has 3 test routes to validate the rate limiter:
  - `/token`: this route needs an API_KEY in the header. 
The token needs to be in the list of tokens configured in the file `configs/config.env`.
The limiter will use the token and the limit specified in the config file to the specific token to validate the access.
  - `/ip`: this route doesn't need an API_KEY in the header. The limiter will use the IP and the limit specified
in the config file to validate the access.
  - `/both`: this route has limiter configured to block both by ip and token, following the rules of the challenge.
If a token is not sent, it will use the block by ip. If a token is sent, then validates the limiter by token.
    - `TIP`:  You can try to run this route without token to validate the limit by ip and when it blocks your ip, you can send a valid
    token, to validate the limit by token.
- In any route, you can check the logs to see which `origin` is being validated: token or ip, 
and also the key, which is the specific token or your ip.
  - When the access is valid, you will se in the logs the count of access.
  - When the token or the ip is blocked, you will see in the logs the message `access is blocked` and the specific origin and key.
  - When the access is blocked, you will receive a response with status code 429 and the message 
`you have reached the maximum number of requests or actions allowed within a certain time frame`.

### Contact
- [Email](mailto:mateusmatinato@gmail.com)
- [LinkedIN](https://linkedin.com/in/mateusmatinato)
