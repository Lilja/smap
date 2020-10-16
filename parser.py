import storm_parser
import shutil
import os


def parse(item):
    return {
        "host": item["host"],
        "user": item["options"].get("user"),
        "port": item["options"].get("port"),
        "hostname": item["options"].get("hostname"),
        "proxyjump": item["options"].get("proxyjump"),
        "localforward": format_localforward_text(item["options"].get("localforward")),
    }


def add_proxyjump_graph(item, items, proxyjump_graph_fmt):
    proxyjump_graph = get_proxyjump_graph(items, item)
    formatted_proxyjump_graph = format_proxyjump_graph(
        proxyjump_graph, proxyjump_graph_fmt
    )

    item["proxyjump_graph"] = formatted_proxyjump_graph

    return item


def format_localforward_text(fw):
    if fw:
        if isinstance(fw, list):
            forwards = []
            for forward in fw:
                c = forward.split(" ")
                c.insert(1, " ↔ ")
                forwards.append("".join(c))
            return " ".join(forwards)
        else:
            c = fw.split(" ")
            c.insert(1, " ↔ ")
            return "".join(c)


def find_host(items, _host):
    for item in items:
        if item["host"] == _host:
            return item


def get_proxyjump_graph(_items, _item):
    def recurision(items, item):
        if not item["proxyjump"]:
            return [item["host"]]
        else:
            new_host = find_host(items, item["proxyjump"])
            return [item["host"]] + recurision(items, new_host)

    return recurision(_items, _item)[::-1]


def format_proxyjump_graph(p, proxyjump_graph_fmt):
    if len(p) == 1:
        return ''
    else:
        p = [
            proxyjump_graph_fmt.replace('[item]', x)
            for x in p
        ]
        return " → ".join(p)


def load(proxyjump_graph_fmt):
    data = list([
        x
        for x in storm_parser.ConfigParser().load()
        if x['type'] == 'entry' and x['host'] != '*'
    ])

    data = [
        parse(x)
        for x in data
    ]
    data = [
        add_proxyjump_graph(x, data, proxyjump_graph_fmt)
        for x in data
    ]
    return sorted(data, key=lambda x: x["host"], reverse=True)


def ssh(host):
    cmd = shutil.which('ssh')
    os.execv(cmd, [cmd, host])
