package cryptoadp_test

import (
	"testing"
	cryptoadp "video-provider/user-service/adapters/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestNewBCryptPasswordHasher(t *testing.T) {
	hasher := cryptoadp.NewBCryptPasswordHasher()
	require.NotNil(t, hasher)
}

func TestBcryptPasswordHasher_Hash(t *testing.T) {
	cases := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password hash",
			password: "MySecurePassword123!!",
			wantErr:  false,
		},
		{
			name:     "Hash with special characters",
			password: "P@ssw0rd!@#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "Very long password",
			password: "VeryLongPasswordVeryLongPasswordVeryLongPasswordVeryLongPassword123!!",
			wantErr:  false,
		},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hasher := cryptoadp.NewBCryptPasswordHasher()

			hash, err := hasher.Hash(c.password)

			if c.wantErr {
				require.Error(t, err)
				assert.Nil(t, hash)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, hash)
				assert.NotEmpty(t, hash)

				// Verify the hash is a valid bcrypt hash
				err := bcrypt.CompareHashAndPassword(hash, []byte(c.password))
				assert.NoError(t, err)
			}
		})
	}
}

func TestBcryptPasswordHasher_CompareHashAndPassword(t *testing.T) {
	hasher := cryptoadp.NewBCryptPasswordHasher()
	password := "TestPassword123!!"

	// First hash the password
	hash, err := hasher.Hash(password)
	require.NoError(t, err)

	cases := []struct {
		name     string
		hash     []byte
		password []byte
		wantErr  bool
	}{
		{
			name:     "Matching password and hash",
			hash:     hash,
			password: []byte(password),
			wantErr:  false,
		},
		{
			name:     "Different password",
			hash:     hash,
			password: []byte("WrongPassword123!!"),
			wantErr:  true,
		},
		{
			name:     "Empty password attempt",
			hash:     hash,
			password: []byte(""),
			wantErr:  true,
		},
		{
			name:     "Case sensitive password",
			hash:     hash,
			password: []byte("testpassword123!!"), // lowercase version
			wantErr:  true,
		},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := hasher.CompareHashAndPassword(c.hash, c.password)

			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBcryptPasswordHasher_HashConsistency(t *testing.T) {
	hasher := cryptoadp.NewBCryptPasswordHasher()
	password := "ConsistencyTestPassword123!!"

	// Hash the same password twice
	hash1, err1 := hasher.Hash(password)
	require.NoError(t, err1)

	hash2, err2 := hasher.Hash(password)
	require.NoError(t, err2)

	// The hashes should be different due to salt
	assert.NotEqual(t, hash1, hash2)

	// But both should match the same password
	assert.NoError(t, hasher.CompareHashAndPassword(hash1, []byte(password)))
	assert.NoError(t, hasher.CompareHashAndPassword(hash2, []byte(password)))
}

func TestBcryptPasswordHasher_InvalidHashFormat(t *testing.T) {
	hasher := cryptoadp.NewBCryptPasswordHasher()
	password := []byte("SomePassword123!!")

	invalidHashes := []struct {
		name string
		hash []byte
	}{
		{
			name: "Invalid hash format",
			hash: []byte("not-a-valid-bcrypt-hash"),
		},
		{
			name: "Empty hash",
			hash: []byte(""),
		},
		{
			name: "Random bytes",
			hash: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}

	for _, ih := range invalidHashes {
		t.Run(ih.name, func(t *testing.T) {
			err := hasher.CompareHashAndPassword(ih.hash, password)
			require.Error(t, err)
		})
	}
}
