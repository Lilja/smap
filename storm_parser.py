# -*- coding: utf-8 -*-

from glob import glob
from os.path import expanduser
from os.path import exists
from os.path import isabs
from os.path import join as os_join
from operator import itemgetter
import re

from paramiko.config import SSHConfig
import six


class StormConfig(SSHConfig):
    def parse(self, file_obj):
        """
        Read an OpenSSH config from the given file object.

        @param file_obj: a file-like object to read the config file from
        @type file_obj: file
        """
        order = 1
        host = {"host": ['*'], "config": {}, }
        for line in file_obj:
            line = line.rstrip('\n').lstrip()
            if line == '':
                self._config.append({
                    'type': 'empty_line',
                    'value': line,
                    'host': '',
                    'order': order,
                })
                order += 1
                continue

            if line.startswith('#'):
                self._config.append({
                    'type': 'comment',
                    'value': line,
                    'host': '',
                    'order': order,
                })
                order += 1
                continue

            if '=' in line:
                # Ensure ProxyCommand gets properly split
                if line.lower().strip().startswith('proxycommand'):
                    proxy_re = re.compile(r"^(proxycommand)\s*=*\s*(.*)", re.I)
                    match = proxy_re.match(line)
                    key, value = match.group(1).lower(), match.group(2)
                else:
                    key, value = line.split('=', 1)
                    key = key.strip().lower()
            else:
                # find first whitespace, and split there
                i = 0
                while (i < len(line)) and not line[i].isspace():
                    i += 1
                if i == len(line):
                    raise Exception('Unparsable line: %r' % line)
                key = line[:i].lower()
                value = line[i:].lstrip()
            if key == 'host':
                self._config.append(host)
                value = value.split()
                host = {
                    key: value,
                    'config': {},
                    'type': 'entry',
                    'order': order
                }
                order += 1
            elif key in ['identityfile', 'localforward', 'remoteforward']:
                if key in host['config']:
                    host['config'][key].append(value)
                else:
                    host['config'][key] = [value]
            elif key not in host['config']:
                host['config'].update({key: value})
        self._config.append(host)


class ConfigParser(object):
    """
    Config parser for ~/.ssh/config files.
    """

    def __init__(self, ssh_config_file=None):
        if not ssh_config_file:
            ssh_config_file = self.get_default_ssh_config_file()

        self.extra_config_files = self.recursively_find_config_files(ssh_config_file, [])

        self.defaults = {}

        self.ssh_config_file = ssh_config_file

        self.config_data = []

    def get_default_ssh_config_file(self):
        return expanduser("~/.ssh/config")

    def _absify(self, path):
        if not isabs(path):
            return os_join(
                expanduser("~/.ssh"),
                path
            )
        else:
            return path

    def _find_config_files(self, o):
        for line in o.read().splitlines():
            if "Include" in line:
                path = self._absify(line.split(" ")[1].strip())
                if '*' in path:
                    for g in glob(path):
                        yield g
                else:
                    yield path

    def recursively_find_config_files(self, config_file, parsed_files):
        if not exists(config_file):
            raise ValueError(
                f"Cannot parse {config_file}, it does not exist"
            )
        with open(config_file, 'r') as o:
            if config_file in parsed_files:
                raise ValueError(
                    f"Cannot parse {config_file}, it's already been parsed. "
                    "Check your ssh config for recursion."
                )
            parsed_files.append(config_file)
            _files = list(self._find_config_files(o))
            new_files = []
            for file in _files:
                new_files.append(
                    self.recursively_find_config_files(file, parsed_files)
                )

            if new_files:
                return new_files
            else:
                return config_file

    def load(self):
        config = StormConfig()

        with open(self.ssh_config_file) as fd:
            config.parse(fd)

        for config_file in self.extra_config_files:
            with open(config_file) as fd:
                config.parse(fd)

        for entry in config.__dict__.get("_config"):
            if entry.get("host") == ["*"]:
                self.defaults.update(entry.get("config"))

            if entry.get("type") in ["comment", "empty_line"]:
                self.config_data.append(entry)
                continue

            host_item = {
                'host': entry["host"][0],
                'options': entry.get("config"),
                'type': 'entry',
                'order': entry.get("order", 0),
            }

            if len(entry["host"]) > 1:
                host_item.update({
                    'host': " ".join(entry["host"]),
                })

            # minor bug in paramiko.SSHConfig that duplicates
            # "Host *" entries.
            if entry.get("config") and len(entry.get("config")) > 0:
                self.config_data.append(host_item)

        return self.config_data

    def search_host(self, search_string):
        results = []
        for host_entry in self.config_data:
            if host_entry.get("type") != 'entry':
                continue
            if host_entry.get("host") == "*":
                continue

            searchable_information = host_entry.get("host")
            for key, value in six.iteritems(host_entry.get("options")):
                if isinstance(value, list):
                    value = " ".join(value)
                if isinstance(value, int):
                    value = str(value)

                searchable_information += " " + value

            if search_string in searchable_information:
                results.append(host_entry)

        return results

    def dump(self):
        if len(self.config_data) < 1:
            return

        file_content = ""
        self.config_data = sorted(self.config_data, key=itemgetter("order"))

        for host_item in self.config_data:
            if host_item.get("type") in ['comment', 'empty_line']:
                file_content += host_item.get("value") + "\n"
                continue
            host_item_content = "Host {0}\n".format(host_item.get("host"))
            for key, value in six.iteritems(host_item.get("options")):
                if isinstance(value, list):
                    sub_content = ""
                    for value_ in value:
                        sub_content += "    {0} {1}\n".format(
                            key, value_
                        )
                    host_item_content += sub_content
                else:
                    host_item_content += "    {0} {1}\n".format(
                        key, value
                    )
            file_content += host_item_content

        return file_content
