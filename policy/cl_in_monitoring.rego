package closedloop.inner.monitoring

import rego.v1

default cpu := false

cpu if {
    some cpu_idx
    data.Monitoringpolicies.Data[cpu_idx] == "cpu"
    cpu_idx_tresh := sprintf("%v%v", [cpu_idx, "-thresholdvalue"])
    cpu_idx_kind := sprintf("%v%v", [cpu_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[cpu_idx_kind] == "inferior"
    to_number(input.cpu) < to_number(data.Monitoringpolicies.Tresholdvalue[cpu_idx_tresh])
}

cpu if {
    some cpu_idx
    data.Monitoringpolicies.Data[cpu_idx] == "cpu"
    cpu_idx_tresh := sprintf("%v%v", [cpu_idx, "-thresholdvalue"])
    cpu_idx_kind := sprintf("%v%v", [cpu_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[cpu_idx_kind] == "superior"
    to_number(input.cpu) > to_number(data.Monitoringpolicies.Tresholdvalue[cpu_idx_tresh])
}

cpu if {
    some cpu_idx
    data.Monitoringpolicies.Data[cpu_idx] == "cpu"
    cpu_idx_tresh := sprintf("%v%v", [cpu_idx, "-thresholdvalue"])
    cpu_idx_kind := sprintf("%v%v", [cpu_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[cpu_idx_kind] == "uniform"
    to_number(input.cpu) == to_number(data.Monitoringpolicies.Tresholdvalue[cpu_idx_tresh])
}

default memory := false

memory if {
    some memory_idx
    data.Monitoringpolicies.Data[memory_idx] == "memory"
    memory_idx_tresh := sprintf("%v%v", [memory_idx, "-thresholdvalue"])
    memory_idx_kind := sprintf("%v%v", [memory_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[memory_idx_kind] == "inferior"
    to_number(input.memory) < to_number(data.Monitoringpolicies.Tresholdvalue[memory_idx_tresh])
}

memory if {
    some memory_idx
    data.Monitoringpolicies.Data[memory_idx] == "memory"
    memory_idx_tresh := sprintf("%v%v", [memory_idx, "-thresholdvalue"])
    memory_idx_kind := sprintf("%v%v", [memory_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[memory_idx_kind] == "superior"
    to_number(input.memory) > to_number(data.Monitoringpolicies.Tresholdvalue[memory_idx_tresh])
}

memory if {
    some memory_idx
    data.Monitoringpolicies.Data[memory_idx] == "memory"
    memory_idx_tresh := sprintf("%v%v", [memory_idx, "-thresholdvalue"])
    memory_idx_kind := sprintf("%v%v", [memory_idx, "-thresholdkind"])
    data.Monitoringpolicies.Tresholdkind[memory_idx_kind] == "uniform"
    to_number(input.memory) == to_number(data.Monitoringpolicies.Tresholdvalue[memory_idx_tresh])
}

