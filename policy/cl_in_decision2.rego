package closedloop.inner.decision

import rego.v1

default metric := "none"
# default monitoring := {"cpu":false,"memory":false}

default monitoring.cpu := false
monitoring.cpu := input.cpu
default monitoring.memory := false
monitoring.memory := input.memory

to_bool("high") := true
to_bool("low") := true
to_bool("true") := true
to_bool(true) := true

to_bool(value):=false if{
	value != "true"
    value != true
    value != "low"
    value != "high"
}
to_bool(false) := false

metric := result if {
    data.Decisionpolicies.Decisiontype == "Priority"
    to_bool(monitoring.cpu) == true
    to_bool(monitoring.memory) == true
    result := data.Decisionpolicies.Priorityspec.Priorityrank["rank-1"]
}

metric := result if {
    data.Decisionpolicies.Decisiontype == "Priority"
    to_bool(monitoring.cpu) == true
    to_bool(monitoring.memory) == false
    result := "cpu"
}

metric := result if {
    data.Decisionpolicies.Decisiontype == "Priority"
    to_bool(monitoring.cpu) == false
    to_bool(monitoring.memory) == true
    result := "memory"
}

action := result if {
    input[metric] == "low"
    result := "increase"
}

action := result if {
    input[metric] == "high"
    result := "decrease"
}

action := result if {
    input[metric] == "ok"
    result := "ok"
}

