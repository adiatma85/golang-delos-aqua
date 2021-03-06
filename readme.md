# Preparation
1. I use Go 1.17.8
2. Framework --> ``Gin`` & ``Gorm``
3. Test Framework --> Suite from ``Testify``

# Code Architecture
 ```
 internal
    - api
        - handler       (handler func for gin framework)
        - middleware    (middleware func for gin framework)
        - router        (router)
    - pkg
        - config        (app configuration)
        - db            (database configuration)
        - models        (models)
        - repository    (CRUD logic for resources)
        - validator     (struct for gin binding)
pkg
    - crypto            (for crypt, like password crypt and jwt)
    - helpers           (util)
    - response          (to standarize response to client)
```

# Instruction to Start
1. Copy ``.env.example`` to ``.env``. You can use ``cp .env.example .env``
2. Run docker storage with ``docker-compose -f docker-compose-storage.yml up -d``
3. Run ``go run ./main.go``

# List of Enpoints
    (Default at localhost:5000, but you can change the port number if you want in .env)
    - Farm
        - /api/v1/farm --> [GET] Get All Farm
            - body
                - (none)
            - expected response
                - [200] Return the list of all Farm
                - [404] If there is no anything in farms table, then it return no found
        - /api/v1/farm --> [POST]
            - body (JSON)
                - name [REQUIRED, String]
            - expected response
                - [200] Return the new created farm
                - [409] If there is another resource that already exist in storage and both of them are identical
        - /api/v1/farm/:id --> [GET]
            - body
                - (none)
            - param
                - id --> used to identify what resource that must be taken
            - expected response
                - [200] Return the instance of existed farm
                - [404] No instance exist with inserted id
        - /api/v1/farm --> [PUT]
            - body (JSON)
                - id [OPTIONAL]
                - name [OPTIONAL]
            - expected response
                - [200] This is happen when there is no existing instance yet, so we instead create new instance with payload
                - [204] This is happen when you update particular resource
        - /api/v1/farm/:id [DELETE]
            - body
                - (none)
            - param
                - id --> used to identify what resource that must be taken
            - expected response
                - [204] Return no content, but can be considered as success
                - [404] No instance exist with inserted id
    
    - Pond
        - /api/v1/pond --> [GET] Get All Pond
            - body
                - (none)
            - expected response
                - [200] Return the list of all pond
                - [404] If there is no anything in ponds table, then it return no found
        - /api/v1/pond --> [POST]
            - body (JSON)
                - name [REQUIRED, String]
            - expected response
                - [200] Return the new created pond
                - [409] If there is another resource that already exist in storage and both of them are identical
        - /api/v1/pond/:id --> [GET]
            - body
                - (none)
            - param
                - id --> used to identify what resource that must be taken
            - expected response
                - [200] Return the instance of existed pond
                - [404] No instance exist with inserted id
        - /api/v1/pond --> [PUT]
            - body (JSON)
                - id [OPTIONAL]
                - name [OPTIONAL]
            - expected response
                - [200] This is happen when there is no existing instance yet, so we instead create new instance with payload
                - [204] This is happen when you update particular resource
        - /api/v1/pond/:id [DELETE]
            - body
                - (none)
            - param
                - id --> used to identify what resource that must be taken
            - expected response
                - [204] Return no content, but can be considered as success
                - [404] No instance exist with inserted id
    
    - Record
        - /api/v1/records --> [GET]
            - body
                - (none)
            - expected response
                - [200] Return the list of all traffic records
                - [404] If there is no anything in traffic records table, then it return no found
# How to Test?
1. Run docker storage with ``docker-compose -f docker-compose-storage_test.yml up -d``
2. Rung ``go test ./test/repository`` for repository test or ``go test ./test/handler`` for handler test