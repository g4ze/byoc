# This software enables you to deploy your containers with a single click.

Checkout all docs and implementation [here](https://g4ze.github.io/byoc/)

Read on LinkedIn [here](https://www.linkedin.com/posts/guptanilay1_bring-your-own-container-activity-7222729928209809408-VaSP?utm_source=share&utm_medium=member_desktop)

Read the tech blog [here](https://nilaygupta.hashnode.dev/bring-your-own-container)






## Database setup etc
### in the pkg/database dir
setup the db `docker run --name byoc_postgres -e POSTGRES_PASSWORD=Welcome -p 5432:5432 postgres:latest`
First time cmd:
[run db push to synchronize your schema with your database. It will also create the database if it doesn't exist.]
`go run github.com/steebchen/prisma-client-go db push`

Everytime model is changed, migrate your database and re-generate your prisma code:
`go run github.com/steebchen/prisma-client-go migrate dev --name add_comment_model`

Access the databse using:
`docker exec -it psql-dev psql -U postgres -d postgres`


`mkdir data`
`
docker run -e POSTGRES_PASSWORD=mysecretpassword  -p 5432:5432 -v ./data:/var/lib/postgresql/data --name byoc-postgres -d postgres
`
```
docker run --name byoc_postgres \
  -e POSTGRES_PASSWORD=Welcome \
  -p 5432:5432 \
  -v byoc_pg:/var/lib/postgresql/data \
  postgres:latest
```

`connect to DBpsql -h localhost -U postgres`