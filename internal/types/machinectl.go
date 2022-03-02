package types

import "strings"

type (
	MachineInfo struct {
		Machine   string      `json:"machine"`
		Class     string      `json:"class"`
		Service   string      `json:"service"`
		Os        string      `json:"os"`
		Version   interface{} `json:"version"`
		Addresses string      `json:"addresses"`
	}
	List []MachineInfo
)

func (l List) ToGroups(bastion string) Answer {
	if len(l) < 1 {
		return nil
	}
	var (
		r = map[string]Group{
			"_meta": {
				HostVariables: map[string]StringList{},
			},
		}
	)
	if bastion != "" {
		r["bastion"] = Group{
			Hosts: []string{"bastion"},
		}
		r["_meta"].HostVariables["bastion"] = StringList{
			"ansible_host": bastion,
		}
	}
	for _, el := range l {
		local := strings.Split(el.Machine, "-")
		if len(local) == 1 {
			local = strings.Split(el.Machine, ".")
		}
		// append to group
		groupName := local[0]
		if group, ok := r[groupName]; !ok {
			r[groupName] = Group{
				Hosts: []string{el.Machine},
			}
		} else {
			group.Hosts = append(group.Hosts, el.Machine)
			r[groupName] = group
		}
		// generate metadata
		var machineMetadata = StringList{}
		if s := strings.Split(el.Addresses, "\n"); len(s) > 1 {
			machineMetadata["ansible_host"] = s[0]
		} else {
			machineMetadata["ansible_host"] = el.Addresses
		}
		if bastion != "" {
			machineMetadata["ansible_ssh_common_args"] = "-o ProxyCommand=\"ssh -W %h:%p -q " + bastion + "\""
		}
		r["_meta"].HostVariables[el.Machine] = machineMetadata
	}
	return r
}
