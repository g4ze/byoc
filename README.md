# This software enables you to deploy your containers with a single click.

## Endpoints:
| Endpoint | Description |
|----------|-------------|
| /api/v1/register | Register New User |
| /api/v1/login | Login of a User |
| /api/v1/createcluster | Create a cluster for user |

## Database setup etc
### in the pkg/database dir
setup the db `docker run --name pgsql-dev -e POSTGRES_PASSWORD=Welcome -p 5432:5432 postgres:latest`
First time cmd:
[run db push to synchronize your schema with your database. It will also create the database if it doesn't exist.]
`go run github.com/steebchen/prisma-client-go db push`

Everytime model is changed, migrate your database and re-generate your prisma code:
`go run github.com/steebchen/prisma-client-go migrate dev --name add_comment_model`


`mkdir data`
`
docker run -e POSTGRES_PASSWORD=mysecretpassword  -p 5432:5432 -v ./data:/var/lib/postgresql/data --name byoc-postgres -d postgres
`
