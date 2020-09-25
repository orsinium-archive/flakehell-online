from pathlib import Path
from textwrap import dedent

from astroid import MANAGER
from flakehell._patched import FlakeHellApplication
from flakehell.formatters import JSONFormatter

# save flakehell config
path = Path("pyproject.toml").write_text(config)  # noqa: F821

# save source code
path = Path("code.py")
path.write_text(dedent(text))  # noqa: F821


# clear astroid (pylint) cache
MANAGER.astroid_cache.clear()


class Formatter(JSONFormatter):
    def after_init(self):
        self._out = []
        return super().after_init()

    def _write(self, output: str) -> None:
        self._out.append(output)


class App(FlakeHellApplication):
    def make_formatter(self, formatter_class=None):
        self.formatter = Formatter(self.options)


app = App()
code = 0
try:
    app.run(["code.py"])
    app.exit()
except SystemExit as err:
    code = int(err.args[0])
