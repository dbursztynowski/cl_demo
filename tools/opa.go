package tools

import (
	"bytes"
	"encoding/json"
	"strings"

	//"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//curl -X POST $OPA_URL/v1/data/closedloop/inner -d '{"input":{"memory":51}}'
//curl -X PUT $OPA_URL/v1/data -d "@data.json"
//curl -X PUT $OPA_URL/v1/policies/closedloop/inner --data-binary "@cl_in.rego"

//{"result":{"cpu":false,"decision":"memory","memory":true,"monitoring":{"cpu":false,"memory":true}}}

const OPA_URL = "http://192.168.49.2:32633"
const INNER_POLICY_MONITORING_IN = "/v1/data/closedloop/inner/monitoring"
const INNER_POLICY_DECISION_IN = "/v1/data/closedloop/inner/decision"
const INNER_POLICY_EXECUTION_IN = "/v1/data/closedloop/inner/execution"

const INNER_POLICY_MONITORING_OUT = "/v1/data/closedloop/outer/monitoring"
const INNER_POLICY_DECISION_OUT = "/v1/data/closedloop/outer/decision"
const INNER_POLICY_EXECUTION_OUT = "/v1/data/closedloop/otuer/execution"

func OpaPolicyFoundPolicy(opaUrl string, apiPolicy string, policyId string, l VerbosityLog) bool {
	l.V(1).Info("Enter PolicyFoundPolicy")
	defer func() {
		if x := recover(); x != nil {
			l.V(1).Info("run time panic: ", "error", x)
		}
	}()

	//curl -X GET $OPA_URL/v1/policies

	posturl := strings.Trim(opaUrl, "/") + "/" + strings.Trim(apiPolicy, "/")

	fmt.Println("posturl: " + posturl)
	fmt.Println()
	fmt.Println("Policy List from OPA")

	r, err := http.NewRequest("GET", posturl, nil)
	if err != nil {
		fmt.Printf("Error on GET : %s\n", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error on GET : %s\n", err)
	}

	defer res.Body.Close()

	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		fmt.Printf("Error on ReadAll(response.Body) : %s\n", err2)
	}

	// Convert to string Debug
	content := string(body)

	type Police struct {
		Id string
	}
	type Wynik struct {
		Result []Police
	}

	fmt.Println("Result List from OPA: >")
	fmt.Println("Result List from OPA: " + content)

	found := false
	forSearch := policyId
	messageJson := content
	var wynik Wynik
	json.Unmarshal([]byte(messageJson), &wynik)
	fmt.Printf("Policies : %+v\n", wynik.Result)
	for _, policy := range wynik.Result {
		if forSearch == policy.Id {
			fmt.Printf("Policy found: %+v\n", policy.Id)
			found = found || true
		}
	}
	fmt.Printf("fOUND : %+v\n", found)

	return found
}

func OpaPolicyWritePolicy(opaUrl string, apiPolicy string, policyName string, ruleBody string, l VerbosityLog) string {
	l.V(1).Info("Enter PolicyWritePolicy")
	defer func() {
		if x := recover(); x != nil {
			l.V(1).Info("run time panic: ", "error", x)
		}
	}()

	//curl -X PUT $OPA_URL/v1/policies/policy/closedloop/inner/monitoring --data-binary "@cl_in_monitoring2.rego"
	posturl := strings.Trim(opaUrl, "/") + "/" + strings.Trim(apiPolicy, "/") + "/" + strings.Trim(policyName, "/")

	// posturl = OPA_URL + "/v1/policies/closedloop/inner/monitoring"
	fmt.Println("for write posturl: " + posturl)
	fmt.Println()
	fmt.Println("Policy Write to OPA")
	fmt.Println(ruleBody)

	body := []byte(ruleBody)

	r, err := http.NewRequest("PUT", posturl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error on PUT : %s\n", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error on PUT : %s\n", err)
	}

	defer res.Body.Close()

	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		fmt.Printf("Error on ReadAll(response.Body) : %s\n", err2)
	}

	// Convert to string Debug
	content := string(body)
	fmt.Println("Result from OPA: >")
	fmt.Println("Result from OPA: " + content)

	return content
}
func OpaPolicyDeletePolicy(opaUrl string, apiPolicy string, policyName string, dataName string, l VerbosityLog) string {
	l.V(1).Info("Enter PolicyDeletePolicy")
	defer func() {
		if x := recover(); x != nil {
			l.V(1).Info("run time panic: ", "error", x)
		}
	}()

	//curl -X DELETE $OPA_URL/v1/policies/closedloop/inner/monitoring
	posturl := opaUrl + apiPolicy
	posturl = OPA_URL + "/v1/policies/closedloop/inner/monitoring"
	fmt.Println("posturl: " + posturl)
	fmt.Println()
	fmt.Println("Policy Delete from OPA")

	//body := []byte("{}")

	r, err := http.NewRequest(http.MethodDelete, posturl, nil) // bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error on DELETE : %s\n", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)

	fmt.Printf("res.StatusCode = %v\n", res.StatusCode)
	if err != nil {
		fmt.Printf("Error on DELETE : %s\n", err)
	}

	defer res.Body.Close()

	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		fmt.Printf("Error on ReadAll(response.Body) : %s\n", err2)
	}

	// Convert to string Debug
	content := string(body)
	fmt.Println("Result from OPA: >")
	fmt.Println("Result from OPA: " + content)

	return content
}

func OpaPolicyReadPolicy(opaUrl string, apiPolicy string, policyName string, dataName string, l VerbosityLog) string {
	l.V(1).Info("Enter PolicyReadPolicy")
	defer func() {
		if x := recover(); x != nil {
			l.V(1).Info("run time panic: ", "error", x)
		}
	}()

	//curl -X GET $OPA_URL/v1/policies/closedloop/inner/monitoring

	posturl := opaUrl + apiPolicy
	posturl = "v1/policies/closedloop/inner/monitoring"
	fmt.Println("posturl: " + posturl)
	fmt.Println()
	fmt.Println("Policy Read from OPA")

	r, err := http.NewRequest("GET", posturl, nil)
	if err != nil {
		fmt.Printf("Error on GET : %s\n", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error on GET : %s\n", err)
	}

	defer res.Body.Close()

	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		fmt.Printf("Error on ReadAll(response.Body) : %s\n", err2)
	}

	// Convert to string Debug
	content := string(body)
	fmt.Println("Result from OPA: >")
	fmt.Println("Result from OPA: " + content)

	return content
}

func OpaPolicyFoundData(opaUrl string, apiData string, dataName string, l VerbosityLog) bool {
	l.V(1).Info("Enter PolicyFoundData")
	defer func() {
		if x := recover(); x != nil {
			l.V(1).Info("run time panic: ", "error", x)
		}
	}()
	dane := strings.TrimSpace(OpaPolicyReadData(opaUrl, apiData, dataName, l))
	found := true
	if dane == "{\"result\":{}}" || dane == "{}" {
		found = false
	}
	return found
}

func OpaPolicyReadData(opaUrl string, apiData string, dataName string, l VerbosityLog) string {
	l.V(1).Info("Enter PolicyReadData")
	defer func() {
		if x := recover(); x != nil {
			l.V(1).Info("run time panic: ", "error", x)
		}
	}()

	//curl -X PUT $OPA_URL/v1/data --data-binary "@data.json"

	posturl := strings.Trim(opaUrl, "/") + "/" + strings.Trim(apiData, "/") + "/" + strings.Trim(dataName, "/")
	fmt.Println("for read data posturl: " + posturl)
	fmt.Println()
	fmt.Println("Data Read from OPA")

	r, err := http.NewRequest(http.MethodGet, posturl, nil)
	if err != nil {
		fmt.Printf("Error on GET : %s\n", err)
	}

	r.Header.Add("Content-Type", "application/json")

	//client := &http.Client{}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Printf("Error on GET : %s\n", err)
	}

	defer res.Body.Close()

	body, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		fmt.Printf("Error on ReadAll(response.Body) : %s\n", err2)
	}

	// Convert to string Debug
	fmt.Printf("body: %v", string(body))
	content := string(body)
	//content := "{\"result\":{\"Decisionpolicies\":{\"Decisiontype\":\"Priority\",\"Priorityspec\":{\"Priorityrank\":{\"rank-1\":\"cpu\",\"rank-2\":\"memory\"},\"Prioritytype\":\"Basic\",\"Time\":\"2023-12-01 21:51:58.427048\"}},\"Monitoringpolicies\":{\"Data\":{\"MonitoringData-1\":\"memory\",\"MonitoringData-2\":\"cpu\"},\"Time\":\"2023-12-01 21:51:58.427048\",\"Tresholdkind\":{\"MonitoringData-1-thresholdkind\":\"inferior\",\"MonitoringData-2-thresholdkind\":\"inferior\"},\"Tresholdvalue\":{\"MonitoringData-1-thresholdvalue\":50,\"MonitoringData-2-thresholdvalue\":5}},\"closedloop\":{\"inner\":{\"decision\":{\"metric\":\"none\",\"monitoring\":{\"cpu\":false,\"memory\":false}},\"monitoring\":{\"cpu\":\"ok\",\"memory\":\"ok\"}}},\"dupa\":{\"allow\":true,\"greeting\":\"hello from pod 'opa-5c4555668-dncjf'\"}}}"
	//content := "{\"Decisionpolicies\":{\"Decisiontype\":\"Priority\",\"Priorityspec\":{\"Priorityrank\":{\"rank-1\":\"cpu\",\"rank-2\":\"memory\"},\"Prioritytype\":\"Basic\",\"Time\":\"2023-12-01 21:51:58.427048\"}},\"Monitoringpolicies\":{\"Data\":{\"MonitoringData-1\":\"memory\",\"MonitoringData-2\":\"cpu\"},\"Time\":\"2023-12-01 21:51:58.427048\",\"Tresholdkind\":{\"MonitoringData-1-thresholdkind\":\"inferior\",\"MonitoringData-2-thresholdkind\":\"inferior\"},\"Tresholdvalue\":{\"MonitoringData-1-thresholdvalue\":50,\"MonitoringData-2-thresholdvalue\":5}}}"

	fmt.Println("Result from OPA: >")
	fmt.Println("Result from OPA: " + content)

	l.V(1).Info("End PolicyReadData")
	return content
}

func OpaPolicyDeleteData(opaUrl string, apiData string, policyName string, dataName string, l VerbosityLog) {

	l.V(1).Info("Enter PolicyDeleteData")
	defer func() {
		if x := recover(); x != nil {
			l.V(1).Info("run time panic: ", "error", x)
		}
	}()

	//curl -X PUT $OPA_URL/v1/data --data-binary "@data.json"

	posturl := opaUrl + apiData
	if len(posturl) == 0 {
		return
	}
	fmt.Println("posturl: " + posturl)
	fmt.Println()
	fmt.Println("Data Write to OPA: ")

	body := []byte("{}")

	r, err := http.NewRequest("PUT", posturl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error on PUT : %s\n", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error on PUT : %s\n", err)
	}

	defer res.Body.Close()

	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		fmt.Printf("Error on ReadAll(response.Body) : %s\n", err2)
	}

	// Convert to string Debug
	content := string(body)
	fmt.Println("Result from OPA: " + content)

	return

}

func OpaPolicyWriteData(opaUrl string, apiData string, dataName string, dataValue string, l VerbosityLog) {

	l.V(1).Info("Enter PolicyWriteData")
	defer func() {
		if x := recover(); x != nil {
			l.V(1).Info("run time panic: ", "error", x)
		}
	}()

	//curl -X PUT $OPA_URL/v1/data --data-binary "@data.json"

	posturl := strings.Trim(opaUrl, "/") + "/" + strings.Trim(apiData, "/") + "/" + strings.Trim(dataName, "/")
	if len(posturl) == 0 {
		return
	}

	fmt.Println("for data write posturl: " + posturl)
	fmt.Println()
	fmt.Println("Data Write to OPA: " + dataValue)

	body := []byte(dataValue)

	r, err := http.NewRequest("PUT", posturl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error on PUT : %s\n", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error on PUT : %s\n", err)
	}

	defer res.Body.Close()

	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		fmt.Printf("Error on ReadAll(response.Body) : %s\n", err2)
	}

	// Convert to string Debug
	content := string(body)
	fmt.Println("Result from OPA: " + content)

	return

}

func OpaPolicy(opaUrl string, apiPolicy string, policyName string, message string, l VerbosityLog) string {

	defer func() {
		if x := recover(); x != nil {
			l.V(1).Info("run time panic: ", "error", x)
		}
	}()
	l.V(1).Info("Enter Policy_monitoring")
	posturl := opaUrl + apiPolicy + policyName
	fmt.Println("posturl: " + posturl)
	//input := map[string](string){
	//	"input": message,
	//	}
	//message_json,_ := json.Marshal(input)
	message_json := []byte("{\"input\":" + strings.ReplaceAll(message, "'", "\"") + "}")
	fmt.Println("after ReplaceAll")
	fmt.Println("Message JSON to OPA: " + string(message_json))

	body := []byte(message_json)

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error on POST : %s\n", err)
		return ""
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error on POST : %s\n", err)
	}

	defer res.Body.Close()

	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		fmt.Printf("Error on ReadAll(response.Body) : %s\n", err2)
	}

	// Convert to string Debug
	content := string(body)
	fmt.Println("Result from OPA: " + content)

	return content

}
