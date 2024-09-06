### TODO before production:

[] Make sure ports are set correctly, because starting backend before frontend breaks things.

[] Sort by character frequency, obtained here: https://lingua.mtsu.edu/chinese-computing/statistics/char/list.php?Which=MO


[] Fix login and registration.  Make it actually secure.
    [] Write tests for logins with multiple users
    [] Persist login between sessions.
[] Make email save to database, so I can definitely never send email to anybody who gives me their email address

[] make `docker-compose up` run everything as multiple containers... so it's nice and easy to deploy on a server
