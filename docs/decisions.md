### Who should do the userExists check? The auth middleware or the handler?
- Who should do the userExists check? The auth middleware (just after successfully decoding the jwt token) or the handler?
- IF authmiddleware does it
  - CONS
    - if user asks for the projects route(or any route that does not require user info), the userCheck query will run => extra query
- IF not checking in authmiddleware
  - if i am logged in on one tab and in the other duplicate tab, i delete my account, the in the first tab, i could still access the resource unless i log out or the resource handler does the userExists check
  - also duplicate user query in each handler
> DECISION :- do userExists check in AuthenticationMiddleware to prevent duplicate db-user-fetch-by-id query and also prevent deleted account's from accessing resources