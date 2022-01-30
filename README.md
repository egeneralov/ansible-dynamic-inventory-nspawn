# ansible-dynamic-inventory-nspawn

## example usage

- `export NSPAWN_INVENTORY_COMMAND_LIST="ssh home machinectl list -o json" DEBUG=0`
- `ansible-inventory -i ~/go/bin/ansible-dynamic-inventory-nspawn --list`

```json
{
    "_meta": {
        "hostvars": {
            "etcd-0": {
                "ansible_host": "192.168.88.61"
            },
            "master-0": {
                "ansible_host": "192.168.88.56"
            }
        }
    },
    "all": {
        "children": [
            "etcd",
            "master",
            "ungrouped"
        ]
    },
    "etcd": {
        "hosts": [
            "etcd-0"
        ]
    },
    "master": {
        "hosts": [
            "master-0"
        ]
    }
}
```

- `ansible -i ~/go/bin/ansible-dynamic-inventory-nspawn -m ping all`

```text
master-0 | SUCCESS => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python"
    },
    "changed": false,
    "ping": "pong"
}
etcd-0 | SUCCESS => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python"
    },
    "changed": false,
    "ping": "pong"
}
```

## inside

### configuration

environment variables:

- **NSPAWN_INVENTORY_COMMAND_LIST** - you can override execution command, for example `ssh -l user -p 22 -i /tmp/id_rsa 1.1.1.2 machinectl list -o json`
- **DEBUG** if eq `1` - logs will be present to stderr

### what happens after launch

Executing command `machinectl list -o json` (or overrided), stdout must be like:

```json
[
  {
    "machine": "etcd-0",
    "class": "container",
    "service": "systemd-nspawn",
    "os": "arch",
    "version": null,
    "addresses": "192.168.88.61"
  }
]
```

And packs in *ansible-inventory* format:

```json
{
  "_meta": {
    "hostvars": {
      "etcd-0": {
        "ansible_host": "192.168.88.61"
      }
    }
  },
  "etcd": {
    "hosts": [
      "etcd-0"
    ]
  }
}
```
