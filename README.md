# go-chapper

## Installation

### CLI

Run the `./server install` command.

### Manual installation

The manual installation provides more control over the installation process but requires
a lot more setup, like creating directories, config files and more...

## Running

To run your Chapper instance use the `./server run --config path/to/your/config.toml`
command.

## TODOs

-   Finish bridge
-   Use null package (store/database)
-   Make error ctx more robust/thought trough
-   Create invite template
-   Rework server handlers/services (and more?)
-   Rework TURN New function
-   Handle signals gracefully

## Ideas

-   Add AutoTLS support (ACME Manager) and ask this in the install CLI

## All Todos

### Authentication

-   [x] Password Argon2 Hashing
-   [ ] Authentication Routes
    -   [x] Login
    -   [x] Register
    -   [ ] Refresh
    -   [ ] Register 2FA code
    -   [ ] Enter 2FA code
-   [ ] Password Reset

### Avatars

-   [x] Generate Avatars based on Username
    -   [ ] Save Avatars on disk
    -   [ ] JPEG or SVG?
    -   [ ] Resizer
-   [ ] Custom Avatar support
    -   [ ] Routes
    -   [ ] Resizer
    -   [ ] Constraints

### Profile

-   [ ] Custom Avatar
-   [ ] Username change?
-   [ ] Password change
-   [ ] 2FA add/delete
-   [ ] Change e-mail
-   [ ] Privacy settings
    -   [ ] What data is public?
    -   [ ] Who can add me as a friend?
    -   [ ] ...

### Config

-   [x] Add TOML config support
    -   [ ] Validation
    -   [ ] Default values

### Virtual Servers

-   [ ] Add virtual server support
    -   [ ] Invite System
    -   [ ] Routes (CRUD Actions)
    -   [ ] Keep track which virtual servers the user is on

### Rooms

-   [ ] Add Voice Rooms
    -   [ ] Session management
    -   [ ] Routes
    -   [ ] Key exchange
    -   [ ] Admin controls (Mute, kick user, etc)
-   [ ] Add Text Rooms
    -   [ ] Session management
    -   [ ] Routes
    -   [ ] Key exchange
    -   [ ] Admin controls
    -   [ ] Multimedia message support

### 1 on 1 Room

-   [ ] Add Voice Room
    -   [ ] Session management
    -   [ ] Routes
    -   [ ] Key exchange
    -   [ ] Admin controls (Mute, kick user, etc)
-   [ ] Add Text Room
    -   [ ] Session management
    -   [ ] Routes
    -   [ ] Key exchange
    -   [ ] Admin controls
    -   [ ] Multimedia message support

### Media Handler

-   [ ] Routes
    -   [ ] Upload
    -   [ ] Original
    -   [ ] Preview
-   [ ] Saving media on disk
-   [ ] Images
    -   [ ] Validation / Constraints
    -   [ ] Resizing
    -   [ ] Quality
    -   [ ] Save in different file formats
-   [ ] Videos
    -   [ ] Validation / Constraints
    -   [ ] Resizing
    -   [ ] Quality
    -   [ ] Save in different file formats (Transcoding)

### Broadcaster / Session Handler / TURN / STUN

-   [ ] Alot

### Scheduler?

-   [ ] How do we handle background jobs?
-   [ ] Microservice approach?
