This is a repository for a programming exercise in Go.


**Instructions for running**

In these steps, the '$' is not meant to be typed, it indicates what to type at
your terminal prompt.

1. Clone the repo on your system in your go workspace.
2. From inside the repository directory, build and install via $go install
3. Run the program as a standalone command: $jc-assignment
4. The program listens on localhost:8080 by default. The port can be changed
   via the command line flag --port.
   For example, to run on port 1234: $jc-assignment --port 1234.
5. You can run $jc-assignment -h to see the usage information.
6. You can run the unit tests: $go test

Notes:



Assumptions:
1. I'm assuming that the stats data returned by the /stats call is
   not persisted across restarts of the server.


