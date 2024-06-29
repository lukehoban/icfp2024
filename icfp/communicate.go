package icfp

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CommunicateToken(s string) ([]Expr, error) {
	request, err := http.NewRequest("POST", "https://boundvariable.space/communicate", strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+"3b4c0eaa-bdc2-42ff-8e28-4a652814bd73")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	byts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	expr := Parse(string(byts))
	return expr, nil
}

func CommunicateString(s string) (string, error) {
	tok := StringToToken(s)
	ret, err := CommunicateToken(string(tok))
	if err != nil {
		return "", err
	}
	expr, rest := CombineToExpr(ret)
	if len(rest) > 0 {
		return "", fmt.Errorf("didn't use all input! %v", rest)
	}
	res := Eval(expr)
	out, ok := res.(String)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", res)
	}
	return string(out), nil
}
