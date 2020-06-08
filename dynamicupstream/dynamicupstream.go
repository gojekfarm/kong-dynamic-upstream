package dynamicupstream

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var PrefixExpression = regexp.MustCompile(`/(.+?)/\$upstream`)

//DynamicUpstream provides Regex based routing to inferred Upstream
type DynamicUpstream struct {
	regex    *regexp.Regexp
	regexStr string
	port     int
	Prefix   string
}

//NewUpstream creates a new dynamic upstream
func New(pathExp string, port string) (*DynamicUpstream, error) {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	prefixMatches := PrefixExpression.FindStringSubmatch(pathExp)
	if len(prefixMatches) < 1 {
		return nil, fmt.Errorf("Upstream doesnt have a prefix")
	}

	regexStr := strings.ReplaceAll(pathExp, "$upstream", "(.+?)")

	re, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, err
	}
	return &DynamicUpstream{regex: re, port: intPort, Prefix: fmt.Sprintf("/%s/", prefixMatches[1]), regexStr: regexStr}, nil
}

func (du *DynamicUpstream) MatchPrefix(path string) bool {
	return strings.HasPrefix(path, du.Prefix)
}

func (du *DynamicUpstream) Target(kong KongPDK) error {

	path, err := kong.RequestPath()
	if err != nil {
		return err
	}
	if du.MatchPrefix(path) {

		hostPathSegment := du.regex.FindStringSubmatch(path)[1]

		err = kong.SetUpstreamTarget(hostPathSegment, du.port)
		if err != nil {
			return err
		}

		return kong.SetUpstreamTargetRequestPath(strings.TrimPrefix(path, fmt.Sprintf("%s%s", du.Prefix, hostPathSegment)))
	}

	return fmt.Errorf("Path %s, didnt match prefix: %s", path, du.Prefix)
}
