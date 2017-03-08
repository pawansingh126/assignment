# assignment

This is a demo assignment as an introduction to Golang and Redis.

Go Routines, Channels and how to work with Redis.

Frontend is the initial interface with listens to web requests.
It also invokes two *go routines* namely Bidder() and Reducer().

Frontend module, takes request form client.go. One such request is taking multipart POST data 
containing test.txt. Uploaded file typically contains different URLs in differeny lines.
There is one more GET request frontend receives which enqires about the status of previously 
sent POST request. Frontend fetches status from Redis using request id and replies accordingly.

Frontend creates a request id for the POST data and saves it with marshaled data in Redis
with request status message as "Bidding".Then it add the request id to the bidder channel.

Bidder module is a "Go Routine" which is constantly looking for request id added to bidder channel.
If it finds bidder channel not empty, it takes out that request id and retrieves all the data from
Redis. Then it hits each and every URL to get their HTTP statuscode and again save everything in Redis
with request status message as "Reducing". Once done with this, it add the request id to the reducer channel.

Reducer module is agaib a "Go Routine" which is constantly looking for request id added to reducer channel.
If it finds reducer channel not empty, it also takes out that request id and retrieves all the data from
Redis. Then according to HTTP statuscodes received for different URLs it groups those URLs accordinf to
the Statuscode received. Finally it saves all the data in Redis with request status message as "Completed".

--- END ---
