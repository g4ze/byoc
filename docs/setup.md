# Setup

As our fortunes seem, the project is not deployed for public use due to shortage of finances. We support open source initiative and hence urge you to go through a few simple and easy steps to get your own instance up and running!

## Prerequisites

### AWS
- AWS secret access id and keys.
- IAM user with ECS admin policy.
- Decide on AWS Region (based on wherever you are situated)
- ECS_VPC- any VPC you make in your region.
- Subnets, based on your needs, make any amount. The project has three subnet variables in its template, add subnets based on what you desire. Make sure to change `/pkg/core/service-ops.go` ; `pkg/core/load-balancer-ops.go` and `pkg/handlers/health.go` w.r.t. how many subnets you add. 
!!! tip
    It is recommended to have three subnets if you're new to the project and don't want to tinker around.
### Docker 
- Docker cli, engine and running demon
- Docker compose
!!! note
    We've added compose for local dev and usage, suggested to use swarm for more reliable use cases. 
### Golang
Go version `go1.22.5` (not required if developing on fe)
### Frontend
Any Javascript/typescript runtime, preferrably node.

## Dev setup
For Frontend based development, it is recommended to build and run the backend using `docker compose`.

### Backend
Most straight forward to quicly start BE would be to use docker compose. For ease of setup, everything is in place.

We do need to setup all unset variables in the .env.template file.
```
cp .env.template .env
```
Now set all the required variables.


Just execute the commands below in project root folder (BYOC/)

- start the containers:
```
docker compose up
```
- for first time running, we need to push and create initial client in the db:
```
docker exec -it byoc_app bash
cd ../app/pkg/database
go run github.com/steebchen/prisma-client-go db push
```
That's it. Everything good to go now!

!!! note
    if you want to work on the BE, just start it using `go cmd/main.go` for BE dev. Do a compose build to check if eveything is in place and you're good to go.

### Frontend
The `fe` folder acts as the root flder for our frontend written in NEXTjs. It also acts as a proxy server for proxying custom paths to the deployment DNS.

set env using `cp fe/.env.template fe/.env`

Since our backend is serving on port 80 (if deployed from docker compose), do check in `fe/.env` that `NEXT_PUBLIC_BE_URL=http://localhost:80`
In the frontend `fe` folder:
```
npm i
npm run dev
```

!!! tip
    If you're not using docker compose to run the backend, backend would be served on port `2001` and hence you should either delete `NEXT_PUBLIC_BE_URL=http://localhost:80` or set it to `http://localhost:2001`
