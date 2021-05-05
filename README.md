[![Build Status](https://travis-ci.com/Distributed-Systems-CSGY9223/fk652-sl5202-final-project.svg?token=xN9AKxnyzevD59kRtRmw&branch=master)](https://travis-ci.com/Distributed-Systems-CSGY9223/fk652-sl5202-final-project)

[![Go Report Card](https://goreportcard.com/badge/github.com/sujaylokesh/NotTwitter)](https://goreportcard.com/report/github.com/sujaylokesh/NotTwitter)

# Distributed Systems: Final Project
## Our Project Details
Our project is built using the Gin-Gonic go framework, gRPCs, and etcd Raft.

Run "go test -v" in our web folder to run all the tests. 
  * Note that you need the backend, auth, and web applications running first.
  * Backend needs to be reset at the start of every test, for correctness on storage lengths and to avoid other problems.

Similarly run "go test -v" in our backend folder to see Raft tests
  * Recommended to run the Raft tests one at a time, but they also work sequentially
---
etcd is a distributed key-value store for distributed applications. In our project, we use etcd's implementation of the raft consensus algorithm to manage an available replicated log. 

We chose CoreOs' etcd over hashicorp mainly due to the fact that etcd had done a way better job in terms of documentation and was relatively easy to understand and implement. Understanding and documentation of hashicorp was not easy to read and with the lack of included examples drove us to go with etcd.

At the time of this project, Hashicorp's raft build was failing which made us lose confidence in it. 

---

Current Structure as of Part 3:
  1. Backend
     * main.go
     * repository (repository.user and repository.article)
       * repository functions are mainly the backend stuff
       * interacts with the Raft storage to read info and send proposals to add new info
         * for demo purposes, the Raft node contacted is rotated through the 3 node cluster
           * this is to show that contacting any Raft node in the cluster will work
         * data is replicated through all Raft nodes
         * followers redirect proposals to the leader
       * also retrieves article lists and sorts them chronologically
     * storage (storage.user and storage.article)
       * AS OF PART 3 these files are mainly used to define storage structs
         * Raft nodes actually store the information, but wrap around these storage file structs
       * for part 1 and 2 any information is stored in memory in these files.
         * storage.article stores tweet info
         * storage.user stores usernames, password, and maps follow relations between users
         * there is also a write lock for each storage data structure, to handle concurrent reads and writes
         * each user has their own article list, follow map, and lock for overall better and faster concurrency handling
     * raft (raft.go)
       * Implements etcd Raft to work with our storage needs
       * Defines additional functions to initialize a 3 node cluster, retry proposals, and other essential Raft functions
  2. Web
     * main.go
     * img
       * images used for front end
     * templates
       * contains the html files tha get rendered
     * middleware.auth.go
       * communicates with the auth server via gRPCs
     * routes.go
       * handles the flow of redirecting web pages, passing through auth security checks first
     * handlers (handlers.user and handlers.article)
       * handler functions define which web pages to serve up, and pass along any relevant info to load up into those pages
     * service (service.user and service.article)
       * service functions usually call backend functions located in models to retrieve info and handle user input
       * they call certain handler functions to display the appropriate web pages using the processed data
  3. Auth
     * main.go
     * auth.go
     * defines the security checks used during routing, and sets cookie logged_in status
  4. Commonpb
     * contains the protobuf3 files used for gRPCs
  5. Vendor
     * contains the necessary packages to run our project

Web, backend, and auth are 3 seperate servers.
  * web handles the client side
  * Web communicates to backend and auth via gRPCs
  * Once authorized, web will contact backend via gRPCs
                      
User and Articles are the 2 divided sections
  * User actions and web pages include log in, registration, and following/unfollowing other users
  * Article actions include making and viewing tweets in various ways
  
---
Instructions for running the project:
  1. Download all our github source code into some folder
  2. In 3 seperate terminals do the following after switching to our project folder
     * web, backend, and auth must each be ran in their own terminal
     * Run the following commands (Windows)
       * "cd web" -> "go build -o web.exe" -> "web.exe"
       * "cd backend" -> "go build -o backend.exe" -> "backend.exe"
       * "cd auth" -> "go build -o auth.exe" -> "auth.exe"
     * Run go build and the commands accordingly to your machine
       * for example in Bash it would be "go build -o web" -> "./web"
     * Alternatively you can try running the already made .exe files
     * If firewall messages pop up, you can allow it
  3. Go to localhost:8080 on a web browser
     * Wait for Raft to initialize on the backend first
       * Raft is ready when heartbeat messages keep printing on the backend terminal
        
Navigating through the website:
  * Our web pages are very simple and straightforward
  * On the index page (first page that should show up) click sign up and register for an account
    * Passwords must contain at least one  number and one uppercase and lowercase letter, and at least 8 or more characters
    * Usernames must be unique (user1, user2, user3 are reserved as test data)
  * After successful registration, you are automatically logged in
  * Use the navbar at the top to browse through our simple features
    * We only allow adding tweets, following/unfollowing others, and viewing/searching tweets for now
    * Clicking the logo on the navbar will show all tweets (mainly for testing and demo purposes)
  * Play around with it :)

---
The project for this course will be to develop (in stages), a distributed, reliable backend in support of a mildly complex (think twitter) social media application.

The development of this project is divided into 3 stages!

## Summary

The stages are split up as follows:

1. A monolithic web application, with all logic and persistence happening on process.  This process exposes a simple social media website.
1. The monolithic web application is now split into several services, with a *stateless* web server and *at least one* backend service, which persists data and communicates with the stateless web server via gRPC.
1. The backend service(s) are now stateless, persisting their state in a raft, replicated data store.

These segments are due March 25th, April 8th, and April 29th, respectively.

\newpage

## Stage 1

### Deliverable

#### Summary

A simple web application, comprised of a web server written in Go, serving html files from a single machine.  For this stage of the project, you do not need to persist data in files or a database – instead, keep everything in memory.  Do not use a database for this application.

#### Features

This application needs to have a small number of very clear functions.  The minimum set are as follows:

1. Creating an account, with username and password
1. Logging in as a given user, given username and password
1. Users can follow other users. This action must be reversible.
1. Users can create posts that are associated with their identity.
1. Users can view some feed composed only of content generated by users they follow.

Note that these operations *should* be accomplished efficiently, but *need* not be.  You won't lose points for O(n) algorithms, where an O(logn) solution is possible.  However, you're selling yourself short ;)

If you want to build a system other than this twitter example, that is totally OK.  Just speak to the professor/TA first, to get approval.  We just want to make sure that your application is of comparable difficulty and testability.

#### Tests

You *will* be expected to write tests for your functions!! If you encounter any issues with testing, don't hesitate to ping your TA!

If you want to use fancy testing frameworks (like [ginkgo](https://github.com/onsi/ginkgo)), that's awesome!  Not required, but might be fun.

To be clear we're not setting an arbitrary code coverage metric, but tests are super important, so they'll factor into your team's grade.  If you're nervous about whether your tests are good enough, speak to your TA about it.  We just want to drive home that tests are super important, especially in this phase, as these tests will be your canaries when you split off the backend in step 2.

\newpage

#### Frameworks and Vendoring

On that note, if you want to use frameworks such as buffalo or ginkgo, go ahead, but make sure you practice [vendoring them as dependencies](https://github.com/golang/go/wiki/modules)!  It'll make it easier on your teammates, and easier for me to grade!  That being said, you definitely will not *need* to use any of these frameworks.  As our system is relatively small, they won't save you a ton of time.  It's also worth learning how stl packages like `net/http` and `testing` work, as that's what all of these frameworks are built on.

#### Structure

As friendly advice, you should try to structure your application as follows.  You don't *need* to, but it will definitely make stage 2 easier! This is because you will be decoupling your modules in the second phase, and if you don't structure your project well, like for example, dump the whole code in one or two files, it'll be hard to decouple the modules later. 

        cmd
        |-- web
        |   `-- web.go --> build target (the code that actually runs the server)
        web
        |-- web.go
        |-- config.go --> config for web, determines port to bind to, etc.
        |-- cmd       --> implementation of build target
        `-- auth
            |-- auth.go   --> authentication module, creates a new auth object,
            |                 starts it, stops it.
            |-- config.go --> token validity duration, hash difficulties, etc.
            |-- errors.go --> errors that auth object can reply with
            |-- *.go      --> implement auth; use contexts, generated
            |                 proto, and storage packages
            |-- *_test.go --> don't forget tests!
            |-- storage
            |   |-- storage.go --> storage interface for auth
            |   `-- memory
            |       |-- config.go
            |       `-- memory.go --> *threadsafe* implementation,
                                      using maps and lists
            `-- authpb
                `-- models.proto --> used to generate all data
                                     primitives used by auth module

The auth package is enumerated as an example.  All of auth's sister packages should look pretty similar, and `web` should compose these packages, hiding access to them behind http handlers, as was shown in class.

If you want, you can totally abandon this structure and follow your own path!! It's totally up to you, but please pick something sensible!! :) 

\newpage

### Resources

The following resources might be helpful:

1. This [open source book](https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/), which explains how to build a go web server from fundamentals.  Note that there is way more information in this link than you'll need, but it's a great overview.  Some of the snippets in Chapter 3 might be especially useful.
1. This [blog post](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html) explains a popular method of organizing projects, which is quite similar to the one your TA demonstrated in-class!
1. ~~This [exploration of Kubernetes source](https://developer.ibm.com/opentech/2017/06/21/tour-kubernetes-source-code-part-one-kubectl-api-server/), which your TA found on the day he was writing this doc, is really cool, and gives you an intimate sense of how the project works.~~
1. Since the above link doesn't work anymore, if you want to learn about Kubernetes, use [this link](https://jvns.ca/blog/2017/06/04/learning-about-kubernetes/).
1. This [example travis script](https://gist.github.com/y0ssar1an/df2dab474520c4086926f672c52db139) for those of you who want to play around with automated builds and tests using Travis!
1. For leveraging contexts, check out [this blog post](https://blog.golang.org/context), which goes over the subject way better than I could!!
1. This might be a bit late to slide in, but [this talk](https://blog.golang.org/advanced-go-concurrency-patterns) on concurrency patterns is great.a
1. For faking http responses, [this package](https://github.com/jarcoal/httpmock) is great.

\newpage

## Stage 2

### Deliverable

In the next stage, you will split off your backend into a separate service.  It can be a monolith (i.e. one service that performs all of the tasks related to your system's state management), or a set of services that each perform one small task, as described in-class.

This service (set) must communicate with the web server using an RPC Framework like gRPC (strongly recommended!) or Thrift.  If you run into any issues getting these frameworks set up, don't hesitate to contact the TA!

At this point, your webpage service (i.e. the service that is rendering templates / serving web pages) should be *stateless*.  In practical terms, you should be able to horizontally scale your web service without causing any concurrency problems on the backend, or UX problems on the front end.  The Web Service should no longer be persisting any data (even session data!!), but instead be fully reliant on the backend service.  If this isn't the case, you still have work to do :)

Depending on how you solved stage 1, this could be REALLY challenging -- you might need to totally re-write your service, in some cases.  If you followed the model your TA showed in-class, however, this might be the least time-intensive stage.  If you have extra time and want a challenge, you might want to implement more features, or split the backend into more distinct services.

### Resources

The following resources might be helpful:

1. [This document](https://grpc.io/docs/quickstart/go.html) is a super concise explanation of installing gRPC and protobufs (the transport serialization used by gRPC).  It also references an example, which is a great resource in and of itself for learning protobuf syntax.  In particular, it contains snippets for creating, registering, and dialling with gRPC servers/clients, which you can directly use in the project.
1. Review the example rpc slides, and the accompanying repo, which are posted on piazza
1. Adam's [gen proto script](https://raw.githubusercontent.com/adamsanghera/go-websub/master/gen-proto-go.sh) :)
1. Your TA! If you have any issues!

## Stage 3

### Summary

In stage 3, we are finally going to give our system a much-needed upgrade -- persistence!!  And not just any kind of persistence -- we're talking about highly-available, replicated persistence!

To accomplish this, we're going to use one of the two popular open source Raft implementations out there, either the one from [CoreOS](https://github.com/etcd-io/etcd/tree/master/raft) or [Hashicorp](https://github.com/hashicorp/raft).  You can choose to use whichever one you prefer; in your final presentation, you to give some rationale for why you picked the one you did! :)

### Deliverable

This is a relatively small step, but it requires learning how your raft implmentation of choice looks like!! It should end up being more of an ops problem than a dev problem.  Again, if you've followed the example architecture, it shouldn't be a huge headache.  If you find that it is, contact your TA!!

You'll need to:

1. Write code that stands up raft nodes (on your local machine is ok!!)
1. Replace your tried and tested storage implementations with a new one, which wraps a raft client for your raft implementation of choice.  This will live alongside your old implementation in the `storage` directory, if you're following the recommended structure.
   1. This involves code that establishes a connection with your raft cluster.  Now, you'll find that your storage implementation's `config`, `start` and `shutdown` actually do some heavy lifting!!

### Details

Each service should have its *own client* wrapper for contacting the raft cluster, but it is totally OK if they are all targeting the same cluster.  If you want, you can target multiple clusters -- good luck getting all those to run on your laptop, though :).

Don't worry about hosting this solution anywhere other than your laptop.  If you want to host it on AWS, heroku, or what have you (just to show off!), that's totally OK.  Just make sure that you have a configuration that will also work on your local machine.  One hack for accomplishing this in the recommended structure, is to build out `cmd` directory as follows:

\newpage

        cmd
        |-- local
        |   |-- web
        |   `-- backend
        `-- aws
            |-- web
            `-- backend

In the above example, your config initializations under `local` would look very different from your config inits under `prod`.  Alternatively, you can read config objects values from environment variables, have different env-setting shell scripts for aws/local, and keep a single `cmd` folder.  It's really up to you!

