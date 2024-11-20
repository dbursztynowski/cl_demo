package tools

import (
	"context"
	"encoding/json"
	"errors"
	"sort"

	"fmt"
	"reflect"

	"strings"
	"text/template"

	"github.com/tidwall/gjson"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type target struct {
	NameCR string
	KindCR string
}

type myMessage struct {
	Target  target
	Metric  string
	Path    []string
	Message string
}

func MessageDataCollect(client client.Client, ctx context.Context, req ctrl.Request, input string, inputmessage string, nameR string, kindR string, l VerbosityLog) string {

	input_map := gjson.Get(input, "@this")
	result := input
	//adJson := new(ADJson) //use it after solving bugs
	//adJson.SetJson(input)
	l.V(1).Info("MessageDataCollect input: " + input)
	input_map.ForEach(func(key, value gjson.Result) bool {
		l.V(1).Info(key.String() + ": " + value.String())
		//adJson.SetElement(key.String(), CollectValue(client, ctx, req, value.String(), nameR, kindR, l))
		result = strings.Replace(result, value.String(), CollectValue(client, ctx, req, value.String(), nameR, kindR, l), 1)
		l.V(1).Info(key.String() + ": " + result)
		return true // keep iterating
	})

	//println("input string from data collection: " + adJson.jsonString)
	return result
}

func CollectValue(client client.Client, ctx context.Context, req ctrl.Request, where string, nameR string, kindR string, l VerbosityLog) string {
	var path []string
	var cr []string
	var nameCR string
	var kindCR string

	if where[0] == '#' {
		path = strings.Split(where[1:], ".")
		nameCR = nameR
		kindCR = kindR

	} else {
		pos := strings.Index(where, "cr:")
		if pos >= 0 {
			cre := strings.Split(where[pos+3:], "#")
			cr = strings.Split(cre[0], ".")
			path = strings.Split(cre[1], ".")
			nameCR = cr[0]
			kindCR = cr[1]

		}
	}
	clcm := Cr_type(nameCR, kindCR)
	if clcm == nil || path == nil {
		return ""
	}

	err := client.Get(ctx, types.NamespacedName{Name: req.Name[:strings.LastIndex(req.Name, "-")] + "-" + strings.ToLower(nameCR),
		Namespace: req.Namespace}, clcm)

	if err != nil {
		l.Error(err, "Error getting "+req.Name)
	} else {
		//println("Is OK with " + req.Name)
		v := MyReflectValue{reflect.ValueOf(clcm).Elem()}
		pathJson := ""
		for _, node := range path {
			node = strings.Title(node)
			l.V(1).Info("node: " + node)
			l.V(1).Info(v.XFieldByName(node).Value.String())
			if v.XFieldByName(node).Value != reflect.ValueOf(Invalid{0}) {
				l.V(1).Info("xvalue: " + v.XFieldByName(node).String())
				v = v.XFieldByName(node)
			} else {
				if pathJson == "" {
					pathJson = strings.ToLower(node)
				} else {
					pathJson += "." + strings.ToLower(node)
				}
			}

		}
		l.V(1).Info("pathJson: " + pathJson)
		message := strings.Replace(v.String(), "'", "\"", -1)
		l.V(1).Info("v.String(): " + message)
		l.V(1).Info("gjson.Get(v.String(), pathJson).String(): " + gjson.Get(message, pathJson).String())
		return gjson.Get(message, pathJson).String()
	}

	return ""
}

func MessageInputConvert(jsonSchemaStr string,
	inputmessage string,
	input string, l VerbosityLog) string {
	fmt.Println("jsonSchemaStr:" + jsonSchemaStr)
	var schema JSONSchema
	err := json.Unmarshal([]byte(jsonSchemaStr), &schema)
	if err != nil {
		l.Error(err, "Error parsing JSON schema")
		return inputmessage
	}
	templateStr := generateTemplateFromSchema(schema, "")
	/* 	fmt.Println("Generated Template:\n", templateStr)
	   	fmt.Println("inputmessage: ",inputmessage)
	*/var data map[string]interface{}
	err = json.Unmarshal([]byte(inputmessage), &data)
	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	fmt.Printf("data: %v", data)
	tmpl, err := template.New("jsonTemplate").Parse(templateStr)
	if err != nil {
		l.Error(err, "Error creating template")
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		l.Error(err, "Error executing template")
	}

	fmt.Println("Rendered JSON:\n", result.String())

	message := result.String()
	return message
}

func MessageOutputConvert(client client.Client, ctx context.Context, req ctrl.Request, nameCR string, kindCR string,
	jsonSchemaStr string,
	opaanswer string,
	result string, l VerbosityLog) map[target][]myMessage {

	//		fmt.Println("out jsonSchemaStr:" + jsonSchemaStr)
	var schema JSONSchema
	err := json.Unmarshal([]byte(jsonSchemaStr), &schema)
	if err != nil {
		l.Error(err, "Error parsing JSON schema out")
	}
	templateStr := generateTemplateFromSchema(schema, "")
	/* 		fmt.Println("Generated Template out:\n", templateStr)
	   		fmt.Println("opaanswer in out: ",opaanswer)
	*/
	var data map[string]interface{}
	err = json.Unmarshal([]byte(opaanswer), &data)
	if err != nil {
		fmt.Println(err)
		return map[target][]myMessage{}
	}

	tmpl, err := template.New("jsonTemplate").Parse(templateStr)
	if err != nil {
		l.Error(err, "Error creating template out")
	}

	var results strings.Builder
	err = tmpl.Execute(&results, data)
	if err != nil {
		l.Error(err, "Error executing template")
	}

	//		fmt.Println("Rendered JSON out:\n", results.String())

	//messages := DistributeValue(result, results.String())
	messages := MessageDataDistributeMap(client, ctx, req,
		result, results.String(), nameCR, kindCR, l)
	return messages
}

func MessageDataDistributeMap(client client.Client, ctx context.Context, req ctrl.Request,
	targets string, inputmessage string, nameR string, kindR string, l VerbosityLog) map[target][]myMessage {
	var messages = make([]myMessage, 0)
	var messagesMap = make(map[target][]myMessage, 0)
	targets_map := gjson.Get(targets, "@this")

	//	l.V(1).Info("*************************MessageDataCollect targets: " + targets)

	targets_map.ForEach(func(key, value gjson.Result) bool {
		//		l.V(1).Info(key.String() + ": " + value.String())
		m := DistributeValue(key.String(), value.String(), inputmessage, nameR, kindR, l)
		//		fmt.Printf("*************************m from DistributeValue: %v\n", m)

		messages = append(messages, m)
		return true // keep iterating
	})
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Target.NameCR < messages[j].Target.NameCR
	})
	var prev_target string
	if len(messages) > 0 {
		prev_target = messages[0].Target.NameCR
	}
	var messagesForCR = make([]myMessage, 0)
	var m myMessage
	for _, m = range messages {
		/* 		fmt.Printf("*************************messagesForCR: %v\n", messagesForCR)
		   		fmt.Printf("*************************   m: %v\n", m)
		   		fmt.Printf("*************************   prev_target: %v\n", prev_target)
		*/if m.Target.NameCR != prev_target {
			messagesMap[m.Target] = messagesForCR
			messagesForCR = make([]myMessage, 0)

		}
		prev_target = m.Target.NameCR
		messagesForCR = append(messagesForCR, m)
	}
	messagesMap[m.Target] = messagesForCR
	return messagesMap
}

func DistributeValue(metric string, targ string, mess string, nameR string, kindR string, l VerbosityLog) myMessage {

	var path []string
	var cr []string
	var nameCR string
	var kindCR string

	if targ[0] == '#' {
		path = strings.Split(targ[1:], ".")
		nameCR = nameR
		kindCR = kindR

	} else {
		pos := strings.Index(targ, "cr:")
		if pos >= 0 {
			cre := strings.Split(targ[pos+3:], "#")
			cr = strings.Split(cre[0], ".")
			path = strings.Split(cre[1], ".")
			nameCR = cr[0]
			kindCR = cr[0]

		}
	}
	//clcm := Cr_type(nameCR, kindCR)
	//if clcm == nil || path == nil {return myMessage{}}

	/* 	fmt.Printf("*************************  nameCR: %v\n", nameCR)
	   	fmt.Printf("*************************  kindCR: %v\n", kindCR)
	   	fmt.Printf("*************************  path: %v\n", path)
	   	fmt.Printf("*************************  mess: %v\n", mess)
	   	fmt.Printf("*************************  metric: %v\n", metric)
	*/
	return myMessage{Target: target{NameCR: nameCR, KindCR: kindCR}, Path: path, Message: mess, Metric: metric}
}

func ApplyToCR(client client.Client, ctx context.Context, req ctrl.Request, mess []myMessage, l VerbosityLog) error {

	//	if monitoring.Spec.Time == monitoring.Spec.MonitoringPolicies.Time {
	//		return nil
	//	}
	if len(mess) == 0 {
		return errors.New("empty messages")
	}
	clcm := Cr_type(mess[0].Target.NameCR, mess[0].Target.KindCR)
	println("+++++++++++++++++++++++mess.Target.NameCR: " + mess[0].Target.NameCR)
	println("+++++++++++++++++++++++mess.Message: " + mess[0].Message)
	if clcm == nil {
		return errors.New("unknown CR")
	}

	v := MyReflectValue{Value: reflect.ValueOf(clcm).Elem()}

	err := client.Get(ctx, types.NamespacedName{Name: req.Name[:strings.LastIndex(req.Name, "-")] + "-" + strings.ToLower(mess[0].Target.NameCR),
		Namespace: req.Namespace}, clcm)

	if err != nil {
		l.Error(err, "Error getting "+req.Name)
	}
	//	client.Get(ctx, types.NamespacedName{Name: v.FieldByName("Spec").FieldByName("Affix").String() + "-decision", Namespace: v.FieldByName("Namespace").String()}, clcm)

	l.V(1).Info("Send Message to CR: " + "Message=" + fmt.Sprintf("%v", mess[0].Message) + " Time=" + v.FieldByName("Spec").FieldByName("Time").String())

	//to improve. iterate all messages and put in proper path
	v.FieldByName("Spec").FieldByName("Message").SetString(mess[0].Message)
	//v.FieldByName("Spec").FieldByName("Time").SetString(time)

	return client.Update(ctx, clcm)

}
