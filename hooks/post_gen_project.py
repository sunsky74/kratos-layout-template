import os
from subprocess import Popen
from cookiecutter.main import cookiecutter

PROJECT_DIRECTORY = os.path.realpath(os.path.curdir)
app_name = '{{ cookiecutter.project_name }}'


def init_dir():
    """
    Initialises folder on the new project
    """
    dirs_to_create = [
        "logs",
        "docs",
        "docs/api",
        "docs/wiki",
    ]

    for dir_path in dirs_to_create:
        full_path = os.path.join(PROJECT_DIRECTORY, dir_path)
        try:
            os.makedirs(full_path, exist_ok=True)  # exist_ok=True 避免重复创建报错
            print(f"Created directory: {full_path}")
        except OSError as e:
            print(f"Error creating {full_path}: {e}")

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

init_dir()
init_protoc()
init_config()
init_wire()
