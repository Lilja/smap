from iterfzf import iterfzf
from smap.colors import colorize
from smap.parser import load_ssh_config, ssh
from smap.format import format


fmt = '{BOLD}[H]{NC} {GREEN}[u]{NC}@{BLUE}[h]{NC} ([PG])([LF])'
proxyjump_graph_fmt = '{CURSIVE}[item]{NC}'

data = load_ssh_config(proxyjump_graph_fmt)

_iter = [
    colorize(format(x, fmt))
    for x in data
]


lole = iterfzf(_iter, ansi=True)
if lole:
    ssh(lole.split(' ')[0])

