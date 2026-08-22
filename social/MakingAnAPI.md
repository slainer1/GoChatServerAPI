Making a RESTful API
1. To start making a RESTful API, start with your goals and your structure. You must plan meticulously about how the structure of your API is going to be.
A smart structure separates/nests each essence of a RESTful API to be structured. 
    --Example: have an API folder that holds routes and implementations, have
a migrations folder for pure migrations in the database, scripts folder to make your life easier, and modded/configuration folders (or in the root) to customize libraries

2. As much as structure matters, documentation also matters. Each heandler/method (not helper methods) must be thoroughly documented and easily accessible through curl commands.
It is absolutely crucial in a team environment to know what methods and handlers route to. Having an easily accessible UI like Swagger makes this process
much easier to debug errors.

3. Speaking of errors, you must account for errors in each method and write 2-3 edge cases. Each function should throw an error in different scenarios.


Making Routes and methods
1. First start off with your documentation. Define each element you are going to route in the documentation first when working with GOLang.
2. Write your route in the mounting method of the API. 
3. Route that method through a clean NEW program page so each element/feature stays clean
4. The language of RESTApis is usually JSON, so there is a lot of writing and reading JSON from said features
5. If needed to create a local storage to SQL inject into the database, or to hold the data there locally
6. CI/CD this stuff

Commands of REST HTTP verbs:

1. GET / https://example.com → Retrieves specific user 12.
2. POST / https://example.com → Creates a new user.
3. PUT / https://example.com → Replaces user 12 entirely.
4. PATCH / https://example.com → Partially updates user 12.
5. DELETE / https://example.com → Deletes user 12

TIPS and Rules to follow
1. Robust security and AUTH. Encrypt all data in transit and use JWT, OAuth2, rate limit, and CORS
2. Have a predictable URL structure so users can easily navigate the API: nouns not verbs, logical nesting, idempotency
3. Clear data handling and versioning, include the version in the URL path to protect production integrations, picj one camelCase or snake_case, pgination, filtering/sorting
4. Sementic status codes, 401,500,109 etc. You need to know which errors are which
5. Detailed Error payloads (a block of JSON) always return a structured JSON block explaning the failure, not raw backend stack traces
6. Structured logs and documentation (swagger)

Structured logging
1. {Key-value pair: error "," properties: ","}

Email Template Files that get injected
1. Can be -- Locally stored in the server
2. Can be -- Cloud Storage
3. Can be -- SMTP provider template builder
4. Must have a 'retries' mechanism

Caching 
An API must be built to scale if there are many objects interacting with it.
Put yourself in a better position by planning this now rather than later -- scalability. 
Caching can significant help the performance of your server with expensive computations or high freq. data
This is the redis cache
1. User goes to API and checks if it is a cached response if not then goes to database
2. returns the response to the api and caches it
3. sends response to user
4. Overall reduces load. It may not seem like it since every recent response is being cached, but if a response
is requested more frequently the response will be in the cache almost all the time, making the route faster.
5. Invalidate the cache if data is updated

Testing
"if testing your code is difficult, then using your code is difficulr too".
1. Treat your tests as a first user
2. You'll find hidden bugs you would not find
3. fast code feedback loop
4. The testing pyramid unit test -> integration tests -> end to end test
TDS
1. Transport Layer -- routes and such
2. Data Layer -- database logic and queries (local storage)
3. Service layer -- external services

CI/CD
A more recent concept, but has let developers continue push out changes without having to package code
and push it. CI/CD -- continuous integration and continuous deployment allows for community-driven
tests and packages to test linters in your code -- the grammar and vocabulary -- for production ready
code. The continuous deployment aspects include dynamic versioning and the ability to push out a change
and it being live in minutes. This is almost entirely done today with github actions
1. create a .github/workflows directory in the root of the project. 
2. Define as many YAML files to use different tests and different permissions
- Github action permission vocabulary (IMPORTANT): Read/Write: On Pull Request (PR) -- when someone makes a Pull Requests a runner is set out and tests
- On push to main -- anyone with the access to push to main triggers the runner for that specific test
- Read/Write can be dangerous for local runners because it can trigger specific custom tests someone has written exposing secrets
- Limit local runners to push branches and private repos/pull requests. NEVER have a public repo and a local runner PR test.

Conventional Commits and Semantic Versioning: commits need to have conventions and intentional messages to show a good history.
But also for semantic versioning. (Automatic versioning) follow these and the versioning will be correct.
- Fix: patches a bug in the codebase
- Feat: A new feature is added to the codebase
- Breaking Chage: A breaking API change
- other types include angular convention -- chore, ci, docs, style, refactor, perf, test for detailed history
Examples for Breaking Change: 1. (Feat!) or 2. Feat: --> Breaking Change: in the footer 3. feat(api)!: 4. 

Deploying to cloud -- google cloud
1. Set up the name and your environment -- create a budget
2. Go to cloud run
3. click service and create one
4. Click github for CD
5. Steup cloud build
Vocabulary
Payload: Typically a block of JSON, built from a class or a struct that defines what data is going to be transferred
HTTPS methods: POST, PUT, PATCH, GET, DELETE, the language used to deliver payloads
BCRYPT: The industry standard for hashing passwords. 
Middleware: A function that recieves a handler and returns a handler: Mainly used for authorization and authentication
HandlerFunc: Reterns middleware that can be used as a function (for a handler)
Plugin: if a feature is optional in prod or dev then it should be a plugin
Spies/Doubles: Record the interactions of methods with objects

