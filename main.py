from iterfzf import iterfzf
from tabulate import tabulate
from smap.colors import colorize
from smap.parser import load_ssh_config, ssh
from smap.format import format


fmt = '{BOLD}[H]{NC} |{GREEN}[u]{NC}@{BLUE}[h]{NC}| ([PG])|([LF])'
proxyjump_graph_fmt = '{CURSIVE}[item]{NC}'
tabulated = True

data = load_ssh_config(proxyjump_graph_fmt)

_iter = [
    colorize(format(x, fmt))
    for x in data
]

if tabulated:
    new_data = [
        x.split('|')
        for x in _iter
    ]
    _iter = tabulate(
        new_data,
        tablefmt="plain",
        headers=["Host", "User@system", "Proxyjump", "LocalForward"]
    ).splitlines()

lole = iterfzf(_iter, ansi=True, header_lines=1)
if lole:
    ssh(lole.split(' ')[0])

