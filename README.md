# Design decisions
- `No frameworks` are used to make the service for example gin. Mainly `databse drivers` and lightweight http router `chi` is used 

- `Seperation of concern` is kept in mind and thus different componenets in this service can be even further breakdown into multiple microservices. Mainly `ports and adapter pattern` is used and as some compoenets abide a interface, I can add `onion layers` on top of that (basicaaly `dependency injection`)

- `Structured logging` is used to log the events and errors.

- Postgres is used for storing `users`, `followers` and `following` data because it is relational data

- MongoDB is used for storing `stories` and `timeline` data, mainly because the stories structure is not fixed and it can be changed in future when adding more features. The timeline is present in mongodb as all the stories are present in mongodb and it is easy to query the stories and show it to the user (see next point)

- After reading about snapchat snaps and how they built snapDB, `redis(storing timeline data) along with mongodb(storiing stories)` can be a good combination.


# Steps to run
- `docker-compose up` to start the services
- `docker-compose ps` to check health status of the services
- Use this [postman](https://www.postman.com/mission-physicist-26981670/workspace/instagram-stories) collection to test the service
