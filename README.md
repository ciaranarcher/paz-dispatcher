paz-dispatcher
==============
A key element of the PAZ (Post Apocalyptic Zendesk) Dublin hackday project

What is this? 
==============
Create a Zendesk ticket even if you local Internet is down using ham radio.

I can't even...
==============
Two 'workers' included:

* Read from a Redis remotely (on a Pi), and dispatch API calls to create tickets
* Read ticket updates from a Redis locally and shovel them into the remote Redis for dispatch to the remote radio
