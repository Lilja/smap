import re


def format_bracket(
    item, _format,
    remove_if_not_found=False, pre_content='', post_content=''
):
    rules = {
        "[H]": item["host"],
        "[h]": item["hostname"],
        "[u]": item["user"],
        "[P]": item["proxyjump"],
        "[p]": item["port"],
        "[PG]": item["proxyjump_graph"],
        "[LF]": item["localforward"],
    }
    for x in re.findall('\\[\\w+\\]', _format):
        if rules[x]:
            _format = _format.replace(x, pre_content + rules[x] + post_content)
        elif remove_if_not_found:
            _format = _format.replace(x, '')

    return _format


def format(item, _format):
    for x in re.finditer('\\((.*?)\\)', _format):
        full_match = x.group(0)
        contents = x.group(1)
        cd = format_bracket(
            item, contents, True, '(', ')'
        )
        _format = _format.replace(full_match, cd)

    _format = format_bracket(item, _format)

    return _format
