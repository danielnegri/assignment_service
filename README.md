Assignment Service
----------------------

**Steps to build the application using Docker**

- Edit config.yml file
    - Enter required configurations.
- Install Docker
- Run following command from root folder to build Docker image
	- ```docker build -t assignment  .```
- Run following command from root folder to run Docker image		
	- ```docker run -it -p 3000:3000 assignment```


**To run in local machine**

Compile the code by running:

```$ make build```

Start the server by running:

```$ make run```


**API Documentation**

In the web browser, load the following url:

```http://localhost:3000/swagger/index.html```