// User storage models.
package users

// A user account.
type Account struct {
  // A numeric ID of the account.
  // This should not be exposed externally.
  AccountId int32

  // The username of the account.
  // This can be exposed externally and used in URLs.
  Username string

  // A hash of the user's password.
  PasswordHash string
}

