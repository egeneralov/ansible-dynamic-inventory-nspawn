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

func(l List) ToGroups() Answer {
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
	for _, el := range l {
		local := strings.Split(el.Machine, "-")
		if len(local) == 1 {
			// skip host for now
			continue
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
		// append to meta
		r["_meta"].HostVariables[el.Machine] = StringList{
			"ansible_host": el.Addresses,
		}
	}
	return r
}
