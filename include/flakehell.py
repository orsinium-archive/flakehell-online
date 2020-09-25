from pathlib import Path
from random import choice
from string import ascii_lowercase
from textwrap import dedent

from flakehell._patched import FlakeHellApplication
from flakehell.formatters import JSONFormatter


def random_name():
    name = ''
    for _ in range(20):
        name += choice(ascii_lowercase)
    return name + '.py'


# save flakehell config
path = Path("pyproject.toml").write_text(config)  # noqa: F821

# save source code
path = Path(random_name())
path.write_text(dedent(text))  # noqa: F821


class Formatter(JSONFormatter):
    def after_init(self):
        self._out = []
        return super().after_init()

    def _write(self, output: str) -> None:
        self._out.append(output)


class App(FlakeHellApplication):
    def make_formatter(self, formatter_class=None):
        self.formatter = Formatter(self.options)


# run flakehell
app = App()
code = 0
try:
    app.run([str(path)])
    app.exit()
except SystemExit as err:
    code = int(err.args[0])

# remove file
path.unlink()
