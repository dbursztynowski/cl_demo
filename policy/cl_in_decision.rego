package closedloop.inner.decision1

import rego.v1

default decision := "none"
# default monitoring := {"cpu":false,"memory":false}

default monitoring.cpu := false
monitoring.cpu := input.cpu
default monitoring.memory := false
monitoring.memory := input.memory

to_bool("true"):=true
to_bool(true):=true
to_bool(value):=false if{
	value != "true"
    value != true
}
to_bool(false):=false

# to_bool(value):=result if{
#	result :=(value == "true")
# } 

decision := result if {
    data.Decisionpolicies.Decisiontype == "Priority"
    to_bool(monitoring.cpu) == true
    to_bool(monitoring.memory) == true
    result := data.Decisionpolicies.Priorityspec.Priorityrank["rank-1"]
}

decision := result if {
    data.Decisionpolicies.Decisiontype == "Priority"
    to_bool(monitoring.cpu) == true
    to_bool(monitoring.memory) == false
    result := "cpu"
}

decision := result if {
    data.Decisionpolicies.Decisiontype == "Priority"
    to_bool(monitoring.cpu) == false
    to_bool(monitoring.memory) == true
    result := "memory"
}

