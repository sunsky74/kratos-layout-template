import os
from subprocess import Popen
from cookiecutter.main import cookiecutter

PROJECT_DIRECTORY = os.path.realpath(os.path.curdir)
app_name = '{{ cookiecutter.project_name }}'

def init_protoc():
    """
    Initialises protoc on the new project folder
    """
    PROTOC_COMMANDS = [
        ["make", "api"],
    ]

    for command in PROTOC_COMMANDS:
        w = Popen(command, cwd=os.path.join(PROJECT_DIRECTORY))
        w.wait()

def init_config():
    """
    Initialises configs on the new project folder
    """
    PROTOC_COMMANDS = [
        ["make", "config"],
    ]

    for command in PROTOC_COMMANDS:
        w = Popen(command, cwd=os.path.join(PROJECT_DIRECTORY))
        w.wait()

def init_wire():
    """
    Initialises wire on the new project folder
    """
    WIRE_COMMANDS = [
        ["wire", "."],
    ]

    for command in WIRE_COMMANDS:
        w = Popen(command, cwd=os.path.join(PROJECT_DIRECTORY,'cmd',app_name))
        w.wait()

init_protoc()
init_config()
init_wire()
