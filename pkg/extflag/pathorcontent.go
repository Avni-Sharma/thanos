package extflag

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

// PathOrContent is a flag type that defines two flags to fetch bytes. Either from file (*-file flag) or content (* flag).
type PathOrContent struct {
	flagName string

	required bool

	path    *string
	content *string
}

type CmdClause interface {
	Flag(name, help string) *kingpin.FlagClause
}

// RegisterPathOrContent registers PathOrContent flag in kingpinCmdClause.
func RegisterPathOrContent(cmd CmdClause, flagName string, help string, required bool) *PathOrContent {
	fileFlagName := fmt.Sprintf("%s-file", flagName)
	contentFlagName := flagName

	fileHelp := fmt.Sprintf("Path to %s", help)
	fileFlag := cmd.Flag(fileFlagName, fileHelp).PlaceHolder("<file-path>").String()

	contentHelp := fmt.Sprintf("Alternative to '%s' flag (lower priority). Content of %s", fileFlagName, help)
	contentFlag := cmd.Flag(contentFlagName, contentHelp).PlaceHolder("<content>").String()

	return &PathOrContent{
		flagName: flagName,
		required: required,
		path:     fileFlag,
		content:  contentFlag,
	}
}

// Content returns content of the file. Flag that specifies path has priority.
// It returns error if the content is empty and required flag is set to true.
func (p *PathOrContent) Content() ([]byte, error) {
	contentFlagName := p.flagName
	fileFlagName := fmt.Sprintf("%s-file", p.flagName)

	if len(*p.path) > 0 && len(*p.content) > 0 {
		return nil, errors.Errorf("both %s and %s flags set.", fileFlagName, contentFlagName)
	}

	var content []byte
	if len(*p.path) > 0 {
		c, err := ioutil.ReadFile(*p.path)
		if err != nil {
			return nil, errors.Wrapf(err, "loading YAML file %s for %s", *p.path, fileFlagName)
		}
		content = c
	} else {
		content = []byte(*p.content)
	}

	if len(content) == 0 && p.required {
		return nil, errors.Errorf("flag %s or %s is required for running this command and content cannot be empty.", fileFlagName, contentFlagName)
	}

	return content, nil
}
