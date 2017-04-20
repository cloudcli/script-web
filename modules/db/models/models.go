package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"io"
	"crypto/rand"
	"encoding/hex"
)

const (
	AdminRole = 0
	NormalRole = 1
)

const (
	saltLength = 16
)

/**
 * User Model
 */
type User struct {
	gorm.Model
	Email      string `gorm:"type:varchar(200);unique_index"`
	Name       string `gorm:"type:varchar(200);"`
	Phone      string `gorm:"type:varchar(20);"`
	Department string `gorm:"type:varchar(200);"`
	Password   string `gorm:"type:varchar(200);"`
	Salt       string `gorm:"type:varchar(32);"`
	Role       int    `gorm:"default:1;"`
	Avatar     string `gorm:"type:blob;"`
}


func newSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// NewUser for create a new user model
func NewUser() (*User,  error) {
	salt, err := newSalt()
	if err != nil {
		return nil, err
	}

	return &User{
		Salt: hex.EncodeToString(salt),
	}, nil
}

func (u User) IsAdmin() bool {
	return u.Role == AdminRole
}

func (u User) IsNormal() bool {
	return u.Role == NormalRole
}

// EncryptPassword using bcrypt for encrypt password
func (u User) EncryptPassword(password string) ([]byte, error) {
	p := []byte(u.Salt + password)

	hashedPassword, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// IsValidPassword check password is valid
func (u User) IsValidPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(u.Salt + password))
	return err == nil
}

type UserFull struct {
	User
}

/**
 * Comment
 */
type Comment struct {
	gorm.Model
	UserID   uint                       // Comment user id
	RefersTo uint                       // Comment id
	Content  string `gorm:"type:text;"` // Comment content in markdown
}

/**
* Message Box: inform user
 */
type Message struct {
	gorm.Model
	From    uint
	To      uint
	Message string `gorm:"type:text;"`         // Message in json format
	Status  int    `gorm:"type:int;default:0"` // Message type: Enum{APPLY, APPROVAL, REJECT}
	Viewed  bool   `gorm:"default:false"`      // Is message viewed
	Sys     bool   `gorm:"default:false"`      // Is message from system
}

/**
 * Sys Log
 *
 *  Every event produced by SYS or User will create an event log
 */
type EventLog struct {
	gorm.Model
	Message       string `gorm:"type:varchar(1000);"` // message in json format
	MessageStatus int    `gorm:"type:int;default:0"`  // Log type: Enum
}

/**
 * Script
 *
 * Scripts lifecycle
 *
 * Create:
 *   1. User: http post to create new script
 *   	- SYS: Create an Apply record
 *      - SYS: Create Message records and send message to `AssignTo`
 *   2. Admin: review apply
 *      - Approved: http post approved
 * 		- SYS: change Apply record status to APPROVE
 *		- SYS: Create Script and ScriptVersion
 *		- SYS: Create Message records and send message to proposer
 *	- Rejected: http post reject
 * 		- SYS: Change Apply record status to REJECT
 *		- SYS: Create Message records and send message to proposer
 *
 */
type Script struct {
	gorm.Model
	Name          string `gorm:"type:varchar(200);"`        // Script name
	Description   string `gorm:"type:text;"`                // description
	Path          string `gorm:"type:varchar(1000);"`       // Script folder path
	ScriptVersion int    `gorm:"default:1"`                 // Script current version, foreign key of ScriptVersion
								//	Tags          string `gorm:"type:varchar(200);"`  // Tags array eg: tag1,tag2,tag3
	DownloadCount int         `gorm:"default:0"`            // Download times
	Tags          []ScriptTag `gorm:many2many:script_stags` // tags
}

type ScriptTag struct {
	gorm.Model
	Name string `gorm:"type:varchar(200)"`
}

/**
* ScriptVersion
 *
  *  - script has multiple versions and every version is a git commit
   *  - we can use git-cat-file to retrieve file content (https://git-scm.com/docs/git-cat-file)
*/
type ScriptVersion struct {
	gorm.Model
	ScriptId uint                              // Foreign key of Script
	ApplyId  uint                              // Foreign key of Apply
	Version  int                               // version
	Hash     string `gorm:"type:varchar(50);"` // Git commit hash
	Diff     string `gorm:"type:text;"`        // Store diff string compared to previous version code
}

const (
	ApplyTypeNewScript = 0
	ApplyTypeUpdateScript = 1
	ApplyTypeDeleteScript = 2
	ApplyTypeNewFolder = 3
)

const (
	ApplyStatusPending = 0
	ApplyStatusApprove = 1
	ApplyStatusReject = 2
)

/**
  * Apply
   *
    *      - Every
*/
type Apply struct {
	gorm.Model
	Proposer    string `gorm:"type:varchar(200);"`
	AssignTo    string `gorm:"type:varchar(200);"`
	Title       string `gorm:"type:varchar(255);"`
	Description string `gorm:"type:varchar(1000);"`
	Payload     string `gorm:"type:text;"`         // Apply event payload in json
	Type        int    `gorm:"type:int;"`          // Apply biz type: Enum{CREATE_SCRIPT、UPDATE_SCRIPT、CREATE_DIR、DELETE_SCRIPT}
	Status      int    `gorm:"type:int;default:0"` // Apply Status: Enum{PROPOS, APROVE, REJECT}
}

func (a Apply) IsNewScript() bool {
	return a.Type == ApplyTypeNewScript
}

func (a Apply) IsUpdateScript() bool {
	return a.Type == ApplyTypeUpdateScript
}

func (a Apply) IsDeleteScript() bool {
	return a.Type == ApplyTypeDeleteScript
}

func (a Apply) IsNewFolder() bool {
	return a.Type == ApplyTypeNewFolder
}

func (a Apply) IsPendingStatus() bool {
	return a.Status == ApplyStatusPending
}

func (a Apply) IsApproveStatu() bool {
	return a.Status == ApplyStatusApprove
}

func (a Apply) IsRejectStatus() bool {
	return a.Status == ApplyStatusReject
}

const (
	ValueTypeInt = 0
	ValueTypeString = 1
)

/**
 * Dynamic customized {key, value} data
 */
type KvStore struct {
	gorm.Model
	Key         string `gorm:"type:varchar(50);unique_index"`
	StringValue string `gorm:"type:varchar(1000);"`
	IntValue    int
	ValueType   int `gorm:"type:int;"` // Enum{INT, STRING}
}

func (s KvStore) IsInt() bool {
	return s.ValueType == ValueTypeInt
}

func (s KvStore) IsString() bool {
	return s.ValueType == ValueTypeString
}

func (s KvStore) GetString() string {
	return s.StringValue
}

func (s KvStore) GetInt() int {
	return s.IntValue
}
