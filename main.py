from iterfzf import iterfzf
from colors import colorize
from parser import load, ssh
from format import format


fmt = '{BOLD}[H]{NC} {GREEN}[u]{NC}@{BLUE}[h]{NC} ([PG])([LF])'
proxyjump_graph_fmt = '{CURSIVE}[item]{NC}'

data = load(proxyjump_graph_fmt)

_iter = [
    colorize(format(x, fmt))
    for x in data
]


lole = iterfzf(_iter, ansi=True)
if lole:
    ssh(lole.split(' ')[0])

