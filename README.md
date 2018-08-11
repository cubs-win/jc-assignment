**Overview**
This is a repository for a programming exercise in Go. The main() function is in
jc-assignment.go. The HTTP handlers are in jc-handlers.go.  

**Instructions for running**

In these steps, the '$' is not meant to be typed, it indicates what to type at
your terminal prompt.

1. Clone the repo on your system in your go workspace. $git clone https://github.com/cubs-win/jc-assignment.git
2. From inside the repository directory, build and install via $go install
3. Run the program as a standalone command: $jc-assignment
4. The program listens on localhost:8080 by default. The port can be changed
   via the command line flag --port.
   For example, to run on port 1234: $jc-assignment --port 1234.
5. You can run $jc-assignment -h to see the usage information.
6. You can run the unit tests: $go test -v
7. I added a separate executable called hasher in tests/hasher which
   you can run to test that multiple connections are processed simultaneously,
   and that issuing a shutdown does not interrupt existing work.
   To run it, cd into the test/hasher directory, do $go build, and
   execute the program with $./hasher. Before doing so, start the 
   jc-assignment server in a separate terminal window.

**Assumptions**
1. I'm assuming that the stats data returned by the /stats call is
   not persisted across restarts of the server.


