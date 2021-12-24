# Language Study Tool (LST) - API module

Written in Golang using DDD approach.

## Domain
### Entities
- Language
- Expression
- Translation
- Text

## Application
### Entities
- User
- Group
- Slice
- Comment

### Use Cases
#### Auth:
- [x] Register user
- [x] Confirm registration
- [x] Log in

#### Group
- [x] Create group
    - [ ] Search languages
    - [ ] Search users to invite
- [x] List groups
- [ ] Update group name
- [ ] Mark group as deleted
- [ ] Attach users to a group
- [ ] Send invitation to group link
- [ ] Accept invitation to group
- [ ] Reject invitation to group
- [ ] Detach users from a group
- [x] Add slice to group
- [ ] Delete slice from group
- [x] List group slices


#### Slice
- [ ] Get slice
- [ ] Move slice
- [ ] Attach expression
- [ ] Detach expression
- [ ] Attach translation
- [ ] Detach translation
- [ ] Set text
- [ ] Unset text

#### Expression
- [ ] Get
- [x] Search

#### Text
- [ ] Get
- [ ] Search
- [ ] Add correction
- [ ] Accept correction
- [ ] Check similarity of adding and existing texts to prevent duplicates

