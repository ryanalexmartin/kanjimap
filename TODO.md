# TODO before production:

## Zhuyin Highlighter Extension

[] Configure server URL dynamically:
   - Handle both "localhost" and "kanjimap.cargocult.tech"
   - Implement environment-based switching ("dev" or "prod")
[] Word rules with context clues for better Zhuyin suggestions (if LiuChan can do it, so can I)
[] publish extension to firefox store

    [] icon for extension

[] handle live DOM reloads better

    [] sites like DCard should be handled
    [] switching orientations or enabling/disbling should be immediately reflected in the DOM

[] stuff like button line breaks need to be handled better.
[] profile performance and try to optimize
[] pop up definitions? handle add/remove from firefox extension?

### more clients for zhuyin:

    [] Google Chrome
    [] Edge
    [] Safari
    [] Obsidian
    ... Desktop?
    Mobile?

[] We really need to make some good test suites and deployment pipelines if we are going to support so many platforms.

---

## Website frontend
    [] Make the website look a lot better
    [] donation? subscription?  ads?
    [] logo
    [] better name


## Backend

    [] I'm nervous about pushing my local changes, because I'm scared it will wipe the database.
    [] make sure setup.sql runs and migrations are solid.  Python shouldn't be creating tables!

[x] Fix login and registration.  Make it actually secure. (???)

    [] Write tests for logins with multiple users
    [] Ensure that a user can be logged in in multiple sessions at once
    [] Make email save to database, so I can definitely never send email to anybody who gives me their email address


# Bugs
[] Make sure there aren't multiple login requests over and over.  Seems to be an issue right now.

# Future features
[] japanese language support
