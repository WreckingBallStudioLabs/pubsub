package name

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/WreckingBallStudioLabs/pubsub/errorcatalog"
)

//////
// Const, vars, and types.
//////

var nameRegex = regexp.MustCompile(`^v\d+\.[a-zA-Z]+(?:\.[a-zA-Z]+)+(?:\.queue)?$`)

type (
	// Name of a topic, or queue.
	Name string

	// Queue is the name of a queue. Should be in the form of the following
	// example: "v1.meta.created.queue".
	Queue Name

	// Topic is the name of a topic. Should be in the form of the following
	// example: "v1.meta.created".
	Topic Name
)

//////
// Methods.
//////

// Implement the Stringer interface.
func (n Name) String() string {
	return string(n)
}

// Validate the name.
func (n Name) Validate() error {
	if !nameRegex.MatchString(n.String()) {
		return errorcatalog.Get().MustGet(errorcatalog.PubSubErrNameName).NewInvalidError()
	}

	return nil
}

// Parts breaks a Name into its parts.
func (n Name) Parts() []string {
	return nameRegex.FindStringSubmatch(n.String())
}

// ToQueue converts a Name to a Queue, adding the .queue suffix only if it's not
// already there.
func (n Name) ToQueue() Queue {
	if err := n.Validate(); err == nil {
		// Check if n contains the .queue suffix.
		if !strings.Contains(n.String(), ".queue") {
			return Queue(fmt.Sprintf("%s.queue", n.String()))
		}
	}

	return Queue(n.String())
}

// ToTopic converts a Name to a Topic, removing the .queue suffix only if it's
// there.
func (n Name) ToTopic() Topic {
	if err := n.Validate(); err == nil {
		// Check if n contains the .queue suffix.
		if strings.Contains(n.String(), ".queue") {
			return Topic(strings.ReplaceAll(n.String(), ".queue", ""))
		}
	}

	return Topic(n.String())
}

//////
// Factory.
//////

// New creates a new name. It should be in the format of the following
// example: "v1.meta.created" or "v1.meta.created.queue".
func New(name string) (Name, error) {
	n := Name(name)

	if err := n.Validate(); err != nil {
		return "", err
	}

	return n, nil
}

// MustNew creates a new name. It should be in the format of the following
// example: "v1.meta.created" or "v1.meta.created.queue". It panics if the name
// is invalid.
func MustNew(name string) Name {
	n, err := New(name)
	if err != nil {
		panic(err)
	}

	return n
}
