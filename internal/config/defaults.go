package config

const defaultReadme = "# Soft Serve\n\n Welcome! You can configure your Soft Serve server by cloning this repo and pushing changes.\n\n```\ngit clone ssh://{{.Host}}:{{.Port}}/config\n```"

const defaultConfig = `# The name of the server to show in the TUI.
name: Soft Serve

# The host and port to display in the TUI. You may want to change this if your
# server is accessible from a different host and/or port that what it's
# actually listening on (for example, if it's behind a reverse proxy).
host: %s
port: %d

# Access level for anonymous users. Options are: read-write, read-only and
# no-access.
anon-access: %s

# You can grant read-only access to users without private keys. Any password
# will be accepted.
allow-keyless: false

# Customize repo display in the menu. Only repos in this list will appear in
# the TUI.
repos:
  - name: Home
    repo: config
    private: true
    note: "Configuration and content repo for this server"
`

const hasKeyUserConfig = `

# Authorized users. Admins have full access to all repos. Users can read all
# repos and push to their collab-repos.
users:
  - name: Admin
    admin: true
    public-keys:
%s
`

const defaultUserConfig = `
# users:
#   - name: Admin
#     admin: true
#     public-keys:
#       - KEY TEXT`

const exampleUserConfig = `
#   - name: Example User
#     collab-repos:
#       - REPO
#     public-keys:
#       - KEY TEXT
`
