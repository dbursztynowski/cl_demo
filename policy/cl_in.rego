package closedloop.inner

import rego.v1

default cpu := false

cpu if {
    some cpu_idx
    data.Monitoringpolicies.Data[cpu_idx] == "cpu"
    cpu_idx_tresh := sprintf("%v%v", [cpu_idx, "-thresholdvalue"])
    cpu_idx_kind := sprintf("%v%v", [cpu_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[cpu_idx_kind] == "inferior"
    input.cpu > data.Monitoringpolicies.Tresholdvalue[cpu_idx_tresh]
}

cpu if {
    some cpu_idx
    data.Monitoringpolicies.Data[cpu_idx] == "cpu"
    cpu_idx_tresh := sprintf("%v%v", [cpu_idx, "-thresholdvalue"])
    cpu_idx_kind := sprintf("%v%v", [cpu_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[cpu_idx_kind] == "superior"
    input.cpu < data.Monitoringpolicies.Tresholdvalue[cpu_idx_tresh]
}

cpu if {
    some cpu_idx
    data.Monitoringpolicies.Data[cpu_idx] == "cpu"
    cpu_idx_tresh := sprintf("%v%v", [cpu_idx, "-thresholdvalue"])
    cpu_idx_kind := sprintf("%v%v", [cpu_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[cpu_idx_kind] == "uniform"
    input.cpu == data.Monitoringpolicies.Tresholdvalue[cpu_idx_tresh]
}

default memory := false

memory if {
    some memory_idx
    data.Monitoringpolicies.Data[memory_idx] == "memory"
    memory_idx_tresh := sprintf("%v%v", [memory_idx, "-thresholdvalue"])
    memory_idx_kind := sprintf("%v%v", [memory_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[memory_idx_kind] == "inferior"
    input.memory > data.Monitoringpolicies.Tresholdvalue[memory_idx_tresh]
}

memory if {
    some memory_idx
    data.Monitoringpolicies.Data[memory_idx] == "memory"
    memory_idx_tresh := sprintf("%v%v", [memory_idx, "-thresholdvalue"])
    memory_idx_kind := sprintf("%v%v", [memory_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[memory_idx_kind] == "superior"
    input.memory < data.Monitoringpolicies.Tresholdvalue[memory_idx_tresh]
}

memory if {
    some memory_idx
    data.Monitoringpolicies.Data[memory_idx] == "memory"
    memory_idx_tresh := sprintf("%v%v", [memory_idx, "-thresholdvalue"])
    memory_idx_kind := sprintf("%v%v", [memory_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[memory_idx_kind] == "uniform"
    input.memory == data.Monitoringpolicies.Tresholdvalue[memory_idx_tresh]
}

monitoring := {"cpu": cpu, "memory": memory}
default decision := "none"

decision := result if {
    data.Decisionpolicies.Decisiontype == "Priority"
    monitoring.cpu == true
    monitoring.memory == true
    result := data.Decisionpolicies.Priorityspec.Priorityrank["rank-1"]
}

decision := result if {
    data.Decisionpolicies.Decisiontype == "Priority"
    monitoring.cpu == true
    monitoring.memory == false
    result := "cpu"
}

decision := result if {
    data.Decisionpolicies.Decisiontype == "Priority"
    monitoring.cpu == false
    monitoring.memory == true
    result := "memory"
}
