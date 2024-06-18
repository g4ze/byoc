# This software enables you to deploy your containers with a single click.

## Endpoints:
| Endpoint | Description |
|----------|-------------|
| /api/v1/register | Register New User |
| /api/v1/login | Login of a User |
| /api/v1/createcluster | Create a cluster for user |

## Database setup etc
`mkdir data`
`
docker run -e POSTGRES_PASSWORD=mysecretpassword  -p 5432:5432 -v ./data:/var/lib/postgresql/data --name byoc-postgres -d postgres
`
