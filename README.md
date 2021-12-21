# Language Study Tool (LST) - API module

Written in Golang using DDD approach.

## Domain
### Entities
- Language
- Expression
- Translation

## Application
### Entities
- User
- Registrant
- Profile
- Folder

### Use Cases
- Save registrant (signup)
- Save user (confirm-email)
- Delete registrant (confirm-email)
- Find user by email (login)
- Get user by id
- Save profile
- Delete profile
- List languages 
- Save folder
- Delete folder
- Save expression
- Delete expression
- Attach expression to folder
- Detach expression from folder
- Save translation
- Delete translation
- Attach translation
- Detach translation

#### Text Use Cases
- Add public text:
    - After publishing the text it cannot be modified if it has one or more comments/corrections ???
    - Users can add corrections and comments
    - Writer can select correction and apply it
    - Approved text can be published globally and indexed.

